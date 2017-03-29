package compute

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
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
	if !reflect.DeepEqual(createdACL, getACLOutput) {
		t.Fatalf("Created and retrived acls don't match.\n Desired: %s\n Actual: %s", createdACL, getACLOutput)
	}
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
	if updateACLOutput.Enabled != true {
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
		if ruleSpec.Name != expectedName {
			t.Errorf("Expected name '%s', was %s", expectedName, ruleSpec.Name)
		}
		if ruleSpec.Enabled != false {
			t.Errorf("Expected enabled to be 'false', was %s", ruleSpec.Enabled)
		}

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
	if info.Enabled != false {
		t.Errorf("Expected enabled 'false', was %s", info.Enabled)
	}
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
