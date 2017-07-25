package java

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/opc"
)

const JAVA_ACCOUNT = "/Java-%s"
const AUTH_HEADER = "Authorization"
const TENANT_HEADER = "X-ID-TENANT-NAME"
const JAVA_QUALIFIED_NAME = "%s%s/%s"

// Client represents an authenticated java client, with compute credentials and an api client.
type JavaClient struct {
	client     *client.Client
	authHeader *string
}

func NewJavaClient(c *opc.Config) (*JavaClient, error) {
	javaClient := &JavaClient{}
	client, err := client.NewClient(c)
	if err != nil {
		return nil, err
	}
	javaClient.client = client

	javaClient.authHeader = javaClient.getAuthenticationHeader()

	return javaClient, nil
}

func (c *JavaClient) executeRequest(method, path string, body interface{}) (*http.Response, error) {
	req, err := c.client.BuildRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	debugReqString := fmt.Sprintf("HTTP %s Req (%s)", method, path)
	if body != nil {
		req.Header.Set("Content-Type", "application/vnd.com.oracle.oracloud.provisioning.Service+json")
	}
	// Log the request before the authentication header, so as not to leak credentials
	c.client.DebugLogString(debugReqString)

	// Set the authentiation headers
	req.Header.Add(AUTH_HEADER, *c.authHeader)
	req.Header.Add(TENANT_HEADER, *c.client.IdentityDomain)
	c.client.DebugLogString(fmt.Sprintf("Req (%+v)", req))
	resp, err := c.client.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *JavaClient) getAccount() string {
	return fmt.Sprintf(JAVA_ACCOUNT, *c.client.IdentityDomain)
}

// GetQualifiedName returns the fully-qualified name of a java object, e.g. /v1/{account}/{name}
func (c *JavaClient) getQualifiedName(version string, name string) string {
	if name == "" {
		return ""
	}
	if strings.HasPrefix(name, "/Java-") || strings.HasPrefix(name, "v1/") {
		return name
	}
	return fmt.Sprintf(JAVA_QUALIFIED_NAME, version, c.getAccount(), name)
}

// GetUnqualifiedName returns the unqualified name of a Java object, e.g. the {name} part of /v1/{account}/{name}
func (c *JavaClient) getUnqualifiedName(name string) string {
	if name == "" {
		return name
	}
	if !strings.Contains(name, "/") {
		return name
	}

	nameParts := strings.Split(name, "/")
	return strings.Join(nameParts[len(nameParts)-1:], "/")
}

func (c *JavaClient) unqualify(names ...*string) {
	for _, name := range names {
		*name = c.getUnqualifiedName(*name)
	}
}

func (c *JavaClient) getContainerPath(root string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain)
}

func (c *JavaClient) getObjectPath(root, name string) string {
	return fmt.Sprintf(root, *c.client.IdentityDomain, name)
}
