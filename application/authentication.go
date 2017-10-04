package application

import (
	"fmt"
)

// Get the authentication header for the application client
func (c *ApplicationClient) getAuthenticationHeader() *string {
	authHeader := fmt.Sprintf("%s:%s", *c.client.UserName, *c.client.Password)
	return &authHeader
}
