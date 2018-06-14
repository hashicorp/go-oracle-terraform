package lbaas

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// Test the Load Balancer lifecycle the create, get, delete a Load Balancer
// instance and validate the fields are set as expected.
func TestAccLoadBalancerLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	lbClient, err := getLoadBalancerClient()
	if err != nil {
		t.Fatal(err)
	}

	createLoadBalancerInput := &CreateLoadBalancerInput{
		Name:        "acc-test-lb",
		Region:      "uscom-central-1",
		Description: "Terraformed Load Balancer Test",
		Scheme:      LoadBalancerSchemeInternetFacing,
		Disabled:    "FALSE",
	}

	resp, err := lbClient.CreateLoadBalancer(createLoadBalancerInput)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyLoadBalancer(t, lbClient, createLoadBalancerInput.Region, createLoadBalancerInput.Name)

	getInput := &GetLoadBalancerInput{
		Name:   createLoadBalancerInput.Name,
		Region: createLoadBalancerInput.Region,
	}

	resp, err = lbClient.GetLoadBalancer(getInput)
	if err != nil {
		t.Fatal(err)
	}

	expected := &LoadBalancerInfo{
		Name:        createLoadBalancerInput.Name,
		Region:      createLoadBalancerInput.Region,
		Description: createLoadBalancerInput.Description,
		Scheme:      createLoadBalancerInput.Scheme,
		Disabled:    createLoadBalancerInput.Disabled,
	}

	// compare resp to expected
	// TODO there must be a more general way to do this with reflection
	compare(t, "Name", resp.Name, expected.Name)
	compare(t, "Region", resp.Region, expected.Region)
	compare(t, "Description", resp.Description, expected.Description)
	compare(t, "Scheme", string(resp.Scheme), string(expected.Scheme))
	compare(t, "Disabled", string(resp.Disabled), string(expected.Disabled))

}

func compare(t *testing.T, attrName, respValue, expectedValue string) {
	if respValue != expectedValue {
		t.Fatalf("%s %s in response does to match expected value of %s", attrName, respValue, expectedValue)
	}
}

func getLoadBalancerClient() (*LoadBalancerClient, error) {
	client, err := GetLoadBalancerTestClient(&opc.Config{})
	if err != nil {
		return &LoadBalancerClient{}, err
	}
	return client.LoadBalancerClient(), nil
}

func destroyLoadBalancer(t *testing.T, client *LoadBalancerClient, region, name string) {
	input := &DeleteLoadBalancerInput{
		Name:   name,
		Region: region,
	}
	if _, err := client.DeleteLoadBalancer(input); err != nil {
		t.Fatal(err)
	}
}
