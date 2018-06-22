package lbaas

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

var (
	listenerContainerPath = "/vlbrs/%s/%s/listeners"
	listenerResourcePath  = "/vlbrs/%s/%s/listeners/%s"
)

const waitForListenerReadyPollInterval = 10 * time.Second  // 10 seconds
const waitForListenerReadyTimeout = 10 * time.Minute       // 10 minutes
const waitForListenerDeletePollInterval = 10 * time.Second // 10 seconds
const waitForListenerDeleteTimeout = 10 * time.Minute      // 10 minutes

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
	Disabled               LBaaSDisabled              `json:"disabled"`
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
	State                  LBaaSState                 `json:"state"`
	Tags                   []string                   `json:"tags"`
	URI                    string                     `json:"uri"`
	VirtualHosts           []string                   `json:"virtual_hosts"`
}

type CreateListenerInput struct {
	BalancerProtocol     Protocol      `json:"balancer_protocol"`
	Disabled             LBaaSDisabled `json:"disabled,omitempty"`
	Name                 string        `json:"name"`
	OriginServerPool     string        `json:"origin_server_pool,omitempty"`
	OriginServerProtocol Protocol      `json:"origin_server_protocol"`
	PathPrefixes         []string      `json:"path_prefixes,omitempty"`
	Policies             []string      `json:"policies,omitempty"`
	Port                 int           `json:"port"`
	SSLCerts             []string      `json:"ssl_cert,omitempty"`
	Tags                 []string      `json:"tags,omitempty"`
	VirtualHosts         []string      `json:"virtual_hosts,omitempty"`
}

type UpdateListenerInput struct {
	BalancerProtocol     Protocol      `json:"balancer_protocol,omitempty"`
	Disabled             LBaaSDisabled `json:"disabled,omitempty"`
	OriginServerPool     string        `json:"origin_server_pool,omitempty"`
	OriginServerProtocol Protocol      `json:"origin_server_protocol,omitempty"`
	PathPrefixes         []string      `json:"path_prefixes,omitempty"`
	Policies             []string      `json:"policies,omitempty"`
	SSLCerts             []string      `json:"ssl_cert,omitempty"`
	Tags                 []string      `json:"tags,omitempty"`
	VirtualHosts         []string      `json:"virtual_hosts,omitempty"`
}

// CreateListener creates a new listener
func (c *ListenerClient) CreateListener(lb LoadBalancerContext, input *CreateListenerInput) (*ListenerInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForListenerReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForListenerReadyTimeout
	}

	var info ListenerInfo
	if err := c.createResource(lb.Region, lb.Name, &input, &info); err != nil {
		return nil, err
	}

	createdStates := []LBaaSState{LBaaSStateCreationInProgress, LBaaSStateCreated, LBaaSStateHealthy}
	// createdStates := []LBaaSState{LBaaSStateCreated, LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateCreationFailed, LBaaSStateDeletionInProgress, LBaaSStateDeleted, LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkListenerState(&info, createdStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return &info, nil
	}
	// else poll till ready
	err = c.WaitForListenerState(lb, input.Name, createdStates, erroredStates, &info)
	return &info, err
}

// DeleteListener deletes the listener with the specified input
func (c *ListenerClient) DeleteListener(lb LoadBalancerContext, name string) (*ListenerInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForListenerDeletePollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForListenerDeleteTimeout
	}

	var info ListenerInfo
	if err := c.deleteResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}

	deletedStates := []LBaaSState{LBaaSStateDeletionInProgress, LBaaSStateDeleted}
	// deletedStates := []LBaaSState{LBaaSStateDeleted}
	erroredStates := []LBaaSState{LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	deleted, err := c.checkListenerState(&info, deletedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if deleted {
		return &info, nil
	}
	// else poll till deleted
	err = c.WaitForListenerState(lb, name, deletedStates, erroredStates, &info)
	if err != nil && client.WasNotFoundError(err) {
		// resource could not be found, thus deleted
		return nil, nil
	}
	return &info, err
}

// GetListener fetchs the listener details
func (c *ListenerClient) GetListener(lb LoadBalancerContext, name string) (*ListenerInfo, error) {
	var info ListenerInfo
	if err := c.getResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// UpdateListener updated the listener
func (c *ListenerClient) UpdateListener(lb LoadBalancerContext, name string, input *UpdateListenerInput) (*ListenerInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForListenerReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForListenerReadyTimeout
	}

	var info ListenerInfo
	if err := c.updateResource(lb.Region, lb.Name, name, &input, &info); err != nil {
		return nil, err
	}

	// updatedStates := []LBaaSState{LBaaSStateModificationInProgress, LBaaSStateHealthy}
	updatedStates := []LBaaSState{LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateModificaitonFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkListenerState(&info, updatedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return &info, nil
	}
	// else poll till ready
	err = c.WaitForListenerState(lb, name, updatedStates, erroredStates, &info)
	return &info, err
}

// WaitForListenerState waits for the resource to be in one of a set of desired states
func (c *ListenerClient) WaitForListenerState(lb LoadBalancerContext, name string, desiredStates, errorStates []LBaaSState, info *ListenerInfo) error {

	var getErr error
	err := c.client.WaitFor("Listener status update", c.PollInterval, c.Timeout, func() (bool, error) {
		info, getErr = c.GetListener(lb, name)
		if getErr != nil {
			return false, getErr
		}

		return c.checkListenerState(info, desiredStates, errorStates)
	})
	return err
}

// check the State, returns in desired state (true), not ready yet (false) or errored state (error)
func (c *ListenerClient) checkListenerState(info *ListenerInfo, desiredStates, errorStates []LBaaSState) (bool, error) {

	c.client.DebugLogString(fmt.Sprintf("Listener %v state is %v", info.Name, info.State))

	state := LBaaSState(info.State)

	if isStateInLBaaSStates(state, desiredStates) {
		// we're good, return okay
		return true, nil
	}
	if isStateInLBaaSStates(state, errorStates) {
		// not good, return error
		return false, fmt.Errorf("Listener %v in errored state %v", info.Name, info.State)
	}
	// not ready lifecycleTimeout
	return false, nil
}
