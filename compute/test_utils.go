package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"time"

	"github.com/hashicorp/go-oracle-terraform/opc"
)

func newAuthenticatingServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request: %s, %s\n", r.Method, r.URL)

		if r.URL.Path == "/authenticate/" {
			http.SetCookie(w, &http.Cookie{Name: "testAuthCookie", Value: "cookie value"})
			w.WriteHeader(200)
		} else {
			handler(w, r)
		}
	}))
}

func getTestClient(c *opc.Config) (*Client, error) {
	// Build up config with default values if omitted
	if c.APIEndpoint == nil {
		if os.Getenv("OPC_ENDPOINT") == "" {
			panic("OPC_ENDPOINT not set in environment")
		}
		endpoint, err := url.Parse(os.Getenv("OPC_ENDPOINT"))
		if err != nil {
			return nil, err
		}
		c.APIEndpoint = endpoint
	}

	if c.IdentityDomain == nil {
		domain := os.Getenv("OPC_IDENTITY_DOMAIN")
		c.IdentityDomain = &domain
	}

	if c.Username == nil {
		username := os.Getenv("OPC_USERNAME")
		c.Username = &username
	}

	if c.Password == nil {
		password := os.Getenv("OPC_PASSWORD")
		c.Password = &password
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				TLSHandshakeTimeout: 120 * time.Second},
		}
	}

	// Test Client should be able to log debug levels if required
	if os.Getenv("ORACLE_LOG") != "" {
		c.LogLevel = opc.LogDebug
	}

	return NewComputeClient(c)
}

// Returns a stub client with default values, and a custom API Endpoint
func getStubClient(endpoint *url.URL) (*Client, error) {
	domain := "test"
	username := "test"
	password := "test"
	config := &opc.Config{
		IdentityDomain: &domain,
		Username:       &username,
		Password:       &password,
		APIEndpoint:    endpoint,
	}
	return getTestClient(config)
}

func unmarshalRequestBody(t *testing.T, r *http.Request, target interface{}) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err := json.Unmarshal(buf.Bytes(), target)
	if err != nil {
		t.Fatalf("Error marshalling request: %s", err)
	}
	if _, err := io.Copy(os.Stdout, buf); err != nil {
		t.Fatalf("Error copying file: %s", err)
	}
}

// Unused Function
/*func marshalToBytes(target interface{}) []byte {
	marshalled, err := json.Marshal(target)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.Read(marshalled)
	io.Copy(os.Stdout, buf)
	fmt.Println()
	return marshalled
}*/
