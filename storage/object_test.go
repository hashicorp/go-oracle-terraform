package storage

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
)

const (
	_TestObjectName          = "testing-acc-object"
	_TestFileFixturesPath    = "test-fixtures"
	_TestSourceContentLength = 2351
	_TestFileContentLength   = 2350
	_TestContentType         = "text/plain;charset=UTF-8"
	_TestAcceptRanges        = "bytes"
)

var _TestObjectMetadata = map[string]string{"Foo": "bar", "Abc-Def": "XYZ"}

func TestAccObjectLifeCycle_contentSource(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatalf("Error creating storage client: %s", err)
	}
	client := sClient.Objects()

	container, err := sClient.getTestContainer()
	if err != nil {
		t.Fatalf("Error creating test container: %s", err)
	}
	defer deleteContainer(t, sClient, container.Name)
	log.Printf("[DEBUG] Container created: %+v", container)

	// Create body seeker
	body := bytes.NewReader([]byte(_SourceInput))
	input := &CreateObjectInput{
		Name:           _TestObjectName,
		Container:      container.Name,
		ContentType:    _TestContentType,
		ObjectMetadata: _TestObjectMetadata,
		Body:           body,
	}

	object, err := client.CreateObject(input)
	if err != nil {
		t.Fatalf("Error creating object: %s", err)
	}
	defer deleteObject(t, client, object)
	log.Printf("[DEBUG] Created Object: %+v", object)

	// Assert desired, with quantifiable fields
	expected := &ObjectInfo{
		Name:               _TestObjectName,
		AcceptRanges:       _TestAcceptRanges,
		Container:          _ContainerName,
		ContentDisposition: "",
		ContentEncoding:    "",
		ContentLength:      _TestSourceContentLength,
		ContentType:        _TestContentType,
		DeleteAt:           0,
		ID:                 fmt.Sprintf("%s/%s", _ContainerName, _TestObjectName),
		ObjectManifest:     "",
		ObjectMetadata:     _TestObjectMetadata,
	}

	if err := testAssertions(object, expected); err != nil {
		t.Fatal(err)
	}
}

func TestAccObjectLifeCycle_fileSource(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatalf("Error creating storage client: %s", err)
	}
	client := sClient.Objects()

	container, err := sClient.getTestContainer()
	if err != nil {
		t.Fatalf("Error creating test container: %s", err)
	}
	defer deleteContainer(t, sClient, container.Name)
	log.Printf("[DEBUG] Container created: %+v", container)

	// Create body seeker
	body, err := os.Open(_TestFileFixturesPath + "/input.txt")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	input := &CreateObjectInput{
		Name:           _TestObjectName,
		Container:      container.Name,
		ContentType:    _TestContentType,
		ObjectMetadata: _TestObjectMetadata,
		Body:           body,
	}

	object, err := client.CreateObject(input)
	if err != nil {
		t.Fatalf("Error creating object: %s", err)
	}
	defer deleteObject(t, client, object)
	log.Printf("[DEBUG] Created Object: %+v", object)

	// Assert desired, with quantifiable fields
	expected := &ObjectInfo{
		Name:               _TestObjectName,
		AcceptRanges:       _TestAcceptRanges,
		Container:          _ContainerName,
		ContentDisposition: "",
		ContentEncoding:    "",
		ContentLength:      _TestFileContentLength,
		ContentType:        _TestContentType,
		DeleteAt:           0,
		ID:                 fmt.Sprintf("%s/%s", _ContainerName, _TestObjectName),
		ObjectManifest:     "",
		ObjectMetadata:     _TestObjectMetadata,
	}

	if err := testAssertions(object, expected); err != nil {
		t.Fatal(err)
	}
}

func TestAccObjectLifeCycle_contentSourceID(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	sClient, err := getStorageTestClient(&opc.Config{})
	if err != nil {
		t.Fatalf("Error creating storage client: %s", err)
	}
	client := sClient.Objects()

	container, err := sClient.getTestContainer()
	if err != nil {
		t.Fatalf("Error creating test container: %s", err)
	}
	defer deleteContainer(t, sClient, container.Name)
	log.Printf("[DEBUG] Container created: %+v", container)

	// Create body seeker
	body := bytes.NewReader([]byte(_SourceInput))
	input := &CreateObjectInput{
		Name:           _TestObjectName,
		Container:      container.Name,
		ContentType:    _TestContentType,
		ObjectMetadata: _TestObjectMetadata,
		Body:           body,
	}

	object, err := client.CreateObject(input)
	if err != nil {
		t.Fatalf("Error creating object: %s", err)
	}
	defer deleteObject(t, client, object)
	log.Printf("[DEBUG] Created Object: %+v", object)

	// Assert desired, with quantifiable fields
	expected := &ObjectInfo{
		Name:               _TestObjectName,
		AcceptRanges:       _TestAcceptRanges,
		Container:          _ContainerName,
		ContentDisposition: "",
		ContentEncoding:    "",
		ContentLength:      _TestSourceContentLength,
		ContentType:        _TestContentType,
		DeleteAt:           0,
		ID:                 fmt.Sprintf("%s/%s", _ContainerName, _TestObjectName),
		ObjectManifest:     "",
		ObjectMetadata:     _TestObjectMetadata,
	}

	if err = testAssertions(object, expected); err != nil {
		t.Fatal(err)
	}

	getInput := &GetObjectInput{
		ID: object.ID,
	}

	result, err := client.GetObject(getInput)
	if err != nil {
		t.Fatalf("Error Reading Object: %s", err)
	}

	if err := testAssertions(result, expected); err != nil {
		t.Fatal(err)
	}
}

// Get a container for testing objects with
func (c *Client) getTestContainer() (*Container, error) {
	input := &CreateContainerInput{
		Name:         _ContainerName,
		PrimaryKey:   _ContainerPrimaryKey,
		SecondaryKey: _ContainerSecondaryKey,
		MaxAge:       _ContainerMaxAge,
	}

	return c.CreateContainer(input)
}

func testAssertions(result, expected *ObjectInfo) error {
	// Check transient fields first, then clear from result
	if result.Date == "" {
		return fmt.Errorf("Date Expected, got nil")
	}
	result.Date = ""

	if result.Timestamp == "" {
		return fmt.Errorf("Timestamp Expected, got nil")
	}
	result.Timestamp = ""

	if result.Etag == "" {
		return fmt.Errorf("ETag expected, got nil")
	}
	result.Etag = ""

	if result.LastModified == "" {
		return fmt.Errorf("Last modified expected, got nil")
	}
	result.LastModified = ""

	if result.TransactionID == "" {
		return fmt.Errorf("Transaction ID expected, got nil")
	}
	result.TransactionID = ""

	if diff := pretty.Compare(result, expected); diff != "" {
		return fmt.Errorf("Result Diff (-got +want)\n%s", diff)
	}

	return nil
}

func deleteObject(
	t *testing.T,
	client *ObjectClient,
	object *ObjectInfo) {
	input := &DeleteObjectInput{
		Name:      object.Name,
		Container: object.Container,
	}
	if err := client.DeleteObject(input); err != nil {
		t.Fatalf("Error deleting storage object: %s", err)
	}
}

const _SourceInput = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi auctor nisi id sem gravida, quis sollicitudin dolor
maximus. Sed est lectus, mollis sit amet neque eu, pulvinar aliquet turpis. Aenean in euismod erat. Proin pulvinar
ex vel lorem malesuada, sed tincidunt urna posuere. Sed fringilla, elit et faucibus maximus, dui orci blandit lectus,
ullamcorper fringilla felis nisl at nisl. Ut leo elit, semper non dui sit amet, sagittis commodo nulla. Nulla pulvinar
purus a nunc pellentesque scelerisque at id elit. Etiam quis bibendum eros. Etiam erat elit, feugiat non ante tempus,
mattis consectetur purus. Cras nunc nibh, fringilla in imperdiet a, tempus porta nisl. Curabitur nec justo nec leo
luctus scelerisque quis sit amet risus. Curabitur finibus fringilla lacus eu vestibulum. Nunc pellentesque aliquam
semper. Proin nec ligula urna. Donec lobortis aliquam nunc vitae feugiat. Integer blandit risus in gravida facilisis.
Pellentesque vitae lectus sed est pretium finibus. Morbi sed lacus purus. Duis nec condimentum urna. Donec vel velit
purus. Ut a velit risus. Vivamus ac euismod magna, eget convallis quam. Sed tincidunt, nisl nec rhoncus facilisis,
orci mauris commodo leo, ut eleifend nisi nisi sit amet mauris. Ut lacinia viverra rhoncus. Phasellus lacinia eleifend
turpis eu rutrum. Donec sed gravida eros, eget molestie ipsum.In hac habitasse platea dictumst. Duis a libero ante.
Quisque euismod placerat risus sit amet maximus. Praesent malesuada velit nec dui tincidunt rutrum. Proin commodo ex
non consectetur cursus. Pellentesque egestas pharetra mauris, et condimentum nibh rhoncus nec. Morbi hendrerit vel
ligula vel varius. Vestibulum in faucibus metus, eget euismod justo. Cras dolor sem, dictum eget scelerisque at,
scelerisque eu enim. Aliquam vulputate rutrum orci, vitae convallis mauris sollicitudin ut.Quisque eu accumsan massa.
Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Ut nisl magna, vulputate eget
eleifend id, tincidunt eget dolor. Nam pulvinar, dui non pellentesque dignissim, nulla neque iaculis sapien, id commodo
nisi nunc vel turpis. Vivamus eget dapibus lacus. Mauris convallis mi sit amet faucibus placerat. Mauris gravida neque
tortor, vel placerat sem elementum venenatis. Integer eu placerat est. Sed sem massa, volutpat eget augue eget, aliquam
semper sem.`
