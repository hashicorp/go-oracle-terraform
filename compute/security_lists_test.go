package compute

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccSecurityListLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	securityListClient, err := getSecurityListsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Obtained Security List Client")

	name := "test-sec-list"

	createSecurityListInput := CreateSecurityListInput{
		Name:               name,
		OutboundCIDRPolicy: SecurityListPolicyDeny,
		Policy:             SecurityListPolicyPermit,
	}

	createdSecurityList, err := securityListClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security List: %+v", createdSecurityList)
	defer deleteSecurityList(t, securityListClient, name)

	getSecurityListInput := GetSecurityListInput{
		Name: name,
	}
	getSecurityListOutput, err := securityListClient.GetSecurityList(&getSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	if createdSecurityList.Policy != getSecurityListOutput.Policy {
		t.Fatalf("Created and retrived policies don't match.\n Desired: %s\n Actual: %s", createdSecurityList.Policy, getSecurityListOutput.Policy)
	}
	log.Print("Successfully retrieved Security List")

	updateSecurityListInput := UpdateSecurityListInput{
		Name:               name,
		OutboundCIDRPolicy: SecurityListPolicyPermit,
		Policy:             SecurityListPolicyDeny,
	}
	updateSecurityListOutput, err := securityListClient.UpdateSecurityList(&updateSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, updateSecurityListOutput.OutboundCIDRPolicy, SecurityListPolicy("PERMIT"), "Outbound policy not successfully updated.")
	assert.Equal(t, updateSecurityListOutput.FQDN, securityListClient.getQualifiedName(name), "Expected FDQN to be equal to qualified name")

	log.Print("Successfully updated Security List")
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

func deleteSecurityList(t *testing.T, client *SecurityListsClient, name string) {
	deleteInput := DeleteSecurityListInput{
		Name: name,
	}
	if err := client.DeleteSecurityList(&deleteInput); err != nil {
		t.Fatal(err)
	}

	log.Print("Successfully deleted Security List")
}
