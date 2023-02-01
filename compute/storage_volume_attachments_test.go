// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

func TestAccStorageAttachmentsClient_WaitForStorageDetachmentSuccessful(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test"
	server := serverThatDetachesStorageVolumeAfterThreeSeconds(t, name)

	defer server.Close()
	sv, err := getStubStorageAttachmentsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	err = sv.waitForStorageAttachmentToBeDeleted(name, time.Duration(1*time.Second), time.Duration(10*time.Second))
	if err != nil {
		t.Fatalf("Wait for storage attachment to become detach request failed: %s", err)
	}

}

func TestAccStorageAttachmentsClient_WaitForStorageDetachmentTimeout(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test"
	server := serverThatDetachesStorageVolumeAfterThreeSeconds(t, name)

	defer server.Close()
	sv, err := getStubStorageAttachmentsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	err = sv.waitForStorageAttachmentToBeDeleted(name, time.Duration(1*time.Second), time.Duration(3*time.Second))
	if err == nil {
		t.Fatal("Expected timeout error")
	}
}

func TestAccStorageAttachmentsClient_WaitForStorageAttachmentToBeFullyAttachedSuccessful(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test"
	server := serverThatAttachesStorageVolumeAfterThreeSeconds(t, name)

	defer server.Close()
	sv, err := getStubStorageAttachmentsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	info, err := sv.waitForStorageAttachmentToFullyAttach(name, time.Duration(1*time.Second), time.Duration(10*time.Second))
	if err != nil {
		t.Fatalf("Wait for storage attachment to become available request failed: %s", err)
	}

	if info.State != Attached {
		fmt.Println(info)
		t.Fatalf("Status of retrieved storage volume attachment was %s, expected 'attached'", info.State)
	}

}

func TestAccStorageAttachmentsClient_WaitForStorageAttachmentToBeFullyAttachedTimeout(t *testing.T) {
	helper.Test(t, helper.TestCase{})

	name := "test"
	server := serverThatAttachesStorageVolumeAfterThreeSeconds(t, name)

	defer server.Close()
	sv, err := getStubStorageAttachmentsClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	_, err = sv.waitForStorageAttachmentToFullyAttach(name, time.Duration(1*time.Second), time.Duration(3*time.Second))
	if err == nil {
		t.Fatal("Expected timeout error")
	}
}

func serverThatAttachesStorageVolumeAfterThreeSeconds(t *testing.T, name string) *httptest.Server {
	count := 0
	return newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		var status string
		if count < 3 {
			status = "attaching"
		} else {
			status = "attached"
		}
		count++
		svr := fmt.Sprintf(
			"{ \"name\": \"/foo/bar/%s\", \"instance_name\": \"/foo/bar/some-instance\", \"storage_volume_name\": \"/foo/bar/example\", \"state\": \"%s\" }", name, status)

		w.Write([]byte(svr))
	})
}

func serverThatDetachesStorageVolumeAfterThreeSeconds(t *testing.T, name string) *httptest.Server {
	count := 0
	return newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		if count < 3 {
			svr := fmt.Sprintf("{ \"name\": \"/foo/bar/%s\", \"instance_name\": \"/foo/bar/some-instance\", \"storage_volume_name\": \"/foo/bar/example\", \"state\": \"detaching\" }", name)
			w.Write([]byte(svr))
		} else {
			w.WriteHeader(404)
		}
		count++

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
