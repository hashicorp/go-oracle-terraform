package compute

import (
	"log"
	"reflect"
	"testing"

	"fmt"
	"math/rand"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccStorageAttachmentsLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rInt := rand.Int()
	instanceName := fmt.Sprintf("test-acc-stor-att-instance-%d", rInt)
	volumeName := fmt.Sprintf("test-acc-stor-att-volume-%d", rInt)

	instancesClient, storageVolumesClient, attachmentsClient, err := buildStorageAttachmentsClients()
	if err != nil {
		panic(err)
	}

	createInstanceInput := &CreateInstanceInput{
		Name:      instanceName,
		Label:     "test-acc-stor-acc-lifecycle",
		Shape:     "oc3",
		ImageList: "/oracle/public/oel_6.7_apaas_16.4.5_1610211300",
		Storage:   nil,
		BootOrder: nil,
		SSHKeys:   []string{},
		Attributes: map[string]interface{}{
			"attr1": 12,
			"attr2": map[string]interface{}{
				"inner_attr1": "foo",
			},
		},
	}

	info, err := instancesClient.CreateInstance(createInstanceInput)
	if err != nil {
		panic(err)
	}
	defer tearDownInstances(t, instancesClient, info.Name, info.ID)

	createStorageVolumeInput := &CreateStorageVolumeInput{
		Name:       volumeName,
		Size:       "10G",
		Properties: []string{"/oracle/public/storage/default"},
	}
	_, err = storageVolumesClient.CreateStorageVolume(createStorageVolumeInput)
	if err != nil {
		panic(err)
	}

	defer tearDownStorageVolumes(t, storageVolumesClient, volumeName)

	createRequest := &CreateStorageAttachmentInput{
		Index:             1,
		InstanceName:      info.getInstanceName(),
		StorageVolumeName: createStorageVolumeInput.Name,
	}
	createResult, err := attachmentsClient.CreateStorageAttachment(createRequest)
	if err != nil {
		panic(err)
	}

	defer tearDownStorageAttachments(t, attachmentsClient, createResult.Name)

	getResult, err := attachmentsClient.GetStorageAttachment(createResult.Name)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(createResult.Index, getResult.Index) {
		t.Fatalf("Retrieved Storage Volume Attachment did not match Expected. \nDesired: %s \nActual: %s", createResult, getResult)
	}

	log.Printf("Attachment created: %#v\n", getResult)
}

func tearDownStorageAttachments(t *testing.T, attachmentsClient *StorageAttachmentsClient, name string) {
	log.Printf("Deleting Storage Attachment %s", name)
	if err := attachmentsClient.DeleteStorageAttachment(name); err != nil {
		t.Fatalf("Error deleting storage attachment, dangling resources may occur: %v", err)
	}
}

func buildStorageAttachmentsClients() (*InstancesClient, *StorageVolumeClient, *StorageAttachmentsClient, error) {
	instancesClient, err := getInstancesClient()
	if err != nil {
		return instancesClient, nil, nil, err
	}

	storageVolumesClient, err := getStorageVolumeClient()
	if err != nil {
		return instancesClient, nil, nil, err
	}

	storageAttachmentsClient, err := getStorageAttachmentsClient()
	if err != nil {
		return instancesClient, storageVolumesClient, nil, err
	}

	return instancesClient, storageVolumesClient, storageAttachmentsClient, nil
}

func getStorageAttachmentsClient() (*StorageAttachmentsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &StorageAttachmentsClient{}, err
	}

	return client.StorageAttachments(), nil
}
