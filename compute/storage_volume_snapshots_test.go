package compute

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccStorageVolumeSnapshot_Lifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rInt := rand.Int()
	volumeName := fmt.Sprintf("test-acc-storage-snapshot-base-%d", rInt)
	snapshotName := fmt.Sprintf("test-acc-storage-snapshot-%d", rInt)

	vClient, snapClient, err := getSnapshotTestClients()
	if err != nil {
		t.Fatal(err)
	}

	sVolumeInput := &CreateStorageVolumeInput{
		Name:        volumeName,
		Description: "test-acc",
		Size:        "20",
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	sVolume, err := vClient.CreateStorageVolume(sVolumeInput)
	if err != nil {
		t.Fatalf("Create Volume failed: %v", err)
	}

	snapshotInput := &CreateStorageVolumeSnapshotInput{
		Description: "testing-acc",
		Name:        snapshotName,
		Property:    SnapshotPropertyCollocated,
		Tags:        []string{"tag1"},
		Volume:      sVolume.Name,
	}

	snapshot, err := snapClient.CreateStorageVolumeSnapshot(snapshotInput)
	if err != nil {
		// Attempt a destroy here with the supplied name in case the snapshot exists, but isn't ready
		tearDownStorageSnapshot(t, snapClient, vClient, snapshotInput.Name, sVolume.Name)
		t.Fatalf("Create Snapshot failed: %v", err)
	}
	defer tearDownStorageSnapshot(t, snapClient, vClient, snapshot.Name, sVolume.Name)

	getInput := &GetStorageVolumeSnapshotInput{
		Name: snapshot.Name,
	}

	getRes, err := snapClient.GetStorageVolumeSnapshot(getInput)
	if err != nil {
		t.Fatalf("Error getting storage snapshot: %v", err)
	}

	assert.Equal(t, getRes, snapshot, "Mismatch after Create.")
	assert.Equal(t, getRes.FQDN, snapClient.getQualifiedName(snapshot.Name), "Expected FDQN to be equal to qualified name")

}

func getSnapshotTestClients() (*StorageVolumeClient, *StorageVolumeSnapshotClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}

	return client.StorageVolumes(), client.StorageVolumeSnapshots(), nil
}

func tearDownStorageSnapshot(t *testing.T, snapClient *StorageVolumeSnapshotClient, volumeClient *StorageVolumeClient, snapshotName, volumeName string) {
	input := &DeleteStorageVolumeSnapshotInput{
		Name: snapshotName,
	}
	if err := snapClient.DeleteStorageVolumeSnapshot(input); err != nil {
		t.Fatalf("Error deleting Storage Volume Snapshot '%s', dangling resources may occur: %v", snapshotName, err)
	}
	// If snapshot deleted successfully, it's now safe to destroy the volume
	tearDownStorageVolumes(t, volumeClient, volumeName)
}
