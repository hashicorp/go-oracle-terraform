package lbaas

import (
	"fmt"
	"net/http"
)

/*
 * The LBaaSResourceClient is the general client used for the majority of the Load Balancer
 * Service child resources (Listeners, Origin Servier Pools and Policies) which have the common URI
 * format https://{api_endpoint}/{lb_name}/{lb_region}/{resource_type}/{resource_name}?{projection}
 *
 * For SSL Certificates use the SSLCertificateClient
 * For the Load Balancer Service Instance use the LoadBalancerResourceClient
 */

// LBaaSResourceClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type LBaaSResourceClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
	Projection       string
	Accept           string
	ContentType      string
}

// executes the Create requests to the LBaaS API
func (c *LBaaSResourceClient) createResource(lbRegion, lbName string, requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("POST", c.getContainerPath(c.ContainerPath, lbRegion, lbName), c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Update requests to the LBaaS API
func (c *LBaaSResourceClient) updateResource(lbRegion, lbName, name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	resp, err := c.executeRequest("PUT", objectPath, c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Update requests to the LBaaS API specific to updating the Origin Server Pool
// which has a different update style using POST + with an PATCH Method override
func (c *LBaaSResourceClient) updateOriginServerPool(lbRegion, lbName, name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	resp, err := c.executeRequestWithMethodOverride("POST", "PATCH", objectPath, c.Accept, c.ContentType, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Get requests to the LBaaS API
func (c *LBaaSResourceClient) getResource(lbRegion, lbName, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	queryParams := ""
	if c.Projection != "" {
		queryParams = fmt.Sprintf("?projection=%s" + c.Projection)
	}
	resp, err := c.executeRequest("GET", objectPath+queryParams, c.Accept, c.ContentType, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Delete requests to the LBaaS API
func (c *LBaaSResourceClient) deleteResource(lbRegion, lbName, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, lbRegion, lbName, name)
	resp, err := c.executeRequest("DELETE", objectPath, c.Accept, c.ContentType, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// execute a request with a X-HTTP-Method-Override
func (c *Client) executeRequestWithMethodOverride(method, methodOverride, path, accept, contentType string, body interface{}) (*http.Response, error) {

	reqBody, err := c.client.MarshallRequestBody(body)
	if err != nil {
		return nil, err
	}

	req, err := c.client.BuildRequestBody(method, path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-HTTP-Method-Override", methodOverride)
	req.Header.Add("Accept", accept)
	debugReqString := fmt.Sprintf("HTTP %s (%s) Req (%s)", method, methodOverride, path)
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
