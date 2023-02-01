// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_VPNEndpointV2TestName               = "test-acc-vpn-endpoint-v2"
	_VPNEndpointV2TestCustomerVPNGateway = "127.0.0.1"
	_VPNEndpointV2TestPSK                = "asdfasdf"
)

func TestAccVPNEndpointV2sLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	vpnClient, iClient, err := getVPNEndpointV2sClients()
	if err != nil {
		t.Fatal(err)
	}

	ipNetwork, err := createTestIPNetwork(iClient, _IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, iClient, ipNetwork.Name)

	createInput := &CreateVPNEndpointV2Input{
		Name:               _VPNEndpointV2TestName,
		CustomerVPNGateway: _VPNEndpointV2TestCustomerVPNGateway,
		IPNetwork:          ipNetwork.Name,
		PSK:                _VPNEndpointV2TestPSK,
		ReachableRoutes:    []string{"127.0.0.1/24"},
		VNICSets:           []string{"default"},
	}

	// Create a vNIC Set

	createdVPNEndpointV2, err := vpnClient.CreateVPNEndpointV2(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("VPN Endpoint V2 succcessfully created")
	defer destroyVPNEndpointV2(t, vpnClient, _VPNEndpointV2TestName)

	getInput := &GetVPNEndpointV2Input{
		Name: _VPNEndpointV2TestName,
	}
	receivedVPNEndpointV2, err := vpnClient.GetVPNEndpointV2(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("VPN Endpoint V2 successfully fetched")

	if !reflect.DeepEqual(createdVPNEndpointV2.URI, receivedVPNEndpointV2.URI) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdVPNEndpointV2, receivedVPNEndpointV2)
	}

	updateInput := &UpdateVPNEndpointV2Input{
		Name:               _VPNEndpointV2TestName,
		CustomerVPNGateway: _VPNEndpointV2TestCustomerVPNGateway,
		IPNetwork:          ipNetwork.Name,
		PSK:                _VPNEndpointV2TestPSK,
		ReachableRoutes:    []string{"127.0.0.1/10"},
		VNICSets:           []string{"default"},
	}
	updatedVPNEndpointV2, err := vpnClient.UpdateVPNEndpointV2(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("VPN Endpoint V2 succcessfully updated")
	receivedVPNEndpointV2, err = vpnClient.GetVPNEndpointV2(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(updatedVPNEndpointV2.ReachableRoutes, receivedVPNEndpointV2.ReachableRoutes) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", updatedVPNEndpointV2, receivedVPNEndpointV2)
	}
}

func destroyVPNEndpointV2(t *testing.T, vpnClient *VPNEndpointV2sClient, name string) {
	input := &DeleteVPNEndpointV2Input{
		Name: name,
	}
	if err := vpnClient.DeleteVPNEndpointV2(input); err != nil {
		t.Fatal(err)
	}
}

func getVPNEndpointV2sClients() (*VPNEndpointV2sClient, *IPNetworksClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}

	return client.VPNEndpointV2s(), client.IPNetworks(), nil
}
