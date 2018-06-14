package lbaas

import (
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

// GetTestClient obtains a client for testing purposes
func GetTestClient(c *opc.Config) (*Client, error) {
	// Build up config with default values if omitted

	if c.Username == nil {
		username := os.Getenv("OPC_USERNAME")
		c.Username = &username
	}

	if c.Password == nil {
		password := os.Getenv("OPC_PASSWORD")
		c.Password = &password
	}

	if c.APIEndpoint == nil {
		apiEndpoint, err := url.Parse(os.Getenv("OPC_LBAAS_ENDPOINT"))
		if err != nil {
			return nil, err
		}
		c.APIEndpoint = apiEndpoint
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				TLSHandshakeTimeout: 120 * time.Second},
		}
	}

	return NewClient(c)
}

func compare(t *testing.T, attrName, respValue, expectedValue string) {
	if respValue != expectedValue {
		t.Fatalf("%s %s in response does to match expected value of %s", attrName, respValue, expectedValue)
	}
}
