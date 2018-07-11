package application

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_ApplicationContainerTestName             = "testaccapplicationcontainer4"
	_ApplicationContainerRuntimeJava          = "java"
	_ApplicationContainerDeploymentFile       = "./test_files/deployment.json"
	_ApplicationContainerManifestFile         = "./test_files/manifest.json"
	_ApplicationContainerManifestCommand      = "sh target/bin/start"
	_ApplicationContainerManifestNotes        = "notes related to release"
	_ApplicationContainerManifestMode         = "rolling"
	_ApplicationContainerManifestMajorVersion = "7"
	_ApplicationContainerManifestBuild        = "150520.1154"
	_ApplicationContainerManifestCommit       = "d8c2596364d9584050461"
	_ApplicationContainerManifestVersion      = "15.1.0"
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

func TestAccApplicationContainerLifeCycle_ManifestFile(t *testing.T) {
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

func TestAccApplicationContainerLifeCycle_BadInput(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	aClient, err := getApplicationContainerTestClients()
	if err != nil {
		t.Fatal(err)
	}

	manifest := &ManifestAttributes{
		Command: _ApplicationContainerManifestCommand,
		Notes:   _ApplicationContainerManifestNotes,
		Mode:    _ApplicationContainerManifestMode,
		Runtime: Runtime{MajorVersion: _ApplicationContainerManifestMajorVersion},
		Release: Release{
			Build:   _ApplicationContainerManifestBuild,
			Commit:  _ApplicationContainerManifestCommit,
			Version: _ApplicationContainerManifestVersion},
	}

	createApplicationContainerAdditionalFields := CreateApplicationContainerAdditionalFields{
		Name:    _ApplicationContainerTestName,
		Runtime: _ApplicationContainerRuntimeJava,
	}

	createApplicationContainerInput := &CreateApplicationContainerInput{
		AdditionalFields:   createApplicationContainerAdditionalFields,
		ManifestAttributes: manifest,
		Manifest:           _ApplicationContainerManifestFile,
	}

	_, err = aClient.CreateApplicationContainer(createApplicationContainerInput)
	if err == nil {
		t.Fatal("expected error when specifiying both a manifest file and manifest attributes")
	}

	environment := map[string]string{
		"NO_OF_CONNECTIONS": "25",
		"TWITTER_ID":        "JAVA",
	}

	deployment := &DeploymentAttributes{
		Memory:      "2G",
		Instances:   1,
		Envrionment: environment,
	}

	createApplicationContainerInput = &CreateApplicationContainerInput{
		AdditionalFields:     createApplicationContainerAdditionalFields,
		DeploymentAttributes: deployment,
		Manifest:             _ApplicationContainerManifestFile,
		Deployment:           _ApplicationContainerDeploymentFile,
	}
	_, err = aClient.CreateApplicationContainer(createApplicationContainerInput)
	if err == nil {
		t.Fatal("expected error when specifiying both a manifest file and manifest attributes")
	}
}

func TestAccApplicationContainerLifeCycle_ManifestAttr(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	aClient, err := getApplicationContainerTestClients()
	if err != nil {
		t.Fatal(err)
	}

	manifest := &ManifestAttributes{
		Command: _ApplicationContainerManifestCommand,
		Notes:   _ApplicationContainerManifestNotes,
		Mode:    _ApplicationContainerManifestMode,
		Runtime: Runtime{MajorVersion: _ApplicationContainerManifestMajorVersion},
		Release: Release{
			Build:   _ApplicationContainerManifestBuild,
			Commit:  _ApplicationContainerManifestCommit,
			Version: _ApplicationContainerManifestVersion},
	}

	createApplicationContainerAdditionalFields := CreateApplicationContainerAdditionalFields{
		Name:    _ApplicationContainerTestName,
		Runtime: _ApplicationContainerRuntimeJava,
	}

	createApplicationContainerInput := &CreateApplicationContainerInput{
		AdditionalFields:   createApplicationContainerAdditionalFields,
		ManifestAttributes: manifest,
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

func TestAccApplicationContainerLifeCycle_DeploymentAttr(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	aClient, err := getApplicationContainerTestClients()
	if err != nil {
		t.Fatal(err)
	}

	environment := map[string]string{
		"NO_OF_CONNECTIONS": "25",
		"TWITTER_ID":        "JAVA",
	}

	deployment := &DeploymentAttributes{
		Memory:      "2G",
		Instances:   1,
		Envrionment: environment,
	}

	createApplicationContainerAdditionalFields := CreateApplicationContainerAdditionalFields{
		Name:    _ApplicationContainerTestName,
		Runtime: _ApplicationContainerRuntimeJava,
	}

	createApplicationContainerInput := &CreateApplicationContainerInput{
		AdditionalFields:     createApplicationContainerAdditionalFields,
		DeploymentAttributes: deployment,
		Manifest:             _ApplicationContainerManifestFile,
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

	updateApplicationContainerInput := &UpdateApplicationContainerInput{
		Name:       _ApplicationContainerTestName,
		Deployment: _ApplicationContainerDeploymentFile,
	}

	_, err = aClient.UpdateApplicationContainer(updateApplicationContainerInput)
	if err != nil {
		t.Fatal(err)
	}

	applicationContainer, err = aClient.GetApplicationContainer(getInput)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "RUNNING", applicationContainer.Status)
}

func deleteTestApplicationContainer(t *testing.T, client *ContainerClient, name string) {
	deleteInput := DeleteApplicationContainerInput{
		Name: name,
	}
	if err := client.DeleteApplicationContainer(&deleteInput); err != nil {
		t.Fatal(err)
	}

	log.Print("Successfully deleted Application Container")
}

func getApplicationContainerTestClients() (*ContainerClient, error) {
	client, err := getApplicationTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}
	return client.ContainerClient(), nil
}
