package database

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccIPReservationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	resClient, err := getIPResvervationClient()
	if err != nil {
		t.Fatal(err)
	}

	ipReservation := CreateIPReservationInput{
		Name: "test-ip-reservation",
	}

	_, err = resClient.CreateIPReservation(&ipReservation)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyIPReservation(t, resClient, ipReservation.Name)

	receivedRes, err := resClient.GetIPReservation(ipReservation.Name)
	if err != nil {
		t.Fatal(err)
	}
	if receivedRes.Name != ipReservation.Name {
		t.Fatal(fmt.Errorf("Names do not match. Wanted: %s Received: %s", ipReservation.Name, receivedRes.Name))
	}

}

func getIPResvervationClient() (*IPReservationClient, error) {
	client, err := GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return &IPReservationClient{}, err
	}

	return client.IPReservationClient(), nil
}

func destroyIPReservation(t *testing.T, client *IPReservationClient, name string) {
	if err := client.DeleteIPReservation(name); err != nil {
		t.Fatal(err)
	}
}
