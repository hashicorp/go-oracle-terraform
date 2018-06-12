package java

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_ServiceInstanceName                        = "testingjavaserviceinstance1"
	_ServiceInstanceLevel                       = "PAAS"
	_ServiceInstanceSubscriptionType            = "HOURLY"
	_ServiceInstanceDBAUser                     = "sys"
	_ServiceInstanceDBAPassword                 = "Test_String7"
	_ServiceInstanceShape                       = "oc3"
	_ServiceInstanceUpdateShape                 = "oc5"
	_ServiceInstanceVersion                     = "12cRelease212"
	_ServiceInstanceAdminUsername               = "sdk-user"
	_ServiceInstanceAdminPassword               = "Test_String7"
	_ServiceInstancePubKey                      = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
	_ServiceInstanceCloudStorageContainer       = "Storage-a459477/test-java-instance"
	_ServiceInstanceCloudStorageCreateIfMissing = true
	_ServiceInstanceManagedServerCount          = 0
	_ServiceInstanceProvisionOTD                = true
	// Database specific configuration
	_ServiceInstanceDatabaseName            = "testing-java-instance-service1"
	_ServiceInstanceBackupDestinationBoth   = "BOTH"
	_ServiceInstanceDBSID                   = "ORCL"
	_ServiceInstanceDBType                  = "db"
	_ServiceInstanceUsableStorage           = "25"
	_ServiceInstanceDBCloudStorageContainer = "Storage-a459477/test-db-java-instancea"
	_ServiceInstanceEdition                 = "EE"
	_ServiceInstanceDatabaseShape           = "oc3"
	_ServiceInstanceDBVersion               = "12.2.0.1"
)

func TestAccServiceInstanceLifeCycle_Basic(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, dClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	databaseParameter := database.ParameterInput{
		AdminPassword:                   _ServiceInstanceDBAPassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceDBType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceDBCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createDatabaseServiceInstance := &database.CreateServiceInstanceInput{
		Name:             _ServiceInstanceDatabaseName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceDatabaseShape,
		SubscriptionType: _ServiceInstanceSubscriptionType,
		Version:          _ServiceInstanceDBVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        databaseParameter,
	}

	_, err = dClient.CreateServiceInstance(createDatabaseServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDatabaseServiceInstance(t, dClient, _ServiceInstanceDatabaseName)

	wlsConfig := &CreateWLS{
		DBAName:            _ServiceInstanceDBAUser,
		DBAPassword:        _ServiceInstanceDBAPassword,
		DBServiceName:      _ServiceInstanceDatabaseName,
		Shape:              _ServiceInstanceShape,
		ManagedServerCount: _ServiceInstanceManagedServerCount,
		AdminUsername:      _ServiceInstanceAdminUsername,
		AdminPassword:      _ServiceInstanceAdminPassword,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		CloudStorageContainer:             _ServiceInstanceCloudStorageContainer,
		CloudStorageContainerAutoGenerate: _ServiceInstanceCloudStorageCreateIfMissing,
		ServiceName:                       _ServiceInstanceName,
		ServiceLevel:                      _ServiceInstanceLevel,
		Components:                        CreateComponents{WLS: wlsConfig},
		VMPublicKeyText:                   _ServiceInstancePubKey,
		Edition:                           ServiceInstanceEditionSuite,
		ServiceVersion:                    _ServiceInstanceVersion,
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
	assert.Equal(t, _ServiceInstanceName, receivedRes.ServiceName, "Service instance name not expected.")
}

func TestAccServiceInstanceLifeCycle_typeOTD(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, dClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}
	databaseParameter := database.ParameterInput{
		AdminPassword:                   _ServiceInstanceDBAPassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceDBType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceDBCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createDatabaseServiceInstance := &database.CreateServiceInstanceInput{
		Name:             _ServiceInstanceDatabaseName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceShape,
		SubscriptionType: _ServiceInstanceSubscriptionType,
		Version:          _ServiceInstanceDBVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        databaseParameter,
	}

	_, err = dClient.CreateServiceInstance(createDatabaseServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDatabaseServiceInstance(t, dClient, _ServiceInstanceDatabaseName)

	wlsConfig := &CreateWLS{
		DBAName:            _ServiceInstanceDBAUser,
		DBAPassword:        _ServiceInstanceDBAPassword,
		DBServiceName:      _ServiceInstanceDatabaseName,
		Shape:              _ServiceInstanceShape,
		ManagedServerCount: _ServiceInstanceManagedServerCount,
		AdminUsername:      _ServiceInstanceAdminUsername,
		AdminPassword:      _ServiceInstanceAdminPassword,
	}

	otdConfig := &CreateOTD{
		AdminUsername: _ServiceInstanceAdminUsername,
		AdminPassword: _ServiceInstanceAdminPassword,
		Shape:         _ServiceInstanceShape,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		ProvisionOTD:                      _ServiceInstanceProvisionOTD,
		CloudStorageContainer:             _ServiceInstanceCloudStorageContainer,
		CloudStorageContainerAutoGenerate: _ServiceInstanceCloudStorageCreateIfMissing,
		ServiceName:                       _ServiceInstanceName,
		ServiceLevel:                      _ServiceInstanceLevel,
		Components:                        CreateComponents{WLS: wlsConfig, OTD: otdConfig},
		VMPublicKeyText:                   _ServiceInstancePubKey,
		Edition:                           ServiceInstanceEditionSuite,
		ServiceVersion:                    _ServiceInstanceVersion,
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
	assert.Equal(t, _ServiceInstanceName, receivedRes.ServiceName, "Service instance name not expected.")
	assert.NotEmpty(t, receivedRes.OTDRoot, "Expected OTDROot to not be empty")
}

func TestAccServiceInstanceLifeCycle_ScaleUp(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, dClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	databaseParameter := database.ParameterInput{
		AdminPassword:                   _ServiceInstanceDBAPassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceDBType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceDBCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createDatabaseServiceInstance := &database.CreateServiceInstanceInput{
		Name:             _ServiceInstanceDatabaseName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceDatabaseShape,
		SubscriptionType: _ServiceInstanceSubscriptionType,
		Version:          _ServiceInstanceDBVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        databaseParameter,
	}

	_, err = dClient.CreateServiceInstance(createDatabaseServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDatabaseServiceInstance(t, dClient, _ServiceInstanceDatabaseName)

	wlsConfig := &CreateWLS{
		DBAName:            _ServiceInstanceDBAUser,
		DBAPassword:        _ServiceInstanceDBAPassword,
		DBServiceName:      _ServiceInstanceDatabaseName,
		Shape:              _ServiceInstanceShape,
		ManagedServerCount: _ServiceInstanceManagedServerCount,
		AdminUsername:      _ServiceInstanceAdminUsername,
		AdminPassword:      _ServiceInstanceAdminPassword,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		CloudStorageContainer:             _ServiceInstanceCloudStorageContainer,
		CloudStorageContainerAutoGenerate: _ServiceInstanceCloudStorageCreateIfMissing,
		ServiceName:                       _ServiceInstanceName,
		ServiceLevel:                      _ServiceInstanceLevel,
		Components:                        CreateComponents{WLS: wlsConfig},
		VMPublicKeyText:                   _ServiceInstancePubKey,
		Edition:                           ServiceInstanceEditionSuite,
		ServiceVersion:                    _ServiceInstanceVersion,
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

	var hostname string
	for _, instance := range receivedRes.Components.WLS.VMInstances {
		hostname = instance.HostName
	}

	if hostname == "" {
		t.Fatal(fmt.Errorf("Unable to find hostname to scale"))
	}

	wlsComponent := ScaleUpDownWLS{
		Shape: _ServiceInstanceUpdateShape,
		Hosts: []string{hostname},
	}

	component := ScaleUpDownComponent{
		WLS: wlsComponent,
	}

	scaleUpInput := &ScaleUpDownServiceInstanceInput{
		Name:       _ServiceInstanceName,
		Components: component,
	}

	err = siClient.ScaleUpDownServiceInstance(scaleUpInput)
	if err != nil {
		t.Fatal(err)
	}

	receivedRes, err = siClient.GetServiceInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, _ServiceInstanceUpdateShape, receivedRes.Components.WLS.VMInstances[hostname].ShapeID, "Service instance shape not expected.")
}

func TestAccServiceInstanceLifeCycle_Stop(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, dClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	databaseParameter := database.ParameterInput{
		AdminPassword:                   _ServiceInstanceDBAPassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceDBType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceDBCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createDatabaseServiceInstance := &database.CreateServiceInstanceInput{
		Name:             _ServiceInstanceDatabaseName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceDatabaseShape,
		SubscriptionType: _ServiceInstanceSubscriptionType,
		Version:          _ServiceInstanceDBVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        databaseParameter,
	}

	_, err = dClient.CreateServiceInstance(createDatabaseServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDatabaseServiceInstance(t, dClient, _ServiceInstanceDatabaseName)

	wlsConfig := &CreateWLS{
		DBAName:            _ServiceInstanceDBAUser,
		DBAPassword:        _ServiceInstanceDBAPassword,
		DBServiceName:      _ServiceInstanceDatabaseName,
		Shape:              _ServiceInstanceShape,
		ManagedServerCount: _ServiceInstanceManagedServerCount,
		AdminUsername:      _ServiceInstanceAdminUsername,
		AdminPassword:      _ServiceInstanceAdminPassword,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		CloudStorageContainer:             _ServiceInstanceCloudStorageContainer,
		CloudStorageContainerAutoGenerate: _ServiceInstanceCloudStorageCreateIfMissing,
		ServiceName:                       _ServiceInstanceName,
		ServiceLevel:                      _ServiceInstanceLevel,
		Components:                        CreateComponents{WLS: wlsConfig},
		VMPublicKeyText:                   _ServiceInstancePubKey,
		Edition:                           ServiceInstanceEditionSuite,
		ServiceVersion:                    _ServiceInstanceVersion,
	}

	_, err = siClient.CreateServiceInstance(createServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyServiceInstance(t, siClient, _ServiceInstanceName)

	getInput := &GetServiceInstanceInput{
		Name: _ServiceInstanceName,
	}

	desiredStateInput := &DesiredStateInput{
		Name:            _ServiceInstanceName,
		AllServiceHosts: true,
		LifecycleState:  ServiceInstanceLifecycleStateStop,
	}

	err = siClient.UpdateDesiredState(desiredStateInput)
	if err != nil {
		t.Fatal(err)
	}

	receivedRes, err := siClient.GetServiceInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ServiceInstanceStatusStopped, receivedRes.State, "Service instance status not expected.")
}

func TestAccServiceInstanceLifeCycle_RestartHost(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	siClient, dClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	databaseParameter := database.ParameterInput{
		AdminPassword:                   _ServiceInstanceDBAPassword,
		BackupDestination:               _ServiceInstanceBackupDestinationBoth,
		SID:                             _ServiceInstanceDBSID,
		Type:                            _ServiceInstanceDBType,
		UsableStorage:                   _ServiceInstanceUsableStorage,
		CloudStorageContainer:           _ServiceInstanceDBCloudStorageContainer,
		CreateStorageContainerIfMissing: _ServiceInstanceCloudStorageCreateIfMissing,
	}

	createDatabaseServiceInstance := &database.CreateServiceInstanceInput{
		Name:             _ServiceInstanceDatabaseName,
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceDatabaseShape,
		SubscriptionType: _ServiceInstanceSubscriptionType,
		Version:          _ServiceInstanceDBVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        databaseParameter,
	}

	_, err = dClient.CreateServiceInstance(createDatabaseServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDatabaseServiceInstance(t, dClient, _ServiceInstanceDatabaseName)

	wlsConfig := &CreateWLS{
		DBAName:            _ServiceInstanceDBAUser,
		DBAPassword:        _ServiceInstanceDBAPassword,
		DBServiceName:      _ServiceInstanceDatabaseName,
		Shape:              _ServiceInstanceShape,
		ManagedServerCount: _ServiceInstanceManagedServerCount,
		AdminUsername:      _ServiceInstanceAdminUsername,
		AdminPassword:      _ServiceInstanceAdminPassword,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		CloudStorageContainer:             _ServiceInstanceCloudStorageContainer,
		CloudStorageContainerAutoGenerate: _ServiceInstanceCloudStorageCreateIfMissing,
		ServiceName:                       _ServiceInstanceName,
		ServiceLevel:                      _ServiceInstanceLevel,
		Components:                        CreateComponents{WLS: wlsConfig},
		VMPublicKeyText:                   _ServiceInstancePubKey,
		Edition:                           ServiceInstanceEditionSuite,
		ServiceVersion:                    _ServiceInstanceVersion,
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

	var hostname string
	for _, instance := range receivedRes.Components.WLS.VMInstances {
		hostname = instance.HostName
	}

	if hostname == "" {
		t.Fatal(fmt.Errorf("Unable to find hostname to scale"))
	}

	wlsComponent := &DesiredStateHost{
		Hosts: []string{hostname},
	}

	component := &DesiredStateComponent{
		WLS: wlsComponent,
	}

	desiredStateInput := &DesiredStateInput{
		Name:           _ServiceInstanceName,
		Components:     component,
		LifecycleState: ServiceInstanceLifecycleStateRestart,
	}

	err = siClient.UpdateDesiredState(desiredStateInput)
	if err != nil {
		t.Fatal(err)
	}

	receivedRes, err = siClient.GetServiceInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ServiceInstanceStatusReady, receivedRes.State, "Service instance status not expected.")
}

func getServiceInstanceTestClients() (*ServiceInstanceClient, *database.ServiceInstanceClient, error) {
	client, err := getJavaTestClient(&opc.Config{})
	if err != nil {
		return &ServiceInstanceClient{}, &database.ServiceInstanceClient{}, err
	}

	dClient, err := database.GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return &ServiceInstanceClient{}, &database.ServiceInstanceClient{}, err
	}

	return client.ServiceInstanceClient(), dClient.ServiceInstanceClient(), nil
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

func destroyDatabaseServiceInstance(t *testing.T, client *database.ServiceInstanceClient, name string) {
	input := &database.DeleteServiceInstanceInput{
		Name: name,
	}

	if err := client.DeleteServiceInstance(input); err != nil {
		t.Fatal(err)
	}
}
