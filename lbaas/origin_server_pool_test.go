package lbaas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

// Test the Origin Server Pool lifecycle to create, get, delete a Origin Server
// Pool and validate the fields are set as expected.
func TestAccOriginServerPoolLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	// CREATE Parent Load Balancer Service Instance

	var region string
	if region = os.Getenv("OPC_TEST_LBAAS_REGION"); region == "" {
		region = "uscom-central-1"
	}
	lb := createParentLoadBalancer(t, region, "acc-test-lb-server-pool1")

	// CREATE Origin Server Pool

	serverPoolClient, err := getOriginServerPoolClient()
	assert.NoError(t, err)

	healthCheckInfo := HealthCheckInfo{
		Type:                "HTTP",
		Path:                "/health",
		AcceptedReturnCodes: []string{"2xx", "3xx", "4xx", "5xx"},
		Enabled:             "TRUE",
		Interval:            30,
		Timeout:             30,
		HealthyThreshold:    6,
		UnhealthyThreshold:  3,
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
		HealthCheck: &healthCheckInfo,
		Tags:        []string{"tag3", "tag2", "tag1"},
		Status:      "ENABLED",
	}

	_, err = serverPoolClient.CreateOriginServerPool(lb, createOriginServerPoolInput)
	assert.NoError(t, err)

	defer destroyOriginServerPool(t, serverPoolClient, lb, createOriginServerPoolInput.Name)

	// FETCH

	resp, err := serverPoolClient.GetOriginServerPool(lb, createOriginServerPoolInput.Name)
	assert.NoError(t, err)

	expectedURI := fmt.Sprintf("%svlbrs/%s/%s/originserverpools/%s", serverPoolClient.client.APIEndpoint.String(), lb.Region, lb.Name, createOriginServerPoolInput.Name)

	expected := &OriginServerPoolInfo{
		Name: createOriginServerPoolInput.Name,
		OriginServers: []OriginServerInfo{
			OriginServerInfo{
				Status:   "ENABLED",
				Hostname: "example.com",
				Port:     3691,
			},
		},
		HealthCheck: healthCheckInfo,
		Status:      createOriginServerPoolInput.Status,
		Tags:        createOriginServerPoolInput.Tags,
		State:       LBaaSStateHealthy,
		URI:         expectedURI,
	}

	// compare resp to expected
	assert.Equal(t, expected.Name, resp.Name, "Origin Server Pool name should match")
	assert.Equal(t, expected.URI, resp.URI, "Origin Server Pool URI should match")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Origin Server Pool tags to match ")
	assert.ElementsMatch(t, expected.OriginServers, resp.OriginServers, "Origin Servers should match expected OriginServerInfo ")
	assert.Equal(t, expected.HealthCheck, resp.HealthCheck, "Health Check should match expected HealthCheckInfo")

	// UPDATES

	updateOriginServers := []CreateOriginServerInput{
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
	}

	updatedHeathCheck := HealthCheckInfo{
		Type:                "HTTP",
		Path:                "",
		Enabled:             "TRUE",
		AcceptedReturnCodes: []string{"4xx"},
		Interval:            10,
		Timeout:             5,
		HealthyThreshold:    6,
		UnhealthyThreshold:  2,
	}

	updateTags := []string{"TAGA", "TAGB", "TAGC"}

	updateInput := &UpdateOriginServerPoolInput{
		Name:          createOriginServerPoolInput.Name,
		OriginServers: &updateOriginServers,
		HealthCheck:   &updatedHeathCheck,
		Status:        LBaaSStatusDisabled,
		Tags:          &updateTags,
	}

	resp, err = serverPoolClient.UpdateOriginServerPool(lb, createOriginServerPoolInput.Name, updateInput)
	assert.NoError(t, err)

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
		},
		Status:      updateInput.Status,
		HealthCheck: updatedHeathCheck,
		Tags:        updateTags,
		State:       LBaaSStateHealthy,
		URI:         expectedURI,
	}

	assert.Equal(t, expected.Name, resp.Name, "Origin Server Pool name should match")
	assert.Equal(t, expected.URI, resp.URI, "Origin Server Pool URI should match")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Origin Server Pool tags to match ")
	assert.Len(t, resp.OriginServers, 2, "Expected two Origin Servers to be defined")
	assert.ElementsMatch(t, expected.OriginServers, resp.OriginServers, "Origin Servers should match expected OriginServerInfo ")
	assert.Equal(t, expected.HealthCheck, resp.HealthCheck, "Health Check should match expected HealthCheckInfo")

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
