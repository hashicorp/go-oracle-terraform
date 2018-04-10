package database

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

// Test Constants
const (
	_TestSSHServiceType      = "DBaaS"
	_TestSSHCredType         = "SSH"
	_TestSSHComponentType    = "DB"
	_TestSSHParentType       = "SERVICE"
	_TestSSHLastUpdateStatus = "SUCCESS"
	_TestSSHOsUserName       = "opc"
)

func TestAccSSHKeysLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	serviceClient, sshClient, err := getSSHKeyTestClients()
	if err != nil {
		t.Fatal(err)
	}

	var (
		instanceName string
		sInstance    *ServiceInstance
	)
	if v := os.Getenv("OPC_TEST_DB_INSTANCE"); v == "" {
		// First Create a Service Instance
		sInstance, err = serviceClient.createTestServiceInstance()
		if err != nil {
			t.Fatalf("Error creating Service Instance: %v", err)
		}
		defer destroyServiceInstance(t, serviceClient, sInstance.Name)
		instanceName = sInstance.Name
	} else {
		log.Print("Using already created DB Service Instance")
		instanceName = v
	}

	// Create an SSH Key
	input := &CreateSSHKeyInput{
		ServiceInstanceID: instanceName,
		PublicKey:         _TestSSHKeyPublicKey,
	}

	// Create SSH Key
	if _, err = sshClient.CreateSSHKey(input); err != nil {
		t.Fatalf("Error creating SSH Key: %v", err)
	}

	// Get Input
	getInput := &GetSSHKeyInput{
		ServiceInstanceID: instanceName,
	}

	res, err := sshClient.GetSSHKey(getInput)
	if err != nil {
		t.Fatalf("Error retrieving SSH Key: %v", err)
	}

	// Verify un-testable fields are not nil
	if res.LastUpdateTime == "" {
		t.Fatalf("Expected value for LastUpdateTime, got nil")
	}
	res.LastUpdateTime = ""

	if res.LastUpdateMessage == "" {
		t.Fatalf("Expected value for LastUpdateMessage, got nil")
	}
	res.LastUpdateMessage = ""

	if res.ComputeKeyName == "" {
		t.Fatalf("Expected value for ComputeKeyName, got nil")
	}
	res.ComputeKeyName = ""

	expected := &SSHKeyInfo{
		ServiceType:       _TestSSHServiceType,
		ServiceName:       instanceName,
		CredName:          DBSSHKeyName,
		CredType:          _TestSSHCredType,
		ComponentType:     _TestSSHComponentType,
		ParentType:        _TestSSHParentType,
		OsUserName:        _TestSSHOsUserName,
		PublicKey:         _TestSSHKeyPublicKey,
		Description:       _TestSSHKeyDescription,
		LastUpdateStatus:  _TestSSHLastUpdateStatus,
		LastUpdateTime:    "",
		LastUpdateMessage: "",
		ComputeKeyName:    "",
	}

	if diff := pretty.Compare(res, expected); diff != "" {
		t.Fatalf("Diff on SSH Key: (-got +want)\n%s", diff)
	}

}

func getSSHKeyTestClients() (*ServiceInstanceClient, *UtilityClient, error) {
	client, err := GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}
	return client.ServiceInstanceClient(), client.SSHKeys(), nil
}

// Testing SSH Key. Breaking constant convention by placing at bottom of file
// because a super long string is ugly.
const (
	_TestSSHKeyPublicKey   = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC5VoI2s4GYbVP/Gr3OBxZi5twtupnXc1gknedifREIBb5m7MxRmtL9ffQS7ZY1emYKey7oCs/7R/1Ta7UdV2BBFk++np3t1HUbcaFeBADZVdqavAv28waDBV78dUarDw4aTGzDtqnATz1AWYV4mtsrUeBhtdHBMgTnnm7V1U2NzIYArBs43tvLfulcWfwU/2goK6kfmi21fECJYF5uim7Nlqtgf7ynwCoLlJGNJxn6u8xQzTqmOrScbzno8oUbk3+Rnds+El0/P3mavETr2gVQVvGjmrTWoN8j6g+QjhKudU5C3PjI2MaFjdRnLcwabjJYaF1p69o1LXY+DtD42xtQ99wtfhLTAwDELNYHKcV/xJbgRPfBLXPltZys+LUT/RrhvHZ8d7Urp2FuOgl8KGreE/XH6oDVU0MT01SeLiJeAA951lpMTXmWErm5DRBHHQZ9Bq8ZHzsdqxSUUQK2SHhXnQolqlJRTr9m1o6XDA8Xcq+buMI7LT8A2wTYHbr6pmrDy1CgWVOreHuo+gZegKkFuZKsubXzeCubgHdproOlph84OjLOpdMC11XAUdvO47APNkjdEA1wfRQAxpREI9yymItmHdc81fsBBkTyH5teUVQjvPk/VoVganKqNBUyAwbVZZLdkS3lQCxAtgBg4i1XFwNP+AaA18EtZcMca/1Ykw==`
	_TestSSHKeyDescription = `Service user ssh public key which can be used to access the service VM instances`
)
