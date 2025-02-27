/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateCreate", func() {
	var (
		clusterInstance *ClusterInstance

		v *clusterInstanceValidator
	)

	BeforeEach(func() {
		v = &clusterInstanceValidator{}
		clusterInstance = &ClusterInstance{
			Spec: ClusterInstanceSpec{
				TemplateRefs: []TemplateRef{
					{
						Name:      "cluster-template",
						Namespace: "default",
					},
				},
				Nodes: []NodeSpec{
					{
						HostName: "node1",
						TemplateRefs: []TemplateRef{
							{
								Name:      "node-template",
								Namespace: "default",
							},
						},
						Role: "master"},
				},
				ClusterType: ClusterTypeSNO,
			},
		}
	})

	It("should succeed for a valid ClusterInstance", func() {
		warnings, err := v.ValidateCreate(context.Background(), clusterInstance)
		Expect(err).NotTo(HaveOccurred())
		Expect(warnings).To(BeNil())
	})
})

var _ = Describe("ValidateUpdate", func() {
	var (
		ctx                                    context.Context
		v                                      *clusterInstanceValidator
		oldObj, newObj                         runtime.Object
		oldClusterInstance, newClusterInstance *ClusterInstance
	)

	BeforeEach(func() {
		ctx = context.Background()
		v = &clusterInstanceValidator{}

		oldClusterInstance = &ClusterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "site-sno-du-1",
				Namespace:   "site-sno-du-1",
				Annotations: make(map[string]string),
			},
			Spec: ClusterInstanceSpec{
				ClusterName:            "site-sno-du-1",
				PullSecretRef:          corev1.LocalObjectReference{Name: "pullSecretName"},
				ClusterImageSetNameRef: "openshift-test",
				SSHPublicKey:           "ssh-rsa",
				BaseDomain:             "example.com",
				ApiVIPs:                []string{"192.0.2.1", "192.0.2.2"},
				HoldInstallation:       false,
				AdditionalNTPSources:   []string{"NTP.server1", "192.0.2.3"},
				MachineNetwork:         []MachineNetworkEntry{{CIDR: "203.0.113.0/24"}},
				ClusterNetwork:         []ClusterNetworkEntry{{CIDR: "203.0.113.0/24", HostPrefix: 23}},
				ServiceNetwork:         []ServiceNetworkEntry{{CIDR: "203.0.113.0/24"}},
				NetworkType:            "OVNKubernetes",
				ExtraLabels:            map[string]map[string]string{"ManagedCluster": {"group-du-sno": "test", "common": "true", "sites": "site-sno-du-1"}},
				InstallConfigOverrides: "{\"capabilities\":{\"baselineCapabilitySet\": \"None\", \"additionalEnabledCapabilities\": [ \"marketplace\", \"NodeTuning\" ] }}",
				ExtraManifestsRefs:     []corev1.LocalObjectReference{{Name: "foobar1"}, {Name: "foobar2"}},
				TemplateRefs:           []TemplateRef{{Name: "cluster-v1", Namespace: "site-sno-du-1"}},
				Nodes: []NodeSpec{{
					BmcAddress:         "idrac-virtualmedia+https://198.51.100.0/redfish/v1/Systems/System.Embedded.1",
					BmcCredentialsName: BmcCredentialsName{Name: "bmc-secret"},
					BootMACAddress:     "00:00:5E:00:53:00",
					HostName:           "node1",
					Role:               "master",
					BootMode:           "UEFI",
					InstallerArgs:      "[\"--append-karg\", \"nameserver=198.51.100.0\", \"-n\"]",
					TemplateRefs:       []TemplateRef{{Name: "node-template", Namespace: "site-sno-du-1"}},
				}},
			},
		}

	})

	It("should return nil - dummy test", func() {
		oldClusterInstance.Status = ClusterInstanceStatus{
			Conditions: []metav1.Condition{
				{
					Type:   string(ClusterProvisioned),
					Status: metav1.ConditionFalse,
					Reason: string(InProgress),
				},
			},
		}
		oldObj = oldClusterInstance

		newClusterInstance = oldClusterInstance.DeepCopy()
		newObj = newClusterInstance

		newClusterInstance.Spec.ExtraAnnotations = map[string]map[string]string{
			"BareMetalHost": {
				"foo": "bar",
			},
		}

		_, err := v.ValidateUpdate(ctx, oldObj, newObj)
		Expect(err).NotTo(HaveOccurred())
	})

})
