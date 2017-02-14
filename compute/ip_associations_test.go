package compute

import (
	"fmt"
	"testing"
)

func TestIPAssociationLifeCycle(t *testing.T) {
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
	fmt.Printf("Created IP Association: %+v", createdIPAssociation)

	ipAssociationInfo, err := ipac.GetIPAssociation(createdIPAssociation.Name)
	if err != nil {
		t.Fatal(err)
	}
	if createdIPAssociation.URI != ipAssociationInfo.URI {
		t.Fatal("IP Association URIs don't match")
	}
	fmt.Printf("Successfully retrived ip association")

	err = ipac.DeleteIPAssociation(ipAssociationInfo.Name)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully deleted IP Association")

	err = svc.DeleteInstance(createdInstanceName)
	if err != nil {
		panic(err)
	}
	//TODO: Initial Implementation
	//err = svc.WaitForInstanceDeleted(createdInstanceName, 600)
	err = svc.WaitForInstanceDeleted(createdInstanceName, 900)
	if err != nil {
		panic(err)
	}
}

func getIPAssociationsClient() (*IPAssociationsClient, error) {
	authenticatedClient, err := getAuthenticatedClient()
	if err != nil {
		return &IPAssociationsClient{}, err
	}

	return authenticatedClient.IPAssociations(), nil
}
