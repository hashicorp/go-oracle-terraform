package compute

import (
	"fmt"
	"strings"
)

// IPAssociationsClient is a client for the IP Association functions of the Compute API.
type IPAssociationsClient struct {
	*ResourceClient
}

// IPAssociations obtains a IPAssociationsClient which can be used to access to the
// IP Association functions of the Compute API
func (c *Client) IPAssociations() *IPAssociationsClient {
	return &IPAssociationsClient{
		ResourceClient: &ResourceClient{
			Client:              c,
			ResourceDescription: "ip association",
			ContainerPath:       "/ip/association/",
			ResourceRootPath:    "/ip/association",
		}}
}

// IPAssociationInfo describes an existing IP association.
type IPAssociationInfo struct {
	Name        string `json:"name"`
	VCable      string `json:"vcable"`
	ParentPool  string `json:"parentpool"`
	URI         string `json:"uri"`
	Reservation string `json:"reservation"`
}

type CreateIPAssociationInput struct {
	VCable     string `json:"vcable"`
	ParentPool string `json:"parentpool"`
}

// CreateIPAssociation creates a new IP association with the supplied vcable and parentpool.
func (c *IPAssociationsClient) CreateIPAssociation(input *CreateIPAssociationInput) (*IPAssociationInfo, error) {
	input.VCable = c.getQualifiedName(input.VCable)
	input.ParentPool = c.getQualifiedParentPoolName(input.ParentPool)
	var assocInfo IPAssociationInfo
	if err := c.createResource(input, &assocInfo); err != nil {
		return nil, err
	}

	return c.success(&assocInfo)
}

type GetIPAssociationInput struct {
	Name string `json:"name"`
}

// GetIPAssociation retrieves the IP association with the given name.
func (c *IPAssociationsClient) GetIPAssociation(input *GetIPAssociationInput) (*IPAssociationInfo, error) {
	var assocInfo IPAssociationInfo
	if err := c.getResource(input.Name, &assocInfo); err != nil {
		return nil, err
	}

	return c.success(&assocInfo)
}

type DeleteIPAssociationInput struct {
	Name string `json:"name"`
}

// DeleteIPAssociation deletes the IP association with the given name.
func (c *IPAssociationsClient) DeleteIPAssociation(input *DeleteIPAssociationInput) error {
	return c.deleteResource(input.Name)
}

func (c *IPAssociationsClient) getQualifiedParentPoolName(parentpool string) string {
	parts := strings.Split(parentpool, ":")
	pooltype := parts[0]
	name := parts[1]
	return fmt.Sprintf("%s:%s", pooltype, c.getQualifiedName(name))
}

func (c *IPAssociationsClient) unqualifyParentPoolName(parentpool *string) {
	parts := strings.Split(*parentpool, ":")
	pooltype := parts[0]
	name := parts[1]
	*parentpool = fmt.Sprintf("%s:%s", pooltype, c.getUnqualifiedName(name))
}

// Unqualifies identifiers
func (c *IPAssociationsClient) success(assocInfo *IPAssociationInfo) (*IPAssociationInfo, error) {
	c.unqualify(&assocInfo.Name, &assocInfo.VCable)
	c.unqualifyParentPoolName(&assocInfo.ParentPool)
	return assocInfo, nil
}
