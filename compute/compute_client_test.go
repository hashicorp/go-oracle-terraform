package compute

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"
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
