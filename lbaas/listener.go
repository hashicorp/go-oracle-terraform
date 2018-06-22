package lbaas

var (
	listenerContainerPath = "/vlbrs/%s/%s/listeners"
	listenerResourcePath  = "/vlbrs/%s/%s/listeners/%s"
)

// ListenerClient is a client for the Load Balancer Listener resources.
type ListenerClient struct {
	LBaaSResourceClient
}

// ListenerClient returns an ListenerClient which is used to access the
// Load Balancer Listener API
func (c *Client) ListenerClient() *ListenerClient {
	c.ContentType = CONTENT_TYPE_LISTENER_JSON
	c.Accept = CONTENT_TYPE_LISTENER_JSON
	return &ListenerClient{
		LBaaSResourceClient: LBaaSResourceClient{
			Client:           c,
			ContainerPath:    listenerContainerPath,
			ResourceRootPath: listenerResourcePath,
		},
	}
}

type Protocol string

const (
	ProtocolHTTP  Protocol = "HTTP"
	ProtocolHTTPS Protocol = "HTTPS"
)

type OriginServerSourceInheritedFrom string

const (
	OriginServerSourceInheritedFromSelf   OriginServerSourceInheritedFrom = "SELF"
	OriginServerSourceInheritedFromVLBR   OriginServerSourceInheritedFrom = "VLBR"
	OriginServerSourceInheritedFromParent OriginServerSourceInheritedFrom = "PARENT_LISTENER"
)

type EffectiveOriginServersInfo struct {
	OperationDetails                string                          `json:"operation_details"`
	OriginServerPool                string                          `json:"origin_server_pool"`
	OriginServerSourceInheritedFrom OriginServerSourceInheritedFrom `json:"origin_server_source_inherited_from"`
}

type ListenerInfo struct {
	BalancerProtocol       Protocol                   `json:"balancer_protocol"`
	Disabled               LoadBalancerDisabled       `json:"disabled"`
	EffectiveOriginServers EffectiveOriginServersInfo `json:"effective_origin_servers"`
	EffectiveState         LoadBalancerEffectiveState `json:"effective_state"`
	InlinePolicies         []string                   `json:"inline_policies"`
	Name                   string                     `json:"name"`
	OperationDetails       string                     `json:"operation_details"`
	OriginServerPool       string                     `json:"origin_server_pool"`
	OriginServerProtocol   Protocol                   `json:"origin_server_protocol"`
	ParentListener         string                     `json:"parent_listener"`
	PathPrefixes           []string                   `json:"path_prefixes"`
	Policies               []string                   `json:"policies"`
	Port                   int                        `json:"port"`
	SSLCerts               []string                   `json:"ssl_cert"`
	State                  LoadBalancerState          `json:"state"`
	Tags                   []string                   `json:"tags"`
	URI                    string                     `json:"uri"`
	VirtualHosts           []string                   `json:"virtual_hosts"`
}

type CreateListenerInput struct {
	BalancerProtocol     Protocol             `json:"balancer_protocol"`
	Disabled             LoadBalancerDisabled `json:"disabled,omitempty"`
	Name                 string               `json:"name"`
	OriginServerPool     string               `json:"origin_server_pool,omitempty"`
	OriginServerProtocol Protocol             `json:"origin_server_protocol"`
	PathPrefixes         []string             `json:"path_prefixes,omitempty"`
	Policies             []string             `json:"policies,omitempty"`
	Port                 int                  `json:"port"`
	SSLCerts             []string             `json:"ssl_cert,omitempty"`
	Tags                 []string             `json:"tags,omitempty"`
	VirtualHosts         []string             `json:"virtual_hosts,omitempty"`
}

type UpdateListenerInput struct {
	BalancerProtocol     Protocol             `json:"balancer_protocol,omitempty"`
	Disabled             LoadBalancerDisabled `json:"disabled,omitempty"`
	OriginServerPool     string               `json:"origin_server_pool,omitempty"`
	OriginServerProtocol Protocol             `json:"origin_server_protocol,omitempty"`
	PathPrefixes         []string             `json:"path_prefixes,omitempty"`
	Policies             []string             `json:"policies,omitempty"`
	SSLCerts             []string             `json:"ssl_cert,omitempty"`
	Tags                 []string             `json:"tags,omitempty"`
	VirtualHosts         []string             `json:"virtual_hosts,omitempty"`
}

// CreateListener creates a new listener
func (c *ListenerClient) CreateListener(lb LoadBalancerContext, input *CreateListenerInput) (*ListenerInfo, error) {
	var info ListenerInfo
	if err := c.createResource(lb.Region, lb.Name, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// DeleteListener deletes the listener with the specified input
func (c *ListenerClient) DeleteListener(lb LoadBalancerContext, name string) (*ListenerInfo, error) {
	var info ListenerInfo
	if err := c.deleteResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetListener fetchs the listener details
func (c *ListenerClient) GetListener(lb LoadBalancerContext, name string) (*ListenerInfo, error) {
	// TODO
	// Query Parameters projection(optional): string
	// Projections can be specified when retrieving collection of resources as well
	// as when retrieving a specific resource. They are of four types : MINIMAL, CONSOL, FULL, and DETAILED.

	var info ListenerInfo
	if err := c.getResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetListener fetchs the listener details
func (c *ListenerClient) UpdateListener(lb LoadBalancerContext, name string, input *UpdateListenerInput) (*ListenerInfo, error) {
	var info ListenerInfo
	if err := c.updateResource(lb.Region, lb.Name, name, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
