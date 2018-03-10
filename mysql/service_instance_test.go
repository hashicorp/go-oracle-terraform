package mysql

import (
	"fmt"
	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"testing"
)

const (
	_ServiceInstanceName              = "testing-mysql-service-instance"
	_ServiceInstanceDesc              = "MySQL Terraform Test Instance"
	_ServiceInstancePubKey            = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC3QxPp0BFK+ligB9m1FBcFELyvN5EdNUoSwTCe4Zv2b51OIO6wGM/dvTr/yj2ltNA/Vzl9tqf9AUBL8tKjAOk8uukip6G7rfigby+MvoJ9A8N0AC2te3TI+XCfB5Ty2M2OmKJjPOPCd6+OdzhT4cWnPOM+OAiX0DP7WCkO4Kx2kntf8YeTEurTCspOrRjGdo+zZkJxEydMt31asu9zYOTLmZPwLCkhel8vY6SnZhDTNSNkRzxZFv+Mh2VGmqu4SSxfVXr4tcFM6/MbAXlkA8jo+vHpy5sC79T4uNaPu2D8Ed7uC3yDdO3KRVdzZCfWHj4NjixdMs2CtK6EmyeVOPuiYb8/mcTybrb4F/CqA4jydAU6Ok0j0bIqftLyxNgfS31hR1Y3/GNPzly4+uUIgZqmsuVFh5h0L7qc1jMv7wRHphogo5snIp45t9jWNj8uDGzQgWvgbFP5wR7Nt6eS0kaCeGQbxWBDYfjQE801IrwhgMfmdmGw7FFveCH0tFcPm6td/8kMSyg/OewczZN3T62ETQYVsExOxEQl2t4SZ/yqklg+D9oGM+ILTmBRzIQ2m/xMmsbowiTXymjgVmvrWuc638X6dU2fKJ7As4hxs3rA1BA5sOt0XyqfHQhtYrL/Ovb1iV+C7MRhKicTyoNTc7oVcDDG0VW785d8CPqttDi50w=="
	_ServiceInstanceBackupDestination = "NONE"

	_Service_MySQLDBName   = "demo_db"
	_Service_MySQLStorage  = "25"
	_Service_MySQLPort     = "3306"
	_Service_MySQLUser     = "root"
	_Service_MySQLPassword = "MySqlPassword_1"
	_Service_MySQLShape    = "oc3"
)

func TestAccServiceInstanceLifeCycle(t *testing.T) {

	helper.Test(t, helper.TestCase{})

	siClient, err := getServiceInstanceTestClients()
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the input parameters
	serviceParameter := ServiceParameters{
		BackupDestination:  _ServiceInstanceBackupDestination,
		ServiceDescription: _ServiceInstanceDesc,
		ServiceName:        _ServiceInstanceName,
		VMPublicKeyText:    _ServiceInstancePubKey,
	}

	mySQLParameter := MySQLParameters{
		DBName:            _Service_MySQLDBName,
		DBStorage:         _Service_MySQLStorage,
		MysqlPort:         _Service_MySQLPort,
		MysqlUserName:     _Service_MySQLUser,
		MysqlUserPassword: _Service_MySQLPassword,
		Shape:             _Service_MySQLShape,
	}

	componentParameter := ComponentParameters{
		Mysql: mySQLParameter,
	}

	createServiceInstance := &CreateServiceInstanceInput{
		ComponentParameters: componentParameter,
		ServiceParameters:   serviceParameter,
	}

	t.Log("... Creating MySQL Instance")
	_, err = siClient.CreateServiceInstance(createServiceInstance)
	if err != nil {
		t.Fatal(err)
	}

	defer destroyServiceInstance(t, siClient, _ServiceInstanceName)

	getInput := &GetServiceInstanceInput{
		Name: _ServiceInstanceName,
	}

	t.Log("... Verifying MySQL Instance")
	receivedRes, err := siClient.GetServiceInstance(getInput)
	if err != nil {
		t.Fatal(err)
	}

	if receivedRes.ServiceName != _ServiceInstanceName {
		t.Fatal(fmt.Errorf("Names do not match. Wanted: %s Received: %s", _ServiceInstanceName, receivedRes.ServiceName))
	}
}

/*
func createAccessRuleInstance(t *testing.T, siClient *ServiceInstanceClient) {

	t.Log("... Creating Access Rules")
	createAccessRuleInfo := &AccessRuleInfo{
		Description: _Service_AccessRule_Description,
		Destination: _Service_AccessRule_Destination,
		Ports:       _Service_AccessRule_Ports,
		Protocol:    _Service_AccessRule_Protocol,
		RuleName:    _Service_AccessRule_RuleName,
		Source:      _Service_AccessRule_Source,
		Status:      _Service_AccessRule_Status,
	}

	err := siClient.CreateAccessRule(_ServiceInstanceName, createAccessRuleInfo)
	if err != nil {
		t.Fatal(err)
	}

	accessRuleList, err := siClient.GetAccessRules(_ServiceInstanceName)
	if err != nil {
		t.Fatal(err)
	}

	var accessRuleMatch = false
	for _, accessRule := range accessRuleList.AccessRules {
		if accessRule.RuleName == _Service_AccessRule_RuleName {
			accessRuleMatch = true
		}
	}

	t.Log("... Checking Access Rules")
	if !accessRuleMatch {
		t.Fatal(fmt.Errorf("Could not find access rule. Wanted: %s", _Service_AccessRule_RuleName))
	}
}
*/
func getServiceInstanceTestClients() (*ServiceInstanceClient, error) {
	client, err := GetMySQLTestClient(&opc.Config{})
	if err != nil {
		return &ServiceInstanceClient{}, err
	}

	return client.ServiceInstanceClient(), nil
}

func destroyServiceInstance(t *testing.T, client *ServiceInstanceClient, name string) {
	t.Log("... Destroying MySQL Instance")
	if err := client.DeleteServiceInstance(name); err != nil {
		t.Fatal(err)
	}
}
