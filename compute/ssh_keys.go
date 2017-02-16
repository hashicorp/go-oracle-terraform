package compute

// SSHKeysClient is a client for the SSH key functions of the Compute API.
type SSHKeysClient struct {
	ResourceClient
}

// SSHKeys obtains an SSHKeysClient which can be used to access to the
// SSH key functions of the Compute API
func (c *Client) SSHKeys() *SSHKeysClient {
	return &SSHKeysClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "SSH key",
			ContainerPath:       "/sshkey/",
			ResourceRootPath:    "/sshkey",
		}}
}

// SSHKeyInfo describes an existing SSH key.
type SSHKey struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
	URI     string `json:"uri"`
}


// CreateSSHKeyInput defines an SSH key to be created.
type CreateSSHKeyInput struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
}

// CreateSSHKey creates a new SSH key with the given name, key and enabled flag.
func (c *SSHKeysClient) CreateSSHKey(createInput *CreateSSHKeyInput) (*SSHKey, error) {

	var keyInfo SSHKey

	createInput.Name = c.getQualifiedName(createInput.Name)
	if err := c.createResource(&createInput, &keyInfo); err != nil {
		return nil, err
	}

	return c.success(&keyInfo)
}

// GetSSHKeyInput describes the ssh key to get
type GetSSHKeyInput struct {
	Name string `json:name`
}

// GetSSHKey retrieves the SSH key with the given name.
func (c *SSHKeysClient) GetSSHKey(getInput *GetSSHKeyInput) (*SSHKey, error) {
	var keyInfo SSHKey
	if err := c.getResource(getInput.Name, &keyInfo); err != nil {
		return nil, err
	}

	return c.success(&keyInfo)
}

// UpdateSSHKeyInput defines an SSH key to be updated
type UpdateSSHKeyInput struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
}

// UpdateSSHKey updates the key and enabled flag of the SSH key with the given name.
func (c *SSHKeysClient) UpdateSSHKey(updateInput *UpdateSSHKeyInput) (*SSHKey, error) {

	var keyInfo SSHKey
	updateInput.Name = c.getQualifiedName(updateInput.Name)
	if err := c.updateResource(updateInput.Name, updateInput, &keyInfo); err != nil {
		return nil, err
	}

	return c.success(&keyInfo)
}

// DeleteKeyInput describes the ssh key to delete
type DeleteSSHKeyInput struct {
	Name string `json:name`
}

// DeleteSSHKey deletes the SSH key with the given name.
func (c *SSHKeysClient) DeleteSSHKey(deleteInput *DeleteSSHKeyInput) error {
	return c.deleteResource(deleteInput.Name)
}

func (c *SSHKeysClient) success(keyInfo *SSHKey) (*SSHKey, error) {
	c.unqualify(&keyInfo.Name)
	return keyInfo, nil
}
