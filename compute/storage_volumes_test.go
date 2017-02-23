package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

func TestAccStorageVolumeClient_WaitForStorageVolumeOnline(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := serverThatReturnsOnlineStorageVolumeAfterThreeSeconds(t)

	defer server.Close()
	sv, err := getStubStorageVolumeClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	sv.VolumeModificationTimeout = 10
	info, err := sv.waitForStorageVolumeToBecomeAvailable("test")
	if err != nil {
		t.Fatalf("Wait for storage volume online request failed: %s", err)
	}

	if info.Status != "Online" {
		fmt.Println(info)
		t.Fatalf("Status of retrieved storage volume info was %s, expected 'Online'", info.Status)
	}
}

func TestAccStorageVolumeClient_WaitForStorageVolumeOnlineTimeout(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := serverThatReturnsOnlineStorageVolumeAfterThreeSeconds(t)

	defer server.Close()
	sv, err := getStubStorageVolumeClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}
	sv.VolumeModificationTimeout = 3
	_, err = sv.waitForStorageVolumeToBecomeAvailable("test")
	if err == nil {
		t.Fatal("Expected timeout error")
	}
}

func serverThatReturnsOnlineStorageVolumeAfterThreeSeconds(t *testing.T) *httptest.Server {
	count := 0
	return newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		var status string
		if count < 3 {
			status = "Initializing"
		} else {
			status = "Online"
		}
		count++
		svr := fmt.Sprintf(
			"{\"result\":[{\"name\":\"/Compute-test/test/test\",\"size\":\"16G\",\"status\":\"%s\"}]}", status)

		w.Write([]byte(svr))
		w.WriteHeader(200)
	})
}

func getStubStorageVolumeClient(server *httptest.Server) (*StorageVolumeClient, error) {
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	client, err := getStubClient(endpoint)
	if err != nil {
		return nil, err
	}
	return client.StorageVolumes(), nil
}
