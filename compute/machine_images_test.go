package compute

import (
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"fmt"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/hashicorp/go-oracle-terraform/storage"
	"github.com/kylelemons/godebug/pretty"
)

const (
	_MachineImageName     = "testing-machine-image"
	_TestFileFixturesPath = "test-fixtures"
)

func TestAccMachineImageLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	client, err := getMachineImagesClient()
	if err != nil {
		t.Fatal(err)
	}

	rInt := rand.Int()
	machineImageName := fmt.Sprintf("%s-%d", _MachineImageName, rInt)
	machineImageFile := fmt.Sprintf("%s.tar.gz", machineImageName)

	// Create dummy image file for the machine image test
	sClient := getStorageClient(t)
	createDummyMachineImageFile(t, sClient, machineImageFile)
	defer deleteDummyMachineImageFile(t, sClient, machineImageFile)

	account := fmt.Sprintf("/Compute-%s/cloud_storage", *client.Client.client.IdentityDomain)

	createMachineImage := &CreateMachineImageInput{
		Account: account,
		Name:    machineImageName,
		File:    machineImageFile,
	}

	machineImage, err := client.CreateMachineImage(createMachineImage)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyMachineImage(t, client, machineImage.Name)

	getInput := &GetMachineImageInput{
		Account: account,
		Name:    machineImageName,
	}

	receivedMachineImage, err := client.GetMachineImage(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if diff := pretty.Compare(machineImage, receivedMachineImage); diff != "" {
		t.Errorf("Created Machine Image Diff: (-got +want)\n%s", diff)
	}
}

func getMachineImagesClient() (*MachineImagesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &MachineImagesClient{}, err
	}

	return client.MachineImages(), nil
}

func destroyMachineImage(t *testing.T, client *MachineImagesClient, name string) {
	input := &DeleteMachineImageInput{
		Name: name,
	}

	if err := client.DeleteMachineImage(input); err != nil {
		t.Fatal(err)
	}
}

func getStorageClient(t *testing.T) *storage.Client {
	config := &opc.Config{}
	tr := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 120 * time.Second,
	}
	config.HTTPClient = &http.Client{Transport: tr}

	apiEndpoint, _ := url.Parse(os.Getenv("OPC_STORAGE_ENDPOINT"))
	domain := os.Getenv("OPC_IDENTITY_DOMAIN")
	username := os.Getenv("OPC_USERNAME")
	password := os.Getenv("OPC_PASSWORD")

	config.APIEndpoint = apiEndpoint
	config.IdentityDomain = &domain
	config.Username = &username
	config.Password = &password

	sClient, _ := storage.NewStorageClient(config)
	return sClient
}

// create a dummy image file the the storage /compute_images container
func createDummyMachineImageFile(t *testing.T, sClient *storage.Client, name string) {
	oClient := sClient.Objects()
	body, err := os.Open(_TestFileFixturesPath + "/dummy.tar.gz")
	if err != nil {
		t.Fatal(err)
	}

	input := &storage.CreateObjectInput{
		Name:        name,
		Container:   "compute_images",
		ContentType: "application/tar+gzip",
		Body:        body,
	}

	if _, err := oClient.CreateObject(input); err != nil {
		t.Fatal(err)
	}
}

func deleteDummyMachineImageFile(t *testing.T, sClient *storage.Client, name string) {
	oClient := sClient.Objects()
	input := &storage.DeleteObjectInput{
		Name:      name,
		Container: "compute_images",
	}

	if err := oClient.DeleteObject(input); err != nil {
		t.Fatal(err)
	}
}
