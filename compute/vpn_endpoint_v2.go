package compute

import (
	"fmt"
)

const (
	VPNEndpointV2Description   = "vpn endpoint v2"
	VPNEndpointV2ContainerPath = "/vpnendpoint/v2/"
	VPNEndpointV2ResourcePath  = "/vpnendpoint/v2"
)

type VPNEndpointV2sClient struct {
	ResourceClient
}

// VPNEndpointV2s() returns an VPNEndpointV2sClient that can be used to access the
// necessary CRUD functions for VPN Endpoint V2s.
func (c *Client) VPNEndpointV2s() *VPNEndpointV2sClient {
	return &VPNEndpointV2sClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: VPNEndpointV2Description,
			ContainerPath:       VPNEndpointV2ContainerPath,
			ResourceRootPath:    VPNEndpointV2ResourcePath,
		},
	}
}

// VPNEndpointV2Info contains the exported fields necessary to hold all the information about an
// VPN Endpoint V2
type VPNEndpointV2Info struct {
	// IP address of the VPN gateway in your data center through which you want
	// to connect to the Oracle Cloud VPN gateway.
	CustomerVPNGateway string `json:"customer_vpn_gateway"`
	// The Internet Key Exchange (IKE) ID that you have specified. The default
	// value is the public IP address of the cloud gateway.
	IKEIdentifier string `json:"ikeIdentifier"`
	// The name of the IP network on which the cloud gateway is created by VPNaaS.
	IPNetwork string `json:"ipNetwork"`
	// The name of the VPN Endpoint V2
	Name string `json:"name"`
	// Indicates whether Perfect Forward Secrecy (PFS) is required and your third-party device supports PFS.
	PFSFlag bool `json:"pfsFlag"`
	// Settings for Phase 1 of protocol (IKE).
	Phase1Settings Phase1Settings `json:"phase1Settings"`
	// Settings for Phase 2 of protocol (IPSEC).
	Phase2Settings Phase2Settings `json:"phase2Settings"`
	// The pre-shared VPN key.
	PSK string `json:"psk"`
	// List of routes (CIDR prefixes) that are reachable through this VPN tunnel.
	ReachableRoutes []string `json:"reachable_routes"`
	// Current status of the tunnel. The tunnel can be in one of the following states:
	// PENDING: indicates that your VPN connection is being set up.
	// UP: indicates that your VPN connection is established.
	// DOWN: indicates that your VPN connection is down.
	// ERROR: indicates that your VPN connection is in the error state.
	TunnelStatus string `json:"tunnelStatus"`
	// Uniform Resource Identifier for the VPN Endpoint V2
	Uri string `json:"uri"`
	// Comma-separated list of vNIC sets. Traffic is allowed to and
	// from these vNIC sets to the cloud gateway's vNIC set.
	VNICSets []string `json:"vnicSets"`
}

type Phase1Settings struct {
	// Encryption options for IKE. Permissible values are aes128, aes192, aes256.
	Encryption string `json:"encryption"`
	// Authentication options for IKE. Permissible values are sha1, sha2_256, and md5.
	Hash string `json:"hash"`
	// Diffie-Hellman group for both IKE and ESP. It is applicable for ESP only if PFS is enabled.
	// Permissible values are group5, group14, group22, group23, and group24.
	DHGroup string `json:"dhGroup"`
}

type Phase2Settings struct {
	// Encryption options for IKE. Permissible values are aes128, aes192, aes256.
	Encryption string `json:"encryption"`
	// Authentication options for IKE. Permissible values are sha1, sha2_256, and md5.
	Hash string `json:"hash"`
}

type CreateVPNEndpointV2Input struct {
	// Specify the IP address of the VPN gateway in your data center through which you want
	// to connect to the Oracle Cloud VPN gateway. Your gateway device must support route-based
	// VPN and IKE (Internet Key Exchange) configuration using pre-shared keys.
	// Required
	CustomerVPNGateway string `json:"customer_vpn_gateway"`
	// Description of the VPN Endpoint
	Description string `json:"description,omitempty"`
	// The Internet Key Exchange (IKE) ID. If you don't specify a value, the default value is
	// the public IP address of the cloud gateway. You can specify either an alternative IP address,
	// or any text string that you want to use as the IKE ID. If you specify a text string, you must
	// prefix the string with @. For example, if you want to specify the text IKEID-for-VPN1, specify
	// @IKEID-for-VPN1 as the value in request body. If you specify an IP address, don't prefix it with @.
	// The IKE ID is case sensitive and can contain a maximum of 255 ASCII alphanumeric characters
	// including special characters, period (.), hyphen (-), and underscore (_). The IKE ID can't contain
	// embedded space characters.
	// Note: If you specify the IKE ID, ensure that you specify the Peer ID type as Domain Name on the
	// third-party device in your data center. Other Peer ID types, such as email address, firewall
	// identifier or key identifier, aren't supported.
	// Optional
	IKEIdentifier string `json:"ikeIdentifier,omitempty"`
	// Specify the name of the IP network
	// which you want to create the cloud gateway. When you send a request to create a VPN connection,
	// a cloud gateway is created and this is assigned an available IP address from the IP network that
	// you specify. So, the cloud gateway is directly connected to the IP network that you specify.
	// You can only specify a single IP network. All other IP networks with are connected to the
	// specified IP network through an IP network exchange are discovered and added automatically to
	// the VPN connection.
	// Required
	IPNetwork string `json:"ipNetwork"`
	// The name of the VPN Endpoint V2 to create. Object names can only contain alphanumeric,
	// underscore, dash, and period characters. Names are case-sensitive.
	// Required
	Name string `json:"name"`
	// This is enabled (set to true) by default. If your third-party device supports Perfect Forward
	// Secrecy (PFS), set this parameter to true to require PFS.
	// Optional. Default true
	PFSFlag bool `json:"pfsFlag,omitmepty"`
	// Settings for Phase 1 of protocol (IKE).
	// Optional
	Phase1Settings *Phase1Settings `json:"phase1Settings,omitempty"`
	// Settings for Phase 2 of protocol (IPSEC).
	// Optional
	Phase2Settings *Phase2Settings `json:"phase2Settings,omitempty"`
	// Pre-shared VPN key. This secret key is shared between your network gateway and
	// the Oracle Cloud network for authentication. Specify the full path and name of
	// the text file that contains the pre-shared key. Ensure that the permission level
	// of the text file is set to 400. The pre-shared VPN key must not exceed 256 characters.
	// Required
	PSK string `json:"psk"`
	// Specify a list of routes (CIDR prefixes) that are reachable through this VPN tunnel.
	// You can specify a maximum of 20 IP subnet addresses. Specify IPv4 addresses in dot-decimal
	// notation with or without mask.
	// Required
	ReachableRoutes []string `json:"reachable_routes"`
	// An array of tags
	Tags []string `json:"tags"`
	// Comma-separated list of vNIC sets. Traffic is allowed to and from these vNIC sets to the
	// cloud gateway's vNIC set.
	// Required
	VNICSets []string `json:"vnicSets"`
}

// CreateVPNEndpointV2 creates a new VPN Endpoint V2 from an VPNEndpointV2sClient and an input struct.
// Returns a populated Info struct for the VPN Endpoint V2, and any errors
func (c *VPNEndpointV2sClient) CreateVPNEndpointV2(input *CreateVPNEndpointV2Input) (*VPNEndpointV2Info, error) {
	input.Name = c.getQualifiedACMEName(input.Name)
	input.IPNetwork = c.getQualifiedName(input.IPNetwork)
	input.VNICSets = c.getQualifiedList(input.VNICSets)

	return nil, fmt.Errorf("%+v", input)
	var vpnEndpointV2Info VPNEndpointV2Info
	if err := c.createResource(&input, &vpnEndpointV2Info); err != nil {
		return nil, err
	}

	return c.success(&vpnEndpointV2Info)
}

// GetVPNEndpointV2Input specifies the information needed to retrive a VPNEndpointV2
type GetVPNEndpointV2Input struct {
	// The name of the VPN Endpoint V2 to query for. Case-sensitive
	// Required
	Name string `json:"name"`
}

// GetVPNEndpointV2 returns a populated VPNEndpointV2Info struct from an input struct
func (c *VPNEndpointV2sClient) GetVPNEndpointV2(input *GetVPNEndpointV2Input) (*VPNEndpointV2Info, error) {
	input.Name = c.getQualifiedName(input.Name)

	var ipInfo VPNEndpointV2Info
	if err := c.getResource(input.Name, &ipInfo); err != nil {
		return nil, err
	}

	return c.success(&ipInfo)
}

// UpdateVPNEndpointV2Input defines what to update in a VPN Endpoint V2
// Only PSK and ReachableRoutes are updatable
type UpdateVPNEndpointV2Input struct {
	// Specify the IP address of the VPN gateway in your data center through which you want
	// to connect to the Oracle Cloud VPN gateway. Your gateway device must support route-based
	// VPN and IKE (Internet Key Exchange) configuration using pre-shared keys.
	// Required
	CustomerVPNGateway string `json:"customer_vpn_gateway"`
	// Description of the VPNGatewayV2
	Description string `json:"description,omitempty"`
	// The Internet Key Exchange (IKE) ID. If you don't specify a value, the default value is
	// the public IP address of the cloud gateway. You can specify either an alternative IP address,
	// or any text string that you want to use as the IKE ID. If you specify a text string, you must
	// prefix the string with @. For example, if you want to specify the text IKEID-for-VPN1, specify
	// @IKEID-for-VPN1 as the value in request body. If you specify an IP address, don't prefix it with @.
	// The IKE ID is case sensitive and can contain a maximum of 255 ASCII alphanumeric characters
	// including special characters, period (.), hyphen (-), and underscore (_). The IKE ID can't contain
	// embedded space characters.
	// Note: If you specify the IKE ID, ensure that you specify the Peer ID type as Domain Name on the
	// third-party device in your data center. Other Peer ID types, such as email address, firewall
	// identifier or key identifier, aren't supported.
	// Optional
	IKEIdentifier string `json:"ikeIdentifier"`
	// Specify the name of the IP network
	// which you want to create the cloud gateway. When you send a request to create a VPN connection,
	// a cloud gateway is created and this is assigned an available IP address from the IP network that
	// you specify. So, the cloud gateway is directly connected to the IP network that you specify.
	// You can only specify a single IP network. All other IP networks with are connected to the
	// specified IP network through an IP network exchange are discovered and added automatically to
	// the VPN connection.
	// Required
	IPNetwork string `json:"ipNetwork"`
	// The name of the VPN Endpoint V2 to create. Object names can only contain alphanumeric,
	// underscore, dash, and period characters. Names are case-sensitive.
	// Required
	Name string `json:"name"`
	// This is enabled (set to true) by default. If your third-party device supports Perfect Forward
	// Secrecy (PFS), set this parameter to true to require PFS.
	// Optional. Default true
	PFSFlag bool `json:"pfsFlag,omitempty"`
	// Settings for Phase 1 of protocol (IKE).
	// Optional
	Phase1Settings Phase1Settings `json:"phase1Settings,omitempty"`
	// Settings for Phase 2 of protocol (IPSEC).
	// Optional
	Phase2Settings Phase2Settings `json:"phase2Settings,omitempty"`
	// Pre-shared VPN key. This secret key is shared between your network gateway and
	// the Oracle Cloud network for authentication. Specify the full path and name of
	// the text file that contains the pre-shared key. Ensure that the permission level
	// of the text file is set to 400. The pre-shared VPN key must not exceed 256 characters.
	// Required.
	PSK string `json:"psk"`
	// Specify a list of routes (CIDR prefixes) that are reachable through this VPN tunnel.
	// You can specify a maximum of 20 IP subnet addresses. Specify IPv4 addresses in dot-decimal
	// notation with or without mask.
	// Required
	ReachableRoutes []string `json:"reachable_routes"`
	// Array of tags
	Tags []string `json:"tags,omitempty"`
	// Comma-separated list of vNIC sets. Traffic is allowed to and from these vNIC sets to the
	// cloud gateway's vNIC set.
	// Required
	VNICSets []string `json:"vnicSets"`
}

// UpdateVPNEndpointV2 update the VPN Endpoint V2
func (c *VPNEndpointV2sClient) UpdateVPNEndpointV2(updateInput *UpdateVPNEndpointV2Input) (*VPNEndpointV2Info, error) {
	updateInput.Name = c.getQualifiedACMEName(updateInput.Name)
	updateInput.IPNetwork = c.getQualifiedName(updateInput.IPNetwork)
	updateInput.VNICSets = c.getQualifiedList(updateInput.VNICSets)

	var ipInfo VPNEndpointV2Info
	if err := c.updateResource(updateInput.Name, updateInput, &ipInfo); err != nil {
		return nil, err
	}

	return c.success(&ipInfo)
}

type DeleteVPNEndpointV2Input struct {
	// The name of the VPN Endpoint V2 to query for. Case-sensitive
	// Required
	Name string `json:"name"`
}

func (c *VPNEndpointV2sClient) DeleteVPNEndpointV2(input *DeleteVPNEndpointV2Input) error {
	return c.deleteResource(input.Name)
}

// Unqualifies any qualified fields in the VPNEndpointV2Info struct
func (c *VPNEndpointV2sClient) success(info *VPNEndpointV2Info) (*VPNEndpointV2Info, error) {
	c.unqualify(&info.Name)
	c.unqualify(&info.IPNetwork)
	info.VNICSets = c.getUnqualifiedList(info.VNICSets)
	return info, nil
}
