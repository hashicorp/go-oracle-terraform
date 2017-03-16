package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_ImageListEntryTestName        = "test-acc-ip-network-image-list-entry"
  _ImageListEntryTestVersion      = "1"
)

func TestAccImageListEntriesLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

  name := "test-image-list"

	imageListClient, err := getImageListClient()
	if err != nil {
		t.Fatal("Error Creating Image List Client: %+v", err)
	}
	createImageListInput := CreateImageListInput{
		Name:        name,
		Description: "This is the second greatest image list in the world. Period.",
		Default:     1,
	}
	createResult, err := imageListClient.CreateImageList(&createImageListInput)
	if err != nil {
		t.Fatal("Error Creating Image List: %+v", err)
	}
	defer tearDownImageList(t, imageListClient, name)

	createClient, err := getImageListEntriesClient(createResult.Name, "")
	if err != nil {
		t.Fatal(err)
	}
	createInput := &CreateImageListEntryInput{
    MachineImages: []string{"/oracle/public/oel_6.7_apaas_16.4.5_1610211300"},
    Version: 1,
	}
	createdImageListEntry, err := createClient.CreateImageListEntry(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Image List Entry succcessfully created")

  gdc, err := getImageListEntriesClient(createResult.Name, _ImageListEntryTestVersion)
  if err != nil {
		t.Fatal(err)
	}
  defer destroyImageListEntry(t, gdc)

  receivedImageListEntry, err := gdc.GetImageListEntry()
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Image List Entry successfully fetched")
  if !reflect.DeepEqual(createdImageListEntry, receivedImageListEntry) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdImageListEntry, receivedImageListEntry)
	}
}

func destroyImageListEntry(t *testing.T, svc *ImageListEntriesClient) {
  if err := svc.DeleteImageListEntry(); err != nil {
		t.Fatal(err)
	}
}

func getImageListEntriesClient(name, version string) (*ImageListEntriesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}
	return client.ImageListEntries(name, version), nil
}
