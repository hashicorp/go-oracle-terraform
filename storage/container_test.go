package storage

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

const _ContainerName = "test-str-container"
const _ContainerPrimaryKey = "test-url-key"
const _ContainerSecondaryKey = "test-url-key2"
const _ContainerMaxAge = 50
const _ContainerQuotaBytes = 1000000000
const _ContainerQuotaCount = 1000

var _ContainerCustomMetadata = map[string]string{"Abc-Def": "xyz", "Foo": "bar"}

func TestAccContainerLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	client, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatal(err)
	}

	readACLs := []string{"test-read-acl1", "test-read-acl2"}
	writeACLs := []string{"test-write-acl1", "test-write-acl2"}
	allowedOrigins := []string{"allowed-origin-1", "allowed-origin-2"}
	exposedHeaders := []string{"exposed-header-1", "exposed-header-2"}

	createContainerInput := CreateContainerInput{
		Name:           _ContainerName,
		ReadACLs:       readACLs,
		WriteACLs:      writeACLs,
		PrimaryKey:     _ContainerPrimaryKey,
		SecondaryKey:   _ContainerSecondaryKey,
		AllowedOrigins: allowedOrigins,
		ExposedHeaders: exposedHeaders,
		MaxAge:         _ContainerMaxAge,
		QuotaBytes:     _ContainerQuotaBytes,
		QuotaCount:     _ContainerQuotaCount,
		CustomMetadata: _ContainerCustomMetadata,
	}

	createdContainer, err := client.CreateContainer(&createContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Container: %+v", createdContainer)
	defer deleteContainer(t, client, _ContainerName)

	getContainerInput := GetContainerInput{
		Name: _ContainerName,
	}
	container, err := client.GetContainer(&getContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	if diff := pretty.Compare(container.ReadACLs, readACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("ReadACL diff (-got +want)\n%s", diff))
	}
	if diff := pretty.Compare(container.WriteACLs, writeACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("WriteACL diff (-got +want)\n%s", diff))
	}
	if container.PrimaryKey != _ContainerPrimaryKey {
		t.Fatalf(fmt.Sprintf("URLKeys don't match. Wanted: %s Recieved: %s", _ContainerPrimaryKey, container.PrimaryKey))
	}
	if container.SecondaryKey != _ContainerSecondaryKey {
		t.Fatalf(fmt.Sprintf("URLKey2 do not match. Wanted: %s Recieved: %s", _ContainerSecondaryKey, container.SecondaryKey))
	}
	if diff := pretty.Compare(container.AllowedOrigins, allowedOrigins); diff != "" {
		t.Fatalf(fmt.Sprintf("AllowedOrigin diff (-got +want)\n%s", diff))
	}
	if diff := pretty.Compare(container.ExposedHeaders, exposedHeaders); diff != "" {
		t.Fatalf(fmt.Sprintf("ExposedHeader diff (-got +want)\n%s", diff))
	}
	if container.MaxAge != _ContainerMaxAge {
		t.Fatalf(fmt.Sprintf("Max Age do not match Wanted: %d Recieved: %d", _ContainerMaxAge, container.MaxAge))
	}
	if container.QuotaBytes != _ContainerQuotaBytes {
		t.Fatalf(fmt.Sprintf("Quota Bytes do not match Wanted: %d Recieved: %d", _ContainerQuotaBytes, container.QuotaBytes))
	}
	if container.QuotaCount != _ContainerQuotaCount {
		t.Fatalf(fmt.Sprintf("Quota Count do not match Wanted: %d Recieved: %d", _ContainerQuotaCount, container.QuotaCount))
	}
	if !reflect.DeepEqual(container.CustomMetadata, _ContainerCustomMetadata) {
		t.Fatalf(fmt.Sprintf("CustomMetadata do not match Wanted: %v Recieved: %v", _ContainerCustomMetadata, container.CustomMetadata))
	}

	log.Print("Successfully retrieved Container")

	updateReadACLs := []string{"test-read-acl3", "test-read-acl4"}
	updateWriteACLs := []string{"test-write-acl3", "test-write-acl4"}
	updatedAllowedOrigins := []string{"allowed-origin-3", "allowed-origin-4"}
	updatedExposedHeaders := []string{"exposed-header-3", "exposed-header-4"}
	updatedMaxAge := _ContainerMaxAge + 1
	updatedQuotaBytes := _ContainerQuotaBytes + 1
	updatedQuotaCount := _ContainerQuotaCount + 1
	updatedCustomMetaData := map[string]string{"Abc-Def": "123", "Bar": "foo"}
	updatedRemoveCustomMetaData := []string{"Foo"}
	updateContainerInput := UpdateContainerInput{
		Name:                 _ContainerName,
		ReadACLs:             updateReadACLs,
		WriteACLs:            updateWriteACLs,
		SecondaryKey:         _ContainerPrimaryKey,
		AllowedOrigins:       updatedAllowedOrigins,
		ExposedHeaders:       updatedExposedHeaders,
		MaxAge:               updatedMaxAge,
		QuotaBytes:           updatedQuotaBytes,
		QuotaCount:           updatedQuotaCount,
		CustomMetadata:       updatedCustomMetaData,
		RemoveCustomMetadata: updatedRemoveCustomMetaData,
	}

	_, err = client.UpdateContainer(&updateContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Container: %+v", createdContainer)

	container, err = client.GetContainer(&getContainerInput)
	if err != nil {
		t.Fatal(err)
	}
	if container.Name != _ContainerName {
		t.Fatalf(fmt.Sprintf("Names don't match. Wanted: %s Recieved: %s", _ContainerName, container.Name))
	}
	if diff := pretty.Compare(container.ReadACLs, updateReadACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("UpdatedReadACL diff (-got +want)\n%s", diff))
	}
	if diff := pretty.Compare(container.WriteACLs, updateWriteACLs); diff != "" {
		t.Fatalf(fmt.Sprintf("UpdatedWriteACL diff (-got +want)\n%s", diff))
	}
	if container.PrimaryKey != "" {
		t.Fatalf(fmt.Sprintf("Expected URL Key to be empty. Recieved: %s", container.PrimaryKey))
	}
	if container.SecondaryKey != _ContainerPrimaryKey {
		t.Fatalf(fmt.Sprintf("Updated URL Key 2 does not match. Wanted: %s Recieved: %s", _ContainerPrimaryKey, container.SecondaryKey))
	}
	if diff := pretty.Compare(container.AllowedOrigins, updatedAllowedOrigins); diff != "" {
		t.Fatalf(fmt.Sprintf("Updated AllowedOrigin diff (-got +want)\n%s", diff))
	}
	if diff := pretty.Compare(container.ExposedHeaders, updatedExposedHeaders); diff != "" {
		t.Fatalf(fmt.Sprintf("Updated Exposed Headers diff (-got +want)\n%s", diff))
	}
	if container.MaxAge != updatedMaxAge {
		t.Fatalf(fmt.Sprintf("Max Age do not match Wanted: %d Recieved: %d", updatedMaxAge, container.MaxAge))
	}
	if container.QuotaBytes != updatedQuotaBytes {
		t.Fatalf(fmt.Sprintf("Quota Bytes do not match Wanted: %d Recieved: %d", updatedQuotaBytes, container.QuotaBytes))
	}
	if container.QuotaCount != updatedQuotaCount {
		t.Fatalf(fmt.Sprintf("Quota Count do not match Wanted: %d Recieved: %d", updatedQuotaCount, container.QuotaCount))
	}
	if diff := pretty.Compare(container.CustomMetadata, updatedCustomMetaData); diff != "" {
		t.Fatalf(fmt.Sprintf("CustomMetadata diff (-got +want)\n%s", diff))
	}

	log.Print("Successfully retrieved Container")
}

func deleteContainer(t *testing.T, client *Client, name string) {
	deleteInput := DeleteContainerInput{
		Name: name,
	}
	if err := client.DeleteContainer(&deleteInput); err != nil {
		t.Fatal(err)
	}

	log.Print("Successfully deleted Container")
}

func Test_isCustomHeader(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	client, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if client.isCustomHeader("X-Container-Meta-Access-Control-Allow-Origin") {
		t.Fatalf("X-Container-Meta-Access-Control-Allow-Origin shoud be a identified as a standard header")
	}
	if !client.isCustomHeader("X-Container-Meta-Some-Other-Header") {
		t.Fatalf("X-Container-Meta-Some-Other-Header shoud be a identified as a custom header")
	}

}

func Test_updateOrRemoveHeadersStringValue(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	client, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatal(err)
	}

	headers := make(map[string]string)

	client.updateOrRemoveStringValue(headers, "X-Container-Meta-RemoveString", "")
	client.updateOrRemoveStringValue(headers, "X-Container-Meta-UpdateString", "updated")

	// test remove
	if _, ok := headers["X-Remove-Container-Meta-RemoveString"]; !(ok) {
		t.Fatalf("X-Container-Meta-RemoveString was not set")
	}
	// test update
	if val, ok := headers["X-Container-Meta-UpdateString"]; !(ok) || val != "updated" {
		t.Fatalf("X-Container-Meta-UpdateString was not set to 'updated'")
	}
}

func Test_updateOrRemoveHeadersIntValue(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	client, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatal(err)
	}

	headers := make(map[string]string)

	client.updateOrRemoveIntValue(headers, "X-Container-Meta-RemoveInt", 0)
	client.updateOrRemoveIntValue(headers, "X-Container-Meta-UpdateInt", 1)

	// test remove
	if _, ok := headers["X-Remove-Container-Meta-RemoveInt"]; !(ok) {
		t.Fatalf("X-Remove-Container-Meta-RemoveInt was not set")
	}
	//test update
	if val, ok := headers["X-Container-Meta-UpdateInt"]; !(ok) || val != "1" {
		t.Fatalf("X-Container-Meta-UpdateInt was not set to 1")
	}
}
