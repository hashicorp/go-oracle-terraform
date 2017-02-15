package compute

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccIPReservationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	createIPReservation := CreateIPReservationInput{
		ParentPool: "/oracle/public/ippool",
		Permanent:  true,
	}

	iprc, err := getIPReservationsClient()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Obtained IP Reservation Client\n")

	ipReservation, err := iprc.CreateIPReservation(createIPReservation)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully created IP Reservation: %+v\n", ipReservation)

	getIPReservationInput := GetIPReservationInput{
		Name: ipReservation.Name,
	}
	ipReservationInput, err := iprc.GetIPReservation(getIPReservationInput)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully retrieved ip reservation\n")

	if ipReservation.IP != ipReservationInput.IP {
		t.Fatal("Created and retrived IP addresses don't match %s %s\n", ipReservation.IP, ipReservationInput.IP)
	}

	deleteIPReservationInput := DeleteIPReservationInput{
		Name: ipReservation.Name,
	}
	err = iprc.DeleteIPReservation(deleteIPReservationInput)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully deleted IPReservation\n")
}

func getIPReservationsClient() (*IPReservationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &IPReservationsClient{}, err
	}

	return client.IPReservations(), nil
}
