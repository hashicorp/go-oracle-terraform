package lbaas

import (
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
	if err != nil {
		t.Fatal(err)
	}

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
		HealthCheck: healthCheckInfo,
		Status:      createOriginServerPoolInput.Status,
		Tags:        createOriginServerPoolInput.Tags,
	}

	// compare resp to expected
	assert.Equal(t, expected.Name, resp.Name, "Origin Server Pool name should match")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Origin Server Pool tags to match ")
	assert.Len(t, resp.OriginServers, 1, "Expected one Origin Server to be defined")
	assert.Equal(t, expected.OriginServers[0].Hostname, resp.OriginServers[0].Hostname, "Origin Server host name should match")
	assert.Equal(t, expected.OriginServers[0].Port, resp.OriginServers[0].Port, "Origin Server port should match")

	assert.Equal(t, expected.HealthCheck.Type, resp.HealthCheck.Type, "Health Check Type should match")
	assert.Equal(t, expected.HealthCheck.Path, resp.HealthCheck.Path, "Health Check Path should match")
	assert.ElementsMatch(t, expected.HealthCheck.AcceptedReturnCodes, resp.HealthCheck.AcceptedReturnCodes, "Health Check AcceptedReturnCodes should match")
	assert.Equal(t, expected.HealthCheck.Enabled, resp.HealthCheck.Enabled, "Health Check Enabled should match")
	assert.Equal(t, expected.HealthCheck.Interval, resp.HealthCheck.Interval, "Health Check Interval should match")
	assert.Equal(t, expected.HealthCheck.Timeout, resp.HealthCheck.Timeout, "Health Check Timeout should match")
	assert.Equal(t, expected.HealthCheck.HealthyThreshold, resp.HealthCheck.HealthyThreshold, "Health Check HealthyThreshold should match")
	assert.Equal(t, expected.HealthCheck.UnhealthyThreshold, resp.HealthCheck.UnhealthyThreshold, "Health Check UnhealthyThreshold should match")

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
		HealthCheck: updatedHeathCheck,
		Tags:        updateTags,
	}

	assert.Equal(t, expected.Name, resp.Name, "Origin Server Pool name should match")
	assert.ElementsMatch(t, expected.Tags, resp.Tags, "Expected Origin Server Pool tags to match ")
	assert.Len(t, resp.OriginServers, 2, "Expected two Origin Servers to be defined")
	assert.ElementsMatch(t, expected.OriginServers, resp.OriginServers, "Expected Origin Servers to match ")

	assert.Equal(t, expected.HealthCheck.Type, resp.HealthCheck.Type, "Health Check Type should match")
	assert.Equal(t, expected.HealthCheck.Path, resp.HealthCheck.Path, "Health Check Path should match")
	assert.ElementsMatch(t, expected.HealthCheck.AcceptedReturnCodes, resp.HealthCheck.AcceptedReturnCodes, "Health Check AcceptedReturnCodes should match")
	assert.Equal(t, expected.HealthCheck.Enabled, resp.HealthCheck.Enabled, "Health Check Enabled should match")
	assert.Equal(t, expected.HealthCheck.Interval, resp.HealthCheck.Interval, "Health Check Interval should match")
	assert.Equal(t, expected.HealthCheck.Timeout, resp.HealthCheck.Timeout, "Health Check Timeout should match")
	assert.Equal(t, expected.HealthCheck.HealthyThreshold, resp.HealthCheck.HealthyThreshold, "Health Check HealthyThreshold should match")
	assert.Equal(t, expected.HealthCheck.UnhealthyThreshold, resp.HealthCheck.UnhealthyThreshold, "Health Check UnhealthyThreshold should match")

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
