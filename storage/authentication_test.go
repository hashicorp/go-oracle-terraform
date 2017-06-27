package storage

import (
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// Test that the client can obtain an authentication cookie from the authentication endpoint.
func TestAccObtainAuthenticationToken(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	client, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatalf("Authentication failed: %s", err)
	}

	if client.authToken == nil {
		t.Fatal("Authentication token not set")
	}
}
