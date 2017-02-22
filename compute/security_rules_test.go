package compute

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
)


func TestAccSecurityRuleLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	securityListClient, err := getSecurityListsClient()
	if err != nil {
		t.Fatal(err)
	}

	createSecurityListInput := CreateSecurityListInput{
		Name:               "test-sec-list",
		OutboundCIDRPolicy: "DENY",
		Policy:             "PERMIT",
	}
	securityList, err := securityListClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security List: %+v", securityList)

	securityIPListClient, err := getSecurityIPListsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security IP List Client")

	createSecurityIPListInput := CreateSecurityIPListInput{
		Name:         "test-sec-ip-list",
		SecIPEntries: []string{"127.0.0.1", "127.0.0.2"},
	}
	securityIPList, err := securityIPListClient.CreateSecurityIPList(&createSecurityIPListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security IP List: %+v", securityIPList)


	securityRuleClient, err := getSecurityRulesClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Rule Client")

	createSecurityRuleInput := CreateSecurityRuleInput{
		Name:               "test-sec-rule",
		Action: "PERMIT",
		Disabled: false,
	}
	securityRule, err := securityRuleClient.CreateSecurityRule(&createSecurityRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Rule: %+v", securityRule)

	getSecurityRuleInput := GetSecurityRuleInput{
		Name: securityRule.Name,
	}
	getSecurityRuleOutput, err := securityRuleClient.GetSecurityRule(&getSecurityRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	if securityRule.Policy != getSecurityRuleOutput.Policy {
		t.Fatalf("Created and retrived policies don't match.\n Desired: %s\n Actual: %s", securityRule.Policy, getSecurityRuleOutput.Policy)
	}
	log.Printf("Successfully retrieved Security Rule")

	updateSecurityRuleInput := UpdateSecurityRuleInput{
		Name:               securityRule.Name,
		OutboundCIDRPolicy: "PERMIT",
		Policy:             "DENY",
	}
	updateSecurityRuleOutput, err := securityRuleClient.UpdateSecurityRule(&updateSecurityRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	if updateSecurityRuleOutput.OutboundCIDRPolicy != "PERMIT" {
		t.Fatalf("Outbound policy not successfully updated \nDesired: %s \nActual: %s", updateSecurityRuleInput.OutboundCIDRPolicy, updateSecurityRuleOutput.OutboundCIDRPolicy)
	}
	log.Printf("Successfully updated Security Rule")

	deleteSecurityRuleInput := DeleteSecurityRuleInput{
		Name: securityRule.Name,
	}
	err = securityRuleClient.DeleteSecurityRule(&deleteSecurityRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted Security Rule")

	deleteSecurityIPListInput := DeleteSecurityIPListInput{
		Name: securityIPList.Name,
	}
	err = securityIPListClient.DeleteSecurityIPList(&deleteSecurityIPListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted Security IP List")

	deleteSecurityListInput := DeleteSecurityListInput{
		Name: securityList.Name,
	}
	err = securityListClient.DeleteSecurityList(&deleteSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully deleted Security List")
}


// Test that the client can create an instance.
func TestAccSecurityRulesClient_CreateRule(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Wrong HTTP method %s, expected POST", r.Method)
		}

		expectedPath := "/secrule/"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong HTTP URL %v, expected %v", r.URL, expectedPath)
		}

		ruleSpec := &SecurityRuleSpec{}
		unmarshalRequestBody(t, r, ruleSpec)

		if ruleSpec.Name != "/Compute-test/test/test-rule1" {
			t.Errorf("Expected name '/Compute-test/test/test-rule1', was %s", ruleSpec.Name)
		}

		if ruleSpec.SourceList != "seciplist:/Compute-test/test/test-list1" {
			t.Errorf("Expected source list 'seciplist:/Compute-test/test/test-list1', was %s",
				ruleSpec.SourceList)
		}

		if ruleSpec.DestinationList != "seclist:/Compute-test/test/test-list2" {
			t.Errorf("Expected destination list 'seclist:/Compute-test/test/test-list2', was %s",
				ruleSpec.DestinationList)
		}

		if ruleSpec.Application != "/oracle/default-application" {
			t.Errorf("Expected application '/oracle/default-application', was %s", ruleSpec.Application)
		}

		w.Write([]byte(exampleCreateSecurityRuleResponse))
		w.WriteHeader(201)
	})

	defer server.Close()
	client, err := getStubSecurityRulesClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	info, err := client.CreateSecurityRule(
		"test-rule1",
		"seciplist:test-list1",
		"seclist:test-list2",
		"/oracle/default-application",
		"PERMIT",
		false)

	if err != nil {
		t.Fatalf("Create security rule request failed: %s", err)
	}

	if info.SourceList != "seciplist:es_iplist" {
		t.Errorf("Expected source list 'seciplist:es_iplist', was %s", info.SourceList)
	}

	if info.DestinationList != "seclist:allowed_video_servers" {
		t.Errorf("Expected source list 'seclist:allowed_video_servers', was %s", info.DestinationList)
	}

	if info.Application != "video_streaming_udp" {
		t.Errorf("Expected application 'video_streaming_udp', was %s", info.Application)
	}
}

func getStubSecurityRulesClient(server *httptest.Server) (*SecurityRulesClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}

	return client.SecurityRules(), nil
}

func getSecurityRulesClient() (*SecurityRulesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecurityRulesClient{}, err
	}

	return client.SecurityRules(), nil
}

var exampleCreateSecurityRuleResponse = `
{
  "dst_list": "seclist:/Compute-acme/jack.jones@example.com/allowed_video_servers",
  "name": "/Compute-acme/jack.jones@example.com/es_to_videoservers_stream",
  "src_list": "seciplist:/Compute-acme/jack.jones@example.com/es_iplist",
  "uri": "https://api.compute.us0.oraclecloud.com/secrule/Compute-acme/jack.jones@example.com/es_to_videoservers_stream",
  "disabled": false,
  "application": "/Compute-acme/jack.jones@example.com/video_streaming_udp",
  "action": "PERMIT"
}
`
