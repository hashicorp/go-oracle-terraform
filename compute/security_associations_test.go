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

	log.Printf("Instance created: %#v", createdInstance)

	securityListClient, err := getSecurityListsClient()
	if err != nil {
		t.Fatal(err)
	}

	createSecurityListInput := CreateSecurityListInput{
		Name:               "test-sec-list",
		OutboundCIDRPolicy: "DENY",
		Policy:             "PERMIT",
	}

	securityList, err := securityListClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security List: %+v", securityList)

	securityAssociationClient, err := getSecurityAssociationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Association Client")

	createSecurityAssociationInput := CreateSecurityAssociationInput{
		SecList: securityList.Name,
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

	deleteSecurityListInput := DeleteSecurityListInput{
		Name: securityList.Name,
	}
	err = securityListClient.DeleteSecurityList(&deleteSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted Security List")
}

func getSecurityAssociationsClient() (*SecurityAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecurityAssociationsClient{}, err
	}

	return client.SecurityAssociations(), nil
}
