package compute

import (
	"fmt"
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

	//TODO: Initial Implementation
	//createdInstanceName, err = svc.LaunchInstance("test", "test", "oc3", "/oracle/public/oel_6.4_2GB_v1", nil, nil, []string{},
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
	fmt.Printf("Instance created: %#v\n", createdInstanceName)

	//TODO: Initial Implementation
	//instanceInfo, err := svc.WaitForInstanceRunning(createdInstanceName, 120)
	instanceInfo, err := svc.WaitForInstanceRunning(createdInstanceName, 300)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Instance retrieved: %#v\n", instanceInfo)

	vcable = instanceInfo.VCableID

	ipac, err := getIPAssociationsClient()
	if err != nil {
		t.Fatal(err)
	}

	createdIPAssociation, err := ipac.CreateIPAssociation(vcable, parentPool)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Created IP Association: %+v\n", createdIPAssociation)

	ipAssociationInfo, err := ipac.GetIPAssociation(createdIPAssociation.Name)
	if err != nil {
		t.Fatal(err)
	}
	if createdIPAssociation.URI != ipAssociationInfo.URI {
		t.Fatal("IP Association URIs don't match")
	}
	fmt.Printf("Successfully retrived ip association\n")

	err = ipac.DeleteIPAssociation(ipAssociationInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully deleted IP Association\n")

	err = svc.DeleteInstance(createdInstanceName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sent Delete instance req\n")
	//TODO: Initial Implementation
	//err = svc.WaitForInstanceDeleted(createdInstanceName, 600)
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
