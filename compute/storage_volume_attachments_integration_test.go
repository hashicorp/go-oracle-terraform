package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccStorageAttachmentsLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	instanceName := "test-acc-stor-att-instance"
	volumeName := "test-acc-stor-att-volume"
	var attachmentName string
	var instanceInfo InstanceInfo

	instancesClient, storageVolumesClient, attachmentsClient, err := buildStorageAttachmentsClients()
	if err != nil {
		panic(err)
	}

	defer tearDownStorageAttachments(instancesClient, storageVolumesClient, attachmentsClient, &instanceInfo, volumeName, &attachmentName)

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
	instanceInfo = *info

	createStorageVolumeInput := storageVolumesClient.NewStorageVolumeSpec("10G", []string{"/oracle/public/storage/default"}, volumeName)
	err = storageVolumesClient.CreateStorageVolume(createStorageVolumeInput)
	if err != nil {
		panic(err)
	}

	_, err = storageVolumesClient.WaitForStorageVolumeOnline(volumeName, 30)
	if err != nil {
		panic(err)
	}

	createRequest := &CreateStorageAttachmentInput{
		Index:             1,
		InstanceName:      info.Name,
		InstanceId:        info.ID,
		StorageVolumeName: createStorageVolumeInput.Name,
	}
	createResult, err := attachmentsClient.CreateStorageAttachment(createRequest)
	if err != nil {
		panic(err)
	}

	attachmentName = createResult.Name

	getResult, err := attachmentsClient.GetStorageAttachment(attachmentName)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(createResult.Index, getResult.Index) {
		t.Fatalf("Retrieved Storage Volume Attachment did not match Expected. \nDesired: %s \nActual: %s", createResult, getResult)
	}

	log.Printf("Attachment created: %#v\n", getResult)
}

func tearDownStorageAttachments(instancesClient *InstancesClient, volumesClient *StorageVolumeClient, attachmentsClient *StorageAttachmentsClient,
	instanceInfo *InstanceInfo, volumeName string, attachmentName *string) {

	// delete the storage attachment only if it exists
	if *attachmentName != "" {
		log.Printf("Deleting Storage Attachment %s", *attachmentName)
		err := attachmentsClient.DeleteStorageAttachment(*attachmentName)
		if err != nil {
			panic(err)
		}

		err = attachmentsClient.WaitForStorageAttachmentDeleted(*attachmentName, 30)
		if err != nil {
			panic(err)
		}
	}

	qualifiedVolumeName := volumesClient.getQualifiedName(volumeName)
	volume, err := volumesClient.GetStorageVolume(qualifiedVolumeName)
	if volume != nil {
		log.Printf("Deleting Storage Volume %s", volumeName)

		_ = volumesClient.DeleteStorageVolume(qualifiedVolumeName)

		err = volumesClient.WaitForStorageVolumeDeleted(volumeName, 30)
		if err != nil {
			panic(err)
		}
	}

	if instanceInfo != nil {
		log.Printf("Deleting Instance %s", instanceInfo.Name)
		deleteInstanceInput := &DeleteInstanceInput{
			Name: instanceInfo.Name,
			ID:   instanceInfo.ID,
		}
		err = instancesClient.DeleteInstance(deleteInstanceInput)
		if err != nil {
			panic(err)
		}
	}
}

func getStorageAttachmentsClient() (*StorageAttachmentsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &StorageAttachmentsClient{}, err
	}

	return client.StorageAttachments(), nil
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
