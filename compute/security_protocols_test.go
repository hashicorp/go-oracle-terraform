package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_SecurityProtocolTestName        = "test-acc-security-protocol"
	_SecurityProtocolTestDescription = "testing security protocol"
	_SecurityProtocolTestIPProtocol  = "tcp"
)

func TestAccSecurityProtocolsLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	svc, err := getSecurityProtocolsClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateSecurityProtocolInput{
		Name:        _SecurityProtocolTestName,
		Description: _SecurityProtocolTestDescription,
		Tags:        []string{"testing"},
		IPProtocol:  _SecurityProtocolTestIPProtocol,
		SrcPortSet:  []string{"17"},
		DstPortSet:  []string{"18"},
	}

	createdSecurityProtocol, err := svc.CreateSecurityProtocol(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Protocol succcessfully created")
	defer destroySecurityProtocol(t, svc, _SecurityProtocolTestName)

	getInput := &GetSecurityProtocolInput{
		Name: _SecurityProtocolTestName,
	}
	receivedSecurityProtocol, err := svc.GetSecurityProtocol(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Protocol successfully fetched")

	if !reflect.DeepEqual(createdSecurityProtocol, receivedSecurityProtocol) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdSecurityProtocol, receivedSecurityProtocol)
	}

	updateInput := &UpdateSecurityProtocolInput{
		Name:        _SecurityProtocolTestName,
		Description: _SecurityProtocolTestDescription,
		Tags:        []string{"testing"},
		IPProtocol:  _SecurityProtocolTestIPProtocol,
		SrcPortSet:  []string{"20"},
		DstPortSet:  []string{"21"},
	}
	updatedSecurityProtocol, err := svc.UpdateSecurityProtocol(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Protocol succcessfully updated")
	receivedSecurityProtocol, err = svc.GetSecurityProtocol(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(updatedSecurityProtocol, receivedSecurityProtocol) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", updatedSecurityProtocol, receivedSecurityProtocol)
	}
}

func destroySecurityProtocol(t *testing.T, svc *SecurityProtocolsClient, name string) {
	input := &DeleteSecurityProtocolInput{
		Name: name,
	}
	if err := svc.DeleteSecurityProtocol(input); err != nil {
		t.Fatal(err)
	}
}

func getSecurityProtocolsClient() (*SecurityProtocolsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.SecurityProtocols(), nil
}
