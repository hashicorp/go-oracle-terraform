package storage

import (
	"testing"

	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_TestObjectName = "testing-acc-object"
)

func TestAccObjectLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatalf("Error creating storage client: %s", err)
	}
	client := sClient.Objects()

	container, err := sClient.getTestContainer()
	if err != nil {
		t.Fatalf("Error creating test container: %s", err)
	}
	defer deleteContainer(t, sClient, container.Name)
	log.Printf("[DEBUG] Container created: %s", spew.Sdump(container))

	input := &CreateObjectInput{
		Name:      _TestObjectName,
		Container: container.Name,
		//DeleteAt:  0,
	}

	object, err := client.CreateObject(input)
	if err != nil {
		t.Fatalf("Error creating object: %s", err)
	}
	defer deleteObject(t, client, object)

	log.Printf("[DEBUG] Created Object: %s", spew.Sdump(object))

}

// Get a container for testing objects with
func (c *StorageClient) getTestContainer() (*Container, error) {
	input := &CreateContainerInput{
		Name:         _ContainerName,
		PrimaryKey:   _ContainerPrimaryKey,
		SecondaryKey: _ContainerSecondaryKey,
	}

	return c.CreateContainer(input)
}

func deleteObject(
	t *testing.T,
	client *ObjectClient,
	object *ObjectInfo) {
	input := &DeleteObjectInput{
		Name:      object.Name,
		Container: object.Container,
	}
	if err := client.DeleteObject(input); err != nil {
		t.Fatalf("Error deleting storage object: %s", err)
	}
}
