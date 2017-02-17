package compute

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccSecurityListLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	securityListClient, err := getSecurityListsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security List Client")

	createSecurityListInput := CreateSecurityListInput{
		Name:               "test-sec-list",
		OutboundCIDRPolicy: "DENY",
		Policy:             "PERMIT",
	}
	securityList, err := securityListClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security List: %+v", securityList)

	getSecurityListInput := GetSecurityListInput{
		Name: securityList.Name,
	}
	getSecurityListOutput, err := securityListClient.GetSecurityList(&getSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	if securityList.Policy != getSecurityListOutput.Policy {
		t.Fatalf("Created and retrived policies don't match.\n Desired: %s\n Actual: %s", securityList.Policy, getSecurityListOutput.Policy)
	}
	log.Printf("Successfully retrieved Security List")

	updateSecurityListInput := UpdateSecurityListInput{
		Name:               securityList.Name,
		OutboundCIDRPolicy: "PERMIT",
		Policy:             "DENY",
	}
	updateSecurityListOutput, err := securityListClient.UpdateSecurityList(&updateSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	if updateSecurityListOutput.OutboundCIDRPolicy != "PERMIT" {
		t.Fatalf("Outbound policy not successfully updated \nDesired: %s \nActual: %s", updateSecurityListInput.OutboundCIDRPolicy, updateSecurityListOutput.OutboundCIDRPolicy)
	}
	log.Printf("Successfully updated Security List")

	deleteSecurityListInput := DeleteSecurityListInput{
		Name: securityList.Name,
	}
	err = securityListClient.DeleteSecurityList(&deleteSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted Security List")
}

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

		listInfo := &CreateSecurityListInput{}
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

	createSecurityListInput := CreateSecurityListInput{
		Name:               "test-list1",
		OutboundCIDRPolicy: "PERMIT",
		Policy:             "DENY",
	}
	info, err := client.CreateSecurityList(&createSecurityListInput)
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

func getSecurityListsClient() (*SecurityListsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecurityListsClient{}, err
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
