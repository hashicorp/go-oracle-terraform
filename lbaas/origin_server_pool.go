package lbaas

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const waitForOriginServerPoolReadyPollInterval = 10 * time.Second  // 10 seconds
const waitForOriginServerPoolReadyTimeout = 30 * time.Minute       // 30 minutes
const waitForOriginServerPoolDeletePollInterval = 10 * time.Second // 10 seconds
const waitForOriginServerPoolDeleteTimeout = 30 * time.Minute      // 30 minutes

var (
	originserverpoolContainerPath = "/vlbrs/%s/%s/originserverpools"
	originserverpoolResourcePath  = "/vlbrs/%s/%s/originserverpools/%s"
)

// OriginServerPoolClient is a client for the Load Balancer Origin Server Pool resources.
type OriginServerPoolClient struct {
	LBaaSResourceClient
}

// OriginServerPoolClient returns an Client which is used to access the
// Load Balancer Origin Server Pool API
func (c *Client) OriginServerPoolClient() *OriginServerPoolClient {
	OriginServerPoolClient := &OriginServerPoolClient{
		LBaaSResourceClient: LBaaSResourceClient{
			Client:           c,
			ContainerPath:    originserverpoolContainerPath,
			ResourceRootPath: originserverpoolResourcePath,
			Accept:           CONTENT_TYPE_ORIGIN_SERVER_POOL_JSON,
			ContentType:      CONTENT_TYPE_ORIGIN_SERVER_POOL_JSON,
		},
	}

	return OriginServerPoolClient
}

type OriginServerInfo struct {
	Hostname string        `json:"hostname"`
	Port     int           `json:"port"`
	Status   LBaaSDisabled `json:"status"`
}

type CreateOriginServerInput struct {
	Status   LBaaSStatus `json:"status"`
	Hostname string      `json:"hostname"`
	Port     int         `json:"port"`
}

type OriginServerPoolInfo struct {
	Consumers          string             `json:"consumers"`
	HealthCheck        HealthCheckInfo    `json:"health_check"`
	Name               string             `json:"name"`
	OperationDetails   string             `json:"operation_details"`
	OriginServers      []OriginServerInfo `json:"origin_servers"`
	ReasonForDisabling string             `json:"reason_for_disabling"`
	State              LBaaSState         `json:"state"`
	Status             LBaaSStatus        `json:"status"`
	Tags               []string           `json:"tags"`
	URI                string             `json:"uri"`
	VnicSetName        string             `json:"vnic_set_name"`
}

type CreateOriginServerPoolInput struct {
	Name          string                    `json:"name"`
	OriginServers []CreateOriginServerInput `json:"origin_servers,omitempty"`
	Status        LBaaSStatus               `json:"status,omitempty"`
	Tags          []string                  `json:"tags,omitempty"`
	VnicSetName   string                    `json:"vnic_set_name,omitempty"`
}

type UpdateOriginServerPoolInput struct {
	Name          string                    `json:"name"`
	OriginServers []CreateOriginServerInput `json:"origin_servers,omitempty"`
	Status        LBaaSStatus               `json:"status,omitempty"`
	Tags          []string                  `json:"tags,omitempty"`
	VnicSetName   string                    `json:"vnic_set_name,omitempty"`
}

// CreateOriginServerPool creates a new server pool
func (c *OriginServerPoolClient) CreateOriginServerPool(lb LoadBalancerContext, input *CreateOriginServerPoolInput) (*OriginServerPoolInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForOriginServerPoolReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForOriginServerPoolReadyTimeout
	}

	var info OriginServerPoolInfo
	if err := c.createResource(lb.Region, lb.Name, &input, &info); err != nil {
		return nil, err
	}

	createdStates := []LBaaSState{LBaaSStateCreated, LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateCreationFailed, LBaaSStateDeletionInProgress, LBaaSStateDeleted, LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkOriginServerPoolState(&info, createdStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return &info, nil
	}
	// else poll till ready
	err = c.WaitForOriginServerPoolState(lb, input.Name, createdStates, erroredStates, c.PollInterval, c.Timeout, &info)
	return &info, err
}

// DeleteOriginServerPool deletes the server pool with the specified input
func (c *OriginServerPoolClient) DeleteOriginServerPool(lb LoadBalancerContext, name string) (*OriginServerPoolInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForOriginServerPoolDeletePollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForOriginServerPoolDeleteTimeout
	}

	var info OriginServerPoolInfo
	if err := c.deleteResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}

	deletedStates := []LBaaSState{LBaaSStateDeleted}
	erroredStates := []LBaaSState{LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	deleted, err := c.checkOriginServerPoolState(&info, deletedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if deleted {
		return &info, nil
	}
	// else poll till deleted
	err = c.WaitForOriginServerPoolState(lb, name, deletedStates, erroredStates, c.PollInterval, c.Timeout, &info)
	if err != nil && client.WasNotFoundError(err) {
		// resource could not be found, thus deleted
		return nil, nil
	}

	return &info, err
}

// GetOriginServerPool fetchs the server pool details
func (c *OriginServerPoolClient) GetOriginServerPool(lb LoadBalancerContext, name string) (*OriginServerPoolInfo, error) {

	var info OriginServerPoolInfo
	if err := c.getResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// UpdateOriginServerPool fetchs the server pool details
func (c *OriginServerPoolClient) UpdateOriginServerPool(lb LoadBalancerContext, name string, input *UpdateOriginServerPoolInput) (*OriginServerPoolInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForOriginServerPoolReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForOriginServerPoolReadyTimeout
	}

	var info OriginServerPoolInfo
	if err := c.updateOriginServerPool(lb.Region, lb.Name, name, &input, &info); err != nil {
		return nil, err
	}

	updatedStates := []LBaaSState{LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateModificaitonFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkOriginServerPoolState(&info, updatedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return &info, nil
	}
	// else poll till ready
	err = c.WaitForOriginServerPoolState(lb, name, updatedStates, erroredStates, c.PollInterval, c.Timeout, &info)
	return &info, err
}

// WaitForOriginServerPoolState waits for the resource to be in one of a set of desired states
func (c *OriginServerPoolClient) WaitForOriginServerPoolState(lb LoadBalancerContext, name string, desiredStates, errorStates []LBaaSState, pollInterval, timeoutSeconds time.Duration, info *OriginServerPoolInfo) error {

	var getErr error
	err := c.client.WaitFor("Origin Server Pool status update", pollInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetOriginServerPool(lb, name)
		if getErr != nil {
			return false, getErr
		}

		return c.checkOriginServerPoolState(info, desiredStates, errorStates)
	})
	return err
}

// check the State, returns in desired state (true), not ready yet (false) or errored state (error)
func (c *OriginServerPoolClient) checkOriginServerPoolState(info *OriginServerPoolInfo, desiredStates, errorStates []LBaaSState) (bool, error) {

	c.client.DebugLogString(fmt.Sprintf("Origin Server Pool %v state is %v", info.Name, info.State))

	state := LBaaSState(info.State)

	if isStateInLBaaSStates(state, desiredStates) {
		// we're good, return okay
		return true, nil
	}
	if isStateInLBaaSStates(state, errorStates) {
		// not good, return error
		return false, fmt.Errorf("Origin Server Pool %v in errored state %v", info.Name, info.State)
	}
	// not ready lifecycleTimeout
	return false, nil
}
