package lbaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/mitchellh/mapstructure"
)

// Resource Content Types
const CONTENT_TYPE_VLBR_JSON = "application/vnd.com.oracle.oracloud.lbaas.VLBR+json"
const CONTENT_TYPE_LISTENER_JSON = "application/vnd.com.oracle.oracloud.lbaas.Listener+json"
const CONTENT_TYPE_ORIGIN_SERVER_POOL_JSON = "application/vnd.com.oracle.oracloud.lbaas.OriginServerPool+json"
const CONTENT_TYPE_SERVER_CERTIFICATE_JSON = "application/vnd.com.oracle.oracloud.lbaas.ServerCertificate+json"

// Policy Specific Content Types
const CONTENT_TYPE_APP_COOKIE_STICKINESS_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.AppCookieStickinessPolicy+json"
const CONTENT_TYPE_CLOUDGATE_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.CloudGatePolicy+json"
const CONTENT_TYPE_LB_COOKIE_STICKINESS_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.LBCookieStickinessPolicy+json"
const CONTENT_TYPE_LOADBALANCING_MECHANISM_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.LoadBalancingMechanismPolicy+json"
const CONTENT_TYPE_RATE_LIMITING_REQUEST_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.RateLimitingRequestPolicy+json"
const CONTENT_TYPE_REDIRECT_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.RedirectPolicy+json"
const CONTENT_TYPE_RESOURCE_ACCESS_CONTROL_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.ResourceAccessControlPolicy+json"
const CONTENT_TYPE_SET_REQUEST_HEADER_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.SetRequestHeaderPolicy+json"
const CONTENT_TYPE_SSL_NEGOTIATION_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.SSLNegotiationPolicy+json"
const CONTENT_TYPE_TRUSTED_CERTIFICATE_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.TrustedCertPolicy+json"

// LBaaSState common State type for all LBaaS service resources
type LBaaSState string

const (
	LBaaSStateCreationInProgress              LBaaSState = "CREATION_IN_PROGRESS"
	LBaaSStateCreated                         LBaaSState = "CREATED"
	LBaaSStateHealthy                         LBaaSState = "HEALTHY"
	LBaaSStateAdministratorInterventionNeeded LBaaSState = "ADMINISTRATOR_INTERVENTION_NEEDED"
	LBaaSStateDeletionInProgress              LBaaSState = "DELETION_IN_PROGRESS"
	LBaaSStateDeleted                         LBaaSState = "DELETED"
	LBaaSStateModificationInProgress          LBaaSState = "MODIFICATION_IN_PROGRESS"
	LBaaSStateCreationFailed                  LBaaSState = "CREATION_FAILED"
	LBaaSStateModificaitonFailed              LBaaSState = "MODIFICATION_FAILED"
	LBaaSStateDeletionFailed                  LBaaSState = "DELETION_FAILED"
	LBaaSStateAccessDenied                    LBaaSState = "ACCESS_DENIED"
	LBaaSStateAbandon                         LBaaSState = "ABANDON"
	LBaaSStateAutoAbandoned                   LBaaSState = "AUTO_ABANDONED"
	LBaaSStatePause                           LBaaSState = "PAUSE"
	LBaaSStateForcePaused                     LBaaSState = "FORCE_PAUSED"
	LBaaSStateResume                          LBaaSState = "RESUME"
)

// LBaaaSStatus common Status type for all LBaaS service resources
type LBaaSStatus string

const (
	LBaaSStatusEnabled  LBaaSStatus = "ENABLED"
	LBaaSStatusDisabled LBaaSStatus = "DISABLED"
)

// LBaaSDisabled common Disabled State type for all LBaaS service resources
type LBaaSDisabled string

const (
	LBaaSDisabledTrue        LBaaSDisabled = "TRUE"
	LBaaSDisabledFalse       LBaaSDisabled = "FALSE"
	LBaaSDisabledMaintenance LBaaSDisabled = "MAINTENANCE_MODE"
)

// Projections can be specified when retrieving collection of resources as well as when retrieving a specific resource.
// There are of four types : MINIMAL, CONSOLE, FULL, and DETAILED
type QueryProjection string

const (
	QueryMinimal  QueryProjection = "MINIMAL"
	QueryConsol   QueryProjection = "CONSOLE"
	QueryFull     QueryProjection = "FULL"
	QueryDetailed QueryProjection = "DETAILED"
)

// Client implementation for Oracle Cloud Infrastructure Load Balancing Classic */
type Client struct {
	client       *client.Client
	PollInterval time.Duration
	Timeout      time.Duration
}

// NewClient returns a new LBaaSClient
func NewClient(c *opc.Config) (*Client, error) {
	appClient := &Client{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	appClient.client = client

	return appClient, nil
}

func (c *Client) executeRequest(method, path, accept, contentType string, body interface{}) (*http.Response, error) {

	reqBody, err := c.client.MarshallRequestBody(body)
	if err != nil {
		return nil, err
	}

	req, err := c.client.BuildRequestBody(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", accept)
	debugReqString := fmt.Sprintf("HTTP %s Req (%s)", method, path)
	debugReqString = fmt.Sprintf("%s:\nAccept: %+v", debugReqString, accept)
	if body != nil {
		req.Header.Set("Content-Type", contentType)
		debugReqString = fmt.Sprintf("%s:\nContent-Type: %+v\nBody: %+v", debugReqString, contentType, string(reqBody))
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)

	// Set the authentication headers
	req.SetBasicAuth(*c.client.UserName, *c.client.Password)

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) getContainerPath(root, lbRegion, lbName string) string {
	return fmt.Sprintf(root, lbRegion, lbName)
}

func (c *Client) getObjectPath(root, lbRegion, lbName, name string) string {
	return fmt.Sprintf(root, lbRegion, lbName, name)
}

func (c *Client) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	c.client.DebugLogString(fmt.Sprintf("HTTP Resp (%d): %s", resp.StatusCode, buf.String()))
	// JSON decode response into interface
	var tmp interface{}
	dcd := json.NewDecoder(buf)
	if err = dcd.Decode(&tmp); err != nil {
		return err
	}

	// Use mapstructure to weakly decode into the resulting interface
	msdcd, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           iface,
		TagName:          "json",
	})
	if err != nil {
		return err
	}

	if err := msdcd.Decode(tmp); err != nil {
		return err
	}
	return nil
}

// return true if a given LBaaSState is in a List of LBaaSStates
func isStateInLBaaSStates(state LBaaSState, list []LBaaSState) bool {
	for _, s := range list {
		if LBaaSState(s) == state {
			return true
		}
	}
	return false
}
