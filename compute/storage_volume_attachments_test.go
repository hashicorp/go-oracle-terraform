package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"net/url"

	"strings"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

func TestAccStorageAttachmentsClient_WaitForStorageAttachmentToBeCreatedOnline(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test"
	server := serverThatAttachesStorageVolumeAfterThreeSeconds(t, name)

	defer server.Close()
	sv, err := getStubStorageAttachmentsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	info, err := sv.waitForStorageAttachmentToBeCreated(name, 10)
	if err != nil {
		t.Fatalf("Wait for storage attachment to become available request failed: %s", err)
	}

	if strings.ToLower(info.State) != "attached" {
		fmt.Println(info)
		t.Fatalf("Status of retrieved storage volume attachment was %s, expected 'attached'", info.State)
	}

}

func TestAccStorageAttachmentsClient_WaitForStorageAttachmentToBeCreatedTimeout(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test"
	server := serverThatAttachesStorageVolumeAfterThreeSeconds(t, name)

	defer server.Close()
	sv, err := getStubStorageAttachmentsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	_, err = sv.waitForStorageAttachmentToBeCreated(name, 3)
	if err == nil {
		t.Fatal("Expected timeout error")
	}
}

func serverThatAttachesStorageVolumeAfterThreeSeconds(t *testing.T, name string) *httptest.Server {
	count := 0
	return newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		var status string
		if count < 3 {
			status = "Attaching"
		} else {
			status = "Attached"
		}
		count++
		svr := fmt.Sprintf(
			"{ \"name\": \"/foo/bar/%s\", \"instance_name\": \"/foo/bar/some-instance\", \"storage_volume_name\": \"/foo/bar/example\", \"state\": \"%s\" }", name, status)

		w.Write([]byte(svr))
	})
}

func getStubStorageAttachmentsClient(server *httptest.Server) (*StorageAttachmentsClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}
	return client.StorageAttachments(), nil
}
