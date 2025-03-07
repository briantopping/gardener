// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeCoreV1beta1 struct {
	*testing.Fake
}

func (c *FakeCoreV1beta1) BackupBuckets() v1beta1.BackupBucketInterface {
	return newFakeBackupBuckets(c)
}

func (c *FakeCoreV1beta1) BackupEntries(namespace string) v1beta1.BackupEntryInterface {
	return newFakeBackupEntries(c, namespace)
}

func (c *FakeCoreV1beta1) CloudProfiles() v1beta1.CloudProfileInterface {
	return newFakeCloudProfiles(c)
}

func (c *FakeCoreV1beta1) ControllerDeployments() v1beta1.ControllerDeploymentInterface {
	return newFakeControllerDeployments(c)
}

func (c *FakeCoreV1beta1) ControllerInstallations() v1beta1.ControllerInstallationInterface {
	return newFakeControllerInstallations(c)
}

func (c *FakeCoreV1beta1) ControllerRegistrations() v1beta1.ControllerRegistrationInterface {
	return newFakeControllerRegistrations(c)
}

func (c *FakeCoreV1beta1) ExposureClasses() v1beta1.ExposureClassInterface {
	return newFakeExposureClasses(c)
}

func (c *FakeCoreV1beta1) InternalSecrets(namespace string) v1beta1.InternalSecretInterface {
	return newFakeInternalSecrets(c, namespace)
}

func (c *FakeCoreV1beta1) NamespacedCloudProfiles(namespace string) v1beta1.NamespacedCloudProfileInterface {
	return newFakeNamespacedCloudProfiles(c, namespace)
}

func (c *FakeCoreV1beta1) Projects() v1beta1.ProjectInterface {
	return newFakeProjects(c)
}

func (c *FakeCoreV1beta1) Quotas(namespace string) v1beta1.QuotaInterface {
	return newFakeQuotas(c, namespace)
}

func (c *FakeCoreV1beta1) SecretBindings(namespace string) v1beta1.SecretBindingInterface {
	return newFakeSecretBindings(c, namespace)
}

func (c *FakeCoreV1beta1) Seeds() v1beta1.SeedInterface {
	return newFakeSeeds(c)
}

func (c *FakeCoreV1beta1) Shoots(namespace string) v1beta1.ShootInterface {
	return newFakeShoots(c, namespace)
}

func (c *FakeCoreV1beta1) ShootStates(namespace string) v1beta1.ShootStateInterface {
	return newFakeShootStates(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeCoreV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
