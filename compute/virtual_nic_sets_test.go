package compute

import (
	"log"
	"testing"

	"reflect"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccVirtNICSetsLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	svc, err := getVirtNICSetsClient()
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
	// Verify they're the same
	if createdSet.Name != returnedSet.Name {
		t.Fatalf("Mismatched Sets found.\nExpected: %+v\nReceived: %+v", createdSet, returnedSet)
	}

	if createdSet.Description != returnedSet.Description {
		t.Fatalf("Mismatched Sets found.\nExpected: %+v\nReceived: %+v", createdSet, returnedSet)
	}

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

	if !reflect.DeepEqual(updatedSet, returnedSet) {
		t.Fatalf("Mismatched Sets found.\nExpected: %+v\nReceived: %+v", updatedSet, returnedSet)
	}
}

func TestAccVirtNICSetAddNICS(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// Fist, create necessary clients
	instanceSvc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

	svc, err := getVirtNICSetsClient()
	if err != nil {
		t.Fatal(err)
	}

	// Create an instance with multiple vNICs
	// TODO: Remove hardcoded IP Network when the IP Network resource is added
	instanceInput := &CreateInstanceInput{
		Name:      "test-acc-virt-nic-lifecycle",
		Label:     "test",
		Shape:     "oc3",
		ImageList: "/oracle/public/Oracle_Solaris_11.3",
		Networking: map[string]NetworkingInfo{
			"eth0": {
				IPNetwork: "testing-1",
				Vnic:      "eth0",
			},
			"eth1": {
				IPNetwork: "testing-2",
				Vnic:      "eth1",
			},
			"eth2": {
				IPNetwork: "testing-3",
				Vnic:      "eth2",
			},
		},
	}

	createdInstance, err := instanceSvc.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, instanceSvc, createdInstance.Name, createdInstance.ID)

	vNIC1 := createdInstance.Networking["eth0"].Vnic
	vNIC2 := createdInstance.Networking["eth1"].Vnic
	vNIC3 := createdInstance.Networking["eth2"].Vnic

	// Create virtual nic set using two of the created vNICs
	input := &CreateVirtualNICSetInput{
		Name:            "test-acc-nic-set-nics",
		Description:     "testing nic sets",
		Tags:            []string{"tag1", "tag2"},
		VirtualNICNames: []string{vNIC1, vNIC2},
	}

	createdSet, err := svc.CreateVirtualNICSet(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVirtualNICSet(t, svc, createdSet.Name)

	// Get the created set and compare
	getInput := &GetVirtualNICSetInput{
		Name: createdSet.Name,
	}

	returnedSet, err := svc.GetVirtualNICSet(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if len(returnedSet.VirtualNICs) != 2 {
		t.Fatalf("Expected 2 sets of virtual nics, got: %d", len(returnedSet.VirtualNICs))
	}

	// Verify that the vNICs in the returned set are populated
	for _, v := range returnedSet.VirtualNICs {
		if v.Name != "eth0" && v.Name != "eth1" {
			t.Fatalf("Expected vNIC to be either 'eth0' or 'eth1'. Got: %s", v.Name)
		}
		if v.MACAddress == "" {
			t.Fatalf("Empty MAC address found for vNIC: %s", v.Name)
		}
		if v.TransitFlag {
			t.Fatalf("Expected transit flag to be false, got %b for %s", v.TransitFlag, v.Name)
		}
		if v.Uri == "" {
			t.Fatalf("Empty URI returned for vNIC %s", v.Name)
		}
	}

	// Update the set with the third vNIC
	updateInput := &UpdateVirtualNICSetInput{
		Name:            createdSet.Name,
		Description:     createdSet.Description,
		Tags:            createdSet.Tags,
		VirtualNICNames: []string{vNIC1, vNIC2, vNIC3},
	}

	_, err = svc.UpdateVirtualNICSet(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	// Get the updated set and compare
	returnedSet, err = svc.GetVirtualNICSet(getInput)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the vNICs in the returned set are populated
	for _, v := range returnedSet.VirtualNICs {
		if v.Name != "eth0" && v.Name != "eth1" && v.Name != "eth2" {
			t.Fatalf("Expected vNIC to be either 'eth0', 'eth1', or 'eth2'. Got: %s", v.Name)
		}
		if v.MACAddress == "" {
			t.Fatalf("Empty MAC address found for vNIC: %s", v.Name)
		}
		if v.TransitFlag {
			t.Fatalf("Expected transit flag to be false, got %b for %s", v.TransitFlag, v.Name)
		}
		if v.Uri == "" {
			t.Fatalf("Empty URI returned for vNIC %s", v.Name)
		}
	}

	log.Printf("Virtual NIC Set successfully created and updated")
}

func getVirtNICSetsClient() (*VirtNICSetsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}
	return client.VirtNICSets(), nil
}

func deleteVirtualNICSet(t *testing.T, svc *VirtNICSetsClient, name string) {
	input := &DeleteVirtualNICSetInput{
		Name: name,
	}
	if err := svc.DeleteVirtualNICSet(input); err != nil {
		t.Fatalf("Error deleting set: %v", err)
	}
}
