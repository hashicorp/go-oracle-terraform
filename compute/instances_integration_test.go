package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccInstanceLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	svc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

	input := &CreateInstanceInput{
		Name:      "test-acc",
		Label:     "test",
		Shape:     "oc3",
		ImageList: "/oracle/public/oel_6.7_apaas_16.4.5_1610211300",
		Storage:   nil,
		BootOrder: nil,
		SSHKeys:   []string{},
		Attributes: map[string]interface{}{
			"attr1": 12,
			"attr2": map[string]interface{}{
				"inner_attr1": "foo",
			},
		},
	}
	createdInstance, err := svc.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Instance created: %#v\n", createdInstance)
	defer deleteInstance(t, svc, createdInstance.Name, createdInstance.ID)
}

func deleteInstance(t *testing.T, client *InstancesClient, name string, id string) {
	deleteInput := DeleteInstanceInput{
		Name: name,
		ID:   id,
	}
	err := client.DeleteInstance(&deleteInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted Instance")
}

func getInstancesClient() (*InstancesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &InstancesClient{}, err
	}

	return client.Instances(), nil
}
