package compute

import (
	"math/rand"
	"testing"

	"fmt"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

const (
	_TestIPAddressResDesc = "testing-acc"
	_TestIPAddressResName = "ip-res-testing"
	_TestIPAddressResTag  = "testing-tag"
)

func TestAccIPAddressReservationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	ipaClient, err := getIPAddressReservationsTestClients()
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

	ipRes, err := ipaClient.CreateIPAddressReservation(input)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPAddressReservation(t, ipaClient, resName)

	if ipRes.Name != resName {
		t.Fatalf("bad name: %s", ipRes.Name)
	}

	getInput := &GetIPAddressReservationInput{
		Name: resName,
	}

	receivedRes, err := ipaClient.GetIPAddressReservation(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if receivedRes.IPAddress == "" {
		t.Fatal("Expected an IPAddress created, got nil")
	}

	if diff := pretty.Compare(ipRes, receivedRes); diff != "" {
		t.Errorf("Created Reservation Diff: (-got +want)\n%s", diff)
	}

	updateInput := &UpdateIPAddressReservationInput{
		Description:   _TestIPAddressResDesc,
		IPAddressPool: PublicIPAddressPool,
		Name:          resName,
		Tags:          []string{_TestIPAddressResTag},
	}

	updatedRes, err := ipaClient.UpdateIPAddressReservation(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	receivedRes, err = ipaClient.GetIPAddressReservation(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if receivedRes.IPAddress == "" {
		t.Fatal("Expected a public IPAddress, got nil on update")
	}

	if diff := pretty.Compare(updatedRes, receivedRes); diff != "" {
		t.Errorf("Created Reservation Diff: (-got +want)\n%s", diff)
	}
}

func getIPAddressReservationsTestClients() (*IPAddressReservationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.IPAddressReservations(), nil
}

func destroyIPAddressReservation(t *testing.T, client *IPAddressReservationsClient, name string) {
	input := &DeleteIPAddressReservationInput{
		Name: name,
	}
	if err := client.DeleteIPAddressReservation(input); err != nil {
		t.Fatal(err)
	}
}
