package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	logger         opc.Logger
	loglevel       opc.LogLevelType
}

func NewComputeClient(c *opc.Config) (*Client, error) {
	// First create a client
	client := &Client{
		identityDomain: c.IdentityDomain,
		userName:       c.Username,
		password:       c.Password,
		apiEndpoint:    c.APIEndpoint,
		httpClient:     c.HTTPClient,
		loglevel:       c.LogLevel,
	}

	// Setup logger; defaults to stdout
	if c.Logger == nil {
		client.logger = opc.NewDefaultLogger()
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

	// Log the request before the authentication cookie, so as not to leak credentials
	// TODO: FIX ME!
	//c.debugLogReq(req)

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

	// Even though the returned body will be in json form, it's undocumented what
	// fields are actually returned. Once we get documentation of the actual
	// error fields that are possible to be returned we can have stricter error types.
	if resp.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		oracleErr.Message = buf.String()
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

func (c *Client) debugLogReq(req *http.Request) {
	// Don't need to log this if not debugging
	if c.loglevel != opc.LogDebug {
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	c.logger.Log(fmt.Sprintf("DEBUG: HTTP %s Req %s: %s",
		req.Method, req.URL.String(), buf.String()))
}

// Log a string if debug logs are on
func (c *Client) debugLogStr(str string) {
	if c.loglevel != opc.LogDebug {
		return
	}
	c.logger.Log(fmt.Sprintf("[DEBUG]: %s", str))
}

// Retry function
func (c *Client) waitFor(description string, timeoutSeconds int, test func() (bool, error)) error {
	tick := time.Tick(1 * time.Second)

	for i := 0; i < timeoutSeconds; i++ {
		select {
		case <-tick:
			completed, err := test()
			c.debugLogStr(fmt.Sprintf("Waiting for %s (%d/%ds)", description, i, timeoutSeconds))
			if err != nil || completed {
				return err
			}
		}
	}
	return fmt.Errorf("Timeout waiting for %s", description)
}

// Used to determine if the checked resource was found or not.
func WasNotFoundError(e error) bool {
	err, ok := e.(*opc.OracleError)
	if ok {
		return err.StatusCode == 404
	}
	return false
}
