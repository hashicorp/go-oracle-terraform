package compute

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_IPNetworkTestName        = "test-acc-ip-network"
	_IPNetworkTestPrefix      = "10.0.10.0/24"
	_IPNetworkTestPrefixTwo   = "10.0.12.0/24"
	_IPNetworkTestPrefixThree = "10.0.14.0/24"
	_IPNetworkTestDescription = "testing ip network"
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

	assert.Equal(t, createdNetwork, receivedNetwork, "Mismatch found after create.")
	assert.Equal(t, createdNetwork.FQDN, svc.getQualifiedName(createInput.Name), "Expected FDQN to be equal to qualified name")

	// Update prefix, NAPT, and tags
	updateInput := &UpdateIPNetworkInput{
		Name:              _IPNetworkTestName,
		IPAddressPrefix:   _IPNetworkTestPrefixTwo,
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

// Checks IP Networks with IP Network Exchanges
func TestAccIPNetworksWithExchangesLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sec, err := getIPNetworkExchangesClient()
	if err != nil {
		t.Fatal(err)
	}

	createExchangeInput := &CreateIPNetworkExchangeInput{
		Name:        _IPNetworkExchangeTestName,
		Description: _IPNetworkExchangeTestDescription,
		Tags:        []string{"testing"},
	}

	createdNetworkExchange, err := sec.CreateIPNetworkExchange(createExchangeInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Network Exchange succcessfully created")
	defer destroyIPNetworkExchange(t, sec, _IPNetworkExchangeTestName)

	svc, err := getIPNetworksClient()
	if err != nil {
		t.Fatal(err)
	}

	createNetworkInput := &CreateIPNetworkInput{
		Name:              _IPNetworkTestName,
		IPAddressPrefix:   _IPNetworkTestPrefix,
		Description:       _IPNetworkTestDescription,
		IPNetworkExchange: createdNetworkExchange.Name,
		Tags:              []string{"testing"},
	}

	createdNetwork, err := svc.CreateIPNetwork(createNetworkInput)
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

// Creates a generic IP Network with a supplied network prefix (to prevent collisions)
// and returns the resulting IP Network Info
func createTestIPNetwork(svc *IPNetworksClient, prefix string) (*IPNetworkInfo, error) {
	// Create a random name for the IP network
	rand.Seed(time.Now().UTC().UnixNano())
	rName := fmt.Sprintf("test-%d", rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	input := &CreateIPNetworkInput{
		Name:              rName,
		Description:       _IPNetworkTestDescription,
		IPAddressPrefix:   prefix,
		PublicNaptEnabled: false,
	}
	return svc.CreateIPNetwork(input)
}
