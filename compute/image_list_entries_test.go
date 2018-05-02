package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_ImageListEntryTestName    = "test-acc-ip-network-image-list-entry"
	_ImageListEntryTestVersion = 1
)

func TestAccImageListEntriesLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	imageListClient, err := getImageListClient()
	if err != nil {
		t.Fatalf("Error Creating Image List Client: %+v", err)
	}
	createImageListInput := CreateImageListInput{
		Name:        _ImageListEntryTestName,
		Description: "This is the second greatest image list in the world. Period.",
		Default:     1,
	}
	_, err = imageListClient.CreateImageList(&createImageListInput)
	if err != nil {
		t.Fatalf("Error Creating Image List: %+v", err)
	}
	defer tearDownImageList(t, imageListClient, _ImageListEntryTestName)

	entryClient, err := getImageListEntriesClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateImageListEntryInput{
		Name:          _ImageListEntryTestName,
		MachineImages: []string{"/oracle/public/OL_7.2_UEKR4_x86_64-18.1.4-20180209-231028"},
		Version:       1,
	}

	createdImageListEntry, err := entryClient.CreateImageListEntry(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Image List Entry succcessfully created")
	defer destroyImageListEntry(t, entryClient, createdImageListEntry)

	getInput := &GetImageListEntryInput{
		Name:    _ImageListEntryTestName,
		Version: _ImageListEntryTestVersion,
	}
	receivedImageListEntry, err := entryClient.GetImageListEntry(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Image List Entry successfully fetched")
	if !reflect.DeepEqual(createdImageListEntry, receivedImageListEntry) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdImageListEntry, receivedImageListEntry)
	}
}

func destroyImageListEntry(t *testing.T, svc *ImageListEntriesClient, imageListEntry *ImageListEntryInfo) {
	deleteInput := &DeleteImageListEntryInput{
		Name:    imageListEntry.URI,
		Version: imageListEntry.Version,
	}
	if err := svc.DeleteImageListEntry(deleteInput); err != nil {
		t.Fatal(err)
	}
}

func getImageListEntriesClient() (*ImageListEntriesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}
	return client.ImageListEntries(), nil
}
