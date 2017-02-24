package compute

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccSecurityRuleLifeCycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test-acc-sec-rule"

	securityListClient, err := getSecurityListsClient()
	if err != nil {
		t.Fatal(err)
	}
	createSecurityListInput := CreateSecurityListInput{
		Name:               name,
		OutboundCIDRPolicy: "DENY",
		Policy:             "PERMIT",
	}
	createdSecurityList, err := securityListClient.CreateSecurityList(&createSecurityListInput)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteSecurityList(t, securityListClient, name)

	securityIPListClient, err := getSecurityIPListsClient()
	if err != nil {
		t.Fatal(err)
	}
	createSecurityIPListInput := CreateSecurityIPListInput{
		Name:         name,
		SecIPEntries: []string{"127.0.0.1", "127.0.0.2"},
	}
	createdSecurityIPList, err := securityIPListClient.CreateSecurityIPList(&createSecurityIPListInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security IP List: %+v", createdSecurityIPList)

	securityApplicationsClient, err := getSecurityApplicationsClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Applications Client")

	createInput := CreateSecurityApplicationInput{
		Name:     name,
		Protocol: SecurityApplicationProtocol(ICMP),
		ICMPType: SecurityApplicationICMPType(Echo),
	}
	defer deleteSecurityApplication(t, securityApplicationsClient, name)

	createdSecurityApplication, err := securityApplicationsClient.CreateSecurityApplication(&createInput)
	if err != nil {
		t.Fatal(err)
	}

	securityRuleClient, err := getSecurityRulesClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Security Rule Client")

	createSecurityRuleInput := CreateSecurityRuleInput{
		Name:            name,
		Action:          "PERMIT",
		Disabled:        false,
		DestinationList: "seclist:" + createdSecurityList.Name,
		SourceList:      "seciplist:" + createdSecurityIPList.Name,
		Application:     createdSecurityApplication.Name,
	}
	defer deleteSecurityRule(t, securityRuleClient, name)

	createdSecurityRule, err := securityRuleClient.CreateSecurityRule(&createSecurityRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Security Rule: %+v", createdSecurityRule)

	getSecurityRuleInput := GetSecurityRuleInput{
		Name: createdSecurityRule.Name,
	}
	getSecurityRuleOutput, err := securityRuleClient.GetSecurityRule(&getSecurityRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	if createdSecurityRule.Action != getSecurityRuleOutput.Action {
		t.Fatalf("Created and retrived actions don't match.\n Desired: %s\n Actual: %s", createdSecurityRule.Action, getSecurityRuleOutput.Action)
	}
	log.Printf("Successfully retrieved Security Rule")

	updateSecurityRuleInput := UpdateSecurityRuleInput{
		Name:            name,
		Action:          "PERMIT",
		Disabled:        true,
		DestinationList: "seclist:" + createdSecurityList.Name,
		SourceList:      "seciplist:" + createdSecurityIPList.Name,
		Application:     createdSecurityApplication.Name,
	}
	updateSecurityRuleOutput, err := securityRuleClient.UpdateSecurityRule(&updateSecurityRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	if updateSecurityRuleOutput.Disabled != true {
		t.Fatal("Security Rule was not updated to disabled")
	}
	log.Printf("Successfully updated Security Rule")
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

		ruleSpec := &CreateSecurityRuleInput{}
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

	createInput := CreateSecurityRuleInput{
		Name:            "test-rule1",
		Action:          "PERMIT",
		Disabled:        false,
		DestinationList: "seclist:test-list2",
		SourceList:      "seciplist:test-list1",
		Application:     "/oracle/default-application",
	}
	info, err := client.CreateSecurityRule(&createInput)

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

func deleteSecurityRule(t *testing.T, client *SecurityRulesClient, name string) {
	deleteInput := DeleteSecurityRuleInput{
		Name: name,
	}
	err := client.DeleteSecurityRule(&deleteInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted Security Rule")
}
