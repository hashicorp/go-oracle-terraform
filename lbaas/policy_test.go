package lbaas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

// Test the Policy lifecycle to create, get, update and delete a Policy
// and validate the fields are set as expected.
func TestAccPolicyLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// CREATE Parent Load Balancer Service Instance

	lb, lbClient := createParentLoadBalancer(t)
	if lbClient != nil {
		defer destroyLoadBalancer(t, lbClient, lb)
	}

	// CREATE Policy

	policyClient, err := getPolicyClient()
	assert.NoError(t, err)

	createPolicyInput := &CreatePolicyInput{
		Name: "acc-test-policy1",
		Type: "SetRequestHeaderPolicy",
		SetRequestHeaderPolicyInfo: SetRequestHeaderPolicyInfo{
			ActionWhenHeaderExists: "OVERWRITE",
			HeaderName:             "MyHeaderName",
			Value:                  "MyValue",
		},
	}

	_, err = policyClient.CreatePolicy(lb, createPolicyInput)
	assert.NoError(t, err)

	defer destroyPolicy(t, policyClient, lb, createPolicyInput.Name)

	// FETCH

	resp, err := policyClient.GetPolicy(lb, createPolicyInput.Name)
	assert.NoError(t, err)

	expectedURI := fmt.Sprintf("%s/vlbrs/%s/%s/policies/%s", policyClient.client.APIEndpoint.String(), lb.Region, lb.Name, createPolicyInput.Name)

	expected := &PolicyInfo{
		Name:                   createPolicyInput.Name,
		Type:                   createPolicyInput.Type,
		HeaderName:             createPolicyInput.SetRequestHeaderPolicyInfo.HeaderName,
		ActionWhenHeaderExists: createPolicyInput.SetRequestHeaderPolicyInfo.ActionWhenHeaderExists,
		Value:                  createPolicyInput.SetRequestHeaderPolicyInfo.Value,
		State:                  LBaaSStateHealthy,
		URI:                    expectedURI,
	}

	// compare resp to expected
	assert.Equal(t, expected, resp, "Response should match expected PolicyInfo")
	assert.Equal(t, expected.Name, resp.Name, "SetRequestHeaderPolicy Name should match")
	assert.Equal(t, expected.ActionWhenHeaderExists, resp.ActionWhenHeaderExists, "SetRequestHeaderPolicy ActionWhenHeaderExists should match")
	assert.Equal(t, expected.HeaderName, resp.HeaderName, "SetRequestHeaderPolicy HeaderName should match")
	assert.Equal(t, expected.Value, resp.Value, "SetRequestHeaderPolicy Value should match")

	// UPDATE

	updatedHeaderName := "UpdatedHeaderName"
	updatedValue := "UpdatedValue"

	updateInput := &UpdatePolicyInput{
		Name: createPolicyInput.Name,
		SetRequestHeaderPolicyInfo: SetRequestHeaderPolicyInfo{
			HeaderName: updatedHeaderName,
			Value:      updatedValue,
		},
	}

	resp, err = policyClient.UpdatePolicy(lb, createPolicyInput.Name, createPolicyInput.Type, updateInput)
	assert.NoError(t, err)

	expected = &PolicyInfo{
		Name:                   updateInput.Name,
		Type:                   createPolicyInput.Type,
		HeaderName:             updatedHeaderName,
		Value:                  updatedValue,
		ActionWhenHeaderExists: createPolicyInput.SetRequestHeaderPolicyInfo.ActionWhenHeaderExists,
		State:                  LBaaSStateHealthy,
		URI:                    expectedURI,
	}

	assert.Equal(t, expected, resp, "Response should match expected PolicyInfo")

	assert.Equal(t, expected.Name, resp.Name, "SetRequestHeaderPolicy Name should match")
	assert.Equal(t, expected.ActionWhenHeaderExists, resp.ActionWhenHeaderExists, "SetRequestHeaderPolicy ActionWhenHeaderExists should match")
	assert.Equal(t, expected.HeaderName, resp.HeaderName, "SetRequestHeaderPolicy HeaderName should match")
	assert.Equal(t, expected.Value, resp.Value, "SetRequestHeaderPolicy Value should match")

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
