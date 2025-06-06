// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/gardener/gardener/pkg/apis/seedmanagement/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PendingReplicaApplyConfiguration represents an declarative configuration of the PendingReplica type for use
// with apply.
type PendingReplicaApplyConfiguration struct {
	Name    *string                        `json:"name,omitempty"`
	Reason  *v1alpha1.PendingReplicaReason `json:"reason,omitempty"`
	Since   *v1.Time                       `json:"since,omitempty"`
	Retries *int32                         `json:"retries,omitempty"`
}

// PendingReplicaApplyConfiguration constructs an declarative configuration of the PendingReplica type for use with
// apply.
func PendingReplica() *PendingReplicaApplyConfiguration {
	return &PendingReplicaApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *PendingReplicaApplyConfiguration) WithName(value string) *PendingReplicaApplyConfiguration {
	b.Name = &value
	return b
}

// WithReason sets the Reason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reason field is set to the value of the last call.
func (b *PendingReplicaApplyConfiguration) WithReason(value v1alpha1.PendingReplicaReason) *PendingReplicaApplyConfiguration {
	b.Reason = &value
	return b
}

// WithSince sets the Since field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Since field is set to the value of the last call.
func (b *PendingReplicaApplyConfiguration) WithSince(value v1.Time) *PendingReplicaApplyConfiguration {
	b.Since = &value
	return b
}

// WithRetries sets the Retries field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Retries field is set to the value of the last call.
func (b *PendingReplicaApplyConfiguration) WithRetries(value int32) *PendingReplicaApplyConfiguration {
	b.Retries = &value
	return b
}
