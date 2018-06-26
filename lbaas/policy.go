package lbaas

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const waitForPolicyReadyPollInterval = 10 * time.Second  // 10 seconds
const waitForPolicyReadyTimeout = 30 * time.Minute       // 30 minutes
const waitForPolicyDeletePollInterval = 10 * time.Second // 10 seconds
const waitForPolicyDeleteTimeout = 30 * time.Minute      // 30 minutes

var (
	policyContainerPath = "/vlbrs/%s/%s/policies"
	policyResourcePath  = "/vlbrs/%s/%s/policies/%s"
)

// PolicyClient is a client for the Load Balancer Policy resources.
type PolicyClient struct {
	LBaaSResourceClient
}

// PolicyClient returns an PolicyClient which is used to access the
// Load Balancer Policy API
func (c *Client) PolicyClient() *PolicyClient {
	// c.ContentType cannot be generally set for the PolicyClient, instead it is set on each
	// Create or Update request based on the Type of the Policy being created/updated.
	// Accept all Policy Content Types
	c.Accept = strings.Join([]string{
		CONTENT_TYPE_APP_COOKIE_STICKINESS_POLICY_JSON,
		CONTENT_TYPE_CLOUDGATE_POLICY_JSON,
		CONTENT_TYPE_LB_COOKIE_STICKINESS_POLICY_JSON,
		CONTENT_TYPE_LOADBALANCING_MECHANISM_POLICY_JSON,
		CONTENT_TYPE_RATE_LIMITING_REQUEST_POLICY_JSON,
		CONTENT_TYPE_RESOURCE_ACCESS_CONTROL_POLICY_JSON,
		CONTENT_TYPE_REDIRECT_POLICY_JSON,
		CONTENT_TYPE_SSL_NEGOTIATION_POLICY_JSON,
		CONTENT_TYPE_SET_REQUEST_HEADER_POLICY_JSON,
		CONTENT_TYPE_TRUSTED_CERTIFICATE_POLICY_JSON,
	}, ",")

	return &PolicyClient{
		LBaaSResourceClient: LBaaSResourceClient{
			Client:           c,
			ContainerPath:    policyContainerPath,
			ResourceRootPath: policyResourcePath,
		},
	}
}

type PolicyInfo struct {
	Name  string     `json:"name"`
	State LBaaSState `json:"state"`
	Type  string     `json:"type"`
	URI   string     `json:"uri"`

	// ApplicationCookieStickinessPolicy
	AppCookieName string `json:"app_cookie_name"`

	// CloudGatePolicy
	CloudGateApplication                string `json:"cloudgate_application"`
	CloudGatePolicyName                 string `json:"cloudgate_policy_name"`
	IdentityServiceInstanceGuid         string `json:"identity_service_instance_guid"`
	VirtualHostnameForPolicyAttribution string `json:"virtual_hostname_for_policy_attribution"`

	// LoadBalancerCookieStickinessPolicy
	CookieExpirationPeriod int `json:"cookie_expiration_period"`

	// LoadBalancingMechanismPolicy
	LoadBalancingMechanism string `json:"load_balancing_mechanism"`

	// RateLimitingRequestPolicy
	BurstSize                   int    `json:"burst_size"`
	DoNotDelayExcessiveRequests bool   `json:"do_not_delay_excessive_requests"`
	HttpStatusErrorCode         int    `json:"http_status_error_code"`
	LogLevel                    string `json:"log_level"`
	RateLimitingCriteria        string `json:"rate_limiting_criteria"`
	RequestsPerSecond           int    `json:"requests_per_second"`
	StorageSize                 int    `json:"storage_size_in_mb"`
	Zone                        string `json:"zone"`

	// RedirectPolicy
	RedirectURI  string `json:"redirect_uri"`
	ResponseCode int    `json:"response_code"`

	// ResourceAccessControlPolicy
	Disposition      string   `json:"disposition"`
	DeniedClients    []string `json:"denied_clients"`
	PermittedClients []string `json:"permitted_clients"`

	// SetRequestHeaderPolicy
	HeaderName                 string   `json:"header_name"`
	Value                      string   `json:"value"`
	ActionWhenHeaderExists     string   `json:"action_when_hdr_exists"`
	ActionWhenHeaderValueIs    []string `json:"action_when_hdr_value_is"`
	ActionWhenHeaderValueIsNot []string `json:"action_when_hdr_value_is_not"`

	// SSLNegotiationPolicy
	Port                  int      `json:"port"`
	ServerOrderPreference string   `json:"server_order_preference"`
	SSLProtocol           []string `json:"ssl_protocol"`
	SSLCiphers            []string `json:"ssl_ciphers"`

	// TrustedCertificatePolicy
	TrustedCertificate string `json:"cert"`
}

type CreatePolicyInput struct {
	Name string `json:"name"`
	Type string `json:"type"`
	ApplicationCookieStickinessPolicyInfo
	CloudGatePolicyInfo
	LoadBalancerCookieStickinessPolicyInfo
	LoadBalancingMechanismPolicyInfo
	RateLimitingRequestPolicyInfo
	RedirectPolicyInfo
	ResourceAccessControlPolicyInfo
	SetRequestHeaderPolicyInfo
	SSLNegotiationPolicyInfo
	TrustedCertificatePolicyInfo
}

type UpdatePolicyInput struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	ApplicationCookieStickinessPolicyInfo
	CloudGatePolicyInfo
	LoadBalancerCookieStickinessPolicyInfo
	LoadBalancingMechanismPolicyInfo
	RateLimitingRequestPolicyInfo
	RedirectPolicyInfo
	ResourceAccessControlPolicyInfo
	SetRequestHeaderPolicyInfo
	SSLNegotiationPolicyInfo
	TrustedCertificatePolicyInfo
}

type ApplicationCookieStickinessPolicyInfo struct {
	AppCookieName string `json:"app_cookie_name,omitempty"`
}

type CloudGatePolicyInfo struct {
	CloudGateApplication                string `json:"cloudgate_application,omitempty"`
	CloudGatePolicyName                 string `json:"cloudgate_policy_name,omitempty"`
	IdentityServiceInstanceGuid         string `json:"identity_service_instance_guid,omitempty"`
	VirtualHostnameForPolicyAttribution string `json:"virtual_hostname_for_policy_attribution,omitempty"`
}

type LoadBalancerCookieStickinessPolicyInfo struct {
	CookieExpirationPeriod int `json:"cookie_expiration_period,omitempty"`
}

type LoadBalancingMechanismPolicyInfo struct {
	LoadBalancingMechanism string `json:"load_balancing_mechanism,omitempty"`
}

type RateLimitingRequestPolicyInfo struct {
	BurstSize                   int    `json:"burst_size,omitempty"`
	DoNotDelayExcessiveRequests bool   `json:"do_not_delay_excessive_requests,omitempty"`
	HttpStatusErrorCode         int    `json:"http_status_error_code,omitempty"`
	LogLevel                    string `json:"log_level,omitempty"`
	RateLimitingCriteria        string `json:"rate_limiting_criteria,omitempty"`
	RequestsPerSecond           int    `json:"requests_per_second,omitempty"`
	StorageSize                 int    `json:"storage_size_in_mb,omitempty"`
	Zone                        string `json:"zone,omitempty"`
}

type RedirectPolicyInfo struct {
	RedirectURI  string `json:"redirect_uri,omitempty"`
	ResponseCode int    `json:"response_code,omitempty"`
}

type ResourceAccessControlPolicyInfo struct {
	Disposition      string   `json:"disposition,omitempty"`
	DeniedClients    []string `json:"denied_clients,omitempty"`
	PermittedClients []string `json:"permitted_clients,omitempty"`
}

type SetRequestHeaderPolicyInfo struct {
	HeaderName                 string   `json:"header_name,omitempty"`
	Value                      string   `json:"value,omitempty"`
	ActionWhenHeaderExists     string   `json:"action_when_hdr_exists,omitempty"`
	ActionWhenHeaderValueIs    []string `json:"action_when_hdr_value_is,omitempty"`
	ActionWhenHeaderValueIsNot []string `json:"action_when_hdr_value_is_not,omitempty"`
}

type SSLNegotiationPolicyInfo struct {
	Port                  int      `json:"port,omitempty"`
	ServerOrderPreference string   `json:"server_order_preference,omitempty"`
	SSLProtocol           []string `json:"ssl_protocol,omitempty"`
	SSLCiphers            []string `json:"ssl_ciphers,omitempty"`
}

type TrustedCertificatePolicyInfo struct {
	TrustedCertificate string `json:"cert,omitempty"`
}

// CreatePolicy creates a new listener
func (c *PolicyClient) CreatePolicy(lb LoadBalancerContext, input *CreatePolicyInput) (*PolicyInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForPolicyReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForPolicyReadyTimeout
	}

	var info PolicyInfo

	// set the content type based on the policy type
	if input.Type == "" {
		return nil, fmt.Errorf("Policy type for %s is not set", input.Name)
	}
	c.ContentType = fmt.Sprintf("application/vnd.com.oracle.oracloud.lbaas.%s+json", input.Type)
	if err := c.createResource(lb.Region, lb.Name, &input, &info); err != nil {
		return nil, err
	}

	// createdStates := []LBaaSState{LBaaSStateCreationInProgress, LBaaSStateCreated, LBaaSStateHealthy}
	createdStates := []LBaaSState{LBaaSStateCreated, LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateCreationFailed, LBaaSStateDeletionInProgress, LBaaSStateDeleted, LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkPolicyState(&info, createdStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return &info, nil
	}
	// else poll till ready
	err = c.WaitForPolicyState(lb, input.Name, createdStates, erroredStates, c.PollInterval, c.Timeout, &info)
	return &info, err
}

// DeletePolicy deletes the listener with the specified input
func (c *PolicyClient) DeletePolicy(lb LoadBalancerContext, name string) (*PolicyInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForPolicyDeletePollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForPolicyDeleteTimeout
	}

	var info PolicyInfo
	if err := c.deleteResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}

	// deletedStates := []LBaaSState{LBaaSStateDeletionInProgress, LBaaSStateDeleted}
	deletedStates := []LBaaSState{LBaaSStateDeleted}
	erroredStates := []LBaaSState{LBaaSStateDeletionFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	deleted, err := c.checkPolicyState(&info, deletedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if deleted {
		return &info, nil
	}
	// else poll till deleted
	err = c.WaitForPolicyState(lb, name, deletedStates, erroredStates, c.PollInterval, c.Timeout, &info)
	if err != nil && client.WasNotFoundError(err) {
		// resource could not be found, thus deleted
		return nil, nil
	}
	return &info, nil
}

// GetPolicy fetchs the listener details
func (c *PolicyClient) GetPolicy(lb LoadBalancerContext, name string) (*PolicyInfo, error) {

	var info PolicyInfo
	if err := c.getResource(lb.Region, lb.Name, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetPolicy fetchs the listener details
func (c *PolicyClient) UpdatePolicy(lb LoadBalancerContext, name, policyType string, input *UpdatePolicyInput) (*PolicyInfo, error) {

	if c.PollInterval == 0 {
		c.PollInterval = waitForPolicyReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = waitForPolicyReadyTimeout
	}

	c.ContentType = c.getContentTypeForPolicyType(policyType)
	var info PolicyInfo
	if err := c.updateResource(lb.Region, lb.Name, name, &input, &info); err != nil {
		return nil, err
	}

	// updatedStates := []LBaaSState{LBaaSStateModificationInProgress, LBaaSStateHealthy}
	updatedStates := []LBaaSState{LBaaSStateHealthy}
	erroredStates := []LBaaSState{LBaaSStateModificaitonFailed, LBaaSStateAbandon, LBaaSStateAutoAbandoned}

	// check the initial response
	ready, err := c.checkPolicyState(&info, updatedStates, erroredStates)
	if err != nil {
		return nil, err
	}
	if ready {
		return &info, nil
	}
	// else poll till ready
	err = c.WaitForPolicyState(lb, name, updatedStates, erroredStates, c.PollInterval, c.Timeout, &info)
	return &info, err
}

// return the corrent Content Type for the Update request depending on the Policy Type
// of the Policy being updated.
func (c *PolicyClient) getContentTypeForPolicyType(policyType string) string {
	return fmt.Sprintf("application/vnd.com.oracle.oracloud.lbaas.%s+json", policyType)
}

// WaitForPolicyState waits for the resource to be in one of a set of desired states
func (c *PolicyClient) WaitForPolicyState(lb LoadBalancerContext, name string, desiredStates, errorStates []LBaaSState, pollInterval, timeoutSeconds time.Duration, info *PolicyInfo) error {

	var getErr error
	err := c.client.WaitFor("Policy status update", pollInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetPolicy(lb, name)
		if getErr != nil {
			return false, getErr
		}

		return c.checkPolicyState(info, desiredStates, errorStates)
	})
	return err
}

// check the State, returns in desired state (true), not ready yet (false) or errored state (error)
func (c *PolicyClient) checkPolicyState(info *PolicyInfo, desiredStates, errorStates []LBaaSState) (bool, error) {

	c.client.DebugLogString(fmt.Sprintf("Policy %v state is %v", info.Name, info.State))

	state := LBaaSState(info.State)

	if isStateInLBaaSStates(state, desiredStates) {
		// we're good, return okay
		return true, nil
	}
	if isStateInLBaaSStates(state, errorStates) {
		// not good, return error
		return false, fmt.Errorf("Policy %v in errored state %v", info.Name, info.State)
	}
	// not ready lifecycleTimeout
	return false, nil
}
