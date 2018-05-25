package compute

import (
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_SnapshotInstanceTestName    = "test-acc-snapshot"
	_SnapshotInstanceTestLabel   = "test"
	_SnapshotInstanceTestShape   = "oc3"
	_SnapshotInstanceTestImage   = "/oracle/public/OL_7.2_UEKR4_x86_64"
	_SnapshotInstanceTestAccount = "cloud_storage"
)

func TestAccSnapshotLifeCycleBasic(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, iClient, mClient, err := getSnapshotsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// In order to get details on a Snapshot we need to create the following resources
	// - Instance

	instanceInput := &CreateInstanceInput{
		Name:      _SnapshotInstanceTestName,
		Label:     _SnapshotInstanceTestLabel,
		Shape:     _SnapshotInstanceTestShape,
		ImageList: _SnapshotInstanceTestImage,
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	createSnapshotInput := &CreateSnapshotInput{
		Instance: strings.Join([]string{createdInstance.Name, createdInstance.ID}, "/"),
	}
	createdSnapshot, err := sClient.CreateSnapshot(createSnapshotInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownSnapshots(t, sClient, mClient, createdSnapshot.Name, createdSnapshot.MachineImage)

	getInput := &GetSnapshotInput{
		Name: createdSnapshot.Name,
	}

	snapshot, err := sClient.GetSnapshot(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Snapshot Retrieved: %+v", snapshot)
	assert.Equal(t, createdSnapshot.Name, snapshot.Name, "Snapshot Name mismatch!")
}

func TestAccSnapshotLifeCycleMachineImage(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, iClient, mClient, err := getSnapshotsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// In order to get details on a Snapshot we need to create the following resources
	// - Instance

	instanceInput := &CreateInstanceInput{
		Name:      _SnapshotInstanceTestName,
		Label:     _SnapshotInstanceTestLabel,
		Shape:     _SnapshotInstanceTestShape,
		ImageList: _SnapshotInstanceTestImage,
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	createSnapshotInput := &CreateSnapshotInput{
		Instance:     strings.Join([]string{createdInstance.Name, createdInstance.ID}, "/"),
		MachineImage: _SnapshotInstanceTestName,
	}
	createdSnapshot, err := sClient.CreateSnapshot(createSnapshotInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownSnapshots(t, sClient, mClient, createdSnapshot.Name, createdSnapshot.MachineImage)

	getInput := &GetSnapshotInput{
		Name: createdSnapshot.Name,
	}

	snapshot, err := sClient.GetSnapshot(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Snapshot Retrieved: %+v", snapshot)
	assert.Equal(t, createdSnapshot.Name, snapshot.Name, "Snapshot Name mismatch!")
}

func TestAccSnapshotLifeCycleNoDeleteMachineImage(t *testing.T) {
	// Same test as above; different tear-down
	helper.Test(t, helper.TestCase{})

	sClient, iClient, mClient, err := getSnapshotsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// In order to get details on a Snapshot we need to create the following resources
	// - Instance

	instanceInput := &CreateInstanceInput{
		Name:      _SnapshotInstanceTestName,
		Label:     _SnapshotInstanceTestLabel,
		Shape:     _SnapshotInstanceTestShape,
		ImageList: _SnapshotInstanceTestImage,
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	createSnapshotInput := &CreateSnapshotInput{
		Instance:     strings.Join([]string{createdInstance.Name, createdInstance.ID}, "/"),
		MachineImage: _SnapshotInstanceTestName,
	}
	createdSnapshot, err := sClient.CreateSnapshot(createSnapshotInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownSnapshotResource(t, sClient, createdSnapshot.Name, createdSnapshot.MachineImage)
	defer tearDownMachineImage(t, mClient, createdSnapshot.MachineImage)

	getInput := &GetSnapshotInput{
		Name: createdSnapshot.Name,
	}

	snapshot, err := sClient.GetSnapshot(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Snapshot Retrieved: %+v", snapshot)
	if !reflect.DeepEqual(snapshot.Name, createdSnapshot.Name) {
		t.Fatalf("Snapshot Name mismatch! Got: %s Expected: %s", snapshot.Name, createdSnapshot.Name)
	}
}

func TestAccSnapshotLifeCycleDelay(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, iClient, mClient, err := getSnapshotsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// In order to get details on a Snapshot we need to create the following resources
	// - Instance

	instanceInput := &CreateInstanceInput{
		Name:      _SnapshotInstanceTestName,
		Label:     _SnapshotInstanceTestLabel,
		Shape:     _SnapshotInstanceTestShape,
		ImageList: _SnapshotInstanceTestImage,
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}

	createSnapshotInput := &CreateSnapshotInput{
		Instance: strings.Join([]string{createdInstance.Name, createdInstance.ID}, "/"),
		Delay:    SnapshotDelayShutdown,
	}
	createdSnapshot, err := sClient.CreateSnapshot(createSnapshotInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownSnapshots(t, sClient, mClient, createdSnapshot.Name, createdSnapshot.MachineImage)

	// Snapshot creation only finishes after the instance has been deleted.
	deleteInstanceInput := &DeleteInstanceInput{
		Name: createdInstance.Name,
		ID:   createdInstance.ID,
	}
	if err = iClient.DeleteInstance(deleteInstanceInput); err != nil {
		t.Fatalf("Error deleting instance, dangling resources may occur: %v", err)
	}

	getInput := &GetSnapshotInput{
		Name: createdSnapshot.Name,
	}

	snapshot, err := sClient.GetSnapshot(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Snapshot Retrieved: %+v", snapshot)
	assert.Equal(t, createdSnapshot.Name, snapshot.Name, "Snapshot Name mismatch!")
}

func TestAccSnapshotLifeCycleAccount(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, iClient, mClient, err := getSnapshotsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// In order to get details on a Snapshot we need to create the following resources
	// - Instance

	instanceInput := &CreateInstanceInput{
		Name:      _SnapshotInstanceTestName,
		Label:     _SnapshotInstanceTestLabel,
		Shape:     _SnapshotInstanceTestShape,
		ImageList: _SnapshotInstanceTestImage,
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	createSnapshotInput := &CreateSnapshotInput{
		Instance: strings.Join([]string{createdInstance.Name, createdInstance.ID}, "/"),
		Account:  _SnapshotInstanceTestAccount,
	}
	createdSnapshot, err := sClient.CreateSnapshot(createSnapshotInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownSnapshots(t, sClient, mClient, createdSnapshot.Name, createdSnapshot.MachineImage)

	getInput := &GetSnapshotInput{
		Name: createdSnapshot.Name,
	}

	snapshot, err := sClient.GetSnapshot(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Snapshot Retrieved: %+v", snapshot)
	assert.Equal(t, createdSnapshot.Name, snapshot.Name, "Snapshot Name mismatch!")
}

func tearDownSnapshots(t *testing.T, snapshotsClient *SnapshotsClient, machineImagesClient *MachineImagesClient, snapshotName string, machineImageName string) {
	log.Printf("Deleting Snapshot %s", snapshotName)

	deleteRequest := &DeleteSnapshotInput{
		Snapshot:     snapshotName,
		MachineImage: machineImageName,
	}
	if err := snapshotsClient.DeleteSnapshot(machineImagesClient, deleteRequest); err != nil {
		t.Fatalf("Error deleting snapshot, dangling resources may occur: %v", err)
	}
}

func tearDownSnapshotResource(t *testing.T, snapshotsClient *SnapshotsClient, snapshotName string, machineImageName string) {
	log.Printf("Deleting Snapshot %s", snapshotName)

	deleteRequest := &DeleteSnapshotInput{
		Snapshot:     snapshotName,
		MachineImage: machineImageName,
	}
	if err := snapshotsClient.DeleteSnapshotResourceOnly(deleteRequest); err != nil {
		t.Fatalf("Error removing snapshot, dangling resources may occur: %v", err)
	}
}

func tearDownMachineImage(t *testing.T, machineImagesClient *MachineImagesClient, machineImageName string) {
	log.Printf("Deleting Machine Image %s", machineImageName)

	deleteInput := &DeleteMachineImageInput{
		Name: machineImageName,
	}
	if err := machineImagesClient.DeleteMachineImage(deleteInput); err != nil {
		t.Fatalf("Error deleting snapshot, dangling resources may occur: %v", err)
	}
}

func getSnapshotsTestClients() (*SnapshotsClient, *InstancesClient, *MachineImagesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}
	return client.Snapshots(), client.Instances(), client.MachineImages(), nil
}
