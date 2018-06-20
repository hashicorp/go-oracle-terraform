package lbaas

// LoadBalancerScheme Scheme types
type LoadBalancerScheme string

const (
	LoadBalancerSchemeInternetFacing LoadBalancerScheme = "INTERNET_FACING"
	LoadBalancerSchemeInternal       LoadBalancerScheme = "INTERNAL"
)

// LoadBalancerDisabled
type LoadBalancerDisabled string

const (
	LoadBalancerDisabledTrue        LoadBalancerDisabled = "TRUE"
	LoadBalancerDisabledFalse       LoadBalancerDisabled = "FALSE"
	LoadBalancerDisabledMaintenance LoadBalancerDisabled = "MAINTENANCE_MODE"
)

// LoadBalancerEffectiveState
type LoadBalancerEffectiveState string

const (
	LoadBalancerEffectiveStateTrue        LoadBalancerEffectiveState = "TRUE"
	LoadBalancerEffectiveStateFalse       LoadBalancerEffectiveState = "FALSE"
	LoadBalancerEffectiveStateMaintenance LoadBalancerEffectiveState = "MAINTENANCE_MODE"
)

// LoadBalancerState
type LoadBalancerState string

const (
	LoadBalancerStateCreationInProgress     LoadBalancerState = "CREATION_IN_PROGRESS"
	LoadBalancerStateCreated                LoadBalancerState = "CREATED"
	LoadBalancerStateHealthy                LoadBalancerState = "HEALTHY"
	LoadBalancerStateInterventionNeeded     LoadBalancerState = "ADMINISTRATOR_INTERVENTION_NEEDED"
	LoadBalancerStateDeletionInProgress     LoadBalancerState = "DELETION_IN_PROGRESS"
	LoadBalancerStateDeleted                LoadBalancerState = "DELETED"
	LoadBalancerStateModificationInProgress LoadBalancerState = "MODIFICATION_IN_PROGRESS"
	LoadBalancerStateCreationFailed         LoadBalancerState = "CREATION_FAILED"
	LoadBalancerStateModificationFailed     LoadBalancerState = "MODIFICATION_FAILED"
	LoadBalancerStateDeletionFailed         LoadBalancerState = "DELETION_FAILED"
	LoadBalancerStateAccessDenied           LoadBalancerState = "ACCESS_DENIED"
	LoadBalancerStateAbandon                LoadBalancerState = "ABANDON"
	LoadBalancerStatePause                  LoadBalancerState = "PAUSE"
	LoadBalancerStateForcePaused            LoadBalancerState = "FORCE_PAUSED"
	LoadBalancerStateResume                 LoadBalancerState = "RESUME"
)

// HttpMethods
type HttpMethod string

const (
	HttpCOPY      HttpMethod = "COPY"
	HttpDELETE    HttpMethod = "DELETE"
	HttpGET       HttpMethod = "GET"
	HttpHEAD      HttpMethod = "HEAD"
	HttpLOCK      HttpMethod = "LOCK"
	HttpMKCOL     HttpMethod = "MKCOL"
	HttpMOVE      HttpMethod = "MOVE"
	HttpOPTIONS   HttpMethod = "OPTIONS"
	HttpPATCH     HttpMethod = "PATCH"
	HttpPOST      HttpMethod = "POST"
	HttpPROPFIND  HttpMethod = "PROPFIND"
	HttpPROPPATCH HttpMethod = "PROPPATCH"
	HttpPUT       HttpMethod = "PUT"
	HttpUNLOCK    HttpMethod = "UNLOCK"
)

// LoadBalancerInfo specifies the Load Balancer obtained from a GET request
type LoadBalancerInfo struct {
	BalancerVIPs             []string                       `json:"balancer_vips"`
	CanonicalHostName        string                         `json:"canonical_host_name"`
	CloudgateCapable         string                         `json:"cloudgate_capable"`
	ComputeSecurityArtifacts []ComputeSecurityArtifactsInfo `json:"compute_security_artifacts"`
	ComputeSite              string                         `json:"compute_site"`
	CreatedOn                string                         `json:"created_on"`
	Description              string                         `json:"description"`
	Disabled                 LoadBalancerDisabled           `json:"disabled"`
	DisplayName              string                         `json:"display_name"`
	HealthCheck              HealthCheckInfo                `json:"health_check"`
	IPNetworkName            string                         `json:"ip_network_name"`
	IsDisabledEffectively    string                         `json:"is_disabled_effectively"`
	Listeners                []ListenerInfo                 `json:"listeners"`
	ModifiedOn               string                         `json:"modified_on"`
	Name                     string                         `json:"name"`
	Owner                    string                         `json:"owner"`
	PermittedMethods         []string                       `json:"permitted_methods"`
	Region                   string                         `json:"region"`
	RestURIs                 []RestURIInfo                  `json:"rest_uri"`
	Scheme                   LoadBalancerScheme             `json:"scheme"`
	ServerPool               string                         `json:"server_pool"`
	State                    LoadBalancerState              `json:"state"`
	Tags                     []string                       `json:"tags"`
	URI                      string                         `json:"uri"`
}

type ComputeSecurityArtifactsInfo struct {
	AddressType  string `json:"address_type"`
	ArtifactType string `json:"artifact_type"`
	URI          string `json:"uri"`
}

type HealthCheckInfo struct {
	AcceptedReturnCodes []string `json:accepted_return_codes`
	Enabled             string   `json:"enabled"`
	HealthyThreshold    int      `json:"healthy_threshold"`
	Interval            int      `json:"interval"`
	Path                string   `json:"path"`
	Timeout             int      `json:"timeout"`
	Type                string   `json:"type"`
	UnhealthyThreshold  int      `json:"unhealthy_threshold"`
}

type RestURIInfo struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

// CreateLoadBalancerInput specifies the create request for a load balancer service instance
type CreateLoadBalancerInput struct {
	Description        string               `json:"description,omitempty"`
	Disabled           LoadBalancerDisabled `json:"disabled"`
	IPNetworkName      string               `json:"ip_network_name,omitempty"`
	Name               string               `json:"name"`
	ParentLoadBalancer string               `json:"parent_vlbr,omitempty"`
	PermittedClients   []string             `json:"permitted_clients,omitempty"`
	PermittedMethods   []string             `json:"permitted_methods,omitempty"`
	Policies           []string             `json:"policies,omitempty"`
	Region             string               `json:"region"`
	Scheme             LoadBalancerScheme   `json:"scheme"`
	ServerPool         string               `json:"server_pool,omitempty"`
	Tags               []string             `json:"tags,omitempty"`
}

// UpdateLoadBalancerInput specifies the create request for a load balancer service instance
type UpdateLoadBalancerInput struct {
	Description        string               `json:"description,omitempty"`
	Disabled           LoadBalancerDisabled `json:"disabled,omitempty"`
	IPNetworkName      string               `json:"ip_network_name,omitempty"`
	Name               string               `json:"name,omitempty"`
	ParentLoadBalancer string               `json:"parent_vlbr,omitempty"`
	PermittedClients   []string             `json:"permitted_clients,omitempty"`
	PermittedMethods   []string             `json:"permitted_methods,omitempty"`
	Policies           []string             `json:"policies,omitempty"`
	ServerPool         string               `json:"server_pool,omitempty"`
	Tags               []string             `json:"tags,omitempty"`
}

// CreateLoadBalancer creates a new Load Balancer instance
func (c *LoadBalancerClient) CreateLoadBalancer(input *CreateLoadBalancerInput) (*LoadBalancerInfo, error) {
	var info LoadBalancerInfo
	if err := c.createResource(&input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// DeleteLoadBalancer deletes the service instance with the specified input
func (c *LoadBalancerClient) DeleteLoadBalancer(region, name string) (*LoadBalancerInfo, error) {
	var info LoadBalancerInfo
	if err := c.deleteResource(region, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetLoadBalancer fetchs the instance details of the Load Balancer
func (c *LoadBalancerClient) GetLoadBalancer(region, name string) (*LoadBalancerInfo, error) {
	var info LoadBalancerInfo
	if err := c.getResource(region, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// UpdateLoadBalancer fetchs the instance details of the Load Balancer
func (c *LoadBalancerClient) UpdateLoadBalancer(region, name string, input *UpdateLoadBalancerInput) (*LoadBalancerInfo, error) {
	var info LoadBalancerInfo
	if err := c.updateResource(region, name, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
