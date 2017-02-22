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

	name := "test-sec-app-tcp"
	protocol := "all"
	dport := "19336"
	icmpType := ""
	icmpCode := ""
	description := "Terraform Acceptance Test TCP Lifecycle"

	createdSecurityApplication, err := securityApplicationsClient.CreateSecurityApplication(name, protocol, dport, icmpType, icmpCode, description)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Application: %+v", createdSecurityApplication)

	retrieveInput := GetSecurityApplicationInput{
		Name: name,
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
		Name: name,
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

	name := "test-sec-app-icmp"
	protocol := "icmp"
	dport := ""
	icmpType := "echo"
	icmpCode := ""
	description := "Terraform Acceptance Test ICMP Lifecycle"

	createdSecurityApplication, err := securityApplicationsClient.CreateSecurityApplication(name, protocol, dport, icmpType, icmpCode, description)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Application: %+v", createdSecurityApplication)

	retrieveInput := GetSecurityApplicationInput{
		Name: name,
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
		Name: name,
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
