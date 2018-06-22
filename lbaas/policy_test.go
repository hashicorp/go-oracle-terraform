package lbaas

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// Test the Policy lifecycle to create, get, update and delete a Policy
// and validate the fields are set as expected.
func TestAccPolicyLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// CREATE Parent Load Balancer Service Instance

	lbClient, err := getLoadBalancerClient()
	if err != nil {
		t.Fatal(err)
	}

	createLoadBalancerInput := &CreateLoadBalancerInput{
		Name:        "acc-test-lb-policy1",
		Region:      "uscom-central-1",
		Description: "Terraformed Load Balancer Test",
		Scheme:      LoadBalancerSchemeInternetFacing,
		Disabled:    LoadBalancerDisabledFalse,
	}

	_, err = lbClient.CreateLoadBalancer(createLoadBalancerInput)
	if err != nil {
		t.Fatal(err)
	}

	lb := LoadBalancerContext{
		Region: createLoadBalancerInput.Region,
		Name:   createLoadBalancerInput.Name,
	}

	defer destroyLoadBalancer(t, lbClient, lb)

	// CREATE Policy

	policyClient, err := getPolicyClient()
	if err != nil {
		t.Fatal(err)
	}

	createPolicyInput := &CreatePolicyInput{
		Name:       "acc-test-policy1",
		Action:     "OVERWRITE",
		HeaderName: "internal",
		Type:       "SetRequestHeaderPolicy",
		Value:      "http://myurl.example.com",
	}

	_, err = policyClient.CreatePolicy(lb, createPolicyInput)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyPolicy(t, policyClient, lb, createPolicyInput.Name)

	// FETCH

	resp, err := policyClient.GetPolicy(lb, createPolicyInput.Name)
	if err != nil {
		t.Fatal(err)
	}

	expected := &PolicyInfo{
		Name:       createPolicyInput.Name,
		Action:     createPolicyInput.Action,
		HeaderName: createPolicyInput.HeaderName,
		Type:       createPolicyInput.Type,
		Value:      createPolicyInput.Value,
	}

	// compare resp to expected
	compare(t, "Name", resp.Name, expected.Name)
	compare(t, "Action", resp.Action, expected.Action)
	compare(t, "HeaderName", resp.HeaderName, expected.HeaderName)
	compare(t, "Type", resp.Type, expected.Type)
	compare(t, "Value", resp.Value, expected.Value)

	// UPDATE

	// TODO updates throw a HTTP 405 Error "Method not allowed"

	// updateInput := &UpdatePolicyInput{
	// 	Value: "http://myurl.example.com",
	// }
	//
	// resp, err = policyClient.UpdatePolicy(createLoadBalancerInput.Region, createLoadBalancerInput.Name, createPolicyInput.Name, createPolicyInput.Type, updateInput)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// expected = &PolicyInfo{
	// 	Value: updateInput.Value,
	// }
	//
	// compare(t, "Value", resp.Value, expected.Value)

}

func getPolicyClient() (*PolicyClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &PolicyClient{}, err
	}
	return client.PolicyClient(), nil
}

func destroyPolicy(t *testing.T, client *PolicyClient, lb LoadBalancerContext, name string) {
	if _, err := client.DeletePolicy(lb, name); err != nil {
		t.Fatal(err)
	}
}
