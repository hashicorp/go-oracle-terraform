// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

// WARNING: This _will_ leak a DB instance, that needs to be removed either manually,
// or via a manual API call. Useful for local test iterations, as it takes > 1 hour to
// spin up and shut down a DB Test Instance. The resulting test instance
// can be used in utility tests for access rules and ssh keys by populating the
// OPC_TEST_DB_INSTANCE environment variable with the full name of the created service instance.
func TestAccDataBaseClient_CreateStandalone(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	if v := os.Getenv("OPC_TEST_DB"); v == "" {
		t.Skip("Skipping Database Standalone Instance Create")
	}
	client, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}
	sInstance, err := client.createTestServiceInstance()
	if err != nil {
		t.Fatalf("Error creating service instance: %s", err)
	}
	t.Logf("Created Service Instance: %s", sInstance.Name)
}

// This is another helper test function similar to the above TestAccDataBaseClient_CreateStandalone,
// however, this will attempt to destroy the service instance found in the env var
// OPC_TEST_DB_INSTANCE. If no instance is specified at env var, exits cleanly.
func TestAccDataBaseClient_DestroyStandalone(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	instance := os.Getenv("OPC_TEST_DB_INSTANCE")
	if instance == "" {
		t.Skip("No DB Instance to destroy")
	}
	client, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}
	// Attempt to destroy the service instance
	destroyServiceInstance(t, client, instance)
}
