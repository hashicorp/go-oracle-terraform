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
		parentPool    string = "ippool:/oracle/public/ippool"
		instanceImage string = "/oracle/public/oel_6.7_apaas_16.4.5_1610211300"
	)

	svc, err := getInstancesClient()
	if err != nil {
		t.Fatal(err)
	}

	input := &CreateInstanceInput{
		Name:      "testacc-ip-association",
		Label:     "testacc-ip-association",
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

	createdInstance, err := svc.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, svc, createdInstance.Name, createdInstance.ID)

	log.Printf("Instance created: %#v\n", createdInstance)

	vcable := createdInstance.VCableID

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
	log.Print("Successfully retrived ip association")

	deleteIPInput := &DeleteIPAssociationInput{
		Name: ipAssociationInfo.Name,
	}
	err = ipac.DeleteIPAssociation(deleteIPInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Successfully deleted IP Association")

	// Instance deletion should be covered by the deferred cleanup function
}

func getIPAssociationsClient() (*IPAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &IPAssociationsClient{}, err
	}

	return client.IPAssociations(), nil
}
