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
	getResult, err := client.GetImageList(&getInput)
	if err != nil {
		t.Fatalf("Error Getting Image List: %+v", err)
	}

	updateInput := UpdateImageListInput{
		Name:        getResult.Name,
		Description: "Updated Description",
	}

	_, err = client.UpdateImageList(&updateInput)
	if err != nil {
		t.Fatalf("Error Updating Image List: %+v", updateInput)
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
		t.Fatal("Error Deleting Image List: %+v", err)
	}
	log.Printf("Successfully deleted Image List")
}
