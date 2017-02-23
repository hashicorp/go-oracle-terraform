package compute

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
)

func TestAccStorageVolumeClient_WaitForStorageVolumeToBeDeletedSuccessful(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := serverWhereStorageVolumeGetsDeletedAfterThreeSeconds(t)
	name := "test"

	defer server.Close()
	sv, err := getStubStorageVolumeClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	sv.VolumeModificationTimeout = 10
	err = sv.waitForStorageVolumeToBeDeleted(name)
	if err != nil {
		t.Fatalf("Wait for storage volume deleted request failed: %s", err)
	}

	getRequest := &GetStorageVolumeInput{
		Name: name,
	}
	getResponse, err := sv.GetStorageVolume(getRequest)
	if err != nil {
		t.Fatalf("error getting storage volume: %s", err)
	}

	if len(getResponse.Result) > 0 {
		t.Fatal("Expected Storage Volume to be Deleted")
	}
}

func TestAccStorageVolumeClient_WaitForStorageVolumeToBeDeletedTimeout(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := serverWhereStorageVolumeGetsDeletedAfterThreeSeconds(t)
	name := "test"

	defer server.Close()
	sv, err := getStubStorageVolumeClient(server)
	if err != nil {
		t.Fatalf("error getting stub client: %s", err)
	}

	sv.VolumeModificationTimeout = 3
	err = sv.waitForStorageVolumeToBeDeleted(name)
	if err == nil {
		t.Fatal("Expected timeout error")
	}
}

func TestAccStorageVolumeClient_WaitForStorageVolumeToBecomeAvailableSuccessful(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := serverWhereStorageVolumeBecomesAvailableAfterThreeSeconds(t)

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

func TestAccStorageVolumeClient_WaitForStorageVolumeToBecomeAvailableTimeout(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	server := serverWhereStorageVolumeBecomesAvailableAfterThreeSeconds(t)

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

func serverWhereStorageVolumeBecomesAvailableAfterThreeSeconds(t *testing.T) *httptest.Server {
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

func serverWhereStorageVolumeGetsDeletedAfterThreeSeconds(t *testing.T) *httptest.Server {
	count := 0
	return newAuthenticatingServer(func(w http.ResponseWriter, r *http.Request) {
		var status string
		if count < 3 {
			status = "{\"result\":[{\"name\":\"/storage/volume/test\",\"size\":\"16G\",\"status\":\"Deleting\"}]}"
		} else {
			status = "{\"result\": []}"
		}
		count++

		w.Write([]byte(status))
		w.WriteHeader(400)
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
