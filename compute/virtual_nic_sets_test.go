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
