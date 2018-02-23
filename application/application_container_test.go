package application

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_ApplicationContainerTestName       = "testaccapplicationcontainer4"
	_ApplicationContainerRuntimeJava    = "java"
	_ApplicationContainerDeploymentFile = "./test_files/deployment.json"
	_ApplicationContainerManifestFile   = "./test_files/manifest.json"
	_ApplicationContainerNote           = "This is a note"
	_ApplicationContainerUpdatedNote    = "This is an updated note"
)

func TestAccApplicationContainerLifeCycle_Basic(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	aClient, err := getApplicationContainerTestClients()
	if err != nil {
		t.Fatal(err)
	}

	createApplicationContainerAdditionalFields := CreateApplicationContainerAdditionalFields{
		Name:    _ApplicationContainerTestName,
		Runtime: _ApplicationContainerRuntimeJava,
	}

	createApplicationContainerInput := &CreateApplicationContainerInput{
		AdditionalFields: createApplicationContainerAdditionalFields,
	}

	createdApplicationContainer, err := aClient.CreateApplicationContainer(createApplicationContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Application Container: %+v", createdApplicationContainer)
	defer deleteTestApplicationContainer(t, aClient, _ApplicationContainerTestName)

	getInput := &GetApplicationContainerInput{
		Name: _ApplicationContainerTestName,
	}

	applicationContainer, err := aClient.GetApplicationContainer(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Application Container Retrieved: %+v", applicationContainer)
	assert.NotEmpty(t, applicationContainer.Name, "Expected Application Container name not to be empty")
	assert.Equal(t, _ApplicationContainerTestName, applicationContainer.Name, "Expected Application Container and name to match.")
	assert.Equal(t, "RUNNING", applicationContainer.Status)
}

func TestAccApplicationContainerLifeCycle_Manifest(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	aClient, err := getApplicationContainerTestClients()
	if err != nil {
		t.Fatal(err)
	}

	createApplicationContainerAdditionalFields := CreateApplicationContainerAdditionalFields{
		Name:    _ApplicationContainerTestName,
		Runtime: _ApplicationContainerRuntimeJava,
	}

	createApplicationContainerInput := &CreateApplicationContainerInput{
		AdditionalFields: createApplicationContainerAdditionalFields,
		Manifest:         _ApplicationContainerManifestFile,
	}

	createdApplicationContainer, err := aClient.CreateApplicationContainer(createApplicationContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Application Container: %+v", createdApplicationContainer)
	defer deleteTestApplicationContainer(t, aClient, _ApplicationContainerTestName)

	getInput := &GetApplicationContainerInput{
		Name: _ApplicationContainerTestName,
	}

	applicationContainer, err := aClient.GetApplicationContainer(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Application Container Retrieved: %+v", applicationContainer)
	assert.NotEmpty(t, applicationContainer.Name, "Expected Application Container name not to be empty")
	assert.Equal(t, _ApplicationContainerTestName, applicationContainer.Name, "Expected Application Container and name to match.")
	assert.Equal(t, "RUNNING", applicationContainer.Status)
}

func TestAccApplicationContainerLifeCycle_Deployment(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	aClient, err := getApplicationContainerTestClients()
	if err != nil {
		t.Fatal(err)
	}

	createApplicationContainerAdditionalFields := CreateApplicationContainerAdditionalFields{
		Name:    _ApplicationContainerTestName,
		Runtime: _ApplicationContainerRuntimeJava,
	}

	createApplicationContainerInput := &CreateApplicationContainerInput{
		AdditionalFields: createApplicationContainerAdditionalFields,
		Deployment:       _ApplicationContainerDeploymentFile,
		Manifest:         _ApplicationContainerManifestFile,
	}

	createdApplicationContainer, err := aClient.CreateApplicationContainer(createApplicationContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Application Container: %+v", createdApplicationContainer)
	defer deleteTestApplicationContainer(t, aClient, _ApplicationContainerTestName)

	getInput := &GetApplicationContainerInput{
		Name: _ApplicationContainerTestName,
	}

	applicationContainer, err := aClient.GetApplicationContainer(getInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Application Container Retrieved: %+v", applicationContainer)
	assert.NotEmpty(t, applicationContainer.Name, "Expected Application Container name not to be empty")
	assert.Equal(t, _ApplicationContainerTestName, applicationContainer.Name, "Expected Application Container and name to match.")
	assert.Equal(t, "RUNNING", applicationContainer.Status)
}

func deleteTestApplicationContainer(t *testing.T, client *ApplicationContainerClient, name string) {
	deleteInput := DeleteApplicationContainerInput{
		Name: name,
	}
	if err := client.DeleteApplicationContainer(&deleteInput); err != nil {
		t.Fatal(err)
	}

	log.Print("Successfully deleted Application Container")
}

func getApplicationContainerTestClients() (*ApplicationContainerClient, error) {
	client, err := getApplicationTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}
	return client.ApplicationContainerClient(), nil
}
