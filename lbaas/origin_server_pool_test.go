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
		Disabled:    LBaaSDisabledTrue,
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

	// CREATE Origin Server Pool

	serverPoolClient, err := getOriginServerPoolClient()
	if err != nil {
		t.Fatal(err)
	}

	createOriginServerPoolInput := &CreateOriginServerPoolInput{
		Name: "acc-test-server-pool1",
		OriginServers: []CreateOriginServerInput{
			CreateOriginServerInput{
				Status:   "ENABLED",
				Hostname: "example.com",
				Port:     3691,
			},
		},
		Tags:   []string{"tag3", "tag2", "tag1"},
		Status: "ENABLED",
	}

	_, err = serverPoolClient.CreateOriginServerPool(lb, createOriginServerPoolInput)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyOriginServerPool(t, serverPoolClient, lb, createOriginServerPoolInput.Name)

	// FETCH

	resp, err := serverPoolClient.GetOriginServerPool(lb, createOriginServerPoolInput.Name)
	if err != nil {
		t.Fatal(err)
	}

	expected := &OriginServerPoolInfo{
		Name: createOriginServerPoolInput.Name,
		OriginServers: []OriginServerInfo{
			OriginServerInfo{
				Status:   "ENABLED",
				Hostname: "example.com",
				Port:     3691,
			},
		},
		Status: createOriginServerPoolInput.Status,
		Tags:   createOriginServerPoolInput.Tags,
	}

	// compare resp to expected
	compare(t, "Name", resp.Name, expected.Name)
	// TODO compare OriginServers
	// TODO compare Status
	// TODO compare Tags

	// UPDATE

	updateInput := &UpdateOriginServerPoolInput{
		Name: createOriginServerPoolInput.Name,
		OriginServers: []CreateOriginServerInput{
			CreateOriginServerInput{
				Status:   "ENABLED",
				Hostname: "example.com",
				Port:     3691,
			},
			CreateOriginServerInput{
				Status:   "ENABLED",
				Hostname: "example.com",
				Port:     8080,
			},
		},
		Status: LBaaSStatusDisabled,
		Tags:   []string{"TAGA", "TAGB", "TAGC"},
	}

	resp, err = serverPoolClient.UpdateOriginServerPool(lb, createOriginServerPoolInput.Name, updateInput)
	if err != nil {
		t.Fatal(err)
	}

	expected = &OriginServerPoolInfo{
		Name: updateInput.Name,
		OriginServers: []OriginServerInfo{
			OriginServerInfo{
				Status:   "ENABLED",
				Hostname: "example.com",
				Port:     3691,
			},
			OriginServerInfo{
				Status:   "ENABLED",
				Hostname: "example.com",
				Port:     8080,
			},
		}, Status: updateInput.Status,
		Tags: updateInput.Tags,
	}

	compare(t, "Name", resp.Name, expected.Name)
	// TODO compare OriginServers
	// TODO compare Status
	// TODO compare Tags

}

func getOriginServerPoolClient() (*OriginServerPoolClient, error) {
	client, err := GetTestClient(&opc.Config{})
	if err != nil {
		return &OriginServerPoolClient{}, err
	}
	return client.OriginServerPoolClient(), nil
}

func destroyOriginServerPool(t *testing.T, client *OriginServerPoolClient, lb LoadBalancerContext, name string) {
	if _, err := client.DeleteOriginServerPool(lb, name); err != nil {
		t.Fatal(err)
	}
}
