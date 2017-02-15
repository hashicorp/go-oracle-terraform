package compute

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

// Test that the client can create an instance.
func TestAccSecurityListsClient_CreateKey(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Wrong HTTP method %s, expected POST", r.Method)
		}

		expectedPath := "/seclist/"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong HTTP URL %v, expected %v", r.URL, expectedPath)
		}

		listInfo := &SecurityListSpec{}
		unmarshalRequestBody(t, r, listInfo)

		if listInfo.Name != "/Compute-test/test/test-list1" {
			t.Errorf("Expected name 'test-list1', was %s", listInfo.Name)
		}

		if listInfo.Policy != "DENY" {
			t.Errorf("Expected policy 'DENY', was %s", listInfo.Policy)
		}

		if listInfo.OutboundCIDRPolicy != "PERMIT" {
			t.Errorf("Expected outbound CIDR policy 'PERMIT', was %s", listInfo.OutboundCIDRPolicy)
		}

		w.Write([]byte(exampleCreateSecurityListResponse))
		w.WriteHeader(201)
	})

	defer server.Close()
	client, err := getStubSecurityListsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	info, err := client.CreateSecurityList("test-list1", "DENY", "PERMIT")
	if err != nil {
		t.Fatalf("Create security list request failed: %s", err)
	}

	if info.Name != "allowed_video_servers" {
		t.Errorf("Expected name 'allowed_video_servers', was %s", info.Name)
	}
}

func getStubSecurityListsClient(server *httptest.Server) (*SecurityListsClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}

	return client.SecurityLists(), nil
}

var exampleCreateSecurityListResponse = `
{
  "account": "/Compute-acme/default",
  "name": "/Compute-acme/jack.jones@example.com/allowed_video_servers",
  "uri": "https://api.compute.us0.oraclecloud.com/seclist/Compute-acme/jack.jones@example.com/es_list",
  "outbound_cidr_policy": "DENY",
  "policy": "PERMIT"
}
`
