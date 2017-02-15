package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"strings"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

const CMP_USERNAME = "/Compute-%s/%s"

// Client represents an authenticated compute client, with compute credentials and an api client.
type Client struct {
	identityDomain *string
	userName       *string
	password       *string
	apiEndpoint    *url.URL
	httpClient     *http.Client
	authCookie     *http.Cookie
	cookieIssued   time.Time
}

func NewComputeClient(c *opc.Config) (*Client, error) {
	// First create a client
	client := &Client{
		identityDomain: c.IdentityDomain,
		userName:       c.Username,
		password:       c.Password,
		apiEndpoint:    c.APIEndpoint,
		httpClient:     c.HTTPClient,
	}

	if err := client.getAuthenticationCookie(); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) executeRequest(method, path string, body interface{}) (*http.Response, error) {
	// Parse URL Path
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Marshall request body
	var requestBody io.ReadSeeker
	if body != nil {
		marshaled, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(marshaled)
	}

	// Create request
	req, err := http.NewRequest(method, c.formatURL(urlPath), requestBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/oracle-compute-v3+json")
	}

	// If we have an authentication cookie, let's authenticate, refreshing cookie if need be
	if c.authCookie != nil {
		if time.Since(c.cookieIssued).Minutes() > 25 {
			if err := c.getAuthenticationCookie(); err != nil {
				return nil, err
			}
		}
		req.AddCookie(c.authCookie)
	}

	// Execute request with supplied client
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp, nil
	}

	oracleErr := &opc.OracleError{
		StatusCode: resp.StatusCode,
	}

	errDecoder := json.NewDecoder(resp.Body)
	if err := errDecoder.Decode(oracleErr); err != nil {
		return nil, err
	}
	return nil, oracleErr
}

func (c *Client) formatURL(path *url.URL) string {
	return c.apiEndpoint.ResolveReference(path).String()
}

func (c *Client) getUserName() string {
	return fmt.Sprintf(CMP_USERNAME, *c.identityDomain, *c.userName)
}

// From compute_client
// GetObjectName returns the fully-qualified name of an OPC object, e.g. /identity-domain/user@email/{name}
func (c *Client) getQualifiedName(name string) string {
	if strings.HasPrefix(name, "/oracle") || strings.HasPrefix(name, "/Compute-") {
		return name
	}
	return fmt.Sprintf("%s/%s", c.getUserName(), name)
}

func (c *Client) getObjectPath(root, name string) string {
	return fmt.Sprintf("%s%s", root, c.getQualifiedName(name))
}

// GetUnqualifiedName returns the unqualified name of an OPC object, e.g. the {name} part of /identity-domain/user@email/{name}
func (c *Client) getUnqualifiedName(name string) string {
	if name == "" {
		return name
	}
	if strings.HasPrefix(name, "/oracle") {
		return name
	}
	nameParts := strings.Split(name, "/")
	return strings.Join(nameParts[3:], "/")
}

func (c *Client) unqualify(names ...*string) {
	for _, name := range names {
		*name = c.getUnqualifiedName(*name)
	}
}

// Used to determine if the checked resource was found or not.
func WasNotFoundError(e error) bool {
	err, ok := e.(*opc.OracleError)
	if ok {
		return err.StatusCode == 404
	}
	return false
}

// Retry function
func waitFor(description string, timeoutSeconds int, test func() (bool, error)) error {
	tick := time.Tick(1 * time.Second)

	for i := 0; i < timeoutSeconds; i++ {
		select {
		case <-tick:
			completed, err := test()
			if err != nil || completed {
				return err
			}
		}
	}
	return fmt.Errorf("Timeout waiting for %s", description)
}
