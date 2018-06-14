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

	// CREATE

	createLoadBalancerInput := &CreateLoadBalancerInput{
		Name:        "acc-test-lb1",
		Region:      "uscom-central-1",
		Description: "Terraformed Load Balancer Test",
		Scheme:      LoadBalancerSchemeInternetFacing,
		Disabled:    LoadBalancerDisabledFalse,
		Tags:        []string{"tag3", "tag2", "tag1"},
	}

	_, err = lbClient.CreateLoadBalancer(createLoadBalancerInput)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyLoadBalancer(t, lbClient, createLoadBalancerInput.Region, createLoadBalancerInput.Name)

	// FETCH

	resp, err := lbClient.GetLoadBalancer(createLoadBalancerInput.Region, createLoadBalancerInput.Name)
	if err != nil {
		t.Fatal(err)
	}

	expected := &LoadBalancerInfo{
		Name:        createLoadBalancerInput.Name,
		Region:      createLoadBalancerInput.Region,
		Description: createLoadBalancerInput.Description,
		Scheme:      createLoadBalancerInput.Scheme,
		Disabled:    createLoadBalancerInput.Disabled,
		Tags:        createLoadBalancerInput.Tags,
	}

	// compare resp to expected
	compare(t, "Name", resp.Name, expected.Name)
	compare(t, "Region", resp.Region, expected.Region)
	compare(t, "Description", resp.Description, expected.Description)
	compare(t, "Scheme", string(resp.Scheme), string(expected.Scheme))
	compare(t, "Disabled", string(resp.Disabled), string(expected.Disabled))
	// TODO compare(t, "Tags", string(resp.Tags), string(expected.Tags))

	// UPDATE

	// TODO updates throw a HTTP 405 Error "Method not allowed"

	// updateInput := &UpdateLoadBalancerInput{
	// 	Description: "Updated Description",
	// 	Tags:        []string{"TAGA", "TAGB", "TAGC"},
	// 	// TODO add updateable attributes
	// }
	//
	// resp, err = lbClient.UpdateLoadBalancer(createLoadBalancerInput.Region, createLoadBalancerInput.Name, updateInput)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// expected = &LoadBalancerInfo{
	// 	Name:        createLoadBalancerInput.Name,
	// 	Region:      createLoadBalancerInput.Region,
	// 	Description: updateInput.Description,
	// 	Disabled:    createLoadBalancerInput.Disabled,
	// 	Scheme:      createLoadBalancerInput.Scheme,
	// 	Tags:        updateInput.Tags,
	// }
	//
	// compare(t, "Name", resp.Name, expected.Name)
	// compare(t, "Region", resp.Region, expected.Region)
	// compare(t, "Description", resp.Description, expected.Description)
	// compare(t, "Scheme", string(resp.Scheme), string(expected.Scheme))
	// compare(t, "Disabled", string(resp.Disabled), string(expected.Disabled))
	// // TODO compare(t, "Tags", string(resp.Tags), string(expected.Tags))

}

func getLoadBalancerClient() (*LoadBalancerClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &LoadBalancerClient{}, err
	}
	return client.LoadBalancerClient(), nil
}

func destroyLoadBalancer(t *testing.T, client *LoadBalancerClient, region, name string) {
	if _, err := client.DeleteLoadBalancer(region, name); err != nil {
		t.Fatal(err)
	}
}
