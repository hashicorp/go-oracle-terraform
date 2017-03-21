package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccImageListLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-image-list"

	client, err := getImageListClient()
	if err != nil {
		t.Fatal("Error Creating Image List Client: %+v", err)
	}

	createInput := CreateImageListInput{
		Name:        name,
		Description: "This is the greatest image list in the world. Period.",
		Default:     1,
	}
	createResult, err := client.CreateImageList(&createInput)
	if err != nil {
		t.Fatal("Error Creating Image List: %+v", err)
	}

	defer tearDownImageList(t, client, name)

	getInput := GetImageListInput{
		Name: createResult.Name,
	}
	createGetResult, err := client.GetImageList(&getInput)
	if err != nil {
		t.Fatalf("Error Getting Image List: %+v", err)
	}

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

	// we can't compare the entire object because of the additional fields :(
	if createInput.Description != createGetResult.Description {
		t.Fatalf("Created and retrieved Image List don't match.\n Desired: %s\n Actual: %s", createInput.Description, createGetResult.Description)
	}

	// we can't compare the entire object because of the additional fields :(
	if updateInput.Description != updatedGetResult.Description {
		t.Fatalf("Updated and retrieved Image List don't match.\n Desired: %s\n Actual: %s", updateInput.Description, updatedGetResult.Description)
	}

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
