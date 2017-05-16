package compute

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_SnapshotInstanceTestName  = "test-acc-snapshot"
	_SnapshotInstanceTestLabel = "test"
	_SnapshotInstanceTestShape = "oc3"
	_SnapshotInstanceTestImage = "/oracle/public/Oracle_Solaris_11.3"
)

func TestAccSnapshotLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, iClient, err := getSnapshotsTestClients()
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

	time.Sleep(300*time.Second)
	createSnapshotInput := &CreateSnapshotInput{
		//Account:      "cloud_storage",
		Instance: strings.Join([]string{createdInstance.Name, createdInstance.ID}, "/"),
		//MachineImage: _SnapshotInstanceTestName,
	}
	createdSnapshot, err := sClient.CreateSnapshot(createSnapshotInput)
	if err != nil {
		t.Fatal(fmt.Sprintf("Snapshot: %+v", createSnapshotInput))
		t.Fatal(err)
	}
	defer tearDownSnapshots(t, sClient, createdInstance.Name)

	getInput := &GetSnapshotInput{
		Name: createdSnapshot.Name,
	}

	snapshot, err := sClient.GetSnapshot(getInput)
	if err != nil {
		t.Fatal(err)
	}
	// Don't need to tear down the Snapshot, it's attached to the instance
	log.Printf("Snapshot Retrieved: %+v", snapshot)
	if !reflect.DeepEqual(snapshot.Name, createdSnapshot.Name) {
		t.Fatal("Snapshot Name mismatch! Got: %s Expected: %s", snapshot.Name, createdSnapshot.Name)
	}
}

func tearDownSnapshots(t *testing.T, snapshotsClient *SnapshotsClient, name string) {
	log.Printf("Deleting Snapshot %s", name)

	deleteRequest := &DeleteSnapshotInput{
		Name: name,
	}
	if err := snapshotsClient.DeleteSnapshot(deleteRequest); err != nil {
		t.Fatalf("Error deleting snapshot, dangling resources may occur: %v", err)
	}
}

func getSnapshotsTestClients() (*SnapshotsClient, *InstancesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}
	return client.Snapshots(), client.Instances(), nil
}
