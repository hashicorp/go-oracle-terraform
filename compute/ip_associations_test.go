package compute

import (
	"testing"

	"log"

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

	createdInstanceName, err = svc.LaunchInstance("test", "test", "oc3", instanceImage, nil, nil, []string{},
		map[string]interface{}{
			"attr1": 12,
			"attr2": map[string]interface{}{
				"inner_attr1": "foo",
			},
		})
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Instance created: %#v\n", createdInstanceName)

	instanceInfo, err := svc.WaitForInstanceRunning(createdInstanceName, 300)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Instance retrieved: %#v\n", instanceInfo)

	vcable = instanceInfo.VCableID

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

	getInput := &GetIPAssociationInput{
		Name: createdIPAssociation.Name,
	}
	ipAssociationInfo, err := ipac.GetIPAssociation(getInput)
	if err != nil {
		t.Fatal(err)
	}
	if createdIPAssociation.URI != ipAssociationInfo.URI {
		t.Fatal("IP Association URIs don't match")
	}
	log.Printf("Successfully retrived ip association\n")

	deleteInput := &DeleteIPAssociationInput{
		Name: ipAssociationInfo.Name,
	}
	err = ipac.DeleteIPAssociation(deleteInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted IP Association\n")

	err = svc.DeleteInstance(createdInstanceName)
	if err != nil {
		panic(err)
	}
	log.Printf("Sent Delete instance req\n")
	waitErr := svc.WaitForInstanceDeleted(createdInstanceName, 600)
	if waitErr != nil {
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
