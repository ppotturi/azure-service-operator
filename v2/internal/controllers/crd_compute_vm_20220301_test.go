/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package controllers_test

import (
	"testing"

	. "github.com/onsi/gomega"

	compute2022 "github.com/Azure/azure-service-operator/v2/api/compute/v1api20220301"
	network "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101"
	resources "github.com/Azure/azure-service-operator/v2/api/resources/v1api20200601"
	"github.com/Azure/azure-service-operator/v2/internal/testcommon"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

func newVirtualMachine20220301(
	tc *testcommon.KubePerTestContext,
	rg *resources.ResourceGroup,
	networkInterface *network.NetworkInterface,
	secretRef genruntime.SecretReference,
) *compute2022.VirtualMachine {
	adminUsername := "bloom"
	size := compute2022.HardwareProfile_VmSize_Standard_D1_V2

	return &compute2022.VirtualMachine{
		ObjectMeta: tc.MakeObjectMeta("vm"),
		Spec: compute2022.VirtualMachine_Spec{
			Location: tc.AzureRegion,
			Owner:    testcommon.AsOwner(rg),
			HardwareProfile: &compute2022.HardwareProfile{
				VmSize: &size,
			},
			OsProfile: &compute2022.OSProfile{
				AdminUsername: &adminUsername,
				// Specifying AdminPassword here rather than SSH Key to ensure that handling and injection
				// of secrets works.
				AdminPassword: &secretRef,
				ComputerName:  to.Ptr("poppy"),
			},
			StorageProfile: &compute2022.StorageProfile{
				ImageReference: &compute2022.ImageReference{
					Offer:     to.Ptr("UbuntuServer"),
					Publisher: to.Ptr("Canonical"),
					Sku:       to.Ptr("18.04-LTS"),
					Version:   to.Ptr("latest"),
				},
			},
			NetworkProfile: &compute2022.NetworkProfile{
				NetworkInterfaces: []compute2022.NetworkInterfaceReference{{
					Reference: tc.MakeReferenceFromResource(networkInterface),
				}},
			},
		},
	}
}

func Test_Compute_VM_20220301_CRUD(t *testing.T) {
	t.Parallel()

	tc := globalTestContext.ForTest(t)
	rg := tc.CreateTestResourceGroupAndWait()

	vnet := newVMVirtualNetwork(tc, testcommon.AsOwner(rg))
	subnet := newVMSubnet(tc, testcommon.AsOwner(vnet))
	networkInterface := newVMNetworkInterface(tc, testcommon.AsOwner(rg), subnet)
	// Inefficient but avoids triggering the vnet/subnets problem.
	// https://github.com/Azure/azure-service-operator/issues/1944
	tc.CreateResourceAndWait(vnet)
	tc.CreateResourcesAndWait(subnet, networkInterface)
	secret := createVMPasswordSecretAndRef(tc)
	vm := newVirtualMachine20220301(tc, rg, networkInterface, secret)

	tc.CreateResourceAndWait(vm)
	tc.Expect(vm.Status.Id).ToNot(BeNil())
	armId := *vm.Status.Id

	// Perform a simple patch to turn on boot diagnostics
	old := vm.DeepCopy()
	vm.Spec.DiagnosticsProfile = &compute2022.DiagnosticsProfile{
		BootDiagnostics: &compute2022.BootDiagnostics{
			Enabled: to.Ptr(true),
		},
	}

	tc.PatchResourceAndWait(old, vm)
	tc.Expect(vm.Status.DiagnosticsProfile).ToNot(BeNil())
	tc.Expect(vm.Status.DiagnosticsProfile.BootDiagnostics).ToNot(BeNil())
	tc.Expect(vm.Status.DiagnosticsProfile.BootDiagnostics.Enabled).ToNot(BeNil())
	tc.Expect(*vm.Status.DiagnosticsProfile.BootDiagnostics.Enabled).To(BeTrue())

	// Delete VM.
	tc.DeleteResourceAndWait(vm)

	// Ensure that the resource was really deleted in Azure
	exists, retryAfter, err := tc.AzureClient.CheckExistenceWithGetByID(tc.Ctx, armId, string(compute2022.APIVersion_Value))
	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(retryAfter).To(BeZero())
	tc.Expect(exists).To(BeFalse())
}
