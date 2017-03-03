package compute

type VirtNICSetsClient struct {
	ResourceClient
}

func (c *Client) VirtNICSets() *VirtNICSetsClient {
	return &VirtNICSetsClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "Virtual NIC Set",
			ContainerPath:       "/network/v1/vnicset/",
			ResourceRootPath:    "/network/v1/vnicset",
		},
	}
}

// Describes an existing virtual nic set
type VirtualNICSet struct {
	AppliedACLs []string `json:"appliedAcls"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
	Uri         string   `json:"uri"`
	VirtualNICs []VirtualNIC
	// VirtualNICNames should not be used in any fashion except internally
	// The API returns a slice of names, and we use those names to
	// populate the slice of VirtualNICs
	VirtualNICNames []string `json:"vnics"`
}

type CreateVirtualNICSetInput struct {
	Name            string   `json:"name"`
	AppliedACLs     []string `json:"appliedAcls"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	VirtualNICNames []string `json:"vnics"`
}

func (c *VirtNICSetsClient) CreateVirtualNICSet(input *CreateVirtualNICSetInput) (*VirtualNICSet, error) {
	input.Name = c.getQualifiedName(input.Name)
	input.AppliedACLs = c.getQualifiedAcls(input.AppliedACLs)
	qualifiedNics := c.qualifyVnicSet(input.VirtualNICNames)
	if len(qualifiedNics) != 0 {
		input.VirtualNICNames = qualifiedNics
	}

	var virtNicSet VirtualNICSet
	if err := c.createResource(input, &virtNicSet); err != nil {
		return nil, err
	}

	// Populate the VirtNICs field
	if virtNicSet.VirtualNICNames != nil {
		if err := c.populateVirtualNICs(&virtNicSet); err != nil {
			return nil, err
		}
	}

	return c.success(&virtNicSet)
}

type GetVirtualNICSetInput struct {
	Name string `json:"name"`
}

func (c *VirtNICSetsClient) GetVirtualNICSet(input *GetVirtualNICSetInput) (*VirtualNICSet, error) {
	var virtNicSet VirtualNICSet
	if err := c.getResource(input.Name, &virtNicSet); err != nil {
		return nil, err
	}

	// Populate the VirtNICs field
	if virtNicSet.VirtualNICNames != nil {
		if err := c.populateVirtualNICs(&virtNicSet); err != nil {
			return nil, err
		}
	}

	return c.success(&virtNicSet)
}

type UpdateVirtualNICSetInput struct {
	Name            string   `json:"name"`
	AppliedACLs     []string `json:"appliedAcls"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	VirtualNICNames []string `json:"vnics"`
}

func (c *VirtNICSetsClient) UpdateVirtualNICSet(input *UpdateVirtualNICSetInput) (*VirtualNICSet, error) {
	input.Name = c.getQualifiedName(input.Name)
	input.AppliedACLs = c.getQualifiedAcls(input.AppliedACLs)
	// Qualify VirtualNICNames
	qualifiedNames := []string{}
	for _, v := range input.VirtualNICNames {
		qualifiedNames = append(qualifiedNames, c.getQualifiedName(v))
	}
	if len(qualifiedNames) != 0 {
		input.VirtualNICNames = qualifiedNames
	}

	var virtNICSet VirtualNICSet
	if err := c.updateResource(input.Name, input, &virtNICSet); err != nil {
		return nil, err
	}
	// Populate the VirtNICs field
	if virtNICSet.VirtualNICNames != nil {
		if err := c.populateVirtualNICs(&virtNICSet); err != nil {
			return nil, err
		}
	}
	return c.success(&virtNICSet)
}

type DeleteVirtualNICSetInput struct {
	Name string `json:"name"`
}

func (c *VirtNICSetsClient) DeleteVirtualNICSet(input *DeleteVirtualNICSetInput) error {
	return c.deleteResource(input.Name)
}

func (c *VirtNICSetsClient) populateVirtualNICs(info *VirtualNICSet) error {
	for _, name := range info.VirtualNICNames {
		// Get a VirtNIC client
		nicsClient := c.Client.VirtNICs()
		input := &GetVirtualNICInput{
			Name: name,
		}
		// Get info about the virtual NIC
		nic, err := nicsClient.GetVirtualNIC(input)
		if err != nil {
			return err
		}
		// Append the NIC to the result
		info.VirtualNICs = append(info.VirtualNICs, *nic)
	}
	return nil
}

func (c *VirtNICSetsClient) getQualifiedAcls(acls []string) []string {
	qualifiedAcls := []string{}
	for _, acl := range acls {
		qualifiedAcls = append(qualifiedAcls, c.getQualifiedName(acl))
	}
	return qualifiedAcls
}

func (c *VirtNICSetsClient) unqualifyAcls(acls []string) []string {
	unqualifiedAcls := []string{}
	for _, acl := range acls {
		unqualifiedAcls = append(unqualifiedAcls, c.getUnqualifiedName(acl))
	}
	return unqualifiedAcls
}

func (c *VirtNICSetsClient) success(info *VirtualNICSet) (*VirtualNICSet, error) {
	c.unqualify(&info.Name)
	info.AppliedACLs = c.unqualifyAcls(info.AppliedACLs)
	return info, nil
}

func (c *Client) qualifyVnicSet(set []string) []string {
	qualifiedNics := []string{}
	for _, v := range set {
		qualifiedNics = append(qualifiedNics, c.getQualifiedName(v))
	}
	return qualifiedNics
}

func (c *Client) unqualifyVnicSet(set []string) []string {
	unqualifiedNics := make([]string, len(set))
	for k, v := range set {
		unqualifiedNics[k] = c.getUnqualifiedName(v)
	}
	return unqualifiedNics
}
