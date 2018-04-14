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
	_IPReservationPerm      = true
	_IPReservationName      = "testing-ip-res"
	_IPReservationTag       = "testing-tag"
	_IPReservationTagUpdate = "testing-tag-update"
)

func TestAccIPReservationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	iprClient, err := getIPReservationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	rInt := rand.Int()
	resName := fmt.Sprintf("%s-%d", _IPReservationName, rInt)

	createIPReservation := &CreateIPReservationInput{
		Name:       resName,
		ParentPool: PublicReservationPool,
		Permanent:  _IPReservationPerm,
		Tags:       []string{_IPReservationTag},
	}

	ipRes, err := iprClient.CreateIPReservation(createIPReservation)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyIPReservation(t, iprClient, ipRes.Name)

	getInput := &GetIPReservationInput{
		Name: resName,
	}

	receivedRes, err := iprClient.GetIPReservation(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if diff := pretty.Compare(ipRes, receivedRes); diff != "" {
		t.Errorf("Created Reservation Diff: (-got +want)\n%s", diff)
	}

	updateInput := &UpdateIPReservationInput{
		Name:       resName,
		ParentPool: PublicReservationPool,
		Permanent:  _IPReservationPerm,
		Tags:       []string{_IPReservationTagUpdate},
	}

	updatedRes, err := iprClient.UpdateIPReservation(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	receivedRes, err = iprClient.GetIPReservation(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if diff := pretty.Compare(updatedRes, receivedRes); diff != "" {
		t.Errorf("Created Reservation Diff: (-got +want)\n%s", diff)
	}

}

func getIPReservationsTestClients() (*IPReservationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &IPReservationsClient{}, err
	}

	return client.IPReservations(), nil
}

func destroyIPReservation(t *testing.T, client *IPReservationsClient, name string) {
	input := &DeleteIPReservationInput{
		Name: name,
	}

	if err := client.DeleteIPReservation(input); err != nil {
		t.Fatal(err)
	}
}
