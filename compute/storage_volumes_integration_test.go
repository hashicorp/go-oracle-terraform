package compute

import (
	"strconv"
	"testing"

	"github.com/hashicorp/go-oracle-terraform/helper"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

func TestAccStorageVolumeLifecycle(t *testing.T) {
	helper.Test(t, helper.TestCase{})
	defer tearDownStorageVolumes()

	svc, err := getStorageVolumeClient()
	if err != nil {
		t.Fatal(err)
	}

	spec := svc.NewStorageVolumeSpec("10G", []string{"/oracle/public/storage/default"}, "myVolume")
	spec.SetDescription("MyDescription")

	err = svc.CreateStorageVolume(spec)

	if err != nil {
		t.Fatalf("Create volume failed: %s\n", err)
	}

	info, err := svc.WaitForStorageVolumeOnline("myVolume", 30)
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := strconv.Itoa(10 << 30)
	if info.Size != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, info.Size)
	}

	err = svc.UpdateStorageVolume("myVolume", "20G", "redescribe", []string{})
	if err != nil {
		t.Fatal(err)
	}

	info, err = svc.WaitForStorageVolumeOnline("myVolume", 30)
	if err != nil {
		t.Fatal(err)
	}

	expectedSize = strconv.Itoa(20 << 30)
	if info.Size != expectedSize {
		t.Fatalf("Expected storage volume size %s, but was %s", expectedSize, info.Size)
	}
}

func tearDownStorageVolumes() {
	svc, err := getStorageVolumeClient()
	if err != nil {
		panic(err)
	}

	err = svc.DeleteStorageVolume("myVolume")
	if err != nil {
		panic(err)

	}

	err = svc.WaitForStorageVolumeDeleted("myVolume", 30)
	if err != nil {
		panic(err)
	}
}

func getStorageVolumeClient() (*StorageVolumeClient, error) {
	client, err := getTestClient(&opc.Config{})
	if err != nil {
		return &StorageVolumeClient{}, err
	}

	return client.StorageVolumes(), nil
}
