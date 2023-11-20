// Copyright 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/pointer"

	. "github.com/gardener/gardener/pkg/apis/core/v1beta1"
)

var _ = Describe("Seed defaulting", func() {
	var obj *Seed

	BeforeEach(func() {
		obj = &Seed{}
	})

	Describe("SeedSettings defaulting", func() {
		It("should default the seed settings (w/o taints)", func() {
			var excessCapacityReservation = SeedSettingExcessCapacityReservation{
				Configs: []SeedSettingExcessCapacityReservationConfig{
					{Resources: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("2"), corev1.ResourceMemory: resource.MustParse("6Gi")}},
				},
			}

			SetObjectDefaults_Seed(obj)

			Expect(obj.Spec.Settings.DependencyWatchdog).NotTo(BeNil())
			Expect(obj.Spec.Settings.ExcessCapacityReservation).To(PointTo(Equal(excessCapacityReservation)))
			Expect(obj.Spec.Settings.Scheduling.Visible).To(BeTrue())
			Expect(obj.Spec.Settings.VerticalPodAutoscaler.Enabled).To(BeTrue())
			Expect(obj.Spec.Settings.TopologyAwareRouting.Enabled).To(BeFalse())
		})

		It("should default the excessCapacityReservation field when excess capacity reservation is enabled and excess capacity reservation config empty", func() {
			var excessCapacityReservation = SeedSettingExcessCapacityReservation{
				Enabled: pointer.Bool(true),
				Configs: []SeedSettingExcessCapacityReservationConfig{
					{Resources: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("2"), corev1.ResourceMemory: resource.MustParse("6Gi")}},
				},
			}

			obj.Spec.Settings = &SeedSettings{}
			obj.Spec.Settings.ExcessCapacityReservation = &SeedSettingExcessCapacityReservation{Enabled: pointer.Bool(true)}
			obj.Spec.Settings.ExcessCapacityReservation.Enabled = pointer.Bool(true)

			SetObjectDefaults_Seed(obj)

			Expect(obj.Spec.Settings.DependencyWatchdog).NotTo(BeNil())
			Expect(obj.Spec.Settings.ExcessCapacityReservation).To(PointTo(Equal(excessCapacityReservation)))
		})

		It("should not overwrite the already set values for seed settings field", func() {
			var (
				excessCapacityReservation = SeedSettingExcessCapacityReservation{
					Enabled: pointer.Bool(true),
					Configs: []SeedSettingExcessCapacityReservationConfig{
						{Resources: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("4"), corev1.ResourceMemory: resource.MustParse("16Gi")}},
					},
				}
			)

			obj.Spec.Settings = &SeedSettings{
				DependencyWatchdog: &SeedSettingDependencyWatchdog{
					Weeder: &SeedSettingDependencyWatchdogWeeder{Enabled: false},
					Prober: &SeedSettingDependencyWatchdogProber{Enabled: false},
				},
				TopologyAwareRouting: &SeedSettingTopologyAwareRouting{
					Enabled: true,
				},
				ExcessCapacityReservation: &excessCapacityReservation,
				Scheduling:                &SeedSettingScheduling{Visible: true},
				VerticalPodAutoscaler:     &SeedSettingVerticalPodAutoscaler{Enabled: false},
			}

			SetObjectDefaults_Seed(obj)

			Expect(obj.Spec.Settings.DependencyWatchdog.Weeder.Enabled).To(BeFalse())
			Expect(obj.Spec.Settings.DependencyWatchdog.Prober.Enabled).To(BeFalse())
			Expect(obj.Spec.Settings.ExcessCapacityReservation).To(PointTo(Equal(excessCapacityReservation)))
			Expect(obj.Spec.Settings.Scheduling.Visible).To(BeTrue())
			Expect(obj.Spec.Settings.VerticalPodAutoscaler.Enabled).To(BeFalse())
			Expect(obj.Spec.Settings.TopologyAwareRouting.Enabled).To(BeTrue())
		})
	})

	Describe("SeedNetworks defaulting", func() {
		It("should default ipFamilies setting to IPv4 single-stack", func() {
			SetObjectDefaults_Seed(obj)

			Expect(obj.Spec.Networks.IPFamilies).To(ConsistOf(IPFamilyIPv4))
		})

		It("should not overwrite the already set values for ipFamilies setting field", func() {
			obj.Spec.Networks.IPFamilies = []IPFamily{IPFamilyIPv6}
			SetObjectDefaults_Seed(obj)

			Expect(obj.Spec.Networks.IPFamilies).To(ConsistOf(IPFamilyIPv6))
		})
	})

	Describe("SeedSettingDependencyWatchdog defaulting", func() {
		It("should default the settings", func() {
			SetObjectDefaults_Seed(obj)

			Expect(obj.Spec.Settings.DependencyWatchdog.Weeder.Enabled).To(BeTrue())
			Expect(obj.Spec.Settings.DependencyWatchdog.Prober.Enabled).To(BeTrue())
		})

		It("should not overwrite the already set values for SeedSettingDependencyWatchdog", func() {
			var (
				dwdWeederEnabled = false
				dwdProberEnabled = false
			)

			obj.Spec.Settings = &SeedSettings{
				DependencyWatchdog: &SeedSettingDependencyWatchdog{
					Weeder: &SeedSettingDependencyWatchdogWeeder{Enabled: dwdWeederEnabled},
					Prober: &SeedSettingDependencyWatchdogProber{Enabled: dwdProberEnabled},
				},
			}

			Expect(obj.Spec.Settings.DependencyWatchdog.Weeder.Enabled).To(Equal(dwdWeederEnabled))
			Expect(obj.Spec.Settings.DependencyWatchdog.Prober.Enabled).To(Equal(dwdProberEnabled))
		})
	})
})
