package compute

import (
	"fmt"
	"log"
	"math/rand"
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

	iClient, nClient, vnClient, iprClient, ipaClient, err := getIPAssociationTestClients()
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
	createdVnic := createdInstance.Networking["eth0"].Vnic
	getVnicInput := &GetVirtualNICInput{
		Name: createdVnic,
	}

	vNIC, err := vnClient.GetVirtualNIC(getVnicInput)
	if err != nil {
		t.Fatal(err)
	}

	rInt := rand.Int()
	resName := fmt.Sprintf("%s-%d", _TestIPAddressResName, rInt)

	input := &CreateIPAddressReservationInput{
		Description:   _TestIPAddressResDesc,
		IPAddressPool: PrivateIPAddressPool,
		Name:          resName,
		Tags:          []string{_TestIPAddressResTag},
	}

	ipRes, err := iprClient.CreateIPAddressReservation(input)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPAddressReservation(t, iprClient, resName)

	createInput := &CreateIPAddressAssociationInput{
		Name:                 _IPAddressAssociationTestName,
		Description:          _IPAddressAssociationTestDescription,
		IPAddressReservation: ipRes.Name,
		Vnic:                 vNIC.Name,
		Tags:                 []string{"testing"},
	}

	createdIPAddressAssociation, err := ipaClient.CreateIPAddressAssociation(createInput)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPAddressAssociation(t, ipaClient, _IPAddressAssociationTestName)
	log.Print("IP Address Association succcessfully created")

	getInput := &GetIPAddressAssociationInput{
		Name: _IPAddressAssociationTestName,
	}
	receivedIPAddressAssociation, err := ipaClient.GetIPAddressAssociation(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Association successfully fetched")

	if !reflect.DeepEqual(createdIPAddressAssociation, receivedIPAddressAssociation) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdIPAddressAssociation, receivedIPAddressAssociation)
	}
	if receivedIPAddressAssociation.Vnic != vNIC.Name {
		t.Fatalf("Vnic Mismatch found after create.\nExpected: %+v\nReceived: %+v", vNIC.Name, receivedIPAddressAssociation.Vnic)
	}
	if receivedIPAddressAssociation.IPAddressReservation != ipRes.Name {
		t.Fatalf("IPAddressReservation Mismatch found after create.\nExpected: %+v\nReceived: %+v", ipRes.Name, receivedIPAddressAssociation.IPAddressReservation)
	}

	updateInput := &UpdateIPAddressAssociationInput{
		Name:                 _IPAddressAssociationTestName,
		Description:          _IPAddressAssociationTestDescription,
		IPAddressReservation: ipRes.Name,
		Vnic:                 vNIC.Name,
		Tags:                 []string{"testing-updated"},
	}
	updatedIPAddressAssociation, err := ipaClient.UpdateIPAddressAssociation(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Association succcessfully updated")
	receivedIPAddressAssociation, err = ipaClient.GetIPAddressAssociation(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(updatedIPAddressAssociation, receivedIPAddressAssociation) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", updatedIPAddressAssociation, receivedIPAddressAssociation)
	}
}

func destroyIPAddressAssociation(t *testing.T, ipaClient *IPAddressAssociationsClient, name string) {
	input := &DeleteIPAddressAssociationInput{
		Name: name,
	}
	if err := ipaClient.DeleteIPAddressAssociation(input); err != nil {
		t.Fatal(err)
	}
}

func getIPAssociationTestClients() (*InstancesClient, *IPNetworksClient, *VirtNICsClient, *IPAddressReservationsClient, *IPAddressAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return client.Instances(), client.IPNetworks(), client.VirtNICs(), client.IPAddressReservations(), client.IPAddressAssociations(), nil
}
