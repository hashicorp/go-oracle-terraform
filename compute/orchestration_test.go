package compute

import (
	"log"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_OrchestrationTestName          = "test-acc-orchestration2"
	_OrchestrationTestLabel         = "test-acc-orchestration-lbl"
	_OrchestrationInstanceTestLabel = "test"
	_OrchestrationInstanceTestShape = "oc3"
	_OrchestrationInstanceTestImage = "/oracle/public/Oracle_Solaris_11.3"
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
	if orchestration.Name != createdOrchestration.Name || orchestration.Name == "" {
		t.Fatal("orchestration Name mismatch! Got: %q Expected: %q", orchestration.Name, createdOrchestration)
	}
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
	if orchestration.Name != createdOrchestration.Name || orchestration.Name == "" {
		t.Fatal("orchestration Name mismatch! Got: %q Expected: %q", orchestration.Name, createdOrchestration)
	}
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
	if orchestration.DesiredState != OrchestrationDesiredStateActive {
		t.Fatal("orchestration state mismatch! Got: %q Expected: %q", orchestration.DesiredState, OrchestrationDesiredStateActive)
	}

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
	if orchestration.DesiredState != OrchestrationDesiredStateInactive {
		t.Fatal("orchestration state mismatch! Got: %q Expected: %q", orchestration.DesiredState, OrchestrationDesiredStateInactive)
	}
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
