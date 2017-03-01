package compute

import (
	"strings"
)

const WaitForVolumeAttachmentDeleteTimeout = 30
const WaitForVolumeAttachmentReadyTimeout = 30

// StorageAttachmentsClient is a client for the Storage Attachment functions of the Compute API.
type StorageAttachmentsClient struct {
	ResourceClient
}

// StorageAttachments obtains a StorageAttachmentsClient which can be used to access to the
// Storage Attachment functions of the Compute API
func (c *Client) StorageAttachments() *StorageAttachmentsClient {
	return &StorageAttachmentsClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "storage volume attachment",
			ContainerPath:       "/storage/attachment/",
			ResourceRootPath:    "/storage/attachment",
		}}
}

// StorageAttachmentInfo describes an existing storage attachment.
type StorageAttachmentInfo struct {
	Name              string `json:"name"`
	Index             int    `json:"index"`
	InstanceName      string `json:"instance_name"`
	StorageVolumeName string `json:"storage_volume_name"`
	State             string `json:"state"`
}

func (c *StorageAttachmentsClient) success(attachmentInfo *StorageAttachmentInfo) (*StorageAttachmentInfo, error) {
	c.unqualify(&attachmentInfo.Name)
	return attachmentInfo, nil
}

type CreateStorageAttachmentInput struct {
	Index             int    `json:"index"`
	InstanceName      string `json:"instance_name"`
	StorageVolumeName string `json:"storage_volume_name"`
}

// CreateStorageAttachment creates a storage attachment attaching the given volume to the given instance at the given index.
func (c *StorageAttachmentsClient) CreateStorageAttachment(input *CreateStorageAttachmentInput) (*StorageAttachmentInfo, error) {
	input.InstanceName = c.getQualifiedName(input.InstanceName)

	var attachmentInfo *StorageAttachmentInfo
	if err := c.createResource(&input, &attachmentInfo); err != nil {
		return nil, err
	}

	return c.waitForStorageAttachmentToFullyAttach(attachmentInfo.Name, WaitForVolumeAttachmentReadyTimeout)
}

// DeleteStorageAttachment deletes the storage attachment with the given name.
func (c *StorageAttachmentsClient) DeleteStorageAttachment(name string) error {
	if err := c.deleteResource(name); err != nil {
		return err
	}

	return c.waitForStorageAttachmentToBeDeleted(name, WaitForVolumeAttachmentDeleteTimeout)
}

// GetStorageAttachment retrieves the storage attachment with the given name.
func (c *StorageAttachmentsClient) GetStorageAttachment(name string) (*StorageAttachmentInfo, error) {
	var attachmentInfo *StorageAttachmentInfo
	if err := c.getResource(name, &attachmentInfo); err != nil {
		return nil, err
	}

	return c.success(attachmentInfo)
}

// waitForStorageAttachmentToFullyAttach waits for the storage attachment with the given name to be fully attached, or times out.
func (c *StorageAttachmentsClient) waitForStorageAttachmentToFullyAttach(name string, timeoutSeconds int) (*StorageAttachmentInfo, error) {
	var waitResult *StorageAttachmentInfo

	err := c.waitFor("storage attachment to be attached", timeoutSeconds, func() (bool, error) {
		info, err := c.GetStorageAttachment(name)
		if err != nil {
			return false, err
		}

		if info != nil {
			if strings.ToLower(info.State) == "attached" {
				waitResult = info
				return true, nil
			}
		}

		return false, nil
	})

	return waitResult, err
}

// waitForStorageAttachmentToBeDeleted waits for the storage attachment with the given name to be fully deleted, or times out.
func (c *StorageAttachmentsClient) waitForStorageAttachmentToBeDeleted(name string, timeoutSeconds int) error {
	return c.waitFor("storage attachment to be deleted", timeoutSeconds, func() (bool, error) {
		_, err := c.GetStorageAttachment(name)
		if err != nil {
			if WasNotFoundError(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}
