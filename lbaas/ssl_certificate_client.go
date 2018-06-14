package lbaas

import "fmt"

var (
	sslCertificateContainerPath = "/certs"
	sslCertificaetResourcePath  = "/certs/%s/"
)

// SSLCertificateClient is an AuthenticatedClient with some additional information about the resources to be addressed.
type SSLCertificateClient struct {
	*Client
	ContainerPath    string
	ResourceRootPath string
}

// SSLCertificateClient returns an ServiceInstanceClient which is used to access the
// Load Balancer API
func (c *Client) SSLCertificateClient() *SSLCertificateClient {
	c.ContentType = CONTENT_TYPE_SERVER_CERTIFICATE_JSON
	c.Accept = CONTENT_TYPE_SERVER_CERTIFICATE_JSON
	return &SSLCertificateClient{
		Client:           c,
		ContainerPath:    sslCertificateContainerPath,
		ResourceRootPath: sslCertificaetResourcePath,
	}
}

func (c *SSLCertificateClient) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, name)
}

// executes the Create requests to the Load Balancer API
func (c *SSLCertificateClient) createResource(requestBody interface{}, responseBody interface{}) error {
	resp, err := c.executeRequest("POST", c.ContainerPath, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Update requests to the Load Balancer API
func (c *SSLCertificateClient) updateResource(name string, requestBody interface{}, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)
	resp, err := c.executeRequest("POST", objectPath, requestBody)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Get requests to the Load Balancer API
func (c *SSLCertificateClient) getResource(name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)
	resp, err := c.executeRequest("GET", objectPath, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}

// executes the Delete requests to the Load Balancer API
func (c *SSLCertificateClient) deleteResource(name string, responseBody interface{}) error {
	objectPath := c.getObjectPath(c.ResourceRootPath, name)
	resp, err := c.executeRequest("DELETE", objectPath, nil)
	if err != nil {
		return err
	}
	return c.unmarshalResponseBody(resp, responseBody)
}
