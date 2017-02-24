package compute

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"log"
	"reflect"

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

	createResult, err := attachmentsClient.CreateStorageAttachment(1, info, createStorageVolumeInput.Name)
	if err != nil {
		panic(err)
	}

	attachmentName = createResult.Name
	err = attachmentsClient.WaitForStorageAttachmentCreated(attachmentName, 30)
	if err != nil {
		panic(err)
	}

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

// Test that the client can create an instance.
func TestAccStorageAttachmentsClient_GetStorageAttachmentsForInstance(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Wrong HTTP method %s, expected GET", r.Method)
		}

		expectedPath := "/storage/attachment/Compute-test/test/?state=attached&instance_name=/Compute-test/test/test-instance/test-id"
		if r.URL.String() != expectedPath {
			t.Errorf("Wrong HTTP URL %v, expected %v", r.URL, expectedPath)
		}

		w.Write([]byte(exampleGetStorageAttachmentsResponse))
		w.WriteHeader(200)
	})

	defer server.Close()
	client, err := getStubStorageAttachmentsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	instance := &InstanceInfo{
		Name: "test-instance",
		ID:   "test-id",
	}

	_, err = client.GetStorageAttachmentsForInstance(instance)

	if err != nil {
		t.Fatalf("Get security attachments request failed: %s", err)
	}
}

func getStubStorageAttachmentsClient(server *httptest.Server) (*StorageAttachmentsClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}
	return client.StorageAttachments(), nil
}

var exampleGetStorageAttachmentsResponse = `
{
"result": [
  {
    "index": 5,
    "account": null,
    "storage_volume_name": "/Compute-acme/jack.jones@example.com/data",
    "hypervisor": null,
    "uri": "https://api.compute.us0.oraclecloud.com/storage/attachment/Compute-acme/jack.jones@example.com/01fa297e-e7e1-4501-84d3-402ccc33e66d/10bf639f-bb78-462a-b5ac-eeb0474771a0",
    "instance_name": "/Compute-acme/jack.jones@example.com/01fa297e-e7e1-4501-84d3-402ccc33e66d",
    "state": "attached",
    "readonly": false,
    "name": "/Compute-acme/jack.jones@example.com/01fa297e-e7e1-4501-84d3-402ccc33e66d/10bf639f-bb78-462a..."
  },
  {
    "index": 1,
    "account": null,
    "storage_volume_name": "/Compute-acme/jack.jones@example.com/boot",
    "hypervisor": null,
    "uri": "https://api.compute.us0.oraclecloud.com/storage/attachment/Compute-acme/jack.jones@example.com/01fa297e-e7e1-4501-84d3-402ccc33e66d/4aa33097-b085-4484-a909-a6a0a5955c05",
    "instance_name": "/Compute-acme/jack.jones@example.com/01fa297e-e7e1-4501-84d3-402ccc33e66d",
    "state": "attached",
    "readonly": false,
    "name": "/Compute-acme/jack.jones@example.com/01fa297e-e7e1-4501-84d3-402ccc33e66d/4aa33097-b085-4484..."
  }
 ]
}
`

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
