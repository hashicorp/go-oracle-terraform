// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package java

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccIPReservationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	resClient, err := getIPResvervationClient()
	if err != nil {
		t.Fatal(err)
	}

	ipReservation := CreateIPReservationInput{
		Name:             "test-ip-reservation",
		Region:           "uscom-central-1",
		IdentityDomainID: *resClient.client.IdentityDomain,
	}

	_, err = resClient.CreateIPReservation(&ipReservation)
	assert.NoError(t, err)

	defer destroyIPReservation(t, resClient, ipReservation.Name)

	resp, err := resClient.GetIPReservation(ipReservation.Name)
	assert.NoError(t, err)

	assert.Equal(t, ipReservation.Name, resp.Name, "IP Reservation Name should match")
	assert.Equal(t, ipReservation.Region, resp.ComputeSiteName, "IP Reservation Compute Site Name should match region")
	assert.Equal(t, ipReservation.IdentityDomainID, resp.IdentityDomain, "IP Reservation Identity Domain should match")

}

func getIPResvervationClient() (*IPReservationClient, error) {
	client, err := getJavaTestClient(&opc.Config{})
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
