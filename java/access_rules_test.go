package java

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/database"
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

	sClient, dbClient, aClient, err := getAccessRulesTestClients()
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

	_, err = dbClient.CreateServiceInstance(createDatabaseServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDatabaseServiceInstance(t, dbClient, _ServiceInstanceDatabaseName)

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

	_, err = sClient.CreateServiceInstance(createServiceInstance)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyServiceInstance(t, sClient, _ServiceInstanceName)

	rInt := rand.Int()
	testAccessRuleName := fmt.Sprintf("test-acc-rule-%d", rInt)
	// Create an Access Rule that's disabled
	input := &CreateAccessRuleInput{
		ServiceInstanceID: _ServiceInstanceName,
		Description:       _TestAccessRuleDescription,
		Destination:       AccessRuleDestinationWLSAdmin,
		Protocol:          AccessRuleProtocolTCP,
		Ports:             _TestAccessRulePorts,
		Name:              testAccessRuleName,
		Source:            _TestAccessRuleSource,
		Status:            AccessRuleDisabled,
	}

	expected := &AccessRuleInfo{
		Description: _TestAccessRuleDescription,
		Destination: AccessRuleDestinationWLSAdmin,
		Ports:       _TestAccessRulePorts,
		Protocol:    AccessRuleProtocolTCP,
		Name:        testAccessRuleName,
		Source:      _TestAccessRuleSource,
		Status:      AccessRuleDisabled,
		RuleType:    AccessRuleTypeUser,
	}

	// Create Access Rule
	if _, err = aClient.CreateAccessRule(input); err != nil {
		t.Fatalf("Error creating AccessRule: %s", err)
	}
	defer destroyAccessRule(t, aClient, _ServiceInstanceName, testAccessRuleName)

	// Get Access Rule (Create only returns AccessRule name)
	getInput := &GetAccessRuleInput{
		ServiceInstanceID: _ServiceInstanceName,
		Name:              testAccessRuleName,
	}

	// Read Result
	result, err := aClient.GetAccessRule(getInput)
	if err != nil {
		t.Fatalf("Error reading AccessRule: %s", err)
	}

	// Test Assertions
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Fatalf("Diff creating AccessRule: (-got, +want):\n%s", diff)
	}

	// Update Access Rule
	updateInput := &UpdateAccessRuleInput{
		ServiceInstanceID: _ServiceInstanceName,
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

func getAccessRulesTestClients() (*ServiceInstanceClient, *database.ServiceInstanceClient, *UtilityClient, error) {
	dbClient, err := database.GetDatabaseTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}
	client, err := getJavaTestClient(&opc.Config{})
	if err != nil {
		return nil, nil, nil, err
	}

	return client.ServiceInstanceClient(), dbClient.ServiceInstanceClient(), client.AccessRules(), nil
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
