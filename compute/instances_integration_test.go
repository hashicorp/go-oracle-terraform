package compute

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

const (
	_InstanceTestName       = "test-acc"
	_InstanceTestLabel      = "test"
	_InstanceTestShape      = "oc3"
	_InstanceTestImage      = "/oracle/public/OL_7.2_UEKR4_x86_64"
	_InstanceTestImageEntry = 4
	_InstanceTestPublicPool = "ippool:/oracle/public/ippool"
)

func TestAccInstanceLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	iClient, ipaClient, nClient, err := getInstancesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	ipNetwork, err := createTestIPNetwork(nClient, _IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetwork.Name)

	rInt := rand.Int()
	resName := fmt.Sprintf("%s-%d", _TestIPAddressResName, rInt)

	ipresInput := &CreateIPAddressReservationInput{
		Description:   _TestIPAddressResDesc,
		IPAddressPool: PrivateIPAddressPool,
		Name:          resName,
	}

	ipRes, err := ipaClient.CreateIPAddressReservation(ipresInput)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPAddressReservation(t, ipaClient, ipRes.Name)

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
				Nat:       []string{ipRes.Name},
			},
			"eth1": {
				Model: NICDefaultModel,
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
		Name:         _InstanceTestName,
		Label:        _InstanceTestLabel,
		Shape:        _InstanceTestShape,
		ImageList:    _InstanceTestImage,
		DesiredState: "running",
		Entry:        _InstanceTestImageEntry,
		ImageFormat:  "raw",
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
		State:         InstanceRunning,
		Storage:       []StorageAttachment{},
		Tags:          []string{},
		Virtio:        false,
	}

	if err = verifyInstance(expectedInstance, receivedInstance, ipRes, ipNetwork); err != nil {
		t.Fatal(err)
	}

	updateInput := &UpdateInstanceInput{
		Name: createdInstance.Name,
		ID:   createdInstance.ID,
		Tags: []string{"new_tag1", "new_tag2"},
	}

	updatedInstance, err := iClient.UpdateInstance(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	getInput = &GetInstanceInput{
		Name: updatedInstance.Name,
		ID:   updatedInstance.ID,
	}

	receivedInstance, err = iClient.GetInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	// Update Expected instance + verify
	expectedInstance.Tags = []string{"new_tag1", "new_tag2"}
	if err := verifyInstance(expectedInstance, receivedInstance, ipRes, ipNetwork); err != nil {
		t.Fatal(err)
	}

}

// Test that we can shutdown and startup an instance
func TestAccInstanceStopStart(t *testing.T) {
	rInt := rand.Int()

	helper.Test(t, helper.TestCase{})

	// Setup Instance Client
	iClient, sClient, err := getInstanceStartStopTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// Create the bootable storage volume
	volumeName := fmt.Sprintf("%s-volume-%d", _InstanceTestName, rInt)
	volumeInput := &CreateStorageVolumeInput{
		Name:           volumeName,
		Size:           "20",
		ImageList:      _InstanceTestImage,
		ImageListEntry: _InstanceTestImageEntry,
		Bootable:       true,
		Properties:     []string{string(StorageVolumeKindDefault)},
	}

	storageVolume, err := sClient.CreateStorageVolume(volumeInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownStorageVolumes(t, sClient, volumeName)

	// Finally, create Basic Testing Instance with a root storage volume
	input := &CreateInstanceInput{
		Name:      _InstanceTestName,
		Label:     _InstanceTestLabel,
		Shape:     _InstanceTestShape,
		BootOrder: []int{1},
		Storage: []StorageAttachmentInput{
			{
				Index:  1,
				Volume: storageVolume.Name,
			},
		},
	}

	createdInstance, err := iClient.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)
	log.Printf("Instance Created: %#v\n", createdInstance)

	// Instance is created, time to shut it down
	updateInput := &UpdateInstanceInput{
		Name:         createdInstance.Name,
		ID:           createdInstance.ID,
		DesiredState: InstanceDesiredShutdown,
	}

	updatedInstance, err := iClient.UpdateInstance(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	getInput := &GetInstanceInput{
		Name: updatedInstance.Name,
		ID:   updatedInstance.ID,
	}

	info, err := iClient.GetInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if info.DesiredState != InstanceDesiredShutdown {
		t.Fatalf("Instance desired state should be `shutdown`, got: %s", info.DesiredState)
	}

	// Verify the instance is actually shut down
	if info.State != InstanceShutdown {
		t.Fatalf("Instance should be in the `shutdown` state, got: %s", info.State)
	}

	// Instance is verified to be shutdown, spin it back up
	rebootInput := &UpdateInstanceInput{
		Name:         updatedInstance.Name,
		ID:           updatedInstance.ID,
		DesiredState: InstanceDesiredRunning,
	}

	rebootedInstance, err := iClient.UpdateInstance(rebootInput)
	if err != nil {
		t.Fatal(err)
	}

	getInput = &GetInstanceInput{
		Name: rebootedInstance.Name,
		ID:   rebootedInstance.ID,
	}

	info, err = iClient.GetInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if info.DesiredState != InstanceDesiredRunning {
		t.Fatalf("Instance desired state should be `running`, got: %s", info.DesiredState)
	}

	if info.State != InstanceRunning {
		t.Fatalf("Instance should be in the `running` state, got: %s", info.State)
	}
	// All pass, let defer cleanup instance
}

func verifyInstance(expected, received *InstanceInfo, ipRes *IPAddressReservation, ipNetwork *IPNetworkInfo) error {
	// Verify Networking before zero
	if received.Networking["eth1"].Model != "e1000" {
		return fmt.Errorf("Expected Networking model to be e1000, got: %s", received.Networking["eth1"].Model)
	}

	if received.Networking["eth0"].IPNetwork != ipNetwork.Name {
		return fmt.Errorf("Expected IPNetwork name %s, got: %s", ipNetwork.Name, received.Networking["eth0"].IPNetwork)
	}

	if diff := pretty.Compare(received.Networking["eth0"].Nat, []string{ipRes.Name}); diff != "" {
		return fmt.Errorf("Networking Diff: (-got +want)\n%s", diff)
	}

	// Zero the fields we can't statically check for
	received.zeroFields()

	if diff := pretty.Compare(received, expected); diff != "" {
		return fmt.Errorf("Created Instance Diff: (-got +want)\n%s", diff)
	}
	// Verify
	if !received.ReverseDNS {
		return fmt.Errorf("Expected ReverseDNS to have default 'true' value. Got False")
	}

	return nil
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

func getInstancesTestClients() (*InstancesClient, *IPAddressReservationsClient, *IPNetworksClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}

	return client.Instances(), client.IPAddressReservations(), client.IPNetworks(), nil
}

func getInstanceStartStopTestClients() (*InstancesClient, *StorageVolumeClient, error) {

	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}

	return client.Instances(), client.StorageVolumes(), nil
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
