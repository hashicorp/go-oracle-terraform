package database

import (
	"fmt"
	"math/rand"
	"testing"

	"os"

	"log"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

// Testing Parameters Used
const (
	_TestAccessRuleDescription = "testing description"
	_TestAccessRulePorts       = "7000-8000"
	_TestAccessRuleSource      = "PUBLIC-INTERNET"
)

func TestAccAccessRulesLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, aClient, err := getAccessRulesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	var (
		instanceName string
		sInstance    *ServiceInstance
	)
	if v := os.Getenv("OPC_TEST_DB_INSTANCE"); v == "" {
		// First Create a Service Instance
		sInstance, err = sClient.createTestServiceInstance()
		if err != nil {
			t.Fatalf("Error creating Service Instance: %s", err)
		}
		defer destroyServiceInstance(t, sClient, sInstance.Name)
		instanceName = sInstance.Name
	} else {
		log.Print("Using already created DB Service Instance")
		instanceName = v
	}

	rInt := rand.Int()
	testAccessRuleName := fmt.Sprintf("test-acc-rule-%d", rInt)

	// Create an Access Rule that's disabled
	input := &CreateAccessRuleInput{
		ServiceInstanceID: instanceName,
		Description:       _TestAccessRuleDescription,
		Destination:       AccessRuleDefaultDestination,
		Ports:             _TestAccessRulePorts,
		Name:              testAccessRuleName,
		Source:            _TestAccessRuleSource,
		Status:            AccessRuleDisabled,
	}

	expected := &AccessRuleInfo{
		Description: _TestAccessRuleDescription,
		Destination: AccessRuleDefaultDestination,
		Ports:       _TestAccessRulePorts,
		Name:        testAccessRuleName,
		Source:      _TestAccessRuleSource,
		Status:      AccessRuleDisabled,
		RuleType:    AccessRuleTypeUser,
	}

	// Create Access Rule
	if _, err = aClient.CreateAccessRule(input); err != nil {
		t.Fatalf("Error creating AccessRule: %s", err)
	}
	defer destroyAccessRule(t, aClient, instanceName, testAccessRuleName)

	// Get Access Rule (Create only returns AccessRule name)
	getInput := &GetAccessRuleInput{
		ServiceInstanceID: instanceName,
		Name:              testAccessRuleName,
	}

	var (
		result *AccessRuleInfo
	)
	result, err = aClient.GetAccessRule(getInput)
	if err != nil {
		t.Fatalf("Error reading AccessRule: %s", err)
	}

	// Test Assertions
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Fatalf("Diff creating AccessRule: (-got, +want):\n%s", diff)
	}

	// Update Access Rule
	updateInput := &UpdateAccessRuleInput{
		ServiceInstanceID: instanceName,
		Name:              testAccessRuleName,
		Status:            AccessRuleEnabled,
	}

	if _, err = aClient.UpdateAccessRule(updateInput); err != nil {
		t.Fatalf("Error updating AccessRule: %s", err)
	}

	// Re-Read Result
	result, err = aClient.GetAccessRule(getInput)
	if err != nil {
		t.Fatalf("Error reading AccessRule: %s", err)
	}

	// Change expected to match
	expected.Status = AccessRuleEnabled

	// Test Assertions
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Fatalf("Diff creating AccessRule: (-got, +want):\n%s", diff)
	}
}

func TestGetDefaultAccessRule_Basic(t *testing.T) {
	sClient, aClient, err := getAccessRulesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	var (
		instanceName string
		sInstance    *ServiceInstance
	)
	if v := os.Getenv("OPC_TEST_DB_INSTANCE"); v == "" {
		// First Create a Service Instance
		sInstance, err = sClient.createTestServiceInstance()
		if err != nil {
			t.Fatalf("Error creating Service Instance: %s", err)
		}
		defer destroyServiceInstance(t, sClient, sInstance.Name)
		instanceName = sInstance.Name
	} else {
		log.Print("Using already created DB Service Instance")
		instanceName = v
	}

	getInput := &GetDefaultAccessRuleInput{
		ServiceInstanceID: instanceName,
	}

	result, err := aClient.GetDefaultAccessRules(getInput)
	if err != nil {
		t.Fatalf("Error reading AccessRule: %s", err)
	}

	expected := &DefaultAccessRuleInfo{
		EnableSSH:        helper.Bool(true),
		EnableHTTP:       helper.Bool(false),
		EnableHTTPSSL:    helper.Bool(false),
		EnableDBConsole:  helper.Bool(false),
		EnableDBExpress:  helper.Bool(false),
		EnableDBListener: helper.Bool(false),
	}

	if diff := pretty.Compare(result, expected); diff != "" {
		t.Fatalf("Diff on Default Access Rules: (-got, +want):\n%s", diff)
	}
}

func TestUpdateDefaultAccessRule_Basic(t *testing.T) {
	sClient, aClient, err := getAccessRulesTestClients()
	if err != nil {
		t.Fatal(err)
	}

	var (
		instanceName string
		sInstance    *ServiceInstance
	)
	if v := os.Getenv("OPC_TEST_DB_INSTANCE"); v == "" {
		// First Create a Service Instance
		sInstance, err = sClient.createTestServiceInstance()
		if err != nil {
			t.Fatalf("Error creating Service Instance: %s", err)
		}
		defer destroyServiceInstance(t, sClient, sInstance.Name)
		instanceName = sInstance.Name
	} else {
		log.Print("Using already created DB Service Instance")
		instanceName = v
	}

	expected := &DefaultAccessRuleInfo{
		ServiceInstanceID: instanceName,
		EnableSSH:         helper.Bool(false),
		EnableHTTP:        helper.Bool(true),
		EnableHTTPSSL:     helper.Bool(true),
		EnableDBConsole:   helper.Bool(false),
		EnableDBExpress:   helper.Bool(true),
		EnableDBListener:  helper.Bool(false),
	}

	result, err := aClient.UpdateDefaultAccessRules(expected)
	if err != nil {
		t.Fatalf("Error updating AccessRule: %s", err)
	}

	if diff := pretty.Compare(result, expected); diff != "" {
		t.Fatalf("Diff on Default Access Rules: (-got, +want):\n%s", diff)
	}
}

func getAccessRulesTestClients() (*ServiceInstanceClient, *UtilityClient, error) {
	client, err := GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, err
	}
	return client.ServiceInstanceClient(), client.AccessRules(), nil
}

func destroyAccessRule(t *testing.T, client *UtilityClient, serviceInstance, name string) {
	input := &DeleteAccessRuleInput{
		Name:              name,
		ServiceInstanceID: serviceInstance,
	}
	if err := client.DeleteAccessRule(input); err != nil {
		t.Fatalf("Error deleting Access Rule: %s", err)
	}
}

func (c *ServiceInstanceClient) createTestServiceInstance() (*ServiceInstance, error) {
	parameter := ParameterInput{
		AdminPassword:     _ServiceInstancePassword,
		BackupDestination: _ServiceInstanceBackupDestination,
		SID:               _ServiceInstanceDBSID,
		Type:              _ServiceInstanceType,
		UsableStorage:     _ServiceInstanceUsableStorage,
	}

	rInt := rand.Int()
	createServiceInstance := &CreateServiceInstanceInput{
		Name:             fmt.Sprintf("test-acc-instance-%d", rInt),
		Edition:          _ServiceInstanceEdition,
		Level:            _ServiceInstanceLevel,
		Shape:            _ServiceInstanceShape,
		SubscriptionType: _ServiceInstanceSubscription,
		Version:          _ServiceInstanceVersion,
		VMPublicKey:      _ServiceInstancePubKey,
		Parameter:        parameter,
	}

	return c.CreateServiceInstance(createServiceInstance)
}
