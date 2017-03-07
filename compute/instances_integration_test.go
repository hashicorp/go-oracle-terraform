package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

const (
	_InstanceTestName       = "test-acc"
	_InstanceTestLabel      = "test"
	_InstanceTestShape      = "oc3"
	_InstanceTestImage      = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"
	_InstanceTestPublicPool = "ippool:/oracle/public/ippool"
)

func TestAccInstanceLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	iClient, nClient, err := getInstancesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	ipNetwork, err := createTestIPNetwork(nClient, _IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetwork.Name)

	// TODO (@jake): Remove static IPReservation once resource is available for modification
	input := &CreateInstanceInput{
		Name:      _InstanceTestName,
		Label:     _InstanceTestLabel,
		Shape:     _InstanceTestShape,
		ImageList: _InstanceTestImage,
		Storage:   nil,
		BootOrder: nil,
		SSHKeys:   []string{},
		Networking: map[string]NetworkingInfo{
			"eth0": {
				IPNetwork: ipNetwork.Name,
				Nat:       []string{"testing-acc"},
			},
			"eth1": {
				Model: "e1000",
				Nat:   []string{_InstanceTestPublicPool},
			},
		},
		Attributes: map[string]interface{}{
			"attr1": 12,
			"attr2": map[string]interface{}{
				"inner_attr1": "foo",
			},
		},
	}

	createdInstance, err := iClient.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	log.Printf("Instance created: %#v\n", createdInstance)

	getInput := &GetInstanceInput{
		Name: createdInstance.Name,
		ID:   createdInstance.ID,
	}

	receivedInstance, err := iClient.GetInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	expectedInstance := &InstanceInfo{
		Name:        _InstanceTestName,
		Label:       _InstanceTestLabel,
		Shape:       _InstanceTestShape,
		ImageList:   _InstanceTestImage,
		Entry:       1,
		ImageFormat: "raw",
		PlacementRequirements: []string{
			"/system/compute/placement/default",
			"/system/compute/allow_instances",
		},
		Platform:      "linux",
		Priority:      "/oracle/public/default",
		Relationships: []string{},
		ReverseDNS:    true,
		Site:          "",
		SSHKeys:       []string{},
		State:         "running",
		Storage:       []StorageAttachment{},
		Tags:          []string{},
		Virtio:        false,
	}

	// Verify Networking before zero
	if receivedInstance.Networking["eth1"].Model != "e1000" {
		t.Fatalf("Expected Networking model to be e1000, got: %s", receivedInstance.Networking["eth1"].Model)
	}

	if receivedInstance.Networking["eth0"].IPNetwork != ipNetwork.Name {
		t.Fatalf("Expected IPNetwork name %s, got: %s", ipNetwork.Name, receivedInstance.Networking["eth0"].IPNetwork)
	}

	if diff := pretty.Compare(receivedInstance.Networking["eth0"].Nat, []string{"testing-acc"}); diff != "" {
		t.Fatalf("Networking Diff: (-got +want)\n%s", diff)
	}

	// Zero the fields we can't statically check for
	receivedInstance.zeroFields()

	if diff := pretty.Compare(receivedInstance, expectedInstance); diff != "" {
		t.Errorf("Created Instance Diff: (-got +want)\n%s", diff)
	}
	// Verify
	if !receivedInstance.ReverseDNS {
		t.Fatal("Expected ReverseDNS to have default 'true' value. Got False")
	}

}

func tearDownInstances(t *testing.T, svc *InstancesClient, name, id string) {
	input := &DeleteInstanceInput{
		Name: name,
		ID:   id,
	}
	if err := svc.DeleteInstance(input); err != nil {
		t.Fatalf("Error deleting instance, dangling resources may occur: %v", err)
	}
}

func getInstancesTestClients() (*InstancesClient, *IPNetworksClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}

	return client.Instances(), client.IPNetworks(), nil
}

// Zero fields that we cannot check with a static struct
func (i *InstanceInfo) zeroFields() {
	i.ID = ""
	i.Attributes = map[string]interface{}{}
	i.AvailabilityDomain = ""
	i.Domain = ""
	i.Hostname = ""
	i.IPAddress = ""
	i.Networking = map[string]NetworkingInfo{}
	i.StartTime = ""
	i.VCableID = ""
	i.VNC = ""
}
