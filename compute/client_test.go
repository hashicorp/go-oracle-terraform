package compute

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/kylelemons/godebug/pretty"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestClient_qualifyList(t *testing.T) {
	client, server, err := getBlankTestClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	input := []string{
		"foo",
		"bar",
		"baz",
	}

	baseStr := fmt.Sprintf("/Compute-%s/%s", _ClientTestDomain, _ClientTestUser)

	expected := []string{
		fmt.Sprintf("%s/%s", baseStr, "foo"),
		fmt.Sprintf("%s/%s", baseStr, "bar"),
		fmt.Sprintf("%s/%s", baseStr, "baz"),
	}

	result := client.getQualifiedList(input)

	if diff := pretty.Compare(result, expected); diff != "" {
		t.Fatalf("Qualified List Diff: (-got +want)\n%s", diff)
	}
}

func TestClient_retryHTTP(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	endpoint, err := url.Parse("http://foo.bar")
	if err != nil {
		t.Fatal(err)
	}

	client := Client{}
	client.maxRetries = opc.Int(5)
	// Can't use a custom transport, otherwise httpmock won't catch request
	client.httpClient = http.DefaultClient
	client.apiEndpoint = endpoint
	client.logger = opc.NewDefaultLogger()
	client.loglevel = opc.LogLevel()

	httpmock.RegisterResponder("GET", "http://foo.bar/",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(404, "mocked error message"), nil
		},
	)

	req, err := http.NewRequest("GET", "http://foo.bar/", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, reqErr := client.retryRequest(req)
	if reqErr == nil {
		t.Fatalf("Expected error, got none")
	}

	if httpmock.GetTotalCallCount() != 5 {
		t.Fatalf("Expected 5 retries, got: %d", httpmock.GetTotalCallCount())
	}
}
