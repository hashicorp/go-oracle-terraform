// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_IPAddressPrefixSetTestName        = "test-acc-ip-address-prefix-set"
	_IPAddressPrefixSetTestDescription = "testing ip address prefix set"
)

func TestAccIPAddressPrefixSetsLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	svc, err := getIPAddressPrefixSetsClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateIPAddressPrefixSetInput{
		Name:              _IPAddressPrefixSetTestName,
		Description:       _IPAddressPrefixSetTestDescription,
		IPAddressPrefixes: []string{"192.0.0.168/16"},
		Tags:              []string{"testing"},
	}

	createdIPAddressPrefixSet, err := svc.CreateIPAddressPrefixSet(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Prefix Set succcessfully created")
	defer destroyIPAddressPrefixSet(t, svc, _IPAddressPrefixSetTestName)

	getInput := &GetIPAddressPrefixSetInput{
		Name: _IPAddressPrefixSetTestName,
	}
	receivedIPAddressPrefixSet, err := svc.GetIPAddressPrefixSet(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Prefix Set successfully fetched")

	if !reflect.DeepEqual(createdIPAddressPrefixSet, receivedIPAddressPrefixSet) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdIPAddressPrefixSet, receivedIPAddressPrefixSet)
	}

	updateInput := &UpdateIPAddressPrefixSetInput{
		Name:              _IPAddressPrefixSetTestName,
		Description:       _IPAddressPrefixSetTestDescription,
		IPAddressPrefixes: []string{"192.0.0.167/16"},
		Tags:              []string{"testing"},
	}
	updatedIPAddressPrefixSet, err := svc.UpdateIPAddressPrefixSet(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("IP Address Prefix Set succcessfully updated")
	receivedIPAddressPrefixSet, err = svc.GetIPAddressPrefixSet(getInput)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, updatedIPAddressPrefixSet, receivedIPAddressPrefixSet, "Mismatch found after create.")
	assert.Equal(t, updatedIPAddressPrefixSet.FQDN, svc.getQualifiedName(_IPAddressPrefixSetTestName), "Expected FDQN to be equal to qualified name")
}

func destroyIPAddressPrefixSet(t *testing.T, svc *IPAddressPrefixSetsClient, name string) {
	input := &DeleteIPAddressPrefixSetInput{
		Name: name,
	}
	if err := svc.DeleteIPAddressPrefixSet(input); err != nil {
		t.Fatal(err)
	}
}

func getIPAddressPrefixSetsClient() (*IPAddressPrefixSetsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.IPAddressPrefixSets(), nil
}
