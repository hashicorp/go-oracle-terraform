// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccImageListLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-image-list"

	client, err := getImageListClient()
	if err != nil {
		t.Fatalf("Error Creating Image List Client: %+v", err)
	}

	createInput := CreateImageListInput{
		Name:        name,
		Description: "This is the greatest image list in the world. Period.",
		Default:     1,
	}
	createResult, err := client.CreateImageList(&createInput)
	if err != nil {
		t.Fatalf("Error Creating Image List: %+v", err)
	}

	defer tearDownImageList(t, client, name)

	getInput := GetImageListInput{
		Name: createResult.Name,
	}
	createGetResult, err := client.GetImageList(&getInput)
	if err != nil {
		t.Fatalf("Error Getting Image List: %+v", err)
	}

	assert.Equal(t, createInput.Description, createGetResult.Description, "Created and retrieved Image List don't match.")

	updateInput := UpdateImageListInput{
		Name:        name,
		Description: "Updated Description",
	}

	_, err = client.UpdateImageList(&updateInput)
	if err != nil {
		t.Fatalf("Error Updating Image List: %+v", updateInput)
	}

	updatedGetResult, err := client.GetImageList(&getInput)
	if err != nil {
		t.Fatalf("Error Getting Image List: %+v", err)
	}

	assert.Equal(t, updateInput.Description, updatedGetResult.Description, "Updated and retrieved Image List don't match.")
	assert.Equal(t, updatedGetResult.FQDN, client.getQualifiedName(updatedGetResult.Name), "Expected FDQN to be equal to qualified name")

	log.Print("Successfully Updated Image List")
}

func getImageListClient() (*ImageListClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.ImageList(), nil
}

func tearDownImageList(t *testing.T, client *ImageListClient, name string) {
	deleteInput := DeleteImageListInput{
		Name: name,
	}
	err := client.DeleteImageList(&deleteInput)
	if err != nil {
		t.Fatalf("Error Deleting Image List: %+v", err)
	}
	log.Printf("Successfully deleted Image List")
}
