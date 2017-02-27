package compute

import (
	"testing"

	"reflect"

	"log"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_IPNetworkTestName         = "test-acc-ip-network"
	_IPNetworkTestPrefix       = "10.0.10.0/24"
	_IPNetworkTestPrefixUpdate = "10.0.12.0/24"
	_IPNetworkTestDescription  = "testing ip network"
)

func TestAccIPNetworksLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	svc, err := getIPNetworksClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateIPNetworkInput{
		Name:              _IPNetworkTestName,
		IPAddressPrefix:   _IPNetworkTestPrefix,
		Description:       _IPNetworkTestDescription,
		PublicNaptEnabled: false,
		Tags:              []string{"testing"},
	}

	createdNetwork, err := svc.CreateIPNetwork(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Network succcessfully created")
	defer destroyIPNetwork(t, svc, _IPNetworkTestName)

	getInput := &GetIPNetworkInput{
		Name: _IPNetworkTestName,
	}
	receivedNetwork, err := svc.GetIPNetwork(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Network successfully fetched")

	if !reflect.DeepEqual(createdNetwork, receivedNetwork) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdNetwork, receivedNetwork)
	}

	// Update prefix, NAPT, and tags
	updateInput := &UpdateIPNetworkInput{
		Name:              _IPNetworkTestName,
		IPAddressPrefix:   _IPNetworkTestPrefixUpdate,
		Description:       _IPNetworkTestDescription,
		PublicNaptEnabled: true,
		Tags:              []string{"updated"},
	}

	updatedNetwork, err := svc.UpdateIPNetwork(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Network successfully updated")

	receivedNetwork, err = svc.GetIPNetwork(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Network successfully fetched")

	if !reflect.DeepEqual(updatedNetwork, receivedNetwork) {
		t.Fatalf("Mismatch found after update.\nExpected: %+v\nReceived: %+v", createdNetwork, receivedNetwork)
	}

}

func destroyIPNetwork(t *testing.T, svc *IPNetworksClient, name string) {
	input := &DeleteIPNetworkInput{
		Name: name,
	}
	if err := svc.DeleteIPNetwork(input); err != nil {
		t.Fatal(err)
	}
}

func getIPNetworksClient() (*IPNetworksClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.IPNetworks(), nil
}
