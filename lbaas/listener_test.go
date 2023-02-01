// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lbaas

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

// Test the Listener lifecycle to create, get, delete a Listener
// instance and validate the fields are set as expected.
func TestAccListenerLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// CREATE Parent Load Balancer Service Instance

	lb, lbClient := createParentLoadBalancer(t)
	if lbClient != nil {
		defer destroyLoadBalancer(t, lbClient, lb)
	}
	// CREATE Listener

	listenerClient, err := getListenerClient()
	assert.NoError(t, err)

	createListenerInput := &CreateListenerInput{
		Name:                 "acc-test-listener1",
		Port:                 8080,
		BalancerProtocol:     ProtocolHTTP,
		OriginServerProtocol: ProtocolHTTP,
		Disabled:             LBaaSDisabledTrue,
		Tags:                 []string{"tag3", "tag2", "tag1"},
	}

	_, err = listenerClient.CreateListener(lb, createListenerInput)
	assert.NoError(t, err)

	defer destroyListener(t, listenerClient, lb, createListenerInput.Name)

	// FETCH

	resp, err := listenerClient.GetListener(lb, createListenerInput.Name)
	assert.NoError(t, err)

	expected := &ListenerInfo{
		Name:                 createListenerInput.Name,
		Port:                 createListenerInput.Port,
		BalancerProtocol:     createListenerInput.BalancerProtocol,
		OriginServerProtocol: createListenerInput.OriginServerProtocol,
		Tags:                 createListenerInput.Tags,
	}

	// compare resp to expected
	assert.Equal(t, expected.Name, resp.Name, "Listener name should match")
	assert.Equal(t, expected.Port, resp.Port, "Listener port should match")
	assert.Equal(t, expected.BalancerProtocol, resp.BalancerProtocol, "Listener balancer protocol should match")
	assert.Equal(t, expected.OriginServerProtocol, resp.OriginServerProtocol, "Listener origin server protocol should match")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Listener tags to match ")

	// UPDATE

	updateTags := []string{"TAGA", "TAGB", "TAGC"}
	updatePathPrefixes := []string{"/path1", "/path2"}

	updateInput := &UpdateListenerInput{
		Name:                 createListenerInput.Name,
		Port:                 8081,
		BalancerProtocol:     ProtocolHTTPS,
		OriginServerProtocol: ProtocolHTTPS,
		PathPrefixes:         &updatePathPrefixes,
		Tags:                 &updateTags,
	}

	resp, err = listenerClient.UpdateListener(lb, createListenerInput.Name, updateInput)
	assert.NoError(t, err)

	expected = &ListenerInfo{
		Name:                 createListenerInput.Name,
		Port:                 updateInput.Port,
		BalancerProtocol:     updateInput.BalancerProtocol,
		OriginServerProtocol: updateInput.OriginServerProtocol,
		PathPrefixes:         updatePathPrefixes,
		Tags:                 updateTags,
	}

	assert.Equal(t, expected.Name, resp.Name, "Listener name should match")
	assert.Equal(t, expected.Port, resp.Port, "Listener port should match")
	assert.Equal(t, expected.BalancerProtocol, resp.BalancerProtocol, "Listener balancer protocol should match")
	assert.Equal(t, expected.OriginServerProtocol, resp.OriginServerProtocol, "Listener origin server protocol should match")
	assert.ElementsMatch(t, expected.PathPrefixes, resp.PathPrefixes, "Expected Listener path prefixes to match ")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Listener tags to match ")

}

func getListenerClient() (*ListenerClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &ListenerClient{}, err
	}
	return client.ListenerClient(), nil
}

func destroyListener(t *testing.T, client *ListenerClient, lb LoadBalancerContext, name string) {
	if _, err := client.DeleteListener(lb, name); err != nil {
		t.Fatal(err)
	}
}
