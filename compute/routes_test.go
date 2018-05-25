package compute

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_RouteTestAdminDistance   = 1
	_RouteTestDescription     = "acc-testing-route"
	_RouteTestIPAddressPrefix = "10.0.11.0/24"
	_RouteTestName            = "testing-route"
)

func TestAccRoutesLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	rClient, iClient, nClient, vClient, err := getRoutesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// Create a vNIC Set
	ipNetwork, err := createTestIPNetwork(nClient, _IPNetworkTestPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyIPNetwork(t, nClient, ipNetwork.Name)

	instanceInput := &CreateInstanceInput{
		Name:      _RouteTestName,
		Label:     _VirtNicInstanceTestLabel,
		Shape:     _VirtNicInstanceTestShape,
		ImageList: _VirtNicInstanceTestImage,
		Networking: map[string]NetworkingInfo{
			"eth0": {
				IPNetwork: ipNetwork.Name,
				Vnic:      "eth0",
			},
		},
	}

	createdInstance, err := iClient.CreateInstance(instanceInput)
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownInstances(t, iClient, createdInstance.Name, createdInstance.ID)

	vnicSetInput := &CreateVirtualNICSetInput{
		Name:        _RouteTestName,
		Description: _RouteTestDescription,
		VirtualNICs: []string{createdInstance.Networking["eth0"].Vnic},
	}

	createdSet, err := vClient.CreateVirtualNICSet(vnicSetInput)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVirtualNICSet(t, vClient, createdSet.Name)

	// Create the route with the created vNIC set
	routeInput := &CreateRouteInput{
		AdminDistance:   _RouteTestAdminDistance,
		Description:     _RouteTestDescription,
		IPAddressPrefix: _RouteTestIPAddressPrefix,
		Name:            _RouteTestName,
		NextHopVnicSet:  createdSet.Name,
		Tags:            []string{"testing"},
	}

	createdRoute, err := rClient.CreateRoute(routeInput)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteRoute(t, rClient, createdRoute.Name)
	log.Print("Created route successfully")

	getRouteInput := &GetRouteInput{
		Name: _RouteTestName,
	}

	receivedRoute, err := rClient.GetRoute(getRouteInput)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(createdRoute, receivedRoute) {
		t.Fatalf("Mismatch found!\nExpected: %+v\nReceived: %+v", createdRoute, receivedRoute)
	}

	updateRouteInput := &UpdateRouteInput{
		AdminDistance:   _RouteTestAdminDistance + 1,
		Description:     _RouteTestDescription,
		IPAddressPrefix: _RouteTestIPAddressPrefix,
		Name:            _RouteTestName,
		NextHopVnicSet:  createdSet.Name,
		Tags:            []string{"testing"},
	}
	updatedRoute, err := rClient.UpdateRoute(updateRouteInput)
	if err != nil {
		t.Fatal(err)
	}

	receivedRoute, err = rClient.GetRoute(getRouteInput)
	if err != nil {
		t.Fatal(err)
	}

	if receivedRoute.AdminDistance != _RouteTestAdminDistance+1 {
		t.Fatalf("Incorrect Admin Distance found. Expected: %d Recieved: %d", _RouteTestAdminDistance+1, _RouteTestAdminDistance)
	}

	if !reflect.DeepEqual(updatedRoute, receivedRoute) {
		t.Fatalf("Mismatch found!\nExpected: %+v\nReceived: %+v", updatedRoute, receivedRoute)
	}
}

func TestRoutesClient_CreateRoute(t *testing.T) {
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Wrong HTTP Method %s, expected POST", r.Method)
		}

		if r.URL.Path != routesContainerPath {
			t.Errorf("Wrong HTTP Path %q, expected %q", r.URL.Path, routesContainerPath)
		}

		routeInput := &CreateRouteInput{}
		unmarshalRequestBody(t, r, routeInput)

		if routeInput.Name != "/Compute-test/test/test-route" {
			t.Errorf(`Incorrect 'Name''. Got: %q, Expected: "/Compute-test/test/test-route"`, routeInput.Name)
		}

		if routeInput.AdminDistance != _RouteTestAdminDistance {
			t.Errorf("Incorrect 'AdminDistance'. Got: %d Expected: %d", routeInput.AdminDistance, _RouteTestAdminDistance)
		}

		if routeInput.IPAddressPrefix != _RouteTestIPAddressPrefix {
			t.Errorf("Incorrect 'IPAddressPrefix'. Got: %q Expected: %q", routeInput.IPAddressPrefix, _RouteTestIPAddressPrefix)
		}

		if routeInput.NextHopVnicSet != "/Compute-test/test/test-vnic-set" {
			t.Errorf("Incorrect 'NextHopVnicSet'. Got: %q Expected: %q", routeInput.NextHopVnicSet, "/Compute-test/test/test-vnic-set")
		}
		w.Write([]byte(testRouteResponse))
	})

	defer server.Close()
	client, err := getStubRoutesClient(server)
	if err != nil {
		t.Fatalf("error getting stub routes client: %s", err)
	}

	createInput := CreateRouteInput{
		Name:            "test-route",
		AdminDistance:   _RouteTestAdminDistance,
		IPAddressPrefix: _RouteTestIPAddressPrefix,
		NextHopVnicSet:  "test-vnic-set",
	}
	route, err := client.CreateRoute(&createInput)
	if err != nil {
		t.Fatal(err)
	}

	if route.NextHopVnicSet != "test-vnic-set" {
		t.Fatalf("Incorrect response 'NextHopVnicSet'. Got: %q Expected: %q", route.NextHopVnicSet, "test-vnic-set")
	}

	if route.IPAddressPrefix != _RouteTestIPAddressPrefix {
		t.Fatalf("Incorrect response 'IPAddressPrefix'. Got: %q Expected: %q", route.IPAddressPrefix, _RouteTestIPAddressPrefix)
	}
}

var testRouteResponse = fmt.Sprintf(`
{
  "name": "/Compute-acme/jack.jones@example.com/test-route",
  "adminDistance": 1,
  "ipAddressPrefix": %q,
  "nextHopVnicSet": "/Compute-acme/jack.jones@example.com/test-vnic-set"
}
`, _RouteTestIPAddressPrefix)

func getRoutesTestClients() (*RoutesClient, *InstancesClient, *IPNetworksClient, *VirtNICSetsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return client.Routes(), client.Instances(), client.IPNetworks(), client.VirtNICSets(), nil
}

func deleteRoute(t *testing.T, client *RoutesClient, name string) {
	input := &DeleteRouteInput{
		Name: name,
	}
	if err := client.DeleteRoute(input); err != nil {
		t.Fatal(err)
	}
}

func getStubRoutesClient(server *httptest.Server) (*RoutesClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}
	return client.Routes(), nil
}
