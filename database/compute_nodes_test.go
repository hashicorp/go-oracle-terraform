// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"testing"

	"os"

	"log"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccComputeNodes(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, cClient, err := getComputeNodesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	var (
		instanceName string
		sInstance    *ServiceInstance
	)
	if v := os.Getenv("OPC_TEST_DB_INSTANCE"); v == "" {
		// First Create a Service Instance
		sInstance, err = sClient.createTestServiceInstance()
		if err != nil {
			t.Fatalf("Error creating Service Instance: %s", err)
		}
		defer destroyServiceInstance(t, sClient, sInstance.Name)
		instanceName = sInstance.Name
	} else {
		log.Print("Using already created DB Service Instance")
		instanceName = v
	}

	getInput := &GetComputeNodesInput{
		ServiceInstanceID: instanceName,
	}

	var computeNodes *ComputeNodesInfo
	computeNodes, err = cClient.GetComputeNodes(getInput)
	if err != nil {
		t.Fatalf("Error reading Compute Nodes: %s", err)
	}

	assert.Equal(t, "PDB1", computeNodes.Nodes[0].PDBName, "Expected PDB Name to match")
	assert.Equal(t, "ORCL", computeNodes.Nodes[0].SID, "Expected SID to match")
	assert.Equal(t, 1521, computeNodes.Nodes[0].ListenerPort, "Expected ListenerPort to match")
}

func getComputeNodesTestClients() (*ServiceInstanceClient, *UtilityClient, error) {
	client, err := GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}
	return client.ServiceInstanceClient(), client.ComputeNodes(), nil
}
