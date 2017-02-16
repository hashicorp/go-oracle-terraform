package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)


func TestAccSSHKeyLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sshKeyClient, err := getSSHKeysClient()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Obtained SSH Key Client\n")

	createSSHKeyInput := CreateSSHKeyInput {
    //Name:     "test-key",
		//Key:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC7BzZyp8CWN7tfIZiZwWx8H9RO2ClKu0ru/+bGEsUmHmSS7aA+iPBVqeK1Kr2nVkoG/32GaMLfVNRlRgZZGTBTFM5nnytNoo3DC9dnIPuIu95TbF1afGkVFNNyeJkC5bQDDaRDMaYBbPVJTa6bA8v7nmzvElQHPRtdRqZnFx80QHdrgTluqhtrxWDBCYMSm2meL/NU11kijoKfYSReT4lroglSxnkvP0vjUqUSvZ6tI231Ggvxg4TU1TL4OgtNyfQgXK585V05n7IT9iiJHThah2/ZGsb0DZimj/D5LxngciXVOkOR1sDt8pQb7QCxgoxOO3sa1K3pFi5UAJQ10tSyhu0yn0AnRG13NWK6DlLKhLzZM5jhGJeeYYuwCL5fzJojflouHgebOO62gqNANkUcf7cWUBJRWjSAYuXe/C6rJOriZuUkC87QpffpYd2WaJmqnjAaj7NaqOTzk5ltpS39EjMenyXWWw1MPs7eEB/A/Rfol0cHzGqoXaIZAJVaEpW7ePWEj193CqSc6uh1nwAT15rvh63z2l1iPL0CbuF4GwZWsIZ6roirmwPpKY79kAls69EKsa7bydSQuYpbU5otkT20FIbtHmyFMYpJzYM6sQHoljO2AHWmWChkYtqglbFPrQgwIrsAHbJtmzNcmbXLUm1AY+SjZd1UYqPBjFDb7w== your_email@example.com",
		Key:      "Bladh",
		Enabled:  true,
	}
	sshKey, err := sshKeyClient.CreateSSHKey(&createSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully created SSH Key: %+v\n", sshKey)

	getSSHKeyInput := GetSSHKeyInput{
		Name: sshKey.Name,
	}
	getSSHKeyOutput, err := sshKeyClient.GetSSHKey(&getSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	if sshKey.Key != getSSHKeyOutput.Key {
		t.Fatal("Created and retrived keys don't match %s %s\n", sshKey.Key, getSSHKeyOutput.Key)
	}
	fmt.Printf("Successfully retrieved ssh key\n")

	updateSSHKeyInput := UpdateSSHKeyInput{
		Name:     sshKey.Name,
		Key:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC7BzZyp8CWN7tfIZiZwWx8H9RO2ClKu0ru/+bGEsUmHmSS7aA+iPBVqeK1Kr2nVkoG/32GaMLfVNRlRgZZGTBTFM5nnytNoo3DC9dnIPuIu95TbF1afGkVFNNyeJkC5bQDDaRDMaYBbPVJTa6bA8v7nmzvElQHPRtdRqZnFx80QHdrgTluqhtrxWDBCYMSm2meL/NU11kijoKfYSReT4lroglSxnkvP0vjUqUSvZ6tI231Ggvxg4TU1TL4OgtNyfQgXK585V05n7IT9iiJHThah2/ZGsb0DZimj/D5LxngciXVOkOR1sDt8pQb7QCxgoxOO3sa1K3pFi5UAJQ10tSyhu0yn0AnRG13NWK6DlLKhLzZM5jhGJeeYYuwCL5fzJojflouHgebOO62gqNANkUcf7cWUBJRWjSAYuXe/C6rJOriZuUkC87QpffpYd2WaJmqnjAaj7NaqOTzk5ltpS39EjMenyXWWw1MPs7eEB/A/Rfol0cHzGqoXaIZAJVaEpW7ePWEj193CqSc6uh1nwAT15rvh63z2l1iPL0CbuF4GwZWsIZ6roirmwPpKY79kAls69EKsa7bydSQuYpbU5otkT20FIbtHmyFMYpJzYM6sQHoljO2AHWmWChkYtqglbFPrQgwIrsAHbJtmzNcmbXLUm1AY+SjZd1UYqPBjFDb7w== your_email@example.com",
		Enabled:  false,
	}
	updateSSHKeyOutput, err := sshKeyClient.UpdateSSHKey(&updateSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	if updateSSHKeyOutput.Enabled != updateSSHKeyInput.Enabled {
		t.Fatal("Key not successfully updated \nDesired: %s \nActual: %s", updateSSHKeyInput.Key, updateSSHKeyOutput.Key )
	}
	fmt.Printf("Successfully updated ssh key\n")

	deleteSSHKeyInput := DeleteSSHKeyInput{
		Name: sshKey.Name,
	}
	err = sshKeyClient.DeleteSSHKey(&deleteSSHKeyInput)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully deleted SSH Key\n")
}

// Test that the client can create an instance.
func TestAccSSHClient_CreateKey(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Wrong HTTP method %s, expected POST", r.Method)
		}

		expectedPath := "/sshkey/"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong HTTP URL %v, expected %v", r.URL, expectedPath)
		}

		keyInfo := &SSHKey{}
		unmarshalRequestBody(t, r, keyInfo)

		if keyInfo.Name != "/Compute-test/test/test-key1" {
			t.Errorf("Expected name '/Compute-test/test/test-key1', was %s", keyInfo.Name)
		}

		if keyInfo.Enabled != true {
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
		Name: "test-key1",
		Key:    "key",
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

func getSSHKeysClient() (*SSHKeysClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SSHKeysClient{}, err
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
