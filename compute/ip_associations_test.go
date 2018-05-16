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
		parentPool    = "ippool:/oracle/public/ippool"
		instanceImage = "/oracle/public/OL_7.2_UEKR4_x86_64"
	)

	iClient, ipaClient, err := getIPAssociationsTestClients()
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

	createdInstance, err := iClient.CreateInstance(input)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	log.Printf("Instance created: %#v", createdInstance)

	vcable := createdInstance.VCableID

	createInput := &CreateIPAssociationInput{
		VCable:     vcable,
		ParentPool: parentPool,
	}

	createdIPAssociation, err := ipaClient.CreateIPAssociation(createInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Created IP Association: %+v", createdIPAssociation)

	getIPInput := &GetIPAssociationInput{
		Name: createdIPAssociation.Name,
	}
	ipAssociationInfo, err := ipaClient.GetIPAssociation(getIPInput)
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
	err = ipaClient.DeleteIPAssociation(deleteIPInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Successfully deleted IP Association")

	// Instance deletion should be covered by the deferred cleanup function
}

func getIPAssociationsTestClients() (*InstancesClient, *IPAssociationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}

	return client.Instances(), client.IPAssociations(), nil
}
