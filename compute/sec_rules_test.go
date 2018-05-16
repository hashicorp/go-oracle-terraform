package compute

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/stretchr/testify/assert"
)

func TestAccSecRuleLifeCycle(t *testing.T) {
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
	defer deleteSecurityIPList(t, securityIPListClient, name)

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

	createdSecurityApplication, err := securityApplicationsClient.CreateSecurityApplication(&createInput)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteSecurityApplication(t, securityApplicationsClient, name)

	secRuleClient, err := getSecRulesClient()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Obtained Sec Rule Client")

	createSecRuleInput := CreateSecRuleInput{
		Name:            name,
		Action:          "PERMIT",
		Disabled:        false,
		DestinationList: "seclist:" + createdSecurityList.Name,
		SourceList:      "seciplist:" + createdSecurityIPList.Name,
		Application:     createdSecurityApplication.Name,
	}

	createdSecRule, err := secRuleClient.CreateSecRule(&createSecRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Successfully created Sec Rule: %+v", createdSecRule)
	defer deleteSecRule(t, secRuleClient, name)

	getSecRuleInput := GetSecRuleInput{
		Name: createdSecRule.Name,
	}
	getSecRuleOutput, err := secRuleClient.GetSecRule(&getSecRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, createdSecRule, getSecRuleOutput,
		"Created and retrived sec rules don't match.")
	log.Printf("Successfully retrieved Sec Rule")

	updateSecRuleInput := UpdateSecRuleInput{
		Name:            name,
		Action:          "PERMIT",
		Disabled:        true,
		DestinationList: "seclist:" + createdSecurityList.Name,
		SourceList:      "seciplist:" + createdSecurityIPList.Name,
		Application:     createdSecurityApplication.Name,
	}
	updateSecRuleOutput, err := secRuleClient.UpdateSecRule(&updateSecRuleInput)
	if err != nil {
		t.Fatal(err)
	}
	if !updateSecRuleOutput.Disabled {
		t.Fatal("Sec Rule was not updated to disabled")
	}
	log.Printf("Successfully updated Sec Rule")
}

// Test that the client can create an instance.
func TestAccSecRulesClient_CreateRule(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Wrong HTTP method %s, expected POST", r.Method)
		}

		expectedPath := "/secrule/"
		if r.URL.Path != expectedPath {
			t.Errorf("Wrong HTTP URL %v, expected %v", r.URL, expectedPath)
		}

		ruleSpec := &CreateSecRuleInput{}
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

		w.Write([]byte(exampleCreateSecRuleResponse))
		w.WriteHeader(201)
	})

	defer server.Close()
	client, err := getStubSecRulesClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	createInput := CreateSecRuleInput{
		Name:            "test-rule1",
		Action:          "PERMIT",
		Disabled:        false,
		DestinationList: "seclist:test-list2",
		SourceList:      "seciplist:test-list1",
		Application:     "/oracle/default-application",
	}
	info, err := client.CreateSecRule(&createInput)
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

func getStubSecRulesClient(server *httptest.Server) (*SecRulesClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}

	return client.SecRules(), nil
}

func getSecRulesClient() (*SecRulesClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &SecRulesClient{}, err
	}

	return client.SecRules(), nil
}

var exampleCreateSecRuleResponse = `
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

func deleteSecRule(t *testing.T, client *SecRulesClient, name string) {
	deleteInput := DeleteSecRuleInput{
		Name: name,
	}
	err := client.DeleteSecRule(&deleteInput)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Successfully deleted Sec Rule")
}
