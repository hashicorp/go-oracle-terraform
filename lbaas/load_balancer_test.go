// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lbaas

import (
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/stretchr/testify/assert"
)

// Test the Load Balancer lifecycle the create, get, delete a Load Balancer
// instance and validate the fields are set as expected.
func TestAccLoadBalancerLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	lbClient, err := getLoadBalancerClient()
	assert.NoError(t, err)

	var region string
	if region = os.Getenv("OPC_TEST_LBAAS_REGION"); region == "" {
		region = "uscom-central-1"
	}

	// CREATE

	createLoadBalancerInput := &CreateLoadBalancerInput{
		Name:        "acc-test-lb1",
		Region:      region,
		Description: "Terraformed Load Balancer Test",
		Scheme:      LoadBalancerSchemeInternetFacing,
		Disabled:    LBaaSDisabledTrue,
		Tags:        []string{"tag3", "tag2", "tag1"},
	}

	_, err = lbClient.CreateLoadBalancer(createLoadBalancerInput)
	assert.NoError(t, err)

	lb := LoadBalancerContext{
		Region: createLoadBalancerInput.Region,
		Name:   createLoadBalancerInput.Name,
	}

	defer destroyLoadBalancer(t, lbClient, lb)

	// FETCH

	resp, err := lbClient.GetLoadBalancer(lb)
	assert.NoError(t, err)

	expected := &LoadBalancerInfo{
		Name:        createLoadBalancerInput.Name,
		Region:      createLoadBalancerInput.Region,
		Description: createLoadBalancerInput.Description,
		Scheme:      createLoadBalancerInput.Scheme,
		Disabled:    createLoadBalancerInput.Disabled,
		Tags:        createLoadBalancerInput.Tags,
	}

	// compare resp to expected
	assert.Equal(t, expected.Name, resp.Name, "Expected Load Balancer name to match")
	assert.Equal(t, expected.Region, resp.Region, "Expected Load Balancer region to match")
	assert.Equal(t, expected.Description, resp.Description, "Expected Load Balancer description to match")
	assert.Equal(t, expected.Scheme, resp.Scheme, "Expected Load Balancer scheme to match")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Load Balancer tags to match ")

	// UPDATE

	updatedDescription := "Updated Description"
	updatedTags := []string{"TAGA", "TAGB", "TAGC"}

	updateInput := &UpdateLoadBalancerInput{
		Name:        createLoadBalancerInput.Name,
		Description: &updatedDescription,
		Tags:        &updatedTags,
	}

	resp, err = lbClient.UpdateLoadBalancer(lb, updateInput)
	assert.NoError(t, err)

	expected = &LoadBalancerInfo{
		Name:        createLoadBalancerInput.Name,
		Description: updatedDescription,
		Tags:        updatedTags,
	}

	// compare resp to expected
	assert.Equal(t, expected.Name, resp.Name, "Expected Load Balancer name to match")
	assert.Equal(t, expected.Description, resp.Description, "Expected Load Balancer description to match")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Load Balancer tags to match ")

}
