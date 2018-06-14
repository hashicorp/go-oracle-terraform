package lbaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// ResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type ResourceClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
}

func (c *ResourceClient) createResource(region, name string, requestBody interface{}, responseBody interface{}) error {

	resp, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath), requestBody)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) updateResource(region, name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, region, name)
	resp, err := c.executeRequest("POST", objectPath, requestBody)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) getResource(region, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, region, name)
	log.Print(c.ContainerPath)
	log.Print(c.ResourceRootPath)
	log.Print(objectPath)
	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) deleteResource(region, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, region, name)
	resp, err := c.executeRequest("DELETE", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *ResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
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
