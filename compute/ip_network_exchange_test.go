package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_IPNetworkExchangeTestName        = "test-acc-ip-network-exchange"
	_IPNetworkExchangeTestDescription = "testing ip network exchange"
)

func TestAccIPNetworkExchangesLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	svc, err := getIPNetworkExchangesClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateIPNetworkExchangeInput{
		Name:        _IPNetworkExchangeTestName,
		Description: _IPNetworkExchangeTestDescription,
		Tags:        []string{"testing"},
	}

	createdNetworkExchange, err := svc.CreateIPNetworkExchange(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Network Exchange succcessfully created")
	defer destroyIPNetworkExchange(t, svc, _IPNetworkExchangeTestName)

	getInput := &GetIPNetworkExchangeInput{
		Name: _IPNetworkExchangeTestName,
	}
	receivedNetworkExchange, err := svc.GetIPNetworkExchange(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Network Exchange successfully fetched")

	if !reflect.DeepEqual(createdNetworkExchange, receivedNetworkExchange) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdNetworkExchange, receivedNetworkExchange)
	}
}

func destroyIPNetworkExchange(t *testing.T, svc *IPNetworkExchangesClient, name string) {
	input := &DeleteIPNetworkExchangeInput{
		Name: name,
	}
	if err := svc.DeleteIPNetworkExchange(input); err != nil {
		t.Fatal(err)
	}
}

func getIPNetworkExchangesClient() (*IPNetworkExchangesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.IPNetworkExchanges(), nil
}
