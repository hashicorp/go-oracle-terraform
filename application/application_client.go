package application

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

// Client represents an authenticated application client, with credentials and an api client.
type Client struct {
	client *client.Client
}

// NewClient returns a new client for the application resources managed by Oracle
func NewClient(c *opc.Config) (*Client, error) {
	appClient := &Client{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	appClient.client = client

	return appClient, nil
}

func (c *Client) executeCreateUpdateRequest(method, path string, input *CreateApplicationContainerInput) (*http.Response, error) {
	req, err := c.client.BuildMultipartFormRequest(method, path, input)
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

func (c *Client) executeRequest(method, path string, body interface{}) (*http.Response, error) {
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

// BuildMultipartFormRequest builds a new HTTP Request for a multipart form request from specifies attributes
func (c *Client) BuildMultipartFormRequest(method, path string, manifestBody interface{}, parameters map[string]interface{}) (*http.Request, error) {

	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	var (
		part io.Writer
	)

	part, err = writer.CreateFormFile("manifest", "manifest.json")
	if err != nil {
		return nil, err
	}
	manifestBytes, err := c.client.MarshallRequestBody(manifestBody)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(manifestBytes)
	if err != nil {
		return nil, err
	}

	// Add additional parameters to the writer
	for key, val := range parameters {
		if val.(string) != "" {
			_ = writer.WriteField(strings.ToLower(key), val.(string))
		}
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.formatURL(urlPath), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}

func (c *Client) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *Client) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}
