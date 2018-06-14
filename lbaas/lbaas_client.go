package lbaas

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const CONTENT_TYPE_VLBR_JSON = "application/vnd.com.oracle.oracloud.lbaas.VLBR+json"

// Client implementation for Oracle Cloud Infrastructure Load Balancing Classic */
type Client struct {
	client *client.Client
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
	req.Header.Add("Accept", CONTENT_TYPE_VLBR_JSON)
	if body != nil {
		req.Header.Set("Content-Type", CONTENT_TYPE_VLBR_JSON)
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

func (c *Client) getContainerPath(root string) string {
	return fmt.Sprintf(root)
}

func (c *Client) getObjectPath(root, region, name string) string {
	return fmt.Sprintf(root, region, name)
}
