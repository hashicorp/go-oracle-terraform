package compute

import (
	"strconv"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccStorageVolumeLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-storage-volume-lifecycle"

	svc, err := getStorageVolumeClient()
	if err != nil {
		t.Fatal(err)
	}

	createRequest := CreateStorageVolumeInput{
		Name:        name,
		Description: "original description",
		Size:        "10G",
		Properties:  []string{"/oracle/public/storage/default"},
	}
	createResponse, err := svc.CreateStorageVolume(&createRequest)
	if err != nil {
		t.Fatalf("Create volume failed: %s\n", err)
	}

	defer tearDownStorageVolumes(t, svc, name)

	getRequest := &GetStorageVolumeInput{
		Name: name,
	}
	createdResponse, err := svc.GetStorageVolume(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	if createResponse.Size != createdResponse.Size {
		t.Fatalf("Retrieved Storage Volume Size did not match Expected. \nDesired: %s \nActual: %s", createResponse, createdResponse)
	}

	actualSize := createdResponse.Size
	expectedSize := strconv.Itoa(10 << 30)
	if actualSize != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, actualSize)
	}

	updateRequest := UpdateStorageVolumeInput{
		Name:        name,
		Size:        "20G",
		Description: "updated description",
		Properties:  []string{"/oracle/public/storage/default"},
	}
	updateResponse, err := svc.UpdateStorageVolume(&updateRequest)
	if err != nil {
		t.Fatal(err)
	}

	updatedResponse, err := svc.GetStorageVolume(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	if updateResponse.Size != updatedResponse.Size {
		t.Fatalf("Retrieved Storage Volume did not match Expected. \nDesired: %s \nActual: %s", updateResponse, updatedResponse)
	}

	actualSize = updatedResponse.Size
	expectedSize = strconv.Itoa(20 << 30)
	if actualSize != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, actualSize)
	}
}

func tearDownStorageVolumes(t *testing.T, svc *StorageVolumeClient, name string) {
	deleteRequest := &DeleteStorageVolumeInput{
		Name: name,
	}
	if err := svc.DeleteStorageVolume(deleteRequest); err != nil {
		t.Fatalf("Error deleting storage volume, dangling resources may occur: %v", err)
	}
}

func getStorageVolumeClient() (*StorageVolumeClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &StorageVolumeClient{}, err
	}

	return client.StorageVolumes(), nil
}
