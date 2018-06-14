package lbaas

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// Test the Listener lifecycle to create, get, delete a Listener
// instance and validate the fields are set as expected.
func TestAccListenerLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// CREATE Parent Load Balancer Service Instance

	lbClient, err := getLoadBalancerClient()
	if err != nil {
		t.Fatal(err)
	}

	createLoadBalancerInput := &CreateLoadBalancerInput{
		Name:        "acc-test-lb-listener1",
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

	// CREATE Listener

	listenerClient, err := getListenerClient()
	if err != nil {
		t.Fatal(err)
	}

	createListenerInput := &CreateListenerInput{
		Name:                 "acc-test-listener1",
		Port:                 8080,
		BalancerProtocol:     ProtocolHTTP,
		OriginServerProtocol: ProtocolHTTP,
	}

	_, err = listenerClient.CreateListener(createLoadBalancerInput.Region, createLoadBalancerInput.Name, createListenerInput)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyListener(t, listenerClient, createLoadBalancerInput.Region, createLoadBalancerInput.Name, createListenerInput.Name)

	// FETCH

	resp, err := listenerClient.GetListener(createLoadBalancerInput.Region, createLoadBalancerInput.Name, createListenerInput.Name)
	if err != nil {
		t.Fatal(err)
	}

	expected := &ListenerInfo{
		Name:                 createListenerInput.Name,
		Port:                 createListenerInput.Port,
		BalancerProtocol:     createListenerInput.BalancerProtocol,
		OriginServerProtocol: createListenerInput.OriginServerProtocol,
	}

	// compare resp to expected
	compare(t, "Name", resp.Name, expected.Name)
	compare(t, "Port", string(resp.Port), string(expected.Port))
	compare(t, "BalancerProtocol", string(resp.BalancerProtocol), string(expected.BalancerProtocol))
	compare(t, "OriginServerProtocol", string(resp.OriginServerProtocol), string(expected.OriginServerProtocol))

	// UPDATE

	// TODO updates throw a HTTP 405 Error "Method not allowed"

	// updateInput := &UpdateListenerInput{
	// 	Port:                 8081,
	// 	BalancerProtocol:     ProtocolHTTPS,
	// 	OriginServerProtocol: ProtocolHTTPS,
	// }
	//
	// resp, err = listenerClient.UpdateListener(createLoadBalancerInput.Region, createLoadBalancerInput.Name, createListenerInput.Name, updateInput)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// expected = &ListenerInfo{
	// 	Name:                 createListenerInput.Name,
	// 	Port:                 updateInput.Port,
	// 	BalancerProtocol:     updateInput.BalancerProtocol,
	// 	OriginServerProtocol: updateInput.OriginServerProtocol,
	// }
	//
	// compare(t, "Name", resp.Name, expected.Name)
	// compare(t, "Port", string(resp.Port), string(expected.Port))
	// compare(t, "BalancerProtocol", string(resp.BalancerProtocol), string(expected.BalancerProtocol))
	// compare(t, "OriginServerProtocol", string(resp.OriginServerProtocol), string(expected.OriginServerProtocol))
	//
}

func getListenerClient() (*ListenerClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &ListenerClient{}, err
	}
	return client.ListenerClient(), nil
}

func destroyListener(t *testing.T, client *ListenerClient, lbRegion, lbName, name string) {
	if _, err := client.DeleteListener(lbRegion, lbName, name); err != nil {
		t.Fatal(err)
	}
}
