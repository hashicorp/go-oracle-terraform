package compute

import (
	"log"
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
	log.Printf("Obtained IP Reservation Client")

	ipReservation, err := iprc.CreateIPReservation(&createIPReservation)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPReservation(t, iprc, ipReservation.Name)

	log.Printf("Successfully created IP Reservation: %+v", ipReservation)

	getIPReservationInput := GetIPReservationInput{
		Name: ipReservation.Name,
	}
	ipReservationOutput, err := iprc.GetIPReservation(&getIPReservationInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully retrieved ip reservation")

	if ipReservation.IP != ipReservationOutput.IP {
		t.Fatalf("Created and retrived IP addresses don't match %s %s", ipReservation.IP, ipReservationOutput.IP)
	}
}

func getIPReservationsClient() (*IPReservationsClient, error) {
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
