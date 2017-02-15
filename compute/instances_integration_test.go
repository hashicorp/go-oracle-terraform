package compute

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

var createdInstanceName *InstanceName

func TestAccInstanceLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	defer tearDownInstances()

	svc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

	//TODO: Initial Implementation
	//createdInstanceName, err = svc.LaunchInstance("test", "test", "oc3", "/oracle/public/oel_6.4_2GB_v1", nil, nil, []string{},
	createdInstanceName, err = svc.LaunchInstance("test", "test", "oc3", "/oracle/public/oel_6.7_apaas_16.4.5_1610211300", nil, nil, []string{},
		map[string]interface{}{
			"attr1": 12,
			"attr2": map[string]interface{}{
				"inner_attr1": "foo",
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Instance created: %#v\n", createdInstanceName)

	//TODO: Initial Implementation
	//instanceInfo, err := svc.WaitForInstanceRunning(createdInstanceName, 120)
	instanceInfo, err := svc.WaitForInstanceRunning(createdInstanceName, 300)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Instance retrieved: %#v\n", instanceInfo)
}

func tearDownInstances() {
	svc, err := getInstancesClient()
	if err != nil {
		panic(err)
	}

	err = svc.DeleteInstance(createdInstanceName)
	if err != nil {
		panic(err)
	}
	//TODO: Initial Implementation
	//err = svc.WaitForInstanceDeleted(createdInstanceName, 600)
	err = svc.WaitForInstanceDeleted(createdInstanceName, 900)
	if err != nil {
		panic(err)
	}
}

func getInstancesClient() (*InstancesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &InstancesClient{}, err
	}

	return client.Instances(), nil
}
