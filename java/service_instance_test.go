package java

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_ServiceInstanceName                        = "testingjavaserviceinstance"
	_ServiceInstanceLevel                       = "PAAS"
	_ServiceInstanceSubscriptionType            = "HOURLY"
	_ServiceInstanceType                        = "weblogic"
	_ServiceInstanceDBAUser                     = "sys"
	_ServiceInstanceDBAPassword                 = "Test_String7"
	_ServiceInstanceDBServiceName               = "test-service-instance-matthew"
	_ServiceInstanceShape                       = "oc1m"
	_ServiceInstanceVersion                     = "12.2.1"
	_ServiceInstanceAdminUsername               = "sdk-user"
	_ServiceInstanceAdminPassword               = "Test_String7"
	_ServiceInstancePubKey                      = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
	_ServiceInstanceCloudStorageContainer       = "Storage-canonical/test-java-instance"
	_ServiceInstanceCloudStorageCreateIfMissing = true
)

func TestAccServiceInstanceLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	parameter := Parameter{
		Type:          _ServiceInstanceType,
		DBAName:       _ServiceInstanceDBAUser,
		DBAPassword:   _ServiceInstanceDBAPassword,
		DBServiceName: _ServiceInstanceDBServiceName,
		Shape:         _ServiceInstanceShape,
		Version:       _ServiceInstanceVersion,
		AdminUsername: _ServiceInstanceAdminUsername,
		AdminPassword: _ServiceInstanceAdminPassword,
		VMsPublicKey:  _ServiceInstancePubKey,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		CloudStorageContainer:           _ServiceInstanceCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
		ServiceName:                     _ServiceInstanceName,
		Level:                           _ServiceInstanceLevel,
		SubscriptionType:                _ServiceInstanceSubscriptionType,
		Parameters:                      []Parameter{parameter},
	}

	_, err = siClient.CreateServiceInstance(createServiceInstance)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyServiceInstance(t, siClient, _ServiceInstanceName)

	getInput := &GetServiceInstanceInput{
		Name: _ServiceInstanceName,
	}

	receivedRes, err := siClient.GetServiceInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}
	if receivedRes.ServiceName != _ServiceInstanceName {
		t.Fatal(fmt.Errorf("Names do not match. Wanted: %s Received: %s", _ServiceInstanceName, receivedRes.ServiceName))
	}
}

func getServiceInstanceTestClients() (*ServiceInstanceClient, error) {
	client, err := getJavaTestClient(&opc.Config{})
	if err != nil {
		return &ServiceInstanceClient{}, err
	}

	return client.ServiceInstanceClient(), nil
}

func destroyServiceInstance(t *testing.T, client *ServiceInstanceClient, name string) {
	input := &DeleteServiceInstanceInput{
		Name:        name,
		DBAUsername: _ServiceInstanceDBAUser,
		DBAPassword: _ServiceInstanceDBAPassword,
	}

	if err := client.DeleteServiceInstance(input); err != nil {
		t.Fatal(err)
	}
}
