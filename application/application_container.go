package application

import (
	"time"

	"github.com/fatih/structs"
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

type ApplicationSubscriptionType string

const (
	ApplicationSubscriptionTypeHourly ApplicationSubscriptionType = "HOURLY"
	ApplicationSubscriptionTypeMonthly ApplicationSubscriptionType = "MONTHLY"
)

type ApplicationRepository string

const (
	ApplicationRepositoryDockerHub ApplicationRepository = "dockerhub"
)

type ApplicationRuntime string

const (
	ApplicationRuntimeJava ApplicationRuntime = "java"
	ApplicationRuntimeNode ApplicationRuntime = "node"
	ApplicationRuntimePHP ApplicationRuntime = "php"
	ApplicationRuntimePython ApplicationRuntime = "python"
	ApplicationRuntimeRuby ApplicationRuntime = "ruby"
)

type ApplicationContainer struct {
	// ID of the application
	AppId string `json:"appId"`
	// URL of the created application
	AppURL string `json:"appURL"`
	// Creation time of the application
	CreatedTime string `json:"createdTime"`
	// Identity Domain of the application
	IdentityDomain string `json:"identityDomain"`
	// Shows details of all instances currently running.
	Instances []Instance `json:"instances"`
	// Modification time of the application
	LastModifiedTime string `json:"lastModifiedTime"`
	// Shows all deployments currently in progress.
	LatestDeployments []Deployment `json:"latestDeployment"`
	// Name of the application
	Name string `json:"name"`
	// Shows all deployments currently running.
	RunningDeployments []Deployment `json:"runningDeployment"`
	// Status of the application
	Status string `json:"status"`
	// Type of subscription, Hourly or Monthly
	SubscriptionType ApplicationSubscriptionType `json:"subscriptionType"`
	// Web URL of the application
	WebURL string `json:"webURL"`
}

type Instance struct {
	// Instance URL. Use this url to get a description of the application instance.
	InstanceURL string `json:"instanceURL"`
	// Memory of the instance
	Memory string `json:"memory"`
	// Instance Name. Use this name to manage a specific instance.
	Name string `json:"name"`
	// Status of the instance
	Status string `json:"status"`
}

type Deployment struct {
	// Deployment ID. Use this ID to manage a specific deployment.
	DeploymentId string `json:"deploymentId"`
	// Status of the deployment
	DeploymentStatus string `json:"deploymentStatus"`
	// Deployment URL. Use this URL to get a description of the application deployment.
	DeploymentURL string `json:"deploymentURL"`
}

type CreateApplicationContainerInput struct {
	// The additional fields needed for ApplicationContainer
	AdditionalFields CreateApplicationContainerAdditionalFields
	// Name of the optional deployment file, which specifies memory, number of instances, and service bindings
	// Optional
	Deployment string
	// Name of the manifest file, required if this file is not packaged with the application
	// Optional
	Manifest string

}

type CreateApplicationContainerAdditionalFields struct {
	// Location of the application archive file in Oracle Storage Cloud Service, in the format app-name/file-name
	// Optional
	ArchiveURL string
	// Name of the application
	// Required
	Name string
	// Comments on the application deployment
	Notes string
	// Email address to which application deployment status updates are sent.
	NotificationEmail string
	// Repository of the application. The only allowed value is 'dockerhub'.
	// Optional
	Repository ApplicationRepository
	// Runtime environment: java (the default), node, php, python, or ruby
	// Optional
	Runtime ApplicationRuntime
	// Subscription, either hourly (the default) or monthly
	SubscriptionType ApplicationSubscriptionType
}

// Create a new Applicaiton Container from an ApplicationClient and an input struct.
// Returns a populated ApplicationContainer struct for the Application, and any errors
func (c *ApplicationContainerClient) CreateApplicationContainer(input *CreateApplicationContainerInput) (*ApplicationContainer, error) {

	files := make(map[string]string, 0)
	if input.Deployment != "" {
		files["deployment"] = input.Deployment
	}
	if input.Manifest != "" {
		files["manifest"] = input.Manifest
	}
	additionalFields := structs.Map(input.AdditionalFields)

	var applicationContainer *ApplicationContainer
	if err := c.createResource(files, additionalFields, applicationContainer); err != nil {
		return nil, err
	}
	return applicationContainer, nil
}

type GetApplicationContainerInput struct {
	// Name of the application container
	// Required
	Name string `json:"name"`
}

// GetApplicationContainer retrieves the application container with the given name.
func (c *ApplicationContainerClient) GetApplicationContainer(getInput *GetApplicationContainerInput) (*ApplicationContainer, error) {
	var applicationContainer ApplicationContainer
	if err := c.getResource(getInput.Name, &applicationContainer); err != nil {
		return nil, err
	}

	return &applicationContainer, nil
}

type DeleteApplicationContainerInput struct {
	// Name of the application container
	// Required
	Name string `json:"name"`
}

// DeleteApplicationContainer deletes the application container with the given name.
func (c *ApplicationContainerClient) DeleteApplicationContainer(deleteInput *DeleteApplicationContainerInput) error {
	return c.deleteResource(deleteInput.Name)
}