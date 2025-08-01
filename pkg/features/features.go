// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package features

import (
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/component-base/featuregate"
)

const (
	// Every feature gate should add method here following this template:
	//
	// // MyFeature enable Foo.
	// // owner: @username
	// // alpha: v5.X
	// MyFeature featuregate.Feature = "MyFeature"

	// DefaultSeccompProfile defaults the seccomp profile for Gardener managed workload in the seed to RuntimeDefault.
	// owner: @dimityrmirchev
	// alpha: v1.54.0
	DefaultSeccompProfile featuregate.Feature = "DefaultSeccompProfile"

	// UseNamespacedCloudProfile enables the usage of the NamespacedCloudProfile API object
	// nodes.
	// owner: @timuthy @benedictweis @LucaBernstein
	// alpha: v1.92.0
	// beta: v1.112.0
	// GA: v1.125.0
	UseNamespacedCloudProfile featuregate.Feature = "UseNamespacedCloudProfile"

	// ShootCredentialsBinding enables the usage of the CredentialsBindingName API in shoot spec.
	// owner: @vpnachev @dimityrmirchev
	// alpha: v1.98.0
	// beta: v1.107.0
	ShootCredentialsBinding featuregate.Feature = "ShootCredentialsBinding"

	// NewWorkerPoolHash enables a new calculation method for the worker pool hash. The new
	// calculation supports rolling worker pools if `kubeReserved`, `systemReserved`, `evictionHard` or `cpuManagerPolicy`
	// in the `kubelet` configuration are changed. All provider extensions must be upgraded
	// to support this feature first.
	// owner: @MichaelEischer
	// alpha: v1.98.0
	NewWorkerPoolHash featuregate.Feature = "NewWorkerPoolHash"

	// NewVPN enables the new implementation of the VPN (go rewrite) using an IPv6 transfer network.
	// owner: @MartinWeindel @ScheererJ @axel7born @DockToFuture
	// alpha: v1.104.0
	// beta: v1.115.0
	// GA: v1.116.0
	NewVPN featuregate.Feature = "NewVPN"

	// CredentialsRotationWithoutWorkersRollout enables starting the credentials rotation without immediately causing
	// a rolling update of all worker nodes. Instead, the rolling update can be triggered manually by the user at a
	// later point in time of their convenience.
	// owner: @rfranzke
	// alpha: v1.112.0
	// beta: v1.121.0
	CredentialsRotationWithoutWorkersRollout featuregate.Feature = "CredentialsRotationWithoutWorkersRollout"

	// InPlaceNodeUpdates enables setting the update strategy of worker pools to `AutoInPlaceUpdate` or `ManualInPlaceUpdate` in the Shoot API.
	// owner: @acumino @ary1992 @shafeeqes
	// alpha: v1.113.0
	InPlaceNodeUpdates featuregate.Feature = "InPlaceNodeUpdates"

	// IstioTLSTermination enables TLS termination for the Istio Ingress Gateway instead of TLS termination at the kube-apiserver.
	// It allows load-balancing of requests to the kube-apiserver on request level instead of connection level.
	// owner: @oliver-goetz
	// alpha: v1.114.0
	IstioTLSTermination featuregate.Feature = "IstioTLSTermination"

	// CloudProfileCapabilities enables the usage of capabilities in the CloudProfile. Capabilities are used to create a relation between
	// machineTypes and machineImages. It allows to validate worker groups of a shoot ensuring the selected image and machine combination
	// will boot up successfully. Capabilities are also used to determine valid upgrade paths during automated maintenance operations.
	// owner: @roncossek
	// alpha: v1.117.0
	CloudProfileCapabilities featuregate.Feature = "CloudProfileCapabilities"

	// DoNotCopyBackupCredentials disables the copying of Shoot infrastructure credentials as backup credentials when the Shoot is used as a ManagedSeed.
	// Operators are responsible for providing the credentials for backup explicitly.
	// Credentials that were already copied will be labeled with "secret.backup.gardener.cloud/status=previously-managed" and would have to be cleaned up by operators.
	// owner: @dimityrmirchev
	// alpha: v1.121.0
	// beta: v1.123.0
	DoNotCopyBackupCredentials featuregate.Feature = "DoNotCopyBackupCredentials"

	// OpenTelemetryCollector enables the usage of an OpenTelemetry Collector instance in the Control Plane of Shoot clusters.
	// All logs will be routed through the Collector before they reach the Vali instance.
	// owner: @rrhubenov
	// alpha: v1.124.0
	OpenTelemetryCollector featuregate.Feature = "OpenTelemetryCollector"
)

// DefaultFeatureGate is the central feature gate map used by all gardener components.
// On startup, the component needs to register all feature gates that are available for this component via `Add`, e.g.:
//
//	 utilruntime.Must(features.DefaultFeatureGate.Add(features.GetFeatures(
//			features.MyFeatureGateName,
//		)))
//
// With this, every component has its individual set of available feature gates (different to Kubernetes, where all
// components have all feature gates even if irrelevant).
// Additionally, the component needs to set the feature gates' states based on the operator's configuration, e.g.:
//
//	features.DefaultFeatureGate.SetFromMap(o.config.FeatureGates)
//
// For checking whether a given feature gate is enabled (regardless of which component the code is executed in), use:
//
//	features.DefaultFeatureGate.Enabled(features.DefaultSeccompProfile)
//
// With this, code that needs to check a given feature gate's state can be shared across components, e.g. in API
// validation code for Seeds (executed in gardener-apiserver and gardenlet).
// This variable is an alias to the feature gate map in the apiserver library. The library doesn't allow using a custom
// feature gate map for gardener-apiserver. Hence, we reuse it for all our components.
var DefaultFeatureGate = utilfeature.DefaultMutableFeatureGate

// AllFeatureGates is the list of all feature gates.
var AllFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	DefaultSeccompProfile:                    {Default: false, PreRelease: featuregate.Alpha},
	UseNamespacedCloudProfile:                {Default: true, PreRelease: featuregate.GA, LockToDefault: true},
	ShootCredentialsBinding:                  {Default: true, PreRelease: featuregate.Beta},
	NewWorkerPoolHash:                        {Default: false, PreRelease: featuregate.Alpha},
	NewVPN:                                   {Default: true, PreRelease: featuregate.GA, LockToDefault: true},
	CredentialsRotationWithoutWorkersRollout: {Default: true, PreRelease: featuregate.Beta},
	InPlaceNodeUpdates:                       {Default: false, PreRelease: featuregate.Alpha},
	IstioTLSTermination:                      {Default: false, PreRelease: featuregate.Alpha},
	CloudProfileCapabilities:                 {Default: false, PreRelease: featuregate.Alpha},
	DoNotCopyBackupCredentials:               {Default: true, PreRelease: featuregate.Beta},
	OpenTelemetryCollector:                   {Default: false, PreRelease: featuregate.Alpha},
}

// GetFeatures returns a feature gate map with the respective specifications. Non-existing feature gates are ignored.
func GetFeatures(featureGates ...featuregate.Feature) map[featuregate.Feature]featuregate.FeatureSpec {
	out := make(map[featuregate.Feature]featuregate.FeatureSpec)

	for _, fg := range featureGates {
		if spec, ok := AllFeatureGates[fg]; ok {
			out[fg] = spec
		}
	}

	return out
}
