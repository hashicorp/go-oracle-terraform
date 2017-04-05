package compute

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_IPAddressAssociationTestName        = "test-acc-ip-address-association"
	_IPAddressAssociationTestDescription = "testing ip address association"
)

func TestAccIPAddressAssociationsLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	iClient, nClient, vnClient, err := getVirtNICsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	ipNetwork, err := createTestIPNetwork(nClient, _IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetwork.Name)

	// In order to get details on a Virtual NIC we need to create the following resources
	// - IP Network
	// - Instance

	instanceInput := &CreateInstanceInput{
		Name:      _VirtNicInstanceTestName,
		Label:     _VirtNicInstanceTestLabel,
		Shape:     _VirtNicInstanceTestShape,
		ImageList: _VirtNicInstanceTestImage,
		Networking: map[string]NetworkingInfo{
			"eth0": {
				IPNetwork: ipNetwork.Name,
				Vnic:      "eth0",
			},
		},
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	// Use the static "eth0" interface, as we statically created that above
	createdVNIC := createdInstance.Networking["eth0"].Vnic
	getVNICInput := &GetVirtualNICInput{
		Name: createdVNIC,
	}

	vNIC, err := vnClient.GetVirtualNIC(getVNICInput)
	if err != nil {
		t.Fatal(err)
	}

	ipaClient, err := getIPAddressReservationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	resName := fmt.Sprintf("%s-%d", _TestIPAddressResName, helper.RInt())

	input := &CreateIPAddressReservationInput{
		Description:   _TestIPAddressResDesc,
		IPAddressPool: PrivateIPAddressPool,
		Name:          resName,
		Tags:          []string{_TestIPAddressResTag},
	}

	ipRes, err := ipaClient.CreateIPAddressReservation(input)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPAddressReservation(t, ipaClient, resName)

	svc, err := getIPAddressAssociationsClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateIPAddressAssociationInput{
		Name:                 _IPAddressAssociationTestName,
		Description:          _IPAddressAssociationTestDescription,
		IPAddressReservation: ipRes.Name,
		VNIC:                 vNIC.Name,
		Tags:                 []string{"testing"},
	}

	createdIPAddressAssociation, err := svc.CreateIPAddressAssociation(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Association succcessfully created")
	defer destroyIPAddressAssociation(t, svc, _IPAddressAssociationTestName)

	getInput := &GetIPAddressAssociationInput{
		Name: _IPAddressAssociationTestName,
	}
	receivedIPAddressAssociation, err := svc.GetIPAddressAssociation(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Association successfully fetched")

	if !reflect.DeepEqual(createdIPAddressAssociation, receivedIPAddressAssociation) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdIPAddressAssociation, receivedIPAddressAssociation)
	}
	if receivedIPAddressAssociation.VNIC != vNIC.Name {
		t.Fatalf("VNIC Mismatch found after create.\nExpected: %+v\nReceived: %+v", vNIC.Name, receivedIPAddressAssociation.VNIC)
	}
	if receivedIPAddressAssociation.IPAddressReservation != ipRes.Name {
		t.Fatalf("IPAddressReservation Mismatch found after create.\nExpected: %+v\nReceived: %+v", ipRes.Name, receivedIPAddressAssociation.IPAddressReservation)
	}

	updateInput := &UpdateIPAddressAssociationInput{
		Name:                 _IPAddressAssociationTestName,
		Description:          _IPAddressAssociationTestDescription,
		IPAddressReservation: ipRes.Name,
		VNIC:                 vNIC.Name,
		Tags:                 []string{"testing-updated"},
	}
	updatedIPAddressAssociation, err := svc.UpdateIPAddressAssociation(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Association succcessfully updated")
	receivedIPAddressAssociation, err = svc.GetIPAddressAssociation(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(updatedIPAddressAssociation, receivedIPAddressAssociation) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", updatedIPAddressAssociation, receivedIPAddressAssociation)
	}
}

func destroyIPAddressAssociation(t *testing.T, svc *IPAddressAssociationsClient, name string) {
	input := &DeleteIPAddressAssociationInput{
		Name: name,
	}
	if err := svc.DeleteIPAddressAssociation(input); err != nil {
		t.Fatal(err)
	}
}

func getIPAddressAssociationsClient() (*IPAddressAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.IPAddressAssociations(), nil
}
