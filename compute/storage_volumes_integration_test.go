package compute

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccStorageVolumeLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rInt := rand.Int()
	name := fmt.Sprintf("test-acc-storage-volume-lifecycle-%d", rInt)

	createRequest := CreateStorageVolumeInput{
		Name:        name,
		Description: "original description",
		Size:        "10240",
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	updateRequest := UpdateStorageVolumeInput{
		Name:        name,
		Size:        "20480",
		Description: "updated description",
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	testStorageVolume(t, createRequest, updateRequest)
}

func TestAccStorageVolumeBootableLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rInt := rand.Int()
	name := fmt.Sprintf("test-acc-storage-volume-bootable-lifecycle-%d", rInt)

	imageListName := fmt.Sprintf("test-acc-storage-volume-bootable-lifecycle-il-%d", rInt)

	imageListClient, err := getImageListClient()
	if err != nil {
		t.Fatalf("Error building Image List Client: %+v", err)
	}

	input := CreateImageListInput{
		Name:        imageListName,
		Description: "Test from the TestAccStorageVolumeBootableLifecycle",
		Default:     1,
	}
	_, err = imageListClient.CreateImageList(&input)
	if err != nil {
		t.Fatalf("Error Creating Image List: %+v", err)
	}
	defer tearDownImageList(t, imageListClient, imageListName)
	createRequest := CreateStorageVolumeInput{
		Name:        name,
		Description: "original description",
		Size:        "10240",
		ImageList:   imageListName,
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	updateRequest := UpdateStorageVolumeInput{
		Name:        name,
		Size:        "20480",
		Description: "updated description",
		ImageList:   imageListName,
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	testStorageVolume(t, createRequest, updateRequest)
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

func testStorageVolume(t *testing.T, createInput CreateStorageVolumeInput, updateInput UpdateStorageVolumeInput) {
	svc, err := getStorageVolumeClient()
	if err != nil {
		t.Fatal(err)
	}

	createResponse, err := svc.CreateStorageVolume(&createInput)
	if err != nil {
		t.Fatalf("Create volume failed: %s\n", err)
	}

	defer tearDownStorageVolumes(t, svc, createInput.Name)

	getRequest := &GetStorageVolumeInput{
		Name: createInput.Name,
	}
	createdResponse, err := svc.GetStorageVolume(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	if createResponse.Size != createdResponse.Size {
		t.Fatalf("Retrieved Storage Volume Size did not match Expected. \nDesired: %s \nActual: %s", createResponse, createdResponse)
	}

	actualSize := createdResponse.Size
	expectedSize := "10240"
	if actualSize != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, actualSize)
	}

	updateResponse, err := svc.UpdateStorageVolume(&updateInput)
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
	expectedSize = "20480"

	if actualSize != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, actualSize)
	}
}
