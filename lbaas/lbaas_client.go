package lbaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
	"github.com/mitchellh/mapstructure"
)

const CONTENT_TYPE_VLBR_JSON = "application/vnd.com.oracle.oracloud.lbaas.VLBR+json"
const CONTENT_TYPE_LISTENER_JSON = "application/vnd.com.oracle.oracloud.lbaas.Listener+json"
const CONTENT_TYPE_ORIGIN_SERVER_POOL_JSON = "application/vnd.com.oracle.oracloud.lbaas.OriginServerPool+json"
const CONTENT_TYPE_APP_COOKIE_STICKINESS_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.AppCookieStickinessPolicy+json"
const CONTENT_TYPE_LB_COOKIE_STICKINESS_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.LBCookieStickinessPolicy+json"
const CONTENT_TYPE_RESOURCE_ACCESS_CONTROL_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.ResourceAccessControlPolicy+json"
const CONTENT_TYPE_REDIRECT_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.RedirectPolicy+json"
const CONTENT_TYPE_SSL_NEGOTIATION_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.SSLNegotiationPolicy+json"
const CONTENT_TYPE_SET_REQUEST_HEADER_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.SetRequestHeaderPolicy+json"
const CONTENT_TYPE_TRUSTED_CERTIFICATE_POLICY_JSON = "application/vnd.com.oracle.oracloud.lbaas.TrustedCertPolicy+json"
const CONTENT_TYPE_SERVER_CERTIFICATE_JSON = "application/vnd.com.oracle.oracloud.lbaas.ServerCertificate+json"

// Client implementation for Oracle Cloud Infrastructure Load Balancing Classic */
type Client struct {
	client      *client.Client
	ContentType string
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

func (c *Client) executeRequest(method, path string, body interface{}) (*http.Response, error) {

	reqBody, err := c.client.MarshallRequestBody(body)
	if err != nil {
		return nil, err
	}

	req, err := c.client.BuildRequestBody(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Req (%s)", method, path)
	req.Header.Add("Accept", c.ContentType)
	if body != nil {
		req.Header.Set("Content-Type", c.ContentType)
		// Debug the body for database services
		debugReqString = fmt.Sprintf("%s:\nBody: %+v", debugReqString, string(reqBody))
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
