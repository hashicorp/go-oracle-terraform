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
		instanceImage = "/oracle/public/OL_7.2_UEKR4_x86_64"
		name          = "test-sec-association"
	)

	iClient, slClient, saClient, err := getSecurityAssociationsTestClients()
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

	createdInstance, err := iClient.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	log.Printf("Instance created: %#v", createdInstance)

	createSecurityListInput := CreateSecurityListInput{
		Name:               name,
		OutboundCIDRPolicy: "DENY",
		Policy:             "PERMIT",
	}

	createdSecurityList, err := slClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security List: %+v", createdSecurityList)
	defer deleteSecurityList(t, slClient, createdSecurityList.Name)

	createSecurityAssociationInput := CreateSecurityAssociationInput{
		SecList: createdSecurityList.Name,
		VCable:  createdInstance.VCableID,
	}
	createdSecurityAssociation, err = saClient.CreateSecurityAssociation(&createSecurityAssociationInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Association: %+v", createdSecurityAssociation)
	defer deleteSecurityAssociation(t, saClient, createdSecurityAssociation.Name)

	getSecurityAssociationInput := GetSecurityAssociationInput{
		Name: createdSecurityAssociation.Name,
	}
	getSecurityAssociationOutput, err := saClient.GetSecurityAssociation(&getSecurityAssociationInput)
	if err != nil {
		t.Fatal(err)
	}
	if createdSecurityAssociation.VCable != getSecurityAssociationOutput.VCable {
		t.Fatalf("Created and retrived vcables don't match.\n Desired: %s\n Actual: %s", createdSecurityAssociation.VCable, getSecurityAssociationOutput.VCable)
	}
	log.Printf("Successfully retrieved Security Association")
}

func getSecurityAssociationsTestClients() (*InstancesClient, *SecurityListsClient, *SecurityAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}

	return client.Instances(), client.SecurityLists(), client.SecurityAssociations(), nil
}

func deleteSecurityAssociation(t *testing.T, client *SecurityAssociationsClient, name string) {
	deleteInput := DeleteSecurityAssociationInput{
		Name: name,
	}
	if err := client.DeleteSecurityAssociation(&deleteInput); err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted Security Association")
}
