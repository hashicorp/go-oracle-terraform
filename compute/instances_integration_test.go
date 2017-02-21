package compute

import (
	"testing"

	"log"

	"fmt"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

var createdInstance *InstanceInfo

func TestAccInstanceLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	defer func() {
		if err := tearDownInstances(); err != nil {
			log.Printf("[ERR] Error deleting instance: %#v", createdInstance)
			log.Print("[ERR] Dangling resources may occur!")
			t.Fatalf("Error: %v", err)
		}
	}()

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
	createdInstance, err = svc.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	svc.Client.debugLogStr(fmt.Sprintf("Instance created: %#v\n", createdInstance))
}

func tearDownInstances() error {
	svc, err := getInstancesClient()
	if err != nil {
		return err
	}

	input := &DeleteInstanceInput{
		Name: createdInstance.Name,
		ID:   createdInstance.ID,
	}

	return svc.DeleteInstance(input)
}

func getInstancesClient() (*InstancesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &InstancesClient{}, err
	}

	return client.Instances(), nil
}
