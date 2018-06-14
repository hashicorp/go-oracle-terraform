package lbaas

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// Test the Origin Server Pool lifecycle to create, get, delete a Origin Server
// Pool and validate the fields are set as expected.
func TestAccOriginServerPoolLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// CREATE Parent Load Balancer Service Instance

	lbClient, err := getLoadBalancerClient()
	if err != nil {
		t.Fatal(err)
	}

	createLoadBalancerInput := &CreateLoadBalancerInput{
		Name:        "acc-test-lb-server-pool1",
		Region:      "uscom-central-1",
		Description: "Terraformed Load Balancer Test",
		Scheme:      LoadBalancerSchemeInternetFacing,
		Disabled:    LoadBalancerDisabledFalse,
	}

	_, err = lbClient.CreateLoadBalancer(createLoadBalancerInput)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyLoadBalancer(t, lbClient, createLoadBalancerInput.Region, createLoadBalancerInput.Name)

	// CREATE Origin Server Pool

	serverPoolClient, err := getOriginServerPoolClient()
	if err != nil {
		t.Fatal(err)
	}

	createOriginServerPoolInput := &CreateOriginServerPoolInput{
		Name: "acc-test-server-pool1",
		OriginServers: []OriginServerInfo{
			OriginServerInfo{
				Hostname: "example.com",
				Port:     3691,
				Status:   "ENABLED",
			},
		},
		Status: "ENABLED",
		Tags:   []string{"tag3", "tag2", "tag1"},
	}

	// _, err = serverPoolClient.CreateOriginServerPool("uscom-central-1", "my-test-lb", createOriginServerPoolInput)
	_, err = serverPoolClient.CreateOriginServerPool(createLoadBalancerInput.Region, createLoadBalancerInput.Name, createOriginServerPoolInput)
	if err != nil {
		t.Fatal(err)
	}

	// defer destroyOriginServerPool(t, serverPoolClient, "uscom-central-1", "my-test-lb", createOriginServerPoolInput.Name)
	defer destroyOriginServerPool(t, serverPoolClient, createLoadBalancerInput.Region, createLoadBalancerInput.Name, createOriginServerPoolInput.Name)

	// FETCH

	// resp, err := serverPoolClient.GetOriginServerPool("uscom-central-1", "my-test-lb", createOriginServerPoolInput.Name)
	resp, err := serverPoolClient.GetOriginServerPool(createLoadBalancerInput.Region, createLoadBalancerInput.Name, createOriginServerPoolInput.Name)
	if err != nil {
		t.Fatal(err)
	}

	expected := &OriginServerPoolInfo{
		Name: createOriginServerPoolInput.Name,
		// Status: createOriginServerPoolInput.Status,
		// Tags:   createOriginServerPoolInput.Tags,
	}

	// compare resp to expected
	compare(t, "Name", resp.Name, expected.Name)
	// compare(t, "Status", string(resp.Status), string(expected.Status))

	// UPDATE

	// TODO

}

func getOriginServerPoolClient() (*OriginServerPoolClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &OriginServerPoolClient{}, err
	}
	return client.OriginServerPoolClient(), nil
}

func destroyOriginServerPool(t *testing.T, client *OriginServerPoolClient, lbRegion, lbName, name string) {
	if _, err := client.DeleteOriginServerPool(lbRegion, lbName, name); err != nil {
		t.Fatal(err)
	}
}
