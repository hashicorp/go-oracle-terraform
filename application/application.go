package application

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

var (
  ApplicationContainerPath = "/paas/service/apaas/api/v1.1/apps/%s"
	ApplicationResourcePath  = "/paas/service/apaas/api/v1.1/apps/%s/%s"
)

// ServiceInstanceClient is a client for the Service functions of the Java API.
type ApplicationContainerClient struct {
	ResourceClient
	Timeout time.Duration
}

// ApplicationContainerClient obtains an AppllicationClient which can be used to access to the
// Application Container functions of the Application Container API
func (c *ApplicationClient) ApplicationContainerClient() *ApplicationContainerClient {
	return &ApplicationContainerClient{
		ResourceClient: ResourceClient{
			ApplicationClient:       c,
			ContainerPath:    ApplicationContainerPath,
			ResourceRootPath: ApplicationResourcePath,
		}}
}
