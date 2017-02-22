package compute

// SecurityApplicationsClient is a client for the Security Application functions of the Compute API.
type SecurityApplicationsClient struct {
	ResourceClient
}

// SecurityApplications obtains a SecurityApplicationsClient which can be used to access to the
// Security Application functions of the Compute API
func (c *Client) SecurityApplications() *SecurityApplicationsClient {
	return &SecurityApplicationsClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "security application",
			ContainerPath:       "/secapplication/",
			ResourceRootPath:    "/secapplication",
		}}
}

// SecurityApplicationInfo describes an existing security application.
type SecurityApplicationInfo struct {
	Name        string                      `json:"name"`
	Protocol    SecurityApplicationProtocol `json:"protocol"`
	DPort       string                      `json:"dport"`
	ICMPType    SecurityApplicationICMPType `json:"icmptype"`
	ICMPCode    SecurityApplicationICMPCode `json:"icmpcode"`
	Description string                      `json:"description"`
	URI         string                      `json:"uri"`
}

type SecurityApplicationProtocol string

const (
	All    SecurityApplicationProtocol = "all"
	AH     SecurityApplicationProtocol = "ah"
	ESP    SecurityApplicationProtocol = "esp"
	ICMP   SecurityApplicationProtocol = "icmp"
	ICMPV6 SecurityApplicationProtocol = "icmpv6"
	IGMP   SecurityApplicationProtocol = "igmp"
	IPIP   SecurityApplicationProtocol = "ipip"
	GRE    SecurityApplicationProtocol = "gre"
	MPLSIP SecurityApplicationProtocol = "mplsip"
	OSPF   SecurityApplicationProtocol = "ospf"
	PIM    SecurityApplicationProtocol = "pim"
	RDP    SecurityApplicationProtocol = "rdp"
	SCTP   SecurityApplicationProtocol = "sctp"
	TCP    SecurityApplicationProtocol = "tcp"
	UDP    SecurityApplicationProtocol = "udp"
)

type SecurityApplicationICMPCode string

const (
	Admin    SecurityApplicationICMPCode = "admin"
	Df       SecurityApplicationICMPCode = "df"
	Host     SecurityApplicationICMPCode = "host"
	Network  SecurityApplicationICMPCode = "network"
	Port     SecurityApplicationICMPCode = "port"
	Protocol SecurityApplicationICMPCode = "protocol"
)

type SecurityApplicationICMPType string

const (
	Echo        SecurityApplicationICMPType = "echo"
	Reply       SecurityApplicationICMPType = "reply"
	TTL         SecurityApplicationICMPType = "ttl"
	TraceRoute  SecurityApplicationICMPType = "traceroute"
	Unreachable SecurityApplicationICMPType = "unreachable"
)

func (c *SecurityApplicationsClient) success(result *SecurityApplicationInfo) (*SecurityApplicationInfo, error) {
	c.unqualify(&result.Name)
	return result, nil
}

// CreateSecurityApplicationInput describes the Security Application to create
type CreateSecurityApplicationInput struct {
	Name        string                      `json:"name"`
	Protocol    SecurityApplicationProtocol `json:"protocol"`
	DPort       string                      `json:"dport"`
	ICMPCode    SecurityApplicationICMPCode `json:"icmpcode,omitempty"`
	ICMPType    SecurityApplicationICMPType `json:"icmptype,omitempty"`
	Description string                      `json:"description"`
}

// CreateSecurityApplication creates a new security application.
func (c *SecurityApplicationsClient) CreateSecurityApplication(input *CreateSecurityApplicationInput) (*SecurityApplicationInfo, error) {
	input.Name = c.getQualifiedName(input.Name)

	var appInfo SecurityApplicationInfo
	if err := c.createResource(&input, &appInfo); err != nil {
		return nil, err
	}

	return c.success(&appInfo)
}

// GetSecurityApplicationInput describes the Security Application to obtain
type GetSecurityApplicationInput struct {
	Name string `json:"name"`
}

// GetSecurityApplication retrieves the security application with the given name.
func (c *SecurityApplicationsClient) GetSecurityApplication(input *GetSecurityApplicationInput) (*SecurityApplicationInfo, error) {
	var appInfo SecurityApplicationInfo
	if err := c.getResource(input.Name, &appInfo); err != nil {
		return nil, err
	}

	return c.success(&appInfo)
}

// DeleteSecurityApplicationInput  describes the Security Application to delete
type DeleteSecurityApplicationInput struct {
	Name string `json:"name"`
}

// DeleteSecurityApplication deletes the security application with the given name.
func (c *SecurityApplicationsClient) DeleteSecurityApplication(input *DeleteSecurityApplicationInput) error {
	return c.deleteResource(input.Name)
}
