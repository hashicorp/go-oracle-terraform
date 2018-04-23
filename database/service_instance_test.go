package database

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_ServiceInstanceName                        = "testing-db-service-instance"
	_ServiceInstanceEdition                     = "EE"
	_ServiceInstanceLevel                       = "PAAS"
	_ServiceInstanceShape                       = "oc3"
	_ServiceInstanceUpdateShape                 = "oc4"
	_ServiceInstanceSubscription                = "HOURLY"
	_ServiceInstanceVersion                     = "12.2.0.1"
	_ServiceInstancePubKey                      = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
	_ServiceInstancePassword                    = "Test_String7"
	_ServiceInstanceBackupDestination           = "NONE"
	_ServiceInstanceDBSID                       = "ORCL"
	_ServiceInstanceType                        = "db"
	_ServiceInstanceUsableStorage               = "15"
	_ServiceInstanceCloudStorageContainer       = "Storage-a459477/test-database-instance"
	_ServiceInstanceCloudStorageCreateIfMissing = true
	_ServiceInstanceBackupDestinationBoth       = "BOTH"
	_ServiceInstanceDeleteBackup                = true
)

func TestAccServiceInstanceLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	parameter := ParameterInput{
		AdminPassword:     _ServiceInstancePassword,
		BackupDestination: _ServiceInstanceBackupDestination,
		SID:               _ServiceInstanceDBSID,
		Type:              _ServiceInstanceType,
		UsableStorage:     _ServiceInstanceUsableStorage,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		Name:             _ServiceInstanceName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceShape,
		SubscriptionType: _ServiceInstanceSubscription,
		Version:          _ServiceInstanceVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        parameter,
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
	if receivedRes.Name != _ServiceInstanceName {
		t.Fatal(fmt.Errorf("Names do not match. Wanted: %s Received: %s", _ServiceInstanceName, receivedRes.Name))
	}
}

func TestAccServiceInstanceCloudStorage(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	parameter := ParameterInput{
		AdminPassword:                   _ServiceInstancePassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		Name:             _ServiceInstanceName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceShape,
		SubscriptionType: _ServiceInstanceSubscription,
		Version:          _ServiceInstanceVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        parameter,
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
	if receivedRes.Name != _ServiceInstanceName {
		t.Fatal(fmt.Errorf("Names do not match. Wanted: %s Received: %s", _ServiceInstanceName, receivedRes.Name))
	}
	if receivedRes.CloudStorageContainer != _ServiceInstanceCloudStorageContainer {
		t.Fatal(fmt.Errorf("Cloud storage containers do not match. Wanted: %s Received: %s", _ServiceInstanceCloudStorageContainer, receivedRes.CloudStorageContainer))
	}
}

func TestAccServiceInstanceUpdate(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	parameter := ParameterInput{
		AdminPassword:                   _ServiceInstancePassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		Name:             _ServiceInstanceName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceShape,
		SubscriptionType: _ServiceInstanceSubscription,
		Version:          _ServiceInstanceVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        parameter,
	}

	_, err = siClient.CreateServiceInstance(createServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyServiceInstance(t, siClient, _ServiceInstanceName)

	updateInput := &UpdateServiceInstanceInput{
		Name:  _ServiceInstanceName,
		Shape: _ServiceInstanceUpdateShape,
	}
	_, err = siClient.UpdateServiceInstance(updateInput)
	if err != nil {
		t.Fatal(err)
	}

	getInput := &GetServiceInstanceInput{
		Name: _ServiceInstanceName,
	}

	receivedRes, err := siClient.GetServiceInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}
	if receivedRes.Shape != _ServiceInstanceUpdateShape {
		t.Fatal(fmt.Errorf("Shapes do not match. Wanted: %s Received: %s", _ServiceInstanceUpdateShape, receivedRes.Shape))
	}
}

func TestAccServiceInstanceDesiredState(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	parameter := ParameterInput{
		AdminPassword:                   _ServiceInstancePassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		Name:             _ServiceInstanceName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceShape,
		SubscriptionType: _ServiceInstanceSubscription,
		Version:          _ServiceInstanceVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        parameter,
	}

	_, err = siClient.CreateServiceInstance(createServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyServiceInstance(t, siClient, _ServiceInstanceName)

	desiredStateInput := &DesiredStateInput{
		Name:           _ServiceInstanceName,
		LifecycleState: ServiceInstanceLifecycleStateStop,
	}
	_, err = siClient.UpdateDesiredState(desiredStateInput)
	if err != nil {
		t.Fatal(err)
	}

	getInput := &GetServiceInstanceInput{
		Name: _ServiceInstanceName,
	}

	receivedRes, err := siClient.GetServiceInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}
	if receivedRes.Status != ServiceInstanceStopped {
		t.Fatal(fmt.Errorf("Statuses do not match. Wanted: %s Received: %s", ServiceInstanceStopped, receivedRes.Status))
	}
}

func getServiceInstanceTestClients() (*ServiceInstanceClient, error) {
	client, err := GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return &ServiceInstanceClient{}, err
	}

	return client.ServiceInstanceClient(), nil
}

func destroyServiceInstance(t *testing.T, client *ServiceInstanceClient, name string) {
	input := &DeleteServiceInstanceInput{
		Name:         name,
		DeleteBackup: _ServiceInstanceDeleteBackup,
	}

	if err := client.DeleteServiceInstance(input); err != nil {
		t.Fatal(err)
	}
}
