package compute

// SecurityListsClient is a client for the Security List functions of the Compute API.
type SecurityListsClient struct {
	ResourceClient
}

// SecurityLists obtains a SecurityListsClient which can be used to access to the
// Security List functions of the Compute API
func (c *Client) SecurityLists() *SecurityListsClient {
	return &SecurityListsClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "security list",
			ContainerPath:       "/seclist/",
			ResourceRootPath:    "/seclist",
		}}
}

// SecurityListInfo describes an existing security list.
type SecurityListInfo struct {
	Account            string `json:"account"`
	Description        string `json:description`
	Name               string `json:"name"`
	OutboundCIDRPolicy string `json:"outbound_cidr_policy"`
	Policy             string `json:"policy"`
	URI                string `json:"uri"`
}

// CreateSecurityListInput defines a security list to be created.
type CreateSecurityListInput struct {
	Description        string `json:description`
	Name               string `json:"name"`
	OutboundCIDRPolicy string `json:"outbound_cidr_policy"`
	Policy             string `json:"policy"`
}

// CreateSecurityList creates a new security list with the given name, policy and outbound CIDR policy.
func (c *SecurityListsClient) CreateSecurityList(createInput *CreateSecurityListInput) (*SecurityListInfo, error) {

	createInput.Name = c.getQualifiedName(createInput.Name)

	var listInfo SecurityListInfo
	if err := c.createResource(createInput, &listInfo); err != nil {
		return nil, err
	}

	return c.success(&listInfo)
}

// GetSecurityListInput describes the security list you want to get
type GetSecurityListInput struct {
	Name string `json:name`
}

// GetSecurityList retrieves the security list with the given name.
func (c *SecurityListsClient) GetSecurityList(getInput *GetSecurityListInput) (*SecurityListInfo, error) {
	var listInfo SecurityListInfo
	if err := c.getResource(getInput.Name, &listInfo); err != nil {
		return nil, err
	}

	return c.success(&listInfo)
}

// UpdateSecurityListInput defines what to update in a security list
type UpdateSecurityListInput struct {
	Description        string `json:description`
	Name               string `json:"name"`
	OutboundCIDRPolicy string `json:"outbound_cidr_policy"`
	Policy             string `json:"policy"`
}

// UpdateSecurityList updates the policy and outbound CIDR pol
func (c *SecurityListsClient) UpdateSecurityList(updateInput *UpdateSecurityListInput) (*SecurityListInfo, error) {
	updateInput.Name = c.getQualifiedName(updateInput.Name)

	var listInfo SecurityListInfo
	if err := c.updateResource(updateInput.Name, updateInput, &listInfo); err != nil {
		return nil, err
	}

	return c.success(&listInfo)
}

// DeleteSecurityListInput describes the security list to destroy
type DeleteSecurityListInput struct {
	Name string `json:name`
}

// DeleteSecurityList deletes the security list with the given name.
func (c *SecurityListsClient) DeleteSecurityList(deleteInput *DeleteSecurityListInput) error {
	return c.deleteResource(deleteInput.Name)
}

func (c *SecurityListsClient) success(listInfo *SecurityListInfo) (*SecurityListInfo, error) {
	c.unqualify(&listInfo.Name)
	return listInfo, nil
}
