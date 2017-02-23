package compute

type VirtNICsClient struct {
	ResourceClient
}

func (c *Client) VirtNICs() *VirtNICsClient {
	return &VirtNICsClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "Virtual NIC",
			ContainerPath:       "/network/v1/vnic/",
			ResourceRootPath:    "/network/v1/vnic",
		},
	}
}

type VirtualNIC struct {
	Description string   `json:"description"`
	MACAddress  string   `json:"macAddress"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
	TransitFlag bool     `json:"transitFlag"`
	Uri         string   `json:"uri"`
}

// Can only GET a virtual NIC, not update, create, or delete
type GetVirtualNICInput struct {
	Name string `json:"name"`
}

func (c *VirtNICsClient) GetVirtualNIC(input *GetVirtualNICInput) (*VirtualNIC, error) {
	var virtNIC VirtualNIC
	if err := c.getResource(input.Name, &virtNIC); err != nil {
		return nil, err
	}
	return c.success(&virtNIC)
}

func (c *VirtNICsClient) success(info *VirtualNIC) (*VirtualNIC, error) {
	c.unqualify(&info.Name)
	return info, nil
}
