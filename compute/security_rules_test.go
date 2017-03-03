package compute

import (
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const (
	_SecurityRuleTestName        = "test-acc-security-rule"
	_SecurityRuleTestDescription = "testing security rule"
	_SecurityRuleFlowDirection   = "ingress"
)

func TestAccSecurityRulesLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	svc, err := getSecurityRulesClient()
	if err != nil {
		t.Fatal(err)
	}

	createInput := &CreateSecurityRuleInput{
		Name:          _SecurityRuleTestName,
		Description:   _SecurityRuleTestDescription,
		FlowDirection: _SecurityRuleFlowDirection,
		Tags:          []string{"testing"},
	}

	createdRule, err := svc.CreateSecurityRule(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Rule succcessfully created")
	defer destroySecurityRule(t, svc, _SecurityRuleTestName)

	getInput := &GetSecurityRuleInput{
		Name: _SecurityRuleTestName,
	}
	receivedRule, err := svc.GetSecurityRule(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Rule successfully fetched")

	if !reflect.DeepEqual(createdRule.Tags, receivedRule.Tags) {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", createdRule, receivedRule)
	}

	// Update prefix, NAPT, and tags
	updateInput := &UpdateSecurityRuleInput{
		Name:          _SecurityRuleTestName,
		Description:   _SecurityRuleTestDescription,
		FlowDirection: _SecurityRuleFlowDirection,
		Tags:          []string{"updated"},
	}

	updatedRule, err := svc.UpdateSecurityRule(updateInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Rule successfully updated")

	receivedRule, err = svc.GetSecurityRule(getInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Rule successfully fetched")

	if !reflect.DeepEqual(updatedRule.Tags, receivedRule.Tags) {
		t.Fatalf("Mismatch found after update.\nExpected: %+v\nReceived: %+v", createdRule, receivedRule)
	}
}

func TestAccSecurityRulesWithOptionsLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	ssc, err := getVirtNICSetsClient()
	if err != nil {
		t.Fatal(err)
	}

	dstInput := &CreateVirtualNICSetInput{
		Name:        _SecurityRuleTestName + "dst_set",
		Description: _SecurityRuleTestDescription,
	}

	dstSet, err := ssc.CreateVirtualNICSet(dstInput)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVirtualNICSet(t, ssc, dstSet.Name)
	log.Printf("Created NIC Set: %#v", dstSet)

	srcInput := &CreateVirtualNICSetInput{
		Name:        _SecurityRuleTestName + "src_set",
		Description: _SecurityRuleTestDescription,
	}
	srcSet, err := ssc.CreateVirtualNICSet(srcInput)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVirtualNICSet(t, ssc, srcSet.Name)
	log.Printf("Created NIC Set: %#v", srcSet)

	spc, err := getSecurityProtocolsClient()
	if err != nil {
		t.Fatal(err)
	}

	createSecurityProtocolInput := &CreateSecurityProtocolInput{
		Name:        _SecurityProtocolTestName,
		Description: _SecurityProtocolTestDescription,
		Tags:        []string{"testing"},
		IPProtocol:  _SecurityProtocolTestIPProtocol,
		SrcPortSet:  []string{"17"},
		DstPortSet:  []string{"18"},
	}

	createdSecurityProtocol, err := spc.CreateSecurityProtocol(createSecurityProtocolInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Protocol succcessfully created")
	defer destroySecurityProtocol(t, spc, _SecurityProtocolTestName)

	svc, err := getSecurityRulesClient()
	if err != nil {
		t.Fatal(err)
	}

	sic, err := getIPAddressPrefixSetsClient()
	if err != nil {
		t.Fatal(err)
	}

	dstIPAddressPrefixSetInput := &CreateIPAddressPrefixSetInput{
		Name:              _IPAddressPrefixSetTestName + "-dst",
		Description:       _IPAddressPrefixSetTestDescription,
		IPAddressPrefixes: []string{"192.0.0.168/16"},
		Tags:              []string{"testing"},
	}

	dstIPAddressPrefixSet, err := sic.CreateIPAddressPrefixSet(dstIPAddressPrefixSetInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Dst IP Address Prefix Set succcessfully created")
	defer destroyIPAddressPrefixSet(t, sic, dstIPAddressPrefixSet.Name)

	srcIPAddressPrefixSetInput := &CreateIPAddressPrefixSetInput{
		Name:              _IPAddressPrefixSetTestName + "-src",
		Description:       _IPAddressPrefixSetTestDescription,
		IPAddressPrefixes: []string{"192.0.0.169/16"},
		Tags:              []string{"testing"},
	}

	srcIPAddressPrefixSet, err := sic.CreateIPAddressPrefixSet(srcIPAddressPrefixSetInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Src IP Address Prefix Set succcessfully created")
	defer destroyIPAddressPrefixSet(t, sic, srcIPAddressPrefixSet.Name)

	createInput := &CreateSecurityRuleInput{
		Name:                   _SecurityRuleTestName,
		Description:            _SecurityRuleTestDescription,
		DstVnicSet:             dstSet.Name,
		SrcVnicSet:             srcSet.Name,
		DstIpAddressPrefixSets: []string{dstIPAddressPrefixSet.Name},
		SrcIpAddressPrefixSets: []string{srcIPAddressPrefixSet.Name},
		SecProtocols:           []string{createdSecurityProtocol.Name},
		FlowDirection:          _SecurityRuleFlowDirection,
		Tags:                   []string{"testing"},
	}

	createdRule, err := svc.CreateSecurityRule(createInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Print("Security Rule succcessfully created")
	defer destroySecurityRule(t, svc, _SecurityRuleTestName)

	if dstSet.Name != createdRule.DstVnicSet {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", dstSet.Name, createdRule.DstVnicSet)
	}
	if srcSet.Name != createdRule.SrcVnicSet {
		t.Fatalf("Mismatch found after create.\nExpected: %+v\nReceived: %+v", dstSet.Name, createdRule.DstVnicSet)
	}
}

func destroySecurityRule(t *testing.T, svc *SecurityRuleClient, name string) {
	input := &DeleteSecurityRuleInput{
		Name: name,
	}
	if err := svc.DeleteSecurityRule(input); err != nil {
		t.Fatal(err)
	}
}

func getSecurityRulesClient() (*SecurityRuleClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return nil, err
	}

	return client.SecurityRules(), nil
}
