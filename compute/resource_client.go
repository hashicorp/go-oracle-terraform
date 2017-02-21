package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type ResourceClient struct {
	*Client
	ResourceDescription string
	ContainerPath       string
	ResourceRootPath    string
}

func (c *ResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("POST", c.ContainerPath, requestBody)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) updateResource(name string, requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("PUT", c.getObjectPath(c.ResourceRootPath, name), requestBody)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) getResource(name string, responseBody interface{}) error {
	resp, err := c.executeRequest("GET", c.getObjectPath(c.ResourceRootPath, name), nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) deleteResource(name string) error {
	_, err := c.executeRequest("DELETE", c.getObjectPath(c.ResourceRootPath, name), nil)
	if err != nil {
		return err
	}

	// No errors and no response body to write
	return nil
}

func (c *ResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	c.debugLogStr(fmt.Sprintf("HTTP Resp (%d): %s", resp.StatusCode, buf.String()))
	err := json.Unmarshal(buf.Bytes(), iface)
	if err != nil {
		return fmt.Errorf("Error unmarshalling response body: %s", err)
	}
	return nil
}
