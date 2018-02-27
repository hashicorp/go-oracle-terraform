package application

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// ApplicationClient represents an authenticated application client, with credentials and an api client.
type ApplicationClient struct {
	client *client.Client
}

func NewApplicationClient(c *opc.Config) (*ApplicationClient, error) {
	appClient := &ApplicationClient{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	appClient.client = client

	return appClient, nil
}

func (c *ApplicationClient) executeCreateUpdateRequest(method, path string, files map[string]string, parameters map[string]interface{}) (*http.Response, error) {
	req, err := c.client.BuildMultipartFormRequest(method, path, files, parameters)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Path (%s)", method, path)
	// req.Header.Set("Content-Type", "multipart/form-data")
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)
	c.client.DebugLogString(fmt.Sprintf("Req (%+v)", req))

	// Set the authentication headers
	req.SetBasicAuth(*c.client.UserName, *c.client.Password)
	req.Header.Add("X-ID-TENANT-NAME", *c.client.IdentityDomain)

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ApplicationClient) executeRequest(method, path string, body interface{}) (*http.Response, error) {
	reqBody, err := c.client.MarshallRequestBody(body)
	if err != nil {
		return nil, err
	}

	req, err := c.client.BuildRequestBody(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Path (%s)", method, path)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)
	c.client.DebugLogString(fmt.Sprintf("Req (%+v)", req))

	// Set the authentiation headers
	req.SetBasicAuth(*c.client.UserName, *c.client.Password)
	req.Header.Add("X-ID-TENANT-NAME", *c.client.IdentityDomain)

	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ApplicationClient) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *ApplicationClient) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}
