package compute

import (
	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"log"
	"reflect"
	"testing"
)

func TestAccSecurityApplicationsTCPLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	securityApplicationsClient, err := getSecurityApplicationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Applications List Client")

	createInput := CreateSecurityApplicationInput{
		Name:        "test-sec-app-tcp",
		Description: "Terraform Acceptance Test TCP Lifecycle",
		Protocol:    "all",
		DPort:       "19336",
	}
	createdSecurityApplication, err := securityApplicationsClient.CreateSecurityApplication(&createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Application: %+v", createdSecurityApplication)

	retrieveInput := GetSecurityApplicationInput{
		Name: createInput.Name,
	}
	retrievedSecurityApplication, err := securityApplicationsClient.GetSecurityApplication(&retrieveInput)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(createdSecurityApplication, retrievedSecurityApplication) {
		t.Fatalf("Retrieved Security Application did not match Expected. \nDesired: %s \nActual: %s", createdSecurityApplication, retrievedSecurityApplication)
	}

	log.Printf("Successfully retrieved Security Application")

	deleteInput := DeleteSecurityApplicationInput{
		Name: createInput.Name,
	}
	err = securityApplicationsClient.DeleteSecurityApplication(&deleteInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted Security Application")
}

func TestAccSecurityApplicationsICMPLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	securityApplicationsClient, err := getSecurityApplicationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Applications List Client")

	createInput := CreateSecurityApplicationInput{
		Name:        "test-sec-app-icmp",
		Description: "Terraform Acceptance Test ICMP Lifecycle",
		Protocol:    "icmp",
		ICMPType:    "echo",
	}

	createdSecurityApplication, err := securityApplicationsClient.CreateSecurityApplication(&createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Application: %+v", createdSecurityApplication)

	retrieveInput := GetSecurityApplicationInput{
		Name: createInput.Name,
	}
	retrievedSecurityApplication, err := securityApplicationsClient.GetSecurityApplication(&retrieveInput)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(createdSecurityApplication, retrievedSecurityApplication) {
		t.Fatalf("Retrieved Security Application did not match Expected. \nDesired: %s \nActual: %s", createdSecurityApplication, retrievedSecurityApplication)
	}

	log.Printf("Successfully retrieved Security Application")

	deleteInput := DeleteSecurityApplicationInput{
		Name: createInput.Name,
	}
	err = securityApplicationsClient.DeleteSecurityApplication(&deleteInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted Security Application")
}

func getSecurityApplicationsClient() (*SecurityApplicationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecurityApplicationsClient{}, err
	}

	return client.SecurityApplications(), nil
}
