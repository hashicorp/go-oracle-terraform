package compute

import (
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_SnapshotInstanceTestName    = "test-acc-snapshot"
	_SnapshotInstanceTestLabel   = "test"
	_SnapshotInstanceTestShape   = "oc3"
	_SnapshotInstanceTestImage   = "/oracle/public/JEOS_OL_6.6_10GB_RD-1.2.217-20151201-194209"
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
	if !reflect.DeepEqual(snapshot.Name, createdSnapshot.Name) {
		t.Fatal("Snapshot Name mismatch! Got: %s Expected: %s", snapshot.Name, createdSnapshot.Name)
	}
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
	if !reflect.DeepEqual(snapshot.Name, createdSnapshot.Name) {
		t.Fatal("Snapshot Name mismatch! Got: %s Expected: %s", snapshot.Name, createdSnapshot.Name)
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
	if err := iClient.DeleteInstance(deleteInstanceInput); err != nil {
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
	if !reflect.DeepEqual(snapshot.Name, createdSnapshot.Name) {
		t.Fatal("Snapshot Name mismatch! Got: %s Expected: %s", snapshot.Name, createdSnapshot.Name)
	}
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
	if !reflect.DeepEqual(snapshot.Name, createdSnapshot.Name) {
		t.Fatal("Snapshot Name mismatch! Got: %s Expected: %s", snapshot.Name, createdSnapshot.Name)
	}
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

func getSnapshotsTestClients() (*SnapshotsClient, *InstancesClient, *MachineImagesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}
	return client.Snapshots(), client.Instances(), client.MachineImages(), nil
}
