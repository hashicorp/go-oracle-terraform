package compute

import (
	"fmt"
)

const WaitForVolumeReadyTimeout = 30
const WaitForVolumeDeleteTimeout = 30

// StorageVolumeClient is a client for the Storage Volume functions of the Compute API.
type StorageVolumeClient struct {
	ResourceClient
}

// StorageVolumes obtains a StorageVolumeClient which can be used to access to the
// Storage Volume functions of the Compute API
func (c *Client) StorageVolumes() *StorageVolumeClient {
	return &StorageVolumeClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "storage volume",
			ContainerPath:       "/storage/volume/",
			ResourceRootPath:    "/storage/volume",
		}}

}

// StorageVolumeInfo represents information retrieved from the service about a Storage Volume.
type StorageVolumeInfo struct {
	Account         string   `json:"account,omitempty"`
	Bootable        bool     `json:"bootable,omitempty"`
	Description     string   `json:"description,omitempty"`
	Hypervisor      string   `json:"hypervisor,omitempty"`
	ImageList       string   `json:"imagelist,omitempty"`
	ImageListEntry  int      `json:"imagelist_entry,omitempty"`
	MachineImage    string   `json:"machineimage_name,omitempty"`
	Managed         bool     `json:"managed,omitempty"`
	Name            string   `json:"name"`
	Platform        string   `json:"platform,omitempty`
	Properties      []string `json:"properties,omitempty"`
	Quota           string   `json:"quota,omitempty"`
	ReadOnly        bool     `json:"readonly,omitempty"`
	Shared          bool     `json:"shared,omitempty"`
	Size            string   `json:"size"`
	Snapshot        string   `json:"snapshot,omitempty"`
	SnapshotAccount string   `json:"snapshot_account,omitempty"`
	SnapshotID      string   `json:"snapshot_id,omitempty"`
	Status          string   `json:"status,omitempty"`
	StatusDetail    string   `json:"status_detail,omitempty"`
	StatusTimestamp string   `json:"status_timestamp,omitempty"`
	StoragePool     string   `json:"storage_pool,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	URI             string   `json:"uri,omitempty"`
	WriteCache      bool     `json:"writecache,omitempty"`
}

// StorageVolumeResult represents the body of a response to a query for Storage Volume information.
type StorageVolumeResult struct {
	Result []StorageVolumeInfo `json:"result"`
}

func (c *StorageVolumeClient) getStorageVolumePath(name string) string {
	return c.getObjectPath("/storage/volume", name) + "/"
}

// CreateStorageVolumeInput represents the body of an API request to create a new Storage Volume.
type CreateStorageVolumeInput struct {
	Bootable        bool     `json:"bootable,omitempty"`
	Description     string   `json:"description,omitempty"`
	ImageList       string   `json:"imagelist,omitempty"`
	ImageListEntry  int      `json:"imagelist_entry,omitempty"`
	Name            string   `json:"name"`
	Properties      []string `json:"properties,omitempty"`
	Size            string   `json:"size"`
	Snapshot        string   `json:"snapshot,omitempty"`
	SnapshotAccount string   `json:"snapshot_account,omitempty"`
	SnapshotID      string   `json:"snapshot_id,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

// CreateStorageVolume uses the given CreateStorageVolumeInput to create a new Storage Volume.
func (c *StorageVolumeClient) CreateStorageVolume(input *CreateStorageVolumeInput) (*StorageVolumeInfo, error) {

	input.Name = c.getQualifiedName(input.Name)
	_, err := c.executeRequest("POST", c.ContainerPath, input)
	if err != nil {
		return nil, err
	}

	instanceInfo, err := c.waitForStorageVolumeToBecomeAvailable(input.Name, WaitForVolumeReadyTimeout)
	if err != nil {
		return nil, err
	}

	return instanceInfo, nil
}

// DeleteStorageVolumeInput represents the body of an API request to delete a Storage Volume.
type DeleteStorageVolumeInput struct {
	Name string `json:"name"`
}

// DeleteStorageVolume deletes the specified storage volume.
func (c *StorageVolumeClient) DeleteStorageVolume(input *DeleteStorageVolumeInput) error {
	_, err := c.executeRequest("DELETE", c.getStorageVolumePath(input.Name), nil)
	if err != nil {
		return err
	}

	err = c.waitForStorageVolumeToBeDeleted(input.Name, WaitForVolumeDeleteTimeout)
	if err != nil {
		return err
	}

	return nil
}

// GetStorageVolumeInput represents the body of an API request to obtain a Storage Volume.
type GetStorageVolumeInput struct {
	Name string `json:"name"`
}

// GetStorageVolume gets Storage Volume information for the specified storage volume.
func (c *StorageVolumeClient) GetStorageVolume(input *GetStorageVolumeInput) (*StorageVolumeResult, error) {
	resp, err := c.executeRequest("GET", c.getStorageVolumePath(input.Name), nil)
	if err != nil {
		return nil, err
	}

	var result StorageVolumeResult
	err = c.unmarshalResponseBody(resp, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Result) > 0 {
		c.unqualify(&result.Result[0].Name)
	}
	return &result, nil
}

// UpdateStorageVolumeInput represents the body of an API request to update a Storage Volume.
type UpdateStorageVolumeInput struct {
	Description     string   `json:"description,omitempty"`
	ImageList       string   `json:"imagelist,omitempty"`
	ImageListEntry  int      `json:"imagelist_entry,omitempty"`
	Name            string   `json:"name"`
	Properties      []string `json:"properties"`
	Size            string   `json:"size"`
	Snapshot        string   `json:"snapshot,omitempty"`
	SnapshotAccount string   `json:"snapshot_account,omitempty"`
	SnapshotID      string   `json:"snapshot_id,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

// UpdateStorageVolume updates the specified storage volume, optionally modifying size, description and tags.
func (c *StorageVolumeClient) UpdateStorageVolume(input *UpdateStorageVolumeInput) (*StorageVolumeInfo, error) {
	input.Name = c.getQualifiedName(input.Name)
	path := c.getStorageVolumePath(input.Name)
	_, err := c.executeRequest("PUT", path, input)
	if err != nil {
		return nil, err
	}

	instanceInfo, err := c.waitForStorageVolumeToBecomeAvailable(input.Name, WaitForVolumeReadyTimeout)
	if err != nil {
		return nil, err
	}

	return instanceInfo, nil
}

// waitForStorageVolumeToBecomeAvailable waits until a new Storage Volume is available (i.e. has finished initialising or updating).
func (c *StorageVolumeClient) waitForStorageVolumeToBecomeAvailable(name string, timeoutInSeconds int) (*StorageVolumeInfo, error) {
	var waitResult *StorageVolumeInfo

	err := c.waitFor(
		fmt.Sprintf("storage volume %s to become available", c.getQualifiedName(name)),
		timeoutInSeconds,
		func() (bool, error) {
			getRequest := &GetStorageVolumeInput{
				Name: name,
			}
			result, err := c.GetStorageVolume(getRequest)

			if err != nil {
				return false, err
			}

			if len(result.Result) > 0 {
				waitResult = &result.Result[0]
				if waitResult.Status == "Online" {
					return true, nil
				}
			}

			return false, nil
		})

	return waitResult, err
}

// waitForStorageVolumeToBeDeleted waits until the specified storage volume has been deleted.
func (c *StorageVolumeClient) waitForStorageVolumeToBeDeleted(name string, timeoutInSeconds int) error {
	return c.waitFor(
		fmt.Sprintf("storage volume %s to be deleted", c.getQualifiedName(name)),
		timeoutInSeconds,
		func() (bool, error) {
			getRequest := &GetStorageVolumeInput{
				Name: name,
			}
			result, err := c.GetStorageVolume(getRequest)
			if err != nil {
				return false, err
			}

			return len(result.Result) == 0, nil
		})
}
