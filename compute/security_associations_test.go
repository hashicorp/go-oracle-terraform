package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccSecurityAssociationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	var (
		instanceImage string = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"
	)

	defer func() {
		if err := tearDownInstances(); err != nil {
			log.Printf("Error deleting instance: %#v", createdInstance)
			log.Print("Dangling resources may occur!")
			t.Fatalf("Error: %v", err)
		}
		if err := tearDownSecurityLists(); err != nil {
			log.Printf("Error deleting security list: %#v", createdSecurityList)
			log.Print("Dangling resources may occur!")
			t.Fatalf("Error: %v", err)
		}
	}()

	svc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

	input := &CreateInstanceInput{
		Name:      "testacc-security-association",
		Label:     "testacc-security-association-association",
		Shape:     "oc3",
		ImageList: instanceImage,
		Storage:   nil,
		BootOrder: nil,
		SSHKeys:   []string{},
	}

	createdInstance, err = svc.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Instance created: %#v\n", createdInstance)

	securityListClient, err := getSecurityListsClient()
	if err != nil {
		t.Fatal(err)
	}

	createSecurityListInput := CreateSecurityListInput{
		Name:               "test-sec-list",
		OutboundCIDRPolicy: "DENY",
		Policy:             "PERMIT",
	}

	createdSecurityList, err = securityListClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security List: %+v", createdSecurityList)

	securityAssociationClient, err := getSecurityAssociationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Association Client")

	createSecurityAssociationInput := CreateSecurityAssociationInput{
		SecList: createdSecurityList.Name,
		VCable:  createdInstance.VCableID,
	}
	securityAssociation, err := securityAssociationClient.CreateSecurityAssociation(&createSecurityAssociationInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Association: %+v", securityAssociation)

	getSecurityAssociationInput := GetSecurityAssociationInput{
		Name: securityAssociation.Name,
	}
	getSecurityAssociationOutput, err := securityAssociationClient.GetSecurityAssociation(&getSecurityAssociationInput)
	if err != nil {
		t.Fatal(err)
	}
	if securityAssociation.VCable != getSecurityAssociationOutput.VCable {
		t.Fatalf("Created and retrived vcables don't match.\n Desired: %s\n Actual: %s", securityAssociation.VCable, getSecurityAssociationOutput.VCable)
	}
	log.Printf("Successfully retrieved Security Association")

	deleteSecurityAssociationInput := DeleteSecurityAssociationInput{
		Name: securityAssociation.Name,
	}
	err = securityAssociationClient.DeleteSecurityAssociation(&deleteSecurityAssociationInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted Security Association")
}

func getSecurityAssociationsClient() (*SecurityAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecurityAssociationsClient{}, err
	}

	return client.SecurityAssociations(), nil
}
