package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccInstanceLifeCycle(t *testing.T) {
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
		Networking: map[string]NetworkingInfo{
			"eth0": {},
		},
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
	defer tearDownInstances(t, svc, createdInstance.Name, createdInstance.ID)

	log.Printf("Instance created: %#v\n", createdInstance)
}

func tearDownInstances(t *testing.T, svc *InstancesClient, name, id string) {
	input := &DeleteInstanceInput{
		Name: name,
		ID:   id,
	}
	if err := svc.DeleteInstance(input); err != nil {
		t.Fatalf("Error deleting instance, dangling resources may occur: %v", err)
	}
}

func getInstancesClient() (*InstancesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &InstancesClient{}, err
	}

	return client.Instances(), nil
}
