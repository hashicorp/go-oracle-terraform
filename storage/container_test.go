package storage

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

const _ContainerName = "test-str-container"
const _ContainerPrimaryKey = "test-url-key"
const _ContainerSecondaryKey = "test-url-key2"
const _ContainerMaxAge = 50

func TestAccContainerLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	client, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatal(err)
	}

	readACLs := []string{"test-read-acl1", "test-read-acl2"}
	writeACLs := []string{"test-write-acl1", "test-write-acl2"}
	allowedOrigins := []string{"allowed-origin-1", "allowed-origin-2"}

	createContainerInput := CreateContainerInput{
		Name:           _ContainerName,
		ReadACLs:       readACLs,
		WriteACLs:      writeACLs,
		PrimaryKey:     _ContainerPrimaryKey,
		SecondaryKey:   _ContainerSecondaryKey,
		AllowedOrigins: allowedOrigins,
	}

	createdContainer, err := client.CreateContainer(&createContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Container: %+v", createdContainer)
	defer deleteContainer(t, client, _ContainerName)

	getContainerInput := GetContainerInput{
		Name: _ContainerName,
	}
	container, err := client.GetContainer(&getContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	if diff := pretty.Compare(container.ReadACLs, readACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("ReadACLs do not match Wanted: %+v Recieved: %+v", readACLs, container.ReadACLs))
	}
	if diff := pretty.Compare(container.WriteACLs, writeACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("WriteACLs do not match Wanted: %+v Recieved: %+v", writeACLs, container.WriteACLs))
	}
	if container.PrimaryKey != _ContainerPrimaryKey {
		t.Fatalf(fmt.Sprintf("URLKeys don't match. Wanted: %s Recieved: %s", _ContainerPrimaryKey, container.PrimaryKey))
	}
	if container.SecondaryKey != _ContainerSecondaryKey {
		t.Fatalf(fmt.Sprintf("URLKey2 do not match. Wanted: %s Recieved: %s", _ContainerSecondaryKey, container.SecondaryKey))
	}
	if diff := pretty.Compare(container.AllowedOrigins, allowedOrigins); diff != "" {
		t.Fatalf(fmt.Sprintf("AllowedOrigins do not match Wanted: %+v Recieved: %+v", allowedOrigins, container.AllowedOrigins))
	}

	log.Print("Successfully retrieved Container")

	updateReadACLs := []string{"test-read-acl3", "test-read-acl4"}
	updateWriteACLs := []string{"test-write-acl3", "test-write-acl4"}
	updatedAllowedOrigins := []string{"allowed-origin-3", "allowed-origin-4"}
	updateContainerInput := UpdateContainerInput{
		Name:           _ContainerName,
		ReadACLs:       updateReadACLs,
		WriteACLs:      updateWriteACLs,
		SecondaryKey:   _ContainerPrimaryKey,
		AllowedOrigins: updatedAllowedOrigins,
		MaxAge:         _ContainerMaxAge,
	}

	_, err = client.UpdateContainer(&updateContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Container: %+v", createdContainer)

	container, err = client.GetContainer(&getContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	if container.Name != _ContainerName {
		t.Fatalf(fmt.Sprintf("Names don't match. Wanted: %s Recieved: %s", _ContainerName, container.Name))
	}
	if diff := pretty.Compare(container.ReadACLs, updateReadACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("UpdatedReadACLs do not match Wanted: %+v Recieved: %+v", container.ReadACLs, updateReadACLs))
	}
	if diff := pretty.Compare(container.WriteACLs, updateWriteACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("UpdatedWriteACLs do not match Wanted: %+v Recieved: %+v", container.WriteACLs, updateWriteACLs))
	}
	if container.PrimaryKey != "" {
		t.Fatalf(fmt.Sprintf("Expected URL Key to be empty. Recieved: %s", container.PrimaryKey))
	}
	if container.SecondaryKey != _ContainerPrimaryKey {
		t.Fatalf(fmt.Sprintf("Updated URL Key 2 does not match. Wanted: %s Recieved: %s", _ContainerPrimaryKey, container.SecondaryKey))
	}
	if diff := pretty.Compare(container.AllowedOrigins, updatedAllowedOrigins); diff != "" {
		t.Fatalf(fmt.Sprintf("Updated AllowedOrigins do not match Wanted: %+v Recieved: %+v", updatedAllowedOrigins, container.AllowedOrigins))
	}
	if container.MaxAge != _ContainerMaxAge {
		t.Fatalf(fmt.Sprintf("Max Age do not match Wanted: %s Recieved: %s", _ContainerMaxAge, container.MaxAge))
	}

	log.Print("Successfully retrieved Container")
}

func deleteContainer(t *testing.T, client *StorageClient, name string) {
	deleteInput := DeleteContainerInput{
		Name: name,
	}
	if err := client.DeleteContainer(&deleteInput); err != nil {
		t.Fatal(err)
	}

	log.Print("Successfully deleted Container")
}
