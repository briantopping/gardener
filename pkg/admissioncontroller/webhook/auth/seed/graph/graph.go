// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package graph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	gonumgraph "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/traverse"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	operationsv1alpha1 "github.com/gardener/gardener/pkg/apis/operations/v1alpha1"
	securityv1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
	seedmanagementv1alpha1 "github.com/gardener/gardener/pkg/apis/seedmanagement/v1alpha1"
)

// Interface is used to track resources dependencies.
type Interface interface {
	// Setup registers the event handler functions for the various resource types.
	Setup(ctx context.Context, c cache.Cache) error
	// HasVertex returns true when the given vertex exists in the graph.
	HasVertex(vertexType VertexType, vertexNamespace, vertexName string) bool
	// HasPathFrom returns true when there is a path from <from> to <to>.
	HasPathFrom(fromType VertexType, fromNamespace, fromName string, toType VertexType, toNamespace, toName string) bool
}

type graph struct {
	lock     sync.RWMutex
	logger   logr.Logger
	client   client.Client
	graph    *simple.DirectedGraph
	vertices typeVertexMapping
}

var _ Interface = &graph{}

// New creates a new graph interface for tracking resource dependencies.
func New(logger logr.Logger, client client.Client) *graph {
	return &graph{
		logger:   logger,
		client:   client,
		graph:    simple.NewDirectedGraph(),
		vertices: make(typeVertexMapping),
	}
}

func (g *graph) Setup(ctx context.Context, c cache.Cache) error {
	for _, resource := range []struct {
		obj     client.Object
		setupFn func(context.Context, cache.Informer) error
	}{
		{&gardencorev1beta1.BackupBucket{}, g.setupBackupBucketWatch},
		{&gardencorev1beta1.BackupEntry{}, g.setupBackupEntryWatch},
		{&operationsv1alpha1.Bastion{}, g.setupBastionWatch},
		{&certificatesv1.CertificateSigningRequest{}, g.setupCertificateSigningRequestWatch},
		{&gardencorev1beta1.ControllerInstallation{}, g.setupControllerInstallationWatch},
		{&seedmanagementv1alpha1.Gardenlet{}, g.setupGardenletWatch},
		{&seedmanagementv1alpha1.ManagedSeed{}, g.setupManagedSeedWatch},
		{&gardencorev1beta1.Project{}, g.setupProjectWatch},
		{&gardencorev1beta1.SecretBinding{}, g.setupSecretBindingWatch},
		{&gardencorev1beta1.Seed{}, g.setupSeedWatch},
		{&corev1.ServiceAccount{}, g.setupServiceAccountWatch},
		{&gardencorev1beta1.Shoot{}, g.setupShootWatch},
		{&securityv1alpha1.CredentialsBinding{}, g.setupCredentialsBindingWatch},
	} {
		informer, err := c.GetInformer(ctx, resource.obj)
		if err != nil {
			return err
		}
		if err := resource.setupFn(ctx, informer); err != nil {
			return err
		}
	}

	return nil
}

func (g *graph) HasVertex(vertexType VertexType, vertexNamespace, vertexName string) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()
	_, ok := g.getVertex(vertexType, vertexNamespace, vertexName)
	return ok
}

func (g *graph) HasPathFrom(fromType VertexType, fromNamespace, fromName string, toType VertexType, toNamespace, toName string) bool {
	start := time.Now()
	defer func() {
		metricPathCheckDuration.WithLabelValues(vertexTypes[fromType], vertexTypes[toType]).Observe(time.Since(start).Seconds())
	}()
	g.lock.RLock()
	defer g.lock.RUnlock()

	fromVertex, ok := g.getVertex(fromType, fromNamespace, fromName)
	if !ok {
		return false
	}

	toVertex, ok := g.getVertex(toType, toNamespace, toName)
	if !ok {
		return false
	}

	return (&traverse.DepthFirst{}).Walk(g.graph, fromVertex, func(n gonumgraph.Node) bool {
		return n.ID() == toVertex.ID()
	}) != nil
}

func (g *graph) getOrCreateVertex(vertexType VertexType, namespace, name string) *vertex {
	if v, ok := g.getVertex(vertexType, namespace, name); ok {
		return v
	}
	return g.createVertex(vertexType, namespace, name)
}

func (g *graph) getVertex(vertexType VertexType, namespace, name string) (*vertex, bool) {
	v, ok := g.vertices[vertexType][namespace][name]
	return v, ok
}

func (g *graph) createVertex(vertexType VertexType, namespace, name string) *vertex {
	typedVertices, ok := g.vertices[vertexType]
	if !ok {
		typedVertices = namespaceVertexMapping{}
		g.vertices[vertexType] = typedVertices
	}

	namespacedVertices, ok := typedVertices[namespace]
	if !ok {
		namespacedVertices = map[string]*vertex{}
		typedVertices[namespace] = namespacedVertices
	}

	v := newVertex(vertexType, namespace, name, g.graph.NewNode().ID())
	namespacedVertices[name] = v

	g.graph.AddNode(v)
	g.logger.Info(
		"Added",
		"vertex", fmt.Sprintf("%s (%d)", v, v.ID()),
	)

	return v
}

func (g *graph) deleteVertex(vertexType VertexType, namespace, name string) {
	v, ok := g.getVertex(vertexType, namespace, name)
	if !ok {
		return
	}

	// Now, visit all neighbors of <v> and check if they can also be removed now that <v> will be removed.
	verticesToRemove := []gonumgraph.Node{v}

	// Neighbors to which <v> has an outgoing edge can also be removed if they do not have any outgoing edges (to other
	// vertices) themselves and if they only have one incoming edge (which must be the edge from <v>).
	g.visit(g.graph.From(v.ID()), func(neighbor gonumgraph.Node) {
		if g.graph.From(neighbor.ID()).Len() == 0 && g.graph.To(neighbor.ID()).Len() == 1 {
			verticesToRemove = append(verticesToRemove, neighbor)
		}
	})

	// Neighbors from which <v> has an incoming edge can also be removed if they do not have any incoming edges (from
	// other vertices) themselves and if they only have one outgoing edge (which must be the edge to <v>).
	g.visit(g.graph.To(v.ID()), func(neighbor gonumgraph.Node) {
		if g.graph.To(neighbor.ID()).Len() == 0 && g.graph.From(neighbor.ID()).Len() == 1 {
			verticesToRemove = append(verticesToRemove, neighbor)
		}
	})

	for _, v := range verticesToRemove {
		g.removeVertex(v.(*vertex))
	}
}

func (g *graph) deleteVertexIfIsolated(v *vertex) {
	if g.graph.From(v.ID()).Len() == 0 && g.graph.To(v.ID()).Len() == 0 {
		g.removeVertex(v)
	}
}

func (g *graph) removeVertex(v *vertex) {
	g.graph.RemoveNode(v.ID())
	delete(g.vertices[v.vertexType][v.namespace], v.name)
	if len(g.vertices[v.vertexType][v.namespace]) == 0 {
		delete(g.vertices[v.vertexType], v.namespace)
	}
	g.logger.Info(
		"Removed (with all associated edges)",
		"vertex", fmt.Sprintf("%s (%d)", v, v.ID()),
	)
}

func (g *graph) addEdge(from, to *vertex) {
	g.graph.SetEdge(g.graph.NewEdge(from, to))
	g.logger.Info(
		"Added edge",
		"from", fmt.Sprintf("%s (%d)", from, from.ID()),
		"to", fmt.Sprintf("%s (%d)", to, to.ID()),
	)
}

func (g *graph) deleteAllIncomingEdges(fromVertexType, toVertexType VertexType, toNamespace, toName string) {
	to, ok := g.getVertex(toVertexType, toNamespace, toName)
	if !ok {
		return
	}

	// Now, visit all neighbors of <to> who have an incoming edge to <to> and check whether these vertices can be
	// removed as well.
	var (
		verticesToRemove []gonumgraph.Node
		edgesToRemove    []gonumgraph.Edge
	)

	// Delete all edges from vertices of type <fromVertexType> to <to>. Neighbors from which <to> has an incoming edge
	// can also be removed if they do not have any incoming edges (from other vertices) themselves and if they only have
	// one outgoing edge (which must be the edge to <to>).
	g.visit(g.graph.To(to.ID()), func(neighbor gonumgraph.Node) {
		from, ok := neighbor.(*vertex)
		if !ok || from.vertexType != fromVertexType {
			return
		}

		if g.graph.To(neighbor.ID()).Len() == 0 && g.graph.From(neighbor.ID()).Len() == 1 {
			verticesToRemove = append(verticesToRemove, neighbor)
		} else {
			edgesToRemove = append(edgesToRemove, g.graph.Edge(from.ID(), to.ID()))
		}
	})

	for _, v := range verticesToRemove {
		g.removeVertex(v.(*vertex))
	}

	for _, e := range edgesToRemove {
		g.removeEdge(e)
	}

	// If <to> is now isolated, i.e., has neither incoming nor outgoing edges, then we can delete the vertex as well.
	g.deleteVertexIfIsolated(to)
}

func (g *graph) deleteAllOutgoingEdges(fromVertexType VertexType, fromNamespace, fromName string, toVertexType VertexType) {
	from, ok := g.getVertex(fromVertexType, fromNamespace, fromName)
	if !ok {
		return
	}

	// Now, visit all neighbors of <from> who have an outgoing edge from <from> and check whether these vertices can be
	// removed as well.
	var (
		verticesToRemove []gonumgraph.Node
		edgesToRemove    []gonumgraph.Edge
	)

	// Delete all edges from <from> to vertices of type <toVertexType>. Neighbors to which <from> has an outgoing edge
	// can also be removed if they do not have any outgoing edges (to other vertices) themselves and if they only have
	// one incoming edge (which must be the edge from <from>).
	g.visit(g.graph.From(from.ID()), func(neighbor gonumgraph.Node) {
		to, ok := neighbor.(*vertex)
		if !ok || to.vertexType != toVertexType {
			return
		}

		if g.graph.From(neighbor.ID()).Len() == 0 && g.graph.To(neighbor.ID()).Len() == 1 {
			verticesToRemove = append(verticesToRemove, neighbor)
		} else {
			edgesToRemove = append(edgesToRemove, g.graph.Edge(from.ID(), to.ID()))
		}
	})

	for _, v := range verticesToRemove {
		g.removeVertex(v.(*vertex))
	}

	for _, e := range edgesToRemove {
		g.removeEdge(e)
	}

	// If <from> is now isolated, i.e., has neither incoming nor outgoing edges then we can delete the vertex as well.
	g.deleteVertexIfIsolated(from)
}

func (g *graph) removeEdge(edge gonumgraph.Edge) {
	g.graph.RemoveEdge(edge.From().ID(), edge.To().ID())
	g.logger.Info(
		"Removed edge",
		"from", fmt.Sprintf("%s (%d)", edge.From(), edge.From().ID()),
		"to", fmt.Sprintf("%s (%d)", edge.To(), edge.To().ID()),
	)
}

func (g *graph) visit(nodes gonumgraph.Nodes, visitor func(gonumgraph.Node)) {
	for nodes.Next() {
		if node := nodes.Node(); node != nil {
			visitor(node)
		}
	}
}
