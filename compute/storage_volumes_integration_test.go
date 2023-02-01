// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccStorageVolumeLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rInt := rand.Int()
	name := fmt.Sprintf("test-acc-storage-volume-lifecycle-%d", rInt)

	createRequest := CreateStorageVolumeInput{
		Name:        name,
		Description: "original description",
		Size:        "20",
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	updateRequest := UpdateStorageVolumeInput{
		Name:        name,
		Size:        "30",
		Description: "updated description",
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	testStorageVolume(t, createRequest, updateRequest)
}

func TestAccStorageVolumeBootableLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rInt := rand.Int()
	name := fmt.Sprintf("test-acc-storage-volume-bootable-lifecycle-%d", rInt)

	imageListName := "/oracle/public/OL_7.2_UEKR4_x86_64"

	createRequest := CreateStorageVolumeInput{
		Name:        name,
		Description: "original description",
		Size:        "20",
		ImageList:   imageListName,
		Properties:  []string{string(StorageVolumeKindDefault)},
	}

	updateRequest := UpdateStorageVolumeInput{
		Name:        name,
		Size:        "30",
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

	assert.Equal(t, createResponse.Size, createdResponse.Size,
		"Retrieved Storage Volume Size did not match Expected.")

	actualSize := createdResponse.Size
	expectedSize := "20"
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

	assert.Equal(t, updateResponse.Size, updatedResponse.Size, "Retrieved Storage Volume did not match Expected.")
	assert.Equal(t, updatedResponse.FQDN, svc.getQualifiedName(createInput.Name), "Expected FDQN to be equal to qualified name")

	actualSize = updatedResponse.Size
	expectedSize = "30"

	if actualSize != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, actualSize)
	}
}
