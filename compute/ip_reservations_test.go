package compute

import (
	"fmt"
	"testing"
)

func TestIPReservationLifeCycle(t *testing.T) {
	var (
		parentPool string = "/oracle/public/ippool"
		permanent  bool   = true
	)

	iprc, err := getIPReservationsClient()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Obtained IP Reservation Client\n")

	createdIPReservation, err := iprc.CreateIPReservation(parentPool, permanent, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully created IP Reservation: %+v\n", createdIPReservation)

	ipReservationInfo, err := iprc.GetIPReservation(createdIPReservation.Name)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully retrieved ip reservation\n")

	if createdIPReservation.IP != ipReservationInfo.IP {
		t.Fatal("Created and retrived IP addresses don't match %s %s\n", createdIPReservation.IP, ipReservationInfo.IP)
	}

	err = iprc.DeleteIPReservation(ipReservationInfo.Name)
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
