package compute

import "fmt"

const WaitForInstanceReadyTimeout = 300
const WaitForInstanceDeleteTimeout = 600

// InstancesClient is a client for the Instance functions of the Compute API.
type InstancesClient struct {
	ResourceClient
}

// Instances obtains an InstancesClient which can be used to access to the
// Instance functions of the Compute API
func (c *Client) Instances() *InstancesClient {
	return &InstancesClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "instance",
			ContainerPath:       "/launchplan/",
			ResourceRootPath:    "/instance",
		}}
}

// InstanceInfo represents the Compute API's view of the state of an instance.
type InstanceInfo struct {
	ID          string                 `json:"id"`
	Shape       string                 `json:"shape"`
	ImageList   string                 `json:"imagelist"`
	Name        string                 `json:"name"`
	Label       string                 `json:"label"`
	BootOrder   []int                  `json:"boot_order"`
	SSHKeys     []string               `json:"sshkeys"`
	State       string                 `json:"state"`
	ErrorReason string                 `json:"error_reason"`
	IPAddress   string                 `json:"ip"`
	VCableID    string                 `json:"vcable_id"`
	Attributes  map[string]interface{} `json:"attributes"`
}

func (i *InstanceInfo) getInstanceName() string {
	return fmt.Sprintf(CMP_QUALIFIED_NAME, i.Name, i.ID)
}

type CreateInstanceInput struct {
	Shape      string                 `json:"shape"`
	ImageList  string                 `json:"imagelist"`
	Name       string                 `json:"name"`
	Label      string                 `json:"label"`
	Storage    []StorageAttachment    `json:"storage_attachments"`
	BootOrder  []int                  `json:"boot_order"`
	SSHKeys    []string               `json:"sshkeys"`
	Attributes map[string]interface{} `json:"attributes"`
}

type StorageAttachment struct {
	Index  int    `json:"index"`
	Volume string `json:"volume"`
}

// LaunchPlan defines a launch plan, used to launch instances with the supplied InstanceSpec(s)
type LaunchPlanInput struct {
	Instances []CreateInstanceInput `json:"instances"`
}

type LaunchPlanResponse struct {
	Instances []InstanceInfo `json:"instances"`
}

// LaunchInstance creates and submits a LaunchPlan to launch a new instance.
func (c *InstancesClient) CreateInstance(input *CreateInstanceInput) (*InstanceInfo, error) {
	qualifiedSSHKeys := []string{}
	for _, key := range input.SSHKeys {
		qualifiedSSHKeys = append(qualifiedSSHKeys, c.getQualifiedName(key))
	}

	input.SSHKeys = qualifiedSSHKeys

	qualifiedStorageAttachments := []StorageAttachment{}
	for _, attachment := range input.Storage {
		qualifiedStorageAttachments = append(qualifiedStorageAttachments, StorageAttachment{
			Index:  attachment.Index,
			Volume: c.getQualifiedName(attachment.Volume),
		})
	}
	input.Storage = qualifiedStorageAttachments

	input.Name = fmt.Sprintf(CMP_QUALIFIED_NAME, c.getUserName(), input.Name)

	plan := LaunchPlanInput{Instances: []CreateInstanceInput{*input}}

	var responseBody LaunchPlanResponse
	if err := c.createResource(&plan, &responseBody); err != nil {
		return nil, err
	}

	if len(responseBody.Instances) == 0 {
		return nil, fmt.Errorf("No instance information returned: %#v", responseBody)
	}

	// Call wait for instance ready now, as creating the instance is an eventually consistent operation
	getInput := &GetInstanceInput{
		Name: input.Name,
		ID:   responseBody.Instances[0].ID,
	}

	result, err := c.WaitForInstanceRunning(getInput, WaitForInstanceReadyTimeout)
	if err != nil {
		return nil, err
	}

	// Unqualify instance name
	result.Name = input.Name

	return result, nil
}

type GetInstanceInput struct {
	Name string
	ID   string
}

func (g *GetInstanceInput) String() string {
	return fmt.Sprintf(CMP_QUALIFIED_NAME, g.Name, g.ID)
}

// GetInstance retrieves information about an instance.
func (c *InstancesClient) GetInstance(input *GetInstanceInput) (*InstanceInfo, error) {
	var responseBody InstanceInfo
	if err := c.getResource(input.String(), &responseBody); err != nil {
		return nil, err
	}

	if responseBody.Name == "" {
		return nil, fmt.Errorf("Empty response body when requesting instance %s", input.Name)
	}

	// Overwrite returned name/ID with known name/ID
	// Otherwise the returned name will be the fully qualified name, and the ID will be blank
	responseBody.Name = input.Name
	responseBody.ID = input.ID
	c.unqualify(&responseBody.VCableID)

	// Unqualify SSH Key names
	sshKeyNames := []string{}
	for _, sshKeyRef := range responseBody.SSHKeys {
		sshKeyNames = append(sshKeyNames, c.getUnqualifiedName(sshKeyRef))
	}
	responseBody.SSHKeys = sshKeyNames

	return &responseBody, nil
}

type DeleteInstanceInput struct {
	Name string
	ID   string
}

func (d *DeleteInstanceInput) String() string {
	return fmt.Sprintf(CMP_QUALIFIED_NAME, d.Name, d.ID)
}

// DeleteInstance deletes an instance.
func (c *InstancesClient) DeleteInstance(input *DeleteInstanceInput) error {
	// Call to delete the instance
	if err := c.deleteResource(input.String()); err != nil {
		return err
	}
	// Wait for instance to be deleted
	return c.WaitForInstanceDeleted(input, WaitForInstanceDeleteTimeout)
}

// WaitForInstanceRunning waits for an instance to be completely initialized and available.
func (c *InstancesClient) WaitForInstanceRunning(input *GetInstanceInput, timeoutSeconds int) (*InstanceInfo, error) {
	var info *InstanceInfo
	var getErr error
	err := c.waitFor("instance to be ready", timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetInstance(input)
		if getErr != nil {
			return false, getErr
		}
		if info.State == "error" {
			return false, fmt.Errorf("Error initializing instance: %s", info.ErrorReason)
		}
		if info.State == "running" {
			return true, nil
		}
		return false, nil
	})
	return info, err
}

// WaitForInstanceDeleted waits for an instance to be fully deleted.
func (c *InstancesClient) WaitForInstanceDeleted(input *DeleteInstanceInput, timeoutSeconds int) error {
	return c.waitFor("instance to be deleted", timeoutSeconds, func() (bool, error) {
		var instanceInfo InstanceInfo
		if err := c.getResource(input.String(), &instanceInfo); err != nil {
			if WasNotFoundError(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}
