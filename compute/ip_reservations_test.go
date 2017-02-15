package compute

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

func TestAccIPReservationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	createIPReservation := CreateIPReservationInfo{
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

	getIPReservationInfo := GetIPReservationInfo{
		Name: ipReservation.Name,
	}
	ipReservationInfo, err := iprc.GetIPReservation(getIPReservationInfo)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully retrieved ip reservation\n")

	if ipReservation.IP != ipReservationInfo.IP {
		t.Fatal("Created and retrived IP addresses don't match %s %s\n", ipReservation.IP, ipReservationInfo.IP)
	}

	deleteIPReservationInfo := DeleteIPReservationInfo{
		Name: ipReservation.Name,
	}
	err = iprc.DeleteIPReservation(deleteIPReservationInfo)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully deleted IPReservation\n")
}

func getIPReservationsClient() (*IPReservationsClient, error) {
	authenticatedClient, err := getAuthenticatedClient()
	if err != nil {
		return &IPReservationsClient{}, err
	}

	return authenticatedClient.IPReservations(), nil
}
