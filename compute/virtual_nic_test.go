package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_VirtNicInstanceTestName  = "test-acc-virt-nic"
	_VirtNicInstanceTestLabel = "test"
	_VirtNicInstanceTestShape = "oc3"
	_VirtNicInstanceTestImage = "/oracle/public/Oracle_Solaris_11.3"
)

func TestAccVirtNICLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	ipNetwork, err := createTestIPNetwork(_IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, ipNetwork.Name)

	instanceSvc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

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

	createdInstance, err := instanceSvc.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, instanceSvc, createdInstance.Name, createdInstance.ID)

	svc, err := getVirtNICsClient()
	if err != nil {
		t.Fatal(err)
	}

	// Use the static "eth0" interface, as we statically created that above
	createdVNIC := createdInstance.Networking["eth0"].Vnic
	getInput := &GetVirtualNICInput{
		Name: createdVNIC,
	}

	vNIC, err := svc.GetVirtualNIC(getInput)
	if err != nil {
		t.Fatal(err)
	}
	// Don't need to tear down the VNIC, it's attached to the instance
	log.Printf("VNIC Retrieved: %+v", vNIC)
	if vNIC.Name != createdVNIC || vNIC.Name == "" {
		t.Fatal("VNIC Name mismatch! Got: %s Expected: %s", vNIC.Name, createdVNIC)
	}
}

func getVirtNICsClient() (*VirtNICsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}
	return client.VirtNICs(), nil
}
