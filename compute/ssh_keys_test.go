package compute

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

// Test that the client can create an instance.
func TestAccSSHClient_CreateKey(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "PUT" {
			t.Errorf("Wrong HTTP method %s, expected POST or PUT", r.Method)
		}

		expectedCreatePath := "/sshkey/"
		expectedUpdatePath := "/sshkey/Compute-test/test/test-key1"
		if r.URL.Path != expectedCreatePath && r.URL.Path != expectedUpdatePath {
			t.Errorf("Wrong HTTP URL %v, expected %v or %v",
				r.URL, expectedCreatePath, expectedUpdatePath)
		}

		keyInfo := &SSHKey{}
		unmarshalRequestBody(t, r, keyInfo)

		if keyInfo.Name != "/Compute-test/test/test-key1" {
			t.Errorf("Expected name '/Compute-test/test/test-key1', was %s", keyInfo.Name)
		}

		if !keyInfo.Enabled {
			t.Errorf("Key %s was not enabled", keyInfo.Name)
		}

		if keyInfo.Key != "key" {
			t.Errorf("Expected key 'key', was %s", keyInfo.Key)
		}

		w.Write([]byte(exampleCreateKeyResponse))
		w.WriteHeader(201)
	})

	defer server.Close()
	client, err := getStubSSHKeysClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	createKeyInput1 := CreateSSHKeyInput{
		Name:    "test-key1",
		Key:     "key",
		Enabled: true,
	}
	info, err := client.CreateSSHKey(&createKeyInput1)
	if err != nil {
		t.Fatalf("Create ssh key request failed: %s", err)
	}

	if info.Name != "test-key1" {
		t.Errorf("Expected key 'test-key1, was %s", info.Name)
	}

	expected := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDzU21CEj6JsqIMQAYwNbmZ5P2BVxA..."
	if info.Key != expected {
		t.Errorf("Expected key %s, was %s", expected, info.Key)
	}
}

func getStubSSHKeysClient(server *httptest.Server) (*SSHKeysClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}
	return client.SSHKeys(), nil
}

var exampleCreateKeyResponse = `
{
 "enabled": false,
 "uri": "https://api.compute.us0.oraclecloud.com/sshkey/Compute-test/test/test-key1",
 "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDzU21CEj6JsqIMQAYwNbmZ5P2BVxA...",
 "name": "/Compute-test/test/test-key1"
}
`
