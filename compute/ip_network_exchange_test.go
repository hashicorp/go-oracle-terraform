// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, createdNetworkExchange, receivedNetworkExchange, "Mismatch found after create.")
	assert.Equal(t, createdNetworkExchange.FQDN, svc.getQualifiedName(createInput.Name), "Expected FDQN to be equal to qualified name")

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
