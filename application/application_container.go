package application

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-oracle-terraform/client"
)

const WaitForApplicationContainerRunningTimeout = time.Duration(600 * time.Second)
const WaitForApplicationContainerDeleteTimeout = time.Duration(600 * time.Second)

var (
	ApplicationContainerPath = "/paas/service/apaas/api/v1.1/apps/%s"
	ApplicationResourcePath  = "/paas/service/apaas/api/v1.1/apps/%s/%s"
)

// ApplicationContainerClient is a client for the Application container functions of the Application Container API
type ApplicationContainerClient struct {
	ResourceClient
	Timeout time.Duration
}

// ApplicationContainerClient obtains an AppllicationClient which can be used to access to the
// Application Container functions of the Application Container API
func (c *ApplicationClient) ApplicationContainerClient() *ApplicationContainerClient {
	return &ApplicationContainerClient{
		ResourceClient: ResourceClient{
			ApplicationClient: c,
			ContainerPath:     ApplicationContainerPath,
			ResourceRootPath:  ApplicationResourcePath,
		}}
}

type ApplicationSubscriptionType string

const (
	ApplicationSubscriptionTypeHourly  ApplicationSubscriptionType = "HOURLY"
	ApplicationSubscriptionTypeMonthly ApplicationSubscriptionType = "MONTHLY"
)

type ApplicationRepository string

const (
	ApplicationRepositoryDockerHub ApplicationRepository = "dockerhub"
)

type ApplicationRuntime string

const (
	ApplicationRuntimeJava   ApplicationRuntime = "java"
	ApplicationRuntimeJavaEE ApplicationRuntime = "javaee"
	ApplicationRuntimeNode   ApplicationRuntime = "node"
	ApplicationRuntimePHP    ApplicationRuntime = "php"
	ApplicationRuntimePython ApplicationRuntime = "python"
	ApplicationRuntimeRuby   ApplicationRuntime = "ruby"
	ApplicationRuntimeGo     ApplicationRuntime = "golang"
	ApplicationRuntimeDotnet ApplicationRuntime = "dotnet"
)

type ApplicationContainerAuthType string

const (
	ApplicationContainerAuthTypeBasic ApplicationContainerAuthType = "basic"
	ApplicationContainerAuthTypeOauth ApplicationContainerAuthType = "oauth"
)

type ApplicationContainerStatus string

const (
	ApplicationContainerStatusRunning        ApplicationContainerStatus = "RUNNING"
	ApplicationContainerStatusNew            ApplicationContainerStatus = "NEW"
	ApplicationContainerStatusDestroyPending ApplicationContainerStatus = "DESTROY_PENDING"
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
	LatestDeployment Deployment `json:"latestDeployment"`
	// Name of the application
	Name string `json:"name"`
	// Shows all deployments currently running.
	RunningDeployment Deployment `json:"runningDeployment"`
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
	// Timeout for creating an application container
	Timeout time.Duration
}

type CreateApplicationContainerAdditionalFields struct {
	// Location of the application archive file in Oracle Storage Cloud Service, in the format app-name/file-name
	// Optional
	ArchiveURL string
	// Uses Oracle Identity Cloud Service to control who can access your Java SE 7 or 8, Node.js, or PHP application.
	// This should be ApplicationContainerAuthType but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Allowed values are 'basic' and 'oauth'.
	// Optional
	AuthType string // ApplicationContainerAuthType
	// Name of the application
	// Required
	Name string
	// Comments on the application deployment
	Notes string
	// Email address to which application deployment status updates are sent.
	NotificationEmail string
	// Repository of the application. The only allowed value is 'dockerhub'.
	// This should be ApplicationRepository but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Optional
	Repository string // ApplicationRepository
	// Runtime environment: java (the default), node, php, python, or ruby
	// This should be ApplicationRuntime but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Required
	Runtime string // ApplicationRuntime
	// Subscription, either hourly (the default) or monthly
	// This should be ApplicationSubscriptionType but because of how we need to translate this strut to a map[string]string we are keeping it as a string
	// Optional
	SubscriptionType string //ApplicationSubscriptionType
}

// Create a new Application Container from an ApplicationClient and an input struct.
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

	getInput := &GetApplicationContainerInput{
		Name: input.AdditionalFields.Name,
	}

	if input.Timeout == 0 {
		input.Timeout = WaitForApplicationContainerRunningTimeout
	}

	// Wait for application container to be ready and return the result
	applicationContainerInfo, err := c.WaitForApplicationContainerRunning(getInput, input.Timeout)
	if err != nil {
		return nil, err
	}

	return applicationContainerInfo, nil
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
	// Timeout to delete an application container
	Timeout time.Duration
}

// DeleteApplicationContainer deletes the application container with the given name.
func (c *ApplicationContainerClient) DeleteApplicationContainer(input *DeleteApplicationContainerInput) error {
	// Call to delete the application container
	if err := c.deleteResource(input.Name); err != nil {
		return err
	}

	if input.Timeout == 0 {
		input.Timeout = WaitForApplicationContainerDeleteTimeout
	}

	// Wait for application container to be deleted
	return c.WaitForApplicationContainerDeleted(input, input.Timeout)
}

type UpdateApplicationContainerInput struct {
	// Name of the application container
	Name string
	// The additional fields needed for ApplicationContainer
	AdditionalFields UpdateApplicationContainerAdditionalFields
	// Name of the optional deployment file, which specifies memory, number of instances, and service bindings
	// Optional
	Deployment string
	// Name of the manifest file, required if this file is not packaged with the application
	// Optional
	Manifest string
	// Timeout for creating an application container
	Timeout time.Duration
}

type UpdateApplicationContainerAdditionalFields struct {
	// Location of the application archive file in Oracle Storage Cloud Service, in the format app-name/file-name
	// Optional
	ArchiveURL string
	// Comments on the application deployment
	Notes string
}

// Create a new Application Container from an ApplicationClient and an input struct.
// Returns a populated ApplicationContainer struct for the Application, and any errors
func (c *ApplicationContainerClient) UpdateApplicationContainer(input *UpdateApplicationContainerInput) (*ApplicationContainer, error) {

	files := make(map[string]string, 0)
	if input.Deployment != "" {
		files["deployment"] = input.Deployment
	}
	if input.Manifest != "" {
		files["manifest"] = input.Manifest
	}
	additionalFields := structs.Map(input.AdditionalFields)

	var applicationContainer *ApplicationContainer
	if err := c.updateResource(files, additionalFields, applicationContainer); err != nil {
		return nil, err
	}

	getInput := &GetApplicationContainerInput{
		Name: input.Name,
	}

	if input.Timeout == 0 {
		input.Timeout = WaitForApplicationContainerRunningTimeout
	}

	// Wait for application container to be ready and return the result
	applicationContainerInfo, err := c.WaitForApplicationContainerRunning(getInput, input.Timeout)
	if err != nil {
		return nil, err
	}

	return applicationContainerInfo, nil
}

// WaitForApplicationContainerRunning waits for an application container to be completely initialized and ready.
func (c *ApplicationContainerClient) WaitForApplicationContainerRunning(input *GetApplicationContainerInput, timeoutSeconds time.Duration) (*ApplicationContainer, error) {
	var info *ApplicationContainer
	err := c.client.WaitFor("Waiting for Application container to be ready", timeoutSeconds, func() (bool, error) {
		info, getErr := c.GetApplicationContainer(input)
		if getErr != nil {
			return false, getErr
		}
		c.client.DebugLogString(fmt.Sprintf("Application Container name is %v, Application Contiainer info is %+v", info.Name, info))
		switch s := info.Status; s {
		case string(ApplicationContainerStatusRunning): // Target State
			c.client.DebugLogString("Application Container Running")
			return true, nil
		case string(ApplicationContainerStatusNew):
			c.client.DebugLogString("Application Container New")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown application container state: %s, waiting", s))
			return false, nil
		}
	})
	return info, err
}

// WaitForApplicationContainerDeleted waits for an application container to be fully deleted.
func (c *ApplicationContainerClient) WaitForApplicationContainerDeleted(input *DeleteApplicationContainerInput, timeout time.Duration) error {
	return c.client.WaitFor("application container to be deleted", timeout, func() (bool, error) {
		var (
			info *ApplicationContainer
			err  error
		)
		getApplicationContainerInput := &GetApplicationContainerInput{
			Name: input.Name,
		}
		if info, err = c.GetApplicationContainer(getApplicationContainerInput); err != nil {
			if client.WasNotFoundError(err) {
				// Application Container could not be found, thus deleted
				return true, nil
			}
			// Currently the API returns a 400 with a message containing 404 and we check that here
			if strings.Contains(err.Error(), "404") {
				return true, nil
			}
			// Some other error occurred trying to get application container, exit
			return false, err
		}
		switch s := info.Status; s {
		case string(ApplicationContainerStatusRunning):
			return false, nil
		case string(ApplicationContainerStatusDestroyPending):
			c.client.DebugLogString(fmt.Sprintf("Destroy pending for application container state: %s, waiting", s))
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown application container state: %s, waiting", s))
			return false, fmt.Errorf("Unknown application container state: %s", s)
		}
	})
}
