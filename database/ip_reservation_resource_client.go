package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// IPReservationResourceClient is a client for the IP Reservation functions of the Database API.
type IPReservationResourceClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
}

func (c *IPReservationResourceClient) createResource(requestBody interface{}, responseBody interface{}) error {
	_, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath), requestBody)

	return err
}

func (c *IPReservationResourceClient) getResource(name string, responseBody interface{}) error {
	var objectPath string
	if name != "" {
		objectPath = c.getObjectPath(c.ResourceRootPath, name)
	} else {
		objectPath = c.getContainerPath(c.ContainerPath)
	}

	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}

	return c.unmarshalResponseBody(resp, responseBody)
}

func (c *IPReservationResourceClient) deleteResource(name string) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)

	if _, err := c.executeRequest("DELETE", objectPath, nil); err != nil {
		return err
	}
	return nil
}

func (c *IPReservationResourceClient) unmarshalResponseBody(resp *http.Response, iface interface{}) error {
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
		return fmt.Errorf("Error decoding: %s\n%+v", err.Error(), resp)
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

func (c *IPReservationResourceClient) getContainerPath(root string) string {
	// /paas/api/v1.1/network/{identityDomainId}/services/dbaas/ipreservations
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *IPReservationResourceClient) getObjectPath(root, name string) string {
	// /paas/api/v1.1/network/{identityDomainId}/services/dbaas/ipreservations/{ipResName}
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}
