package lbaas

import (
	"time"
)

var (
	serviceInstanceContainerPath = "/vlbrs"
	serviceInstanceResourcePath  = "/vlbrs/%s/%s/"
)

const waitForServiceInstanceReadyPollInterval = 10 * time.Second
const waitForServiceInstanceReadyTimeout = 300 * time.Second
const waitForServiceInstanceDeletePollInterval = 10 * time.Second
const waitForServiceInstanceDeleteTimeout = 300 * time.Second

// LoadBalancerClient is a client for the Load Balancer service instance.
type LoadBalancerClient struct {
	ResourceClient
	PollInterval time.Duration
	Timeout      time.Duration
}

// LoadBalancerClient returns an ServiceInstanceClient which is used to access the
// Load Balancer API
func (c *Client) LoadBalancerClient() *LoadBalancerClient {
	return &LoadBalancerClient{
		ResourceClient: ResourceClient{
			Client:           c,
			ContainerPath:    serviceInstanceContainerPath,
			ResourceRootPath: serviceInstanceResourcePath,
		},
	}
}

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
	Region                   string                         `json:"region"`
	RestURIs                 []RestURIInfo                  `json:"rest_uri"`
	Scheme                   LoadBalancerScheme             `json:"scheme"`
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
	Enabled            string `json:"enabled"`
	HealthyThreshold   int    `json:"healthy_threshold"`
	Interval           int    `json:"interval"`
	Path               string `json:"path"`
	Timeout            int    `json:"timeout"`
	Type               string `json:"type"`
	UnhealthyThreshold int    `json:"unhealthy_threshold"`
}

type ListenerInfo struct {
	BalancerProtocol     string `json:"balancer_protocol"`
	Disabled             string `json:"disabled"`
	EffectiveState       string `json:"effective_state"`
	Name                 string `json:"name"`
	OriginServerProtocol string `json:"origin_server_protocol"`
	Port                 string `json:"port"`
	URI                  string `json:"uri"`
}

type RestURIInfo struct {
	Type string `json:"type"`
	URI  string `json:"uri"`
}

// CreateLoadBalancerInput specifies the create request for a load balancer service instance
type CreateLoadBalancerInput struct {
	Description        string               `json:"description,omitempty"`
	Disabled           LoadBalancerDisabled `json:"disabled"`
	Name               string               `json:"name"`
	Region             string               `json:"region"`
	Scheme             LoadBalancerScheme   `json:"scheme"`
	ParentLoadBalancer string               `json:"parent_vlbr,omitempty"`
}

// GetLoadBalancerInput request attributes required to Get a Load Balancer instance
type GetLoadBalancerInput struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

// DeleteLoadBalancerInput request attributes to required to Delete a Load Balancer instance
type DeleteLoadBalancerInput struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

// CreateLoadBalancer creates a new Load Balancer instance
func (c *LoadBalancerClient) CreateLoadBalancer(input *CreateLoadBalancerInput) (*LoadBalancerInfo, error) {
	var info LoadBalancerInfo
	if err := c.createResource(input.Region, input.Name, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// DeleteLoadBalancer deletes the service instance with the specified input
func (c *LoadBalancerClient) DeleteLoadBalancer(input *DeleteLoadBalancerInput) (*LoadBalancerInfo, error) {
	var info LoadBalancerInfo
	if err := c.deleteResource(input.Region, input.Name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetLoadBalancer fetchs the instance details of the Load Balancer
func (c *LoadBalancerClient) GetLoadBalancer(input *GetLoadBalancerInput) (*LoadBalancerInfo, error) {
	var info LoadBalancerInfo
	if err := c.getResource(input.Region, input.Name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
