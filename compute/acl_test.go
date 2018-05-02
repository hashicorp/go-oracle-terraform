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

const (
	_ACLTestName        = "test-acc-acl"
	_ACLTestDescription = "testing acl"
)

func TestAccACLLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-acc-acl"

	ACLClient, err := getACLsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained acl Client")

	createACLInput := CreateACLInput{
		Name:        _ACLTestName,
		Enabled:     false,
		Description: _ACLTestDescription,
		Tags:        []string{"tag1"},
	}

	createdACL, err := ACLClient.CreateACL(&createACLInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created acl: %+v", createdACL)
	defer deleteACL(t, ACLClient, name)

	getACLInput := GetACLInput{
		Name: _ACLTestName,
	}
	getACLOutput, err := ACLClient.GetACL(&getACLInput)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, createdACL, getACLOutput, "Created and retrieved acls don't match.")
	log.Print("Successfully retrieved acl")

	updateACLInput := UpdateACLInput{
		Name:        _ACLTestName,
		Enabled:     true,
		Description: _ACLTestDescription,
		Tags:        []string{"tag1"},
	}
	updateACLOutput, err := ACLClient.UpdateACL(&updateACLInput)
	if err != nil {
		t.Fatal(err)
	}
	if !updateACLOutput.Enabled {
		t.Fatal("acl was not updated to enabled")
	}
	log.Print("Successfully updated acl")
}

// Test that the client can create an instance.
func TestAccACLsClient_CreateRule(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Wrong HTTP method %s, expected POST", r.Method)
		}

		expectedPath := "/network/v1/acl/"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong HTTP URL %v, expected %v", r.URL, expectedPath)
		}

		ruleSpec := &CreateACLInput{}
		unmarshalRequestBody(t, r, ruleSpec)

		expectedName := "/Compute-test/test/test-acc-acl"
		assert.Equal(t, expectedName, ruleSpec.Name, "ruleSpec name not expected.")
		assert.False(t, ruleSpec.Enabled, "Expected ruleSpec to not be enabled.")

		w.WriteHeader(201)
		w.Write([]byte(exampleCreateACLResponse))
	})

	defer server.Close()
	client, err := getStubACLsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	createInput := CreateACLInput{
		Name:        _ACLTestName,
		Enabled:     false,
		Description: _ACLTestDescription,
		Tags:        []string{"tag1"},
	}
	info, err := client.CreateACL(&createInput)
	if err != nil {
		t.Fatalf("Create security rule request failed: %s", err)
	}
	assert.False(t, info.Enabled, "Expected `info` to not be enabled.")
}

var exampleCreateACLResponse = `
{
  "name": "/Compute-acme/jack.jones@example.com/es_to_videoservers_stream",
  "enabled": false
}
`

func getStubACLsClient(server *httptest.Server) (*ACLsClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}

	return client.ACLs(), nil
}

func getACLsClient() (*ACLsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &ACLsClient{}, err
	}

	return client.ACLs(), nil
}

func deleteACL(t *testing.T, client *ACLsClient, name string) {
	deleteInput := DeleteACLInput{
		Name: name,
	}
	err := client.DeleteACL(&deleteInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted acl")
}
