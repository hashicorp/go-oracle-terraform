package lbaas

import (
	"fmt"
	"strings"
)

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
		CONTENT_TYPE_LB_COOKIE_STICKINESS_POLICY_JSON,
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
	Action     string `json:"action_when_hdr_exists"`
	HeaderName string `json:"header_name"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Type       string `json:"type"`
	URI        string `json:"uri"`
	Value      string `json:"value"`
}

type CreatePolicyInput struct {
	Action     string `json:"action_when_hdr_exists"`
	HeaderName string `json:"header_name"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Value      string `json:"value"`
}

type UpdatePolicyInput struct {
	Action     string `json:"action_when_hdr_exists,omitempty"`
	HeaderName string `json:"header_name,omitempty"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Value      string `json:"value,omitempty"`
}

// CreatePolicy creates a new listener
func (c *PolicyClient) CreatePolicy(lbRegion, lbName string, input *CreatePolicyInput) (*PolicyInfo, error) {
	var info PolicyInfo
	c.ContentType = CONTENT_TYPE_SET_REQUEST_HEADER_POLICY_JSON
	if err := c.createResource(lbRegion, lbName, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// DeletePolicy deletes the listener with the specified input
func (c *PolicyClient) DeletePolicy(lbRegion, lbName, name string) (*PolicyInfo, error) {
	var info PolicyInfo
	if err := c.deleteResource(lbRegion, lbName, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetPolicy fetchs the listener details
func (c *PolicyClient) GetPolicy(lbRegion, lbName, name string) (*PolicyInfo, error) {
	var info PolicyInfo
	if err := c.getResource(lbRegion, lbName, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetPolicy fetchs the listener details
func (c *PolicyClient) UpdatePolicy(lbRegion, lbName, name, policyType string, input *UpdatePolicyInput) (*PolicyInfo, error) {
	c.ContentType = c.getContentTypeForPolicyType(policyType)
	var info PolicyInfo
	if err := c.updateResource(lbRegion, lbName, name, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// return the corrent Content Type for the Update request depending on the Policy Type
// of the Policy being updated.
func (c *PolicyClient) getContentTypeForPolicyType(policyType string) string {
	return fmt.Sprintf("application/vnd.com.oracle.oracloud.lbaas.%s+json", policyType)
}
