package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccSecurityApplicationsTCPLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-sec-app-tcp"

	securityApplicationsClient, err := getSecurityApplicationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Applications Client")

	createInput := CreateSecurityApplicationInput{
		Name:        name,
		Description: "Terraform Acceptance Test TCP Lifecycle",
		Protocol:    SecurityApplicationProtocol(TCP),
		DPort:       "19336",
	}
	createdSecurityApplication, err := securityApplicationsClient.CreateSecurityApplication(&createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Application: %+v", createdSecurityApplication)
	defer deleteSecurityApplication(t, securityApplicationsClient, name)

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
}

func TestAccSecurityApplicationsICMPLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-sec-app-icmp"
	securityApplicationsClient, err := getSecurityApplicationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Applications Client")

	createInput := CreateSecurityApplicationInput{
		Name:        name,
		Description: "Terraform Acceptance Test ICMP Lifecycle",
		Protocol:    SecurityApplicationProtocol(ICMP),
		ICMPType:    SecurityApplicationICMPType(Echo),
	}

	defer deleteSecurityApplication(t, securityApplicationsClient, name)

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
}

func deleteSecurityApplication(t *testing.T, client *SecurityApplicationsClient, name string) {
	deleteInput := DeleteSecurityApplicationInput{
		Name: name,
	}
	err := client.DeleteSecurityApplication(&deleteInput)
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
