package lbaas

import "fmt"

var (
	loadBalancerContainerPath = "/vlbrs"
	loadBalancerResourcePath  = "/vlbrs/%s/%s/"
)

// LoadBalancerClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type LoadBalancerClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
}

// LoadBalancerClient returns an ServiceInstanceClient which is used to access the
// Load Balancer API
func (c *Client) LoadBalancerClient() *LoadBalancerClient {
	c.ContentType = CONTENT_TYPE_VLBR_JSON
	c.Accept = CONTENT_TYPE_VLBR_JSON
	return &LoadBalancerClient{
		Client:           c,
		ContainerPath:    loadBalancerContainerPath,
		ResourceRootPath: loadBalancerResourcePath,
	}
}

func (c *LoadBalancerClient) getObjectPath(root, region, name string) string {
	return fmt.Sprintf(root, region, name)
}

// executes the Create requests to the Load Balancer API
func (c *LoadBalancerClient) createResource(requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("POST", c.ContainerPath, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Update requests to the Load Balancer API
func (c *LoadBalancerClient) updateResource(region, name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, region, name)
	resp, err := c.executeRequest("PUT", objectPath, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Get requests to the Load Balancer API
func (c *LoadBalancerClient) getResource(region, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, region, name)
	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Delete requests to the Load Balancer API
func (c *LoadBalancerClient) deleteResource(region, name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, region, name)
	resp, err := c.executeRequest("DELETE", objectPath, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}
