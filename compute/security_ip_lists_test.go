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

func TestAccSecurityIPListLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	securityIPListClient, err := getSecurityIPListsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security IP List Client")

	createSecurityIPListInput := CreateSecurityIPListInput{
		Name:         "test-sec-ip-list",
		SecIPEntries: []string{"127.0.0.1", "127.0.0.2"},
	}
	securityIPList, err := securityIPListClient.CreateSecurityIPList(&createSecurityIPListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security IP List: %+v", securityIPList)

	getSecurityIPListInput := GetSecurityIPListInput{
		Name: securityIPList.Name,
	}
	getSecurityIPListOutput, err := securityIPListClient.GetSecurityIPList(&getSecurityIPListInput)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(securityIPList.SecIPEntries, getSecurityIPListOutput.SecIPEntries) {
		t.Fatalf("Created and retrieved security IP entries do not match \nDesired: %s \nActual: %s", securityIPList.SecIPEntries, getSecurityIPListOutput.SecIPEntries)
	}
	log.Printf("Successfully retrieved Security IP List")

	updateSecurityIPListInput := UpdateSecurityIPListInput{
		Name:         securityIPList.Name,
		SecIPEntries: []string{"127.0.0.3", "127.0.0.4"},
	}
	updateSecurityIPListOutput, err := securityIPListClient.UpdateSecurityIPList(&updateSecurityIPListInput)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updateSecurityIPListInput.SecIPEntries, updateSecurityIPListOutput.SecIPEntries) {
		t.Fatalf("Security IP Entry not successfully updated \nDesired: %s \nActual: %s", updateSecurityIPListInput.SecIPEntries[0], updateSecurityIPListOutput.SecIPEntries[0])
	}
	log.Printf("Successfully updated Security IP List")

	deleteSecurityIPListInput := DeleteSecurityIPListInput{
		Name: securityIPList.Name,
	}
	err = securityIPListClient.DeleteSecurityIPList(&deleteSecurityIPListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted Security IP List")
}

// Test that the client can create an instance.
func TestAccSecurityIPListsClient_CreateKey(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Wrong HTTP method %s, expected POST", r.Method)
		}

		expectedPath := "/seciplist/"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong HTTP URL %v, expected %v", r.URL, expectedPath)
		}

		listInfo := &CreateSecurityIPListInput{}
		unmarshalRequestBody(t, r, listInfo)

		if listInfo.Name != "/Compute-test/test/test-list1" {
			t.Errorf("Expected name 'Compute-test/test/test-list1', was %s", listInfo.Name)
		}

		if !reflect.DeepEqual(listInfo.SecIPEntries, []string{"127.0.0.1", "168.10.0.0"}) {
			t.Errorf("Expected entries [127.0.0.1,168.10.0.0], was %s", listInfo.SecIPEntries)
		}

		w.Write([]byte(exampleCreateSecurityIPListResponse))
		w.WriteHeader(201)
	})

	defer server.Close()
	client, err := getStubSecurityIPListsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	createSecurityIPListInput := CreateSecurityIPListInput{
		Name:         "test-list1",
		SecIPEntries: []string{"127.0.0.1", "168.10.0.0"},
	}
	info, err := client.CreateSecurityIPList(&createSecurityIPListInput)
	if err != nil {
		t.Fatalf("Create security ip list request failed: %s", err)
	}

	if info.Name != "es_iplist" {
		t.Errorf("Expected name 'es_iplist', was %s", info.Name)
	}
}

func getStubSecurityIPListsClient(server *httptest.Server) (*SecurityIPListsClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}

	return client.SecurityIPLists(), nil
}

func getSecurityIPListsClient() (*SecurityIPListsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecurityIPListsClient{}, err
	}

	return client.SecurityIPLists(), nil
}

var exampleCreateSecurityIPListResponse = `
{
  "secipentries": [
    "46.16.56.0/21",
    "46.6.0.0/16"
  ],
  "name": "/Compute-acme/jack.jones@example.com/es_iplist",
  "uri": "https://api.compute.us0.oraclecloud.com/seciplist/Compute-acme/jack.jones@example.com/es_iplist"
}
`
