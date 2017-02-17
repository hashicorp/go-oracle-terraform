package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccIPAssociationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	var (
		vcable        string
		parentPool    string = "ippool:/oracle/public/ippool"
		instanceImage string = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"
	)
	svc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

	input := &CreateInstanceInput{
		Name:      "test",
		Label:     "test",
		Shape:     "oc3",
		ImageList: instanceImage,
		Storage:   nil,
		BootOrder: nil,
		SSHKeys:   []string{},
		Attributes: map[string]interface{}{
			"attr1": 12,
			"attr2": map[string]interface{}{
				"inner_attr1": "foo",
			},
		},
	}

	createdInstance, err = svc.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Instance created: %#v\n", createdInstance)

	getInput := &GetInstanceInput{
		Name: createdInstance.Name,
		ID:   createdInstance.ID,
	}

	if err := svc.WaitForInstanceRunning(getInput, 300); err != nil {
		t.Fatal(err)
	}
	log.Print("Instance retrieved")

	vcable = createdInstance.VCableID

	ipac, err := getIPAssociationsClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateIPAssociationInput{
		VCable:     vcable,
		ParentPool: parentPool,
	}

	createdIPAssociation, err := ipac.CreateIPAssociation(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Created IP Association: %+v\n", createdIPAssociation)

	getIPInput := &GetIPAssociationInput{
		Name: createdIPAssociation.Name,
	}
	ipAssociationInfo, err := ipac.GetIPAssociation(getIPInput)
	if err != nil {
		t.Fatal(err)
	}
	if createdIPAssociation.URI != ipAssociationInfo.URI {
		t.Fatal("IP Association URIs don't match")
	}
	log.Printf("Successfully retrived ip association\n")

	deleteIPInput := &DeleteIPAssociationInput{
		Name: ipAssociationInfo.Name,
	}
	err = ipac.DeleteIPAssociation(deleteIPInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted IP Association\n")

	deleteInput := &DeleteInstanceInput{
		Name: createdInstance.Name,
		ID:   createdInstance.ID,
	}
	if err := svc.DeleteInstance(deleteInput); err != nil {
		panic(err)
	}

	log.Print("Sent Delete instance req")

	if err := svc.WaitForInstanceDeleted(deleteInput, 600); err != nil {
		panic(err)
	}
}

func getIPAssociationsClient() (*IPAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &IPAssociationsClient{}, err
	}

	return client.IPAssociations(), nil
}
