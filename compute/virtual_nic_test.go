// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_VirtNicInstanceTestName  = "test-acc-virt-nic"
	_VirtNicInstanceTestLabel = "test"
	_VirtNicInstanceTestShape = "oc3"
	_VirtNicInstanceTestImage = "/oracle/public/OL_7.2_UEKR4_x86_64"
)

func TestAccVirtNICLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	iClient, nClient, vnClient, err := getVirtNICsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	ipNetwork, err := createTestIPNetwork(nClient, _IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetwork.Name)

	// In order to get details on a Virtual NIC we need to create the following resources
	// - IP Network
	// - Instance

	instanceInput := &CreateInstanceInput{
		Name:      _VirtNicInstanceTestName,
		Label:     _VirtNicInstanceTestLabel,
		Shape:     _VirtNicInstanceTestShape,
		ImageList: _VirtNicInstanceTestImage,
		Networking: map[string]NetworkingInfo{
			"eth0": {
				IPNetwork: ipNetwork.Name,
				Vnic:      "eth0",
			},
		},
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	// Use the static "eth0" interface, as we statically created that above
	createdVNIC := createdInstance.Networking["eth0"].Vnic
	getInput := &GetVirtualNICInput{
		Name: createdVNIC,
	}

	vNIC, err := vnClient.GetVirtualNIC(getInput)
	if err != nil {
		t.Fatal(err)
	}
	// Don't need to tear down the VNIC, it's attached to the instance
	log.Printf("VNIC Retrieved: %+v", vNIC)
	assert.NotEmpty(t, vNIC.Name, "Expected VNIC name not to be empty")
	assert.Equal(t, createdVNIC, vNIC.Name, "Expected VNIC and name to match.")
	assert.Equal(t, vNIC.FQDN, vnClient.getQualifiedName(vNIC.Name), "Expected FQDN to match qualifed name")
}

func getVirtNICsTestClients() (*InstancesClient, *IPNetworksClient, *VirtNICsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}
	return client.Instances(), client.IPNetworks(), client.VirtNICs(), nil
}
