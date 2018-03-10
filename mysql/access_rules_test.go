package mysql

import (
	"fmt"
	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
	"testing"
	"time"
)

var _TestAccessRuleName = fmt.Sprintf("test-acc-rule-%d", helper.RInt())

const (
	_Service_AccessRule_Description = "test-mysql-accessrule"
	_Service_AccessRule_Destination = "mysql_MASTER"
	_Service_AccessRule_Ports       = "7000-8000"
	_Service_AccessRule_Protocol    = "tcp"
	_Service_AccessRule_RuleName    = "test-accessrule-1"
	_Service_AccessRule_Source      = "0.0.0.0/24"
	//_Service_AccessRule_Source = "PUBLIC_INTERNET"
	_Service_AccessRule_Status = "enabled"
)

func TestAccAccessRuleLifeCycle(t *testing.T) {

	helper.Test(t, helper.TestCase{})

	sClient, aClient, err := getAccessRulesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	sInstance, err := sClient.createTestServiceInstance()
	if err != nil {
		t.Fatalf("Error creating Service Instance: %s", err)
	}

	instanceName := sInstance.ServiceName
	defer destroyServiceInstance(t, sClient, instanceName)

	//_, aClient, err := getAccessRulesTestClients()
	//instanceName := "sampleInstance"
	//aClient.ServiceInstanceID = instanceName

	createAccessRuleInput := &CreateAccessRuleInput{
		ServiceInstanceID: instanceName,
		Description:       _Service_AccessRule_Description,
		Destination:       _Service_AccessRule_Destination,
		Ports:             _Service_AccessRule_Ports,
		Protocol:          _Service_AccessRule_Protocol,
		RuleName:          _TestAccessRuleName,
		Source:            _Service_AccessRule_Source,
		Status:            _Service_AccessRule_Status,
	}

	expected := &AccessRuleInfo{
		Description: _Service_AccessRule_Description,
		Destination: _Service_AccessRule_Destination,
		Ports:       _Service_AccessRule_Ports,
		Protocol:    _Service_AccessRule_Protocol,
		RuleName:    _TestAccessRuleName,
		Source:      _Service_AccessRule_Source,
		Status:      _Service_AccessRule_Status,
		RuleType:    "USER",
	}

	if err := aClient.CreateAccessRule(createAccessRuleInput); err != nil {
		t.Fatalf("Error creating AccessRule: %s", err)
	}

	defer destroyAccessRule(t, aClient, instanceName, _TestAccessRuleName)

	// Get Access Rule (Create only returns AccessRule name)
	getInput := &GetAccessRuleInput{
		ServiceInstanceID: instanceName,
	}

	allRulesResult, err := aClient.GetAllAccessRules(getInput)
	if err != nil {
		t.Fatalf("Error reading ALL AccessRules : %s", err)
	}

	if len(allRulesResult.AccessRules) == 0 {
		t.Fatalf("Error reading ALL accessRules: Expected at least 1 rule. Got %d", len(allRulesResult.AccessRules))
	}

	// Read Result
	getInput.Name = _TestAccessRuleName
	ruleResult, err := aClient.GetAccessRule(getInput)

	if err != nil {
		t.Fatalf("Error reading AccessRule: %s", err)
	}

	// Test Assertions
	if diff := pretty.Compare(ruleResult, expected); diff != "" {
		t.Fatalf("Diff creating AccessRule: (-got, +want):\n%s", diff)
	}

	// Update Access Rule
	updateInput := &UpdateAccessRuleInput{
		ServiceInstanceID: instanceName,
		Name:              _TestAccessRuleName,
		Status:            "disabled",
	}

	if _, err := aClient.UpdateAccessRule(updateInput); err != nil {
		t.Fatalf("Error updating AccessRule: %s", err)
	}

	// Re-Read Result
	ruleResult, err = aClient.GetAccessRule(getInput)
	if err != nil {
		t.Fatalf("Error reading AccessRule: %s", err)
	}

	// Change expected to match
	expected.Status = "disabled"

	// Test Assertions
	if diff := pretty.Compare(ruleResult, expected); diff != "" {
		t.Fatalf("Diff creating AccessRule: (-got, +want):\n%s", diff)
	}

	// We don't explicitly test the delete, because its called
	// using the defer.

	// Sleep for 2 minutes to prevent failure because the resource is still
	// locked.
	time.Sleep(2 * time.Minute)
}

func getAccessRulesTestClients() (*ServiceInstanceClient, *AccessRulesClient, error) {
	client, err := GetMySQLTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}
	return client.ServiceInstanceClient(), client.AccessRulesClient(), nil
}

func (c *ServiceInstanceClient) createTestServiceInstance() (*ServiceInstance, error) {

	serviceParameter := ServiceParameters{
		BackupDestination:  _ServiceInstanceBackupDestination,
		ServiceDescription: _ServiceInstanceDesc,
		ServiceName:        fmt.Sprintf("test-acc-instance-%d", helper.RInt()),
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

	return c.CreateServiceInstance(createServiceInstance)
}

func destroyAccessRule(t *testing.T, client *AccessRulesClient, serviceInstance, name string) {
	input := &DeleteAccessRuleInput{
		Name:              name,
		ServiceInstanceID: serviceInstance,
	}
	if err := client.DeleteAccessRule(input); err != nil {
		t.Fatalf("Error deleting Access Rule: %s", err)
	}
}
