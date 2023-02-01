// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

const (
	routesDescription   = "IP Network Route"
	routesContainerPath = "/network/v1/route/"
	routesResourcePath  = "/network/v1/route"
)

// RoutesClient specifies the attributes of a route client
type RoutesClient struct {
	ResourceClient
}

// Routes returns a route client
func (c *Client) Routes() *RoutesClient {
	return &RoutesClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: routesDescription,
			ContainerPath:       routesContainerPath,
			ResourceRootPath:    routesResourcePath,
		},
	}
}

// RouteInfo details the attributes for a route
type RouteInfo struct {
	// Admin distance associated with this route
	AdminDistance int `json:"adminDistance"`
	// Description of the route
	Description string `json:"description"`
	// Fully Qualified Domain Name
	FQDN string `json:"name"`
	// CIDR IPv4 Prefix associated with this route
	IPAddressPrefix string `json:"ipAddressPrefix"`
	// Name of the route
	Name string
	// Name of the VNIC set associated with the route
	NextHopVnicSet string `json:"nextHopVnicSet"`
	// Slice of Tags associated with the route
	Tags []string `json:"tags,omitempty"`
	// Uniform resource identifier associated with the route
	URI string `json:"uri"`
}

// CreateRouteInput details the attributes needed to create a route
type CreateRouteInput struct {
	// Specify 0,1, or 2 as the route's administrative distance.
	// If you do not specify a value, the default value is 0.
	// The same prefix can be used in multiple routes. In this case, packets are routed over all the matching
	// routes with the lowest administrative distance.
	// In the case multiple routes with the same lowest administrative distance match,
	// routing occurs over all these routes using ECMP.
	// Optional
	AdminDistance int `json:"adminDistance"`
	// Description of the route
	// Optional
	Description string `json:"description"`
	// The IPv4 address prefix in CIDR format, of the external network (external to the vNIC set)
	// from which you want to route traffic
	// Required
	IPAddressPrefix string `json:"ipAddressPrefix"`
	// Name of the route.
	// Names can only contain alphanumeric, underscore, dash, and period characters. Case-sensitive
	// Required
	Name string `json:"name"`
	// Name of the virtual NIC set to route matching packets to.
	// Routed flows are load-balanced among all the virtual NICs in the virtual NIC set
	// Required
	NextHopVnicSet string `json:"nextHopVnicSet"`
	// Slice of tags to be associated with the route
	// Optional
	Tags []string `json:"tags,omitempty"`
}

// CreateRoute creates the requested route
func (c *RoutesClient) CreateRoute(input *CreateRouteInput) (*RouteInfo, error) {
	input.Name = c.getQualifiedName(input.Name)
	input.NextHopVnicSet = c.getQualifiedName(input.NextHopVnicSet)

	var routeInfo RouteInfo
	if err := c.createResource(&input, &routeInfo); err != nil {
		return nil, err
	}

	return c.success(&routeInfo)
}

// GetRouteInput details the attributes needed to retrive a route
type GetRouteInput struct {
	// Name of the Route to query for. Case-sensitive
	// Required
	Name string `json:"name"`
}

// GetRoute retrieves the specified route
func (c *RoutesClient) GetRoute(input *GetRouteInput) (*RouteInfo, error) {
	input.Name = c.getQualifiedName(input.Name)

	var routeInfo RouteInfo
	if err := c.getResource(input.Name, &routeInfo); err != nil {
		return nil, err
	}
	return c.success(&routeInfo)
}

// UpdateRouteInput details the attributes needed to update a route
type UpdateRouteInput struct {
	// Specify 0,1, or 2 as the route's administrative distance.
	// If you do not specify a value, the default value is 0.
	// The same prefix can be used in multiple routes. In this case, packets are routed over all the matching
	// routes with the lowest administrative distance.
	// In the case multiple routes with the same lowest administrative distance match,
	// routing occurs over all these routes using ECMP.
	// Optional
	AdminDistance int `json:"adminDistance"`
	// Description of the route
	// Optional
	Description string `json:"description"`
	// The IPv4 address prefix in CIDR format, of the external network (external to the vNIC set)
	// from which you want to route traffic
	// Required
	IPAddressPrefix string `json:"ipAddressPrefix"`
	// Name of the route.
	// Names can only contain alphanumeric, underscore, dash, and period characters. Case-sensitive
	// Required
	Name string `json:"name"`
	// Name of the virtual NIC set to route matching packets to.
	// Routed flows are load-balanced among all the virtual NICs in the virtual NIC set
	// Required
	NextHopVnicSet string `json:"nextHopVnicSet"`
	// Slice of tags to be associated with the route
	// Optional
	Tags []string `json:"tags"`
}

// UpdateRoute updates the specified route
func (c *RoutesClient) UpdateRoute(input *UpdateRouteInput) (*RouteInfo, error) {
	input.Name = c.getQualifiedName(input.Name)
	input.NextHopVnicSet = c.getQualifiedName(input.NextHopVnicSet)

	var routeInfo RouteInfo
	if err := c.updateResource(input.Name, &input, &routeInfo); err != nil {
		return nil, err
	}

	return c.success(&routeInfo)
}

// DeleteRouteInput details the route to delete
type DeleteRouteInput struct {
	// Name of the Route to delete. Case-sensitive
	// Required
	Name string `json:"name"`
}

// DeleteRoute deletes the specified route
func (c *RoutesClient) DeleteRoute(input *DeleteRouteInput) error {
	return c.deleteResource(input.Name)
}

func (c *RoutesClient) success(info *RouteInfo) (*RouteInfo, error) {
	info.Name = c.getUnqualifiedName(info.FQDN)
	c.unqualify(&info.NextHopVnicSet)
	return info, nil
}
