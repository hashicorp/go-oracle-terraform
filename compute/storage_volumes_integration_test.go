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
	defer tearDownStorageVolumes(name)

	svc, err := getStorageVolumeClient()
	if err != nil {
		t.Fatal(err)
	}

	createRequest := &CreateStorageVolumeInput{
		Name:        name,
		Description: "original description",
		Size:        "10G",
		Properties:  []string{"/oracle/public/storage/default"},
	}
	err = svc.CreateStorageVolume(createRequest)

	if err != nil {
		t.Fatalf("Create volume failed: %s\n", err)
	}

	getRequest := &GetStorageVolumeInput{
		Name: name,
	}
	createResponse, err := svc.GetStorageVolume(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	actualSize := createResponse.Result[0].Size
	expectedSize := strconv.Itoa(10 << 30)
	if actualSize != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, actualSize)
	}

	updateRequest := &UpdateStorageVolumeInput{
		Name:        name,
		Size:        "20G",
		Description: "updated description",
		Properties:  []string{"/oracle/public/storage/default"},
	}
	err = svc.UpdateStorageVolume(updateRequest)
	if err != nil {
		t.Fatal(err)
	}

	updateResponse, err := svc.GetStorageVolume(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	actualSize = updateResponse.Result[0].Size
	expectedSize = strconv.Itoa(20 << 30)
	if actualSize != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, actualSize)
	}
}

func tearDownStorageVolumes(name string) {
	svc, err := getStorageVolumeClient()
	if err != nil {
		panic(err)
	}

	deleteRequest := &DeleteStorageVolumeInput{
		Name: name,
	}
	err = svc.DeleteStorageVolume(deleteRequest)
	if err != nil {
		panic(err)
	}
}

func getStorageVolumeClient() (*StorageVolumeClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &StorageVolumeClient{}, err
	}

	return client.StorageVolumes(), nil
}
