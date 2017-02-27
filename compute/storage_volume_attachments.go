package compute

import (
	"fmt"
	"log"
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

// StorageAttachmentSpec defines a storage attachment to be created.
type StorageAttachmentSpec struct {
	Index             int    `json:"index"`
	InstanceName      string `json:"instance_name"`
	StorageVolumeName string `json:"storage_volume_name"`
}

// StorageAttachmentInfo describes an existing storage attachment.
type StorageAttachmentInfo struct {
	Name              string `json:"name"`
	Index             int    `json:"index"`
	InstanceName      string `json:"instance_name"`
	StorageVolumeName string `json:"storage_volume_name"`
	State             string `json:"state"`
}

// StorageAttachmentList is a collection of storage attachments attached to a specific instance.
type StorageAttachmentList struct {
	Result []StorageAttachmentInfo `json:"result"`
}

func (c *StorageAttachmentsClient) success(attachmentInfo *StorageAttachmentInfo) (*StorageAttachmentInfo, error) {
	c.unqualify(&attachmentInfo.Name, &attachmentInfo.InstanceName, &attachmentInfo.StorageVolumeName)
	return attachmentInfo, nil
}

type CreateStorageAttachmentInput struct {
	Index             int
	InstanceName      string
	InstanceId        string
	StorageVolumeName string
}

// CreateStorageAttachment creates a storage attachment attaching the given volume to the given instance at the given index.
func (c *StorageAttachmentsClient) CreateStorageAttachment(input *CreateStorageAttachmentInput) (*StorageAttachmentInfo, error) {
	instanceInfo := InstanceInfo{
		Name: input.InstanceName,
		ID:   input.InstanceId,
	}
	storageVolumeName := c.getQualifiedName(input.StorageVolumeName)
	spec := StorageAttachmentSpec{
		Index:             input.Index,
		InstanceName:      c.getQualifiedName(instanceInfo.getInstanceName()),
		StorageVolumeName: storageVolumeName,
	}

	var attachmentInfo StorageAttachmentInfo
	if err := c.createResource(&spec, &attachmentInfo); err != nil {
		return nil, err
	}

	err := c.WaitForStorageAttachmentCreated(attachmentInfo.Name, WaitForVolumeAttachmentReadyTimeout)
	if err != nil {
		return nil, err
	}

	return c.success(&attachmentInfo)
}

// DeleteStorageAttachment deletes the storage attachment with the given name.
func (c *StorageAttachmentsClient) DeleteStorageAttachment(name string) error {
	if err := c.deleteResource(name); err != nil {
		return err
	}

	return c.WaitForStorageAttachmentDeleted(name, WaitForVolumeAttachmentDeleteTimeout)
}

// GetStorageAttachment retrieves the storage attachment with the given name.
func (c *StorageAttachmentsClient) GetStorageAttachment(name string) (*StorageAttachmentInfo, error) {
	var attachmentInfo StorageAttachmentInfo
	if err := c.getResource(name, &attachmentInfo); err != nil {
		return nil, err
	}

	return c.success(&attachmentInfo)
}

// WaitForStorageAttachmentCreated waits for the storage attachment with the given name to be fully attached, or times out.
func (c *StorageAttachmentsClient) WaitForStorageAttachmentCreated(name string, timeoutSeconds int) error {
	return c.waitFor("storage attachment to be attached", timeoutSeconds, func() (bool, error) {
		info, err := c.GetStorageAttachment(name)
		if err != nil {
			return false, err
		}
		if info.State == "attached" {
			return true, nil
		}
		return false, nil
	})
}

// WaitForStorageAttachmentDeleted waits for the storage attachment with the given name to be fully deleted, or times out.
func (c *StorageAttachmentsClient) WaitForStorageAttachmentDeleted(name string, timeoutSeconds int) error {
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

// GetStorageAttachmentsForInstance retrieves all of the storage attachments for the given instance.
func (c *StorageAttachmentsClient) GetStorageAttachmentsForInstance(info *InstanceInfo) (*[]StorageAttachmentInfo, error) {
	return c.getStorageAttachments(
		fmt.Sprintf("instance_name=%s", c.getQualifiedName(info.getInstanceName())),
		"instance",
	)
}

// GetStorageAttachmentsForInstance retrieves all of the storage attachments for the given volume.
func (c *StorageAttachmentsClient) GetStorageAttachmentsForVolume(name string) (*[]StorageAttachmentInfo, error) {
	return c.getStorageAttachments(
		fmt.Sprintf("storage_volume_name=%s", c.getQualifiedName(name)),
		"volume",
	)
}

func (c *StorageAttachmentsClient) getStorageAttachments(query string, description string) (*[]StorageAttachmentInfo, error) {
	queryPath := fmt.Sprintf("/storage/attachment%s/?state=attached&%s",
		c.getUserName(),
		query)
	log.Printf("[DEBUG] Querying for storage attachments: %s", queryPath)

	resp, err := c.executeRequest("GET", queryPath, nil)
	if err != nil {
		return nil, err
	}

	var attachmentList StorageAttachmentList
	if err = c.unmarshalResponseBody(resp, &attachmentList); err != nil {
		return nil, err
	}

	attachments := make([]StorageAttachmentInfo, len(attachmentList.Result))
	for index, attachment := range attachmentList.Result {
		c.unqualify(&attachment.Name, &attachment.InstanceName, &attachment.StorageVolumeName)
		attachments[index] = attachment
	}
	return &attachments, nil
}
