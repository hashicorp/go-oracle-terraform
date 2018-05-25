package compute

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

const (
	_OrchestrationTestName             = "test-acc-orchestration-6"
	_OrchestrationTestLabel            = "test-acc-orchestration-lbl"
	_OrchestrationInstanceTestLabel    = "test"
	_OrchestrationInstanceTestShape    = "oc3"
	_OrchestrationInstanceTestImage    = "/oracle/public/OL_7.2_UEKR4_x86_64"
	_OrchestrationInstanceTestBadImage = "/oracle/public/OL_7.2_UEKR4_x86_64_bad"
)

func TestAccOrchestrationLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	orcClient, err := getOrchestrationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	instanceInput := &CreateInstanceInput{
		Name:      _OrchestrationTestName,
		Label:     _OrchestrationInstanceTestLabel,
		Shape:     _OrchestrationInstanceTestShape,
		ImageList: _OrchestrationInstanceTestImage,
	}

	object := Object{
		Label:         _OrchestrationTestLabel,
		Orchestration: _OrchestrationTestName,
		Template:      instanceInput,
		Type:          OrchestrationTypeInstance,
	}

	input := &CreateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateInactive,
		Objects:      []Object{object},
	}

	createdOrchestration, err := orcClient.CreateOrchestration(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteOrchestration(t, orcClient, createdOrchestration.Name)
	log.Printf("Created NIC Set: %#v", createdOrchestration)

	getInput := &GetOrchestrationInput{
		Name: createdOrchestration.Name,
	}

	orchestration, err := orcClient.GetOrchestration(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("orchestration Retrieved: %+v", orchestration)
	assert.NotEmpty(t, orchestration.Name, "Expected orchestration name not to be empty")
	assert.Equal(t, createdOrchestration.Name, orchestration.Name,
		"Expected orchestration names to match.")
}

func TestAccOrchestrationLifeCycle_active(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	orcClient, err := getOrchestrationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	instanceInput := &CreateInstanceInput{
		Name:      _OrchestrationTestName,
		Label:     _OrchestrationInstanceTestLabel,
		Shape:     _OrchestrationInstanceTestShape,
		ImageList: _OrchestrationInstanceTestImage,
	}

	object := Object{
		Label:         _OrchestrationTestLabel,
		Orchestration: _OrchestrationTestName,
		Template:      instanceInput,
		Type:          OrchestrationTypeInstance,
	}

	input := &CreateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateActive,
		Objects:      []Object{object},
	}

	createdOrchestration, err := orcClient.CreateOrchestration(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteOrchestration(t, orcClient, createdOrchestration.Name)
	log.Printf("Created NIC Set: %#v", createdOrchestration)

	getInput := &GetOrchestrationInput{
		Name: createdOrchestration.Name,
	}

	orchestration, err := orcClient.GetOrchestration(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("orchestration Retrieved: %+v", orchestration)
	assert.NotEmpty(t, orchestration.Name, "Expected orchestration name not to be empty")
	assert.Equal(t, createdOrchestration.Name, orchestration.Name,
		"Expected orchestration names to match.")
}

func TestAccOrchestrationLifeCycle_update(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	orcClient, err := getOrchestrationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	instanceInput := &CreateInstanceInput{
		Name:      _OrchestrationTestName,
		Label:     _OrchestrationInstanceTestLabel,
		Shape:     _OrchestrationInstanceTestShape,
		ImageList: _OrchestrationInstanceTestImage,
	}

	object := Object{
		Label:         _OrchestrationTestLabel,
		Orchestration: _OrchestrationTestName,
		Template:      instanceInput,
		Type:          OrchestrationTypeInstance,
	}

	input := &CreateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateActive,
		Objects:      []Object{object},
	}

	createdOrchestration, err := orcClient.CreateOrchestration(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteOrchestration(t, orcClient, createdOrchestration.Name)
	log.Printf("Created NIC Set: %#v", createdOrchestration)

	getInput := &GetOrchestrationInput{
		Name: createdOrchestration.Name,
	}

	orchestration, err := orcClient.GetOrchestration(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("orchestration Retrieved: %+v", orchestration)
	assert.Equal(t, OrchestrationDesiredStateActive, orchestration.DesiredState, "orchestration state mismatch!")

	updateOrchestrationInput := &UpdateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateInactive,
		Objects:      orchestration.Objects,
		Version:      orchestration.Version,
	}

	_, err = orcClient.UpdateOrchestration(updateOrchestrationInput)
	if err != nil {
		t.Fatal(err)
	}

	orchestration, err = orcClient.GetOrchestration(getInput)
	if err != nil {
		t.Fatal(err)
	}
	// Don't need to tear down the orchestration, it's attached to the instance
	log.Printf("orchestration Retrieved: %+v", orchestration)
	assert.Equal(t, OrchestrationDesiredStateInactive, orchestration.DesiredState, "orchestration state mismatch!")
}

func TestAccOrchestrationLifeCycle_suspend(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	orcClient, err := getOrchestrationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	instanceInput := &CreateInstanceInput{
		Name:      _OrchestrationTestName,
		Label:     _OrchestrationInstanceTestLabel,
		Shape:     _OrchestrationInstanceTestShape,
		ImageList: _OrchestrationInstanceTestImage,
	}

	object := Object{
		Label:         _OrchestrationTestLabel,
		Orchestration: _OrchestrationTestName,
		Template:      instanceInput,
		Type:          OrchestrationTypeInstance,
	}

	input := &CreateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateActive,
		Objects:      []Object{object},
	}

	createdOrchestration, err := orcClient.CreateOrchestration(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteOrchestration(t, orcClient, createdOrchestration.Name)
	log.Printf("Created NIC Set: %#v", createdOrchestration)

	getInput := &GetOrchestrationInput{
		Name: createdOrchestration.Name,
	}

	orchestration, err := orcClient.GetOrchestration(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("orchestration Retrieved: %+v", orchestration)
	assert.Equal(t, OrchestrationDesiredStateActive, orchestration.DesiredState, "orchestration state mismatch!")

	updateOrchestrationInput := &UpdateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateSuspend,
		Objects:      orchestration.Objects,
		Version:      orchestration.Version,
	}

	_, err = orcClient.UpdateOrchestration(updateOrchestrationInput)
	if err != nil {
		t.Fatal(err)
	}

	orchestration, err = orcClient.GetOrchestration(getInput)
	if err != nil {
		t.Fatal(err)
	}
	// Don't need to tear down the orchestration, it's attached to the instance
	log.Printf("orchestration Retrieved: %+v", orchestration)
	assert.Equal(t, OrchestrationDesiredStateSuspend, orchestration.DesiredState, "orchestration state mismatch!")
}

func TestAccOrchestrationLifeCycle_badInstance(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	orcClient, err := getOrchestrationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	instanceInput := &CreateInstanceInput{
		Name:      _OrchestrationTestName,
		Label:     _OrchestrationInstanceTestLabel,
		Shape:     _OrchestrationInstanceTestShape,
		ImageList: _OrchestrationInstanceTestBadImage,
	}

	object := Object{
		Label:         _OrchestrationTestLabel,
		Orchestration: _OrchestrationTestName,
		Template:      instanceInput,
		Type:          OrchestrationTypeInstance,
	}

	input := &CreateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateActive,
		Objects:      []Object{object},
	}

	createdOrchestration, err := orcClient.CreateOrchestration(input)
	assert.Error(t, err, fmt.Sprintf("Orchestration succeded when attempting to create a bad Orchestration: %+v", createdOrchestration))
}

func TestAccOrchestrationLifeCycle_relationship(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	orcClient, err := getOrchestrationsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	instanceInput1 := &CreateInstanceInput{
		Name:      fmt.Sprintf("%s1", _OrchestrationTestName),
		Label:     fmt.Sprintf("%s1", _OrchestrationInstanceTestLabel),
		Shape:     _OrchestrationInstanceTestShape,
		ImageList: _OrchestrationInstanceTestImage,
	}

	instanceInput2 := &CreateInstanceInput{
		Name:      fmt.Sprintf("%s2", _OrchestrationTestName),
		Label:     fmt.Sprintf("%s2", _OrchestrationInstanceTestLabel),
		Shape:     _OrchestrationInstanceTestShape,
		ImageList: _OrchestrationInstanceTestImage,
	}

	object1 := Object{
		Label:         fmt.Sprintf("%s1", _OrchestrationTestLabel),
		Orchestration: _OrchestrationTestName,
		Template:      instanceInput1,
		Type:          OrchestrationTypeInstance,
	}

	relationship := Relationship{
		Type:    OrchestrationRelationshipTypeDepends,
		Targets: []string{object1.Label},
	}

	object2 := Object{
		Label:         fmt.Sprintf("%s2", _OrchestrationTestLabel),
		Orchestration: _OrchestrationTestName,
		Template:      instanceInput2,
		Type:          OrchestrationTypeInstance,
		Relationships: []Relationship{relationship},
	}

	input := &CreateOrchestrationInput{
		Name:         _OrchestrationTestName,
		DesiredState: OrchestrationDesiredStateActive,
		Objects:      []Object{object1, object2},
	}

	createdOrchestration, err := orcClient.CreateOrchestration(input)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteOrchestration(t, orcClient, createdOrchestration.Name)
	log.Printf("Created NIC Set: %#v", createdOrchestration)

	getInput := &GetOrchestrationInput{
		Name: createdOrchestration.Name,
	}

	orchestration, err := orcClient.GetOrchestration(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("orchestration Retrieved: %+v", orchestration)
	assert.NotEmpty(t, orchestration.Name, "Expected orchestration name not to be empty")
	assert.Equal(t, createdOrchestration.Name, orchestration.Name,
		"Expected orchestration names to match.")
	assert.NotNil(t, orchestration.Objects[1].Relationships,
		"Relationship between instances not setup properly")
}

func getOrchestrationsTestClients() (*OrchestrationsClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}
	return client.Orchestrations(), nil
}

func deleteOrchestration(t *testing.T, orcClient *OrchestrationsClient, name string) {
	input := &DeleteOrchestrationInput{
		Name: name,
	}
	if err := orcClient.DeleteOrchestration(input); err != nil {
		t.Fatalf("Error deleting orchestration: %v", err)
	}
}
