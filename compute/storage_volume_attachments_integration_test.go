package compute

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccStorageAttachmentsLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rInt := rand.Int()
	instanceName := fmt.Sprintf("test-acc-stor-att-instance-%d", rInt)
	volumeName := fmt.Sprintf("test-acc-stor-att-volume-%d", rInt)

	instancesClient, storageVolumesClient, attachmentsClient, err := buildStorageAttachmentsClients()
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	defer tearDownInstances(t, instancesClient, info.Name, info.ID)

	createStorageVolumeInput := &CreateStorageVolumeInput{
		Name:       volumeName,
		Size:       "10",
		Properties: []string{"/oracle/public/storage/default"},
	}
	_, err = storageVolumesClient.CreateStorageVolume(createStorageVolumeInput)
	if err != nil {
		t.Fatal(err)
	}

	defer tearDownStorageVolumes(t, storageVolumesClient, volumeName)

	createRequest := &CreateStorageAttachmentInput{
		Index:             1,
		InstanceName:      info.getInstanceName(),
		StorageVolumeName: createStorageVolumeInput.Name,
	}
	createResult, err := attachmentsClient.CreateStorageAttachment(createRequest)
	if err != nil {
		t.Fatal(err)
	}

	defer tearDownStorageAttachments(t, attachmentsClient, createResult.Name)

	getRequest := &GetStorageAttachmentInput{
		Name: createResult.Name,
	}
	getResult, err := attachmentsClient.GetStorageAttachment(getRequest)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, createResult.Index, getResult.Index,
		"Retrieved Storage Volume Attachment did not match Expected.")

	log.Printf("Attachment created: %#v\n", getResult)
}

func tearDownStorageAttachments(t *testing.T, attachmentsClient *StorageAttachmentsClient, name string) {
	log.Printf("Deleting Storage Attachment %s", name)

	deleteRequest := &DeleteStorageAttachmentInput{
		Name: name,
	}
	if err := attachmentsClient.DeleteStorageAttachment(deleteRequest); err != nil {
		t.Fatalf("Error deleting storage attachment, dangling resources may occur: %v", err)
	}
}

func buildStorageAttachmentsClients() (*InstancesClient, *StorageVolumeClient, *StorageAttachmentsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}
	return client.Instances(), client.StorageVolumes(), client.StorageAttachments(), nil
}
