package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

var createdSecurityAssociation *SecurityAssociationInfo

func TestAccSecurityAssociationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	var (
		instanceImage string = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"
		name          string = "test-sec-association"
	)

	svc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

	input := &CreateInstanceInput{
		Name:      name,
		Label:     "testacc-security-association-association",
		Shape:     "oc3",
		ImageList: instanceImage,
		Storage:   nil,
		BootOrder: nil,
		SSHKeys:   []string{},
	}

	createdInstance, err := svc.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, svc, createdInstance.Name, createdInstance.ID)

	log.Printf("Instance created: %#v", createdInstance)

	securityListClient, err := getSecurityListsClient()
	if err != nil {
		t.Fatal(err)
	}

	createSecurityListInput := CreateSecurityListInput{
		Name:               name,
		OutboundCIDRPolicy: "DENY",
		Policy:             "PERMIT",
	}

	createdSecurityList, err := securityListClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security List: %+v", createdSecurityList)
	defer deleteSecurityList(t, securityListClient, createdSecurityList.Name)

	securityAssociationClient, err := getSecurityAssociationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Association Client")

	createSecurityAssociationInput := CreateSecurityAssociationInput{
		SecList: createdSecurityList.Name,
		VCable:  createdInstance.VCableID,
	}
	createdSecurityAssociation, err = securityAssociationClient.CreateSecurityAssociation(&createSecurityAssociationInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Association: %+v", createdSecurityAssociation)
	defer deleteSecurityAssociation(t, securityAssociationClient, createdSecurityAssociation.Name)

	getSecurityAssociationInput := GetSecurityAssociationInput{
		Name: createdSecurityAssociation.Name,
	}
	getSecurityAssociationOutput, err := securityAssociationClient.GetSecurityAssociation(&getSecurityAssociationInput)
	if err != nil {
		t.Fatal(err)
	}
	if createdSecurityAssociation.VCable != getSecurityAssociationOutput.VCable {
		t.Fatalf("Created and retrived vcables don't match.\n Desired: %s\n Actual: %s", createdSecurityAssociation.VCable, getSecurityAssociationOutput.VCable)
	}
	log.Printf("Successfully retrieved Security Association")
}

func getSecurityAssociationsClient() (*SecurityAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecurityAssociationsClient{}, err
	}

	return client.SecurityAssociations(), nil
}

func deleteSecurityAssociation(t *testing.T, client *SecurityAssociationsClient, name string) {
	deleteInput := DeleteSecurityAssociationInput{
		Name: name,
	}
	err := client.DeleteSecurityAssociation(&deleteInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted Security Association")
}
