package compute

import (
	"log"
	"reflect"
	"testing"

	"sort"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccVirtNICSetsLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	_, svc, _, err := getVirtNICSetsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	input := &CreateVirtualNICSetInput{
		Name:        "test-acc-nic-set",
		Description: "testing-nic-set",
	}

	createdSet, err := svc.CreateVirtualNICSet(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVirtualNICSet(t, svc, createdSet.Name)
	log.Printf("Created NIC Set: %#v", createdSet)

	getInput := &GetVirtualNICSetInput{
		Name: createdSet.Name,
	}

	returnedSet, err := svc.GetVirtualNICSet(getInput)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, createdSet.Name, returnedSet.Name, "Mismatched Sets found.")
	assert.Equal(t, createdSet.Description, returnedSet.Description, "Mismatched Sets found.")
	assert.Equal(t, returnedSet.FQDN, svc.getQualifiedName(createdSet.Name), "Expected FDQN to be equal to qualified name")

	// Update the set
	updateInput := &UpdateVirtualNICSetInput{
		Name:        createdSet.Name,
		Description: createdSet.Description,
		AppliedACLs: []string{"default"},
		Tags:        []string{"tag1", "tag2"},
	}

	updatedSet, err := svc.UpdateVirtualNICSet(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Updated NIC Set: %#v", updatedSet)

	// Get the set again to ensure fields are updated
	returnedSet, err = svc.GetVirtualNICSet(getInput)
	if err != nil {
		t.Fatal(err)
	}

	// Need to sort slices, for ordering. The Provider takes care of this at the end-user level, but for
	// testing these can be out of whack when returned from the API
	sort.Strings(updatedSet.Tags)
	sort.Strings(returnedSet.Tags)

	if !reflect.DeepEqual(updatedSet, returnedSet) {
		t.Fatalf("Mismatched Sets found.\nExpected: %+v\nReceived: %+v", updatedSet, returnedSet)
	}
}

func TestAccVirtNICSetsAddNICS(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// Fist, create necessary clients
	iClient, vnClient, nClient, err := getVirtNICSetsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// Create the three IP Networks needed
	ipNetworkOne, err := createTestIPNetwork(nClient, _IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetworkOne.Name)

	ipNetworkTwo, err := createTestIPNetwork(nClient, _IPNetworkTestPrefixTwo)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetworkTwo.Name)

	ipNetworkThree, err := createTestIPNetwork(nClient, _IPNetworkTestPrefixThree)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetworkThree.Name)

	// Create an instance with multiple vNICs
	instanceInput := &CreateInstanceInput{
		Name:      _VirtNicInstanceTestName,
		Label:     _VirtNicInstanceTestLabel,
		Shape:     _VirtNicInstanceTestShape,
		ImageList: _VirtNicInstanceTestImage,
		Networking: map[string]NetworkingInfo{
			"eth0": {
				IPNetwork: ipNetworkOne.Name,
				Vnic:      "eth0",
			},
			"eth1": {
				IPNetwork: ipNetworkTwo.Name,
				Vnic:      "eth1",
			},
			"eth2": {
				IPNetwork: ipNetworkThree.Name,
				Vnic:      "eth2",
			},
		},
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	vNIC1 := createdInstance.Networking["eth0"].Vnic
	vNIC2 := createdInstance.Networking["eth1"].Vnic
	vNIC3 := createdInstance.Networking["eth2"].Vnic

	// Create virtual nic set using two of the created vNICs
	input := &CreateVirtualNICSetInput{
		Name:        "test-acc-nic-set-nics",
		Description: "testing nic sets",
		Tags:        []string{"test-tag"},
		VirtualNICs: []string{vNIC1, vNIC2},
	}

	createdSet, err := vnClient.CreateVirtualNICSet(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVirtualNICSet(t, vnClient, createdSet.Name)

	// Get the created set and compare
	getInput := &GetVirtualNICSetInput{
		Name: createdSet.Name,
	}

	returnedSet, err := vnClient.GetVirtualNICSet(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if len(returnedSet.VirtualNICs) != 2 {
		t.Fatalf("Expected 2 sets of virtual nics, got: %d", len(returnedSet.VirtualNICs))
	}

	// Sort slices
	sort.Strings(createdSet.VirtualNICs)
	sort.Strings(returnedSet.VirtualNICs)

	// Verify that the vNICs in the returned set are populated
	if !reflect.DeepEqual(createdSet, returnedSet) {
		t.Fatalf("Mismatch Found!\nExpected: %+v\nReceived: %+v", createdSet, returnedSet)
	}

	// Update the set with the third vNIC
	updateInput := &UpdateVirtualNICSetInput{
		Name:        createdSet.Name,
		Description: createdSet.Description,
		Tags:        createdSet.Tags,
		VirtualNICs: []string{vNIC1, vNIC2, vNIC3},
	}

	updatedSet, err := vnClient.UpdateVirtualNICSet(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated set and compare
	returnedSet, err = vnClient.GetVirtualNICSet(getInput)
	if err != nil {
		t.Fatal(err)
	}

	// Sort slices
	sort.Strings(updatedSet.VirtualNICs)
	sort.Strings(returnedSet.VirtualNICs)

	// Verify that the vNICs in the returned set are populated
	if !reflect.DeepEqual(updatedSet, returnedSet) {
		t.Fatalf("Mismatch Found!\nExpected: %+v\nReceived: %+v", createdSet, returnedSet)
	}

	log.Printf("Virtual NIC Set successfully created and updated")
}

func getVirtNICSetsTestClients() (*InstancesClient, *VirtNICSetsClient, *IPNetworksClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}
	return client.Instances(), client.VirtNICSets(), client.IPNetworks(), nil
}

func deleteVirtualNICSet(t *testing.T, svc *VirtNICSetsClient, name string) {
	input := &DeleteVirtualNICSetInput{
		Name: name,
	}
	if err := svc.DeleteVirtualNICSet(input); err != nil {
		t.Fatalf("Error deleting set: %v", err)
	}
}
