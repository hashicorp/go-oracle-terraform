//
package mysql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-oracle-terraform/client"
)

// WaitForServiceInstanceReadyPollInterval is the default polling interval value for Creating a service instance and waiting for the instance to be ready
const WaitForServiceInstanceReadyPollInterval = time.Duration(60 * time.Second)

// WaitForServiceInstanceReadyTimeout is the default Timeout value for Creating a service instance and waiting for the instance to be ready
const WaitForServiceInstanceReadyTimeout = time.Duration(3600 * time.Second)

// WaitForServiceInstanceDeletePollInterval is the default polling value for deleting a service instance and waiting for the instance to be completely removed.
const WaitForServiceInstanceDeletePollInterval = time.Duration(60 * time.Second)

// WaitForServiceInstanceDeleteTimeout is the default Timeout value for deleting a service instance and waiting for the instance to be completely removed.
const WaitForServiceInstanceDeleteTimeout = time.Duration(3600 * time.Second)

// API URI Paths for the Container and Root objects.
var (
	ServiceInstanceContainerPath = "/paas/api/v1.1/instancemgmt/%[1]s/services/MySQLCS/instances/"
	ServiceInstanceResourcePath  = "/paas/api/v1.1/instancemgmt/%[1]s/services/MySQLCS/instances/%[2]s"
)

// ServiceInstanceClient is a client for the Service functions of the MySQL API.
type ServiceInstanceClient struct {
	ResourceClient
	Timeout      time.Duration
	PollInterval time.Duration
}

// ServiceInstanceClient obtains an ServiceInstanceClient which can be used to access to the
// Service Instance functions of the Database Cloud API
func (c *MySQLClient) ServiceInstanceClient() *ServiceInstanceClient {
	return &ServiceInstanceClient{
		ResourceClient: ResourceClient{
			MySQLClient:      c,
			ContainerPath:    ServiceInstanceContainerPath,
			ResourceRootPath: ServiceInstanceResourcePath,
		}}
}

// Constants for whether the Enterprise Monitor should be installed
type ServiceInstanceEnterpriseMonitor string

const (
	ServiceInstanceEnterpriseMonitorYes ServiceInstanceEnterpriseMonitor = "Yes"
	ServiceInstanceEnterpriseMonitorNo  ServiceInstanceEnterpriseMonitor = "No"
)

// Constants for the metering frequency for the MySQL CS Service Instance.
type ServiceInstanceMeteringFrequency string

const (
	ServiceInstanceMeteringFrequencyHourly  ServiceInstanceMeteringFrequency = "HOURLY"
	ServiceInstanceMeteringFrequencyMonthly ServiceInstanceMeteringFrequency = "MONTHLY"
)

// Constants for the Backup Destination
type ServiceInstanceBackupDestination string

const (
	ServiceInstanceBackupDestinationBoth ServiceInstanceBackupDestination = "BOTH"
	ServiceInstanceBackupDestinationNone ServiceInstanceBackupDestination = "NONE"
	ServiceInstanceBackupDestinationOSS  ServiceInstanceBackupDestination = "OSS"
)

// Constants for the state of the Service Instance State
type ServiceInstanceState string

const (
	ServiceInstanceReady        ServiceInstanceState = "READY"
	ServiceInstanceInitializing ServiceInstanceState = "INITIALIZING"
	ServiceInstanceStarting     ServiceInstanceState = "STARTING"
	ServiceInstanceStopping     ServiceInstanceState = "STOPPING"
	ServiceInstanceStopped      ServiceInstanceState = "STOPPED"
	ServiceInstanceConfiguring  ServiceInstanceState = "CONFIGURING"
	ServiceInstanceError        ServiceInstanceState = "ERROR"
	ServiceInstanceTerminating  ServiceInstanceState = "TERMINATING"
)

type ActivityLogInfo struct {
	ActivityLogId  string                `json:"activityLogId"`
	AuthDomain     string                `json:"authDomain"`
	AuthUser       string                `json:"authUser"`
	EndDate        string                `json:"endDate"`
	IdentityDomain string                `json:"identityDomain"`
	InitiatedBy    string                `json:"initiatedBy"`
	JobId          string                `json:"jobId"` //TODO: Possible int
	Messages       []ActivityMessageInfo `json:"messages"`
	OperationId    string                `json:"operationId"`
	OperationType  string                `json:"operationType"`
	ServiceId      string                `json:"serviceId"` //TODO: Possible int
	ServiceName    string                `json:"serviceName"`
	StartDate      string                `json:"startDate"`
	Status         string                `json:"status"`
	SummaryMessage string                `json:"summaryMessage"`
	ServiceType    string                `json:"serviceType"` // Not in API
}

type ActivityMessageInfo struct {
	ActivityDate string `json:"activityDate"`
	Messages     string `json:"message"`
}

type AttributeInfo struct {
	DisplayName  string `json:"displayName"`
	Type         string `json:"type"`
	Value        string `json:"value"`
	DisplayValue string `json:"displayValue"`
	IsKeyBinding bool   `json:"isKeyBinding"`
}

// ServiceInstance defines the instance information that is returned from the Get method
// when quering the instance
type ServiceInstance struct {
	ServiceId                    string                   `json:"serviceId"`
	ServiceUuid                  string                   `json:"serviceUuid"` // Not in API
	ServiceLogicalUuid           string                   `json:"serviceLogicalUuid"`
	ServiceName                  string                   `json:"serviceName"`
	ServiceType                  string                   `json:"serviceType"`
	DomainName                   string                   `json:"domainName"`
	ServiceVersion               string                   `json:"serviceVersion"`
	ReleaseVersion               string                   `json:"releaseVersion"`
	BaseReleaseVersion           string                   `json:"baseReleaseVersion"` // Not in API
	MetaVersion                  string                   `json:"metaVersion"`
	ServiceDescription           string                   `json:"serviceDescription"` // Not in API
	ServiceLevel                 string                   `json:"serviceLevel"`
	Subscription                 string                   `json:"subscription"`
	MeteringFrequency            string                   `json:"meteringFrequency"`
	Edition                      string                   `json:"edition"`
	TotalSSDStorage              int                      `json:"totalSSDStorage"`
	Status                       ServiceInstanceState     `json:"state"`
	ServiceStateDisplayName      string                   `json:"serviceStateDisplayName"`
	Clone                        bool                     `json:"clone"`
	Creator                      string                   `json:"creator"`
	CreationDate                 string                   `json:"creationDate"`
	IsBYOL                       bool                     `json:"isBYOL"`
	IsManaged                    bool                     `json:"isManaged"`
	IaasProvider                 string                   `json:"iaasProvider"`
	Attributes                   map[string]AttributeInfo `json:"attributes"`
	Components                   ComponentInfo            `json:"components"`
	ActivityLogs                 []ActivityLogInfo        `json:"activityLogs"`
	LayeringMode                 string                   `json:"layeringMode"`
	ServiceLevelDisplayName      string                   `json:"serviceLevelDisplayName"`
	EditionDisplayName           string                   `json:"editionDisplayName"`
	MeteringFrequencyDisplayName string                   `json:"meteringFrequencyDisplayName"`
	BackupFilePath               string                   `json:"BACKUP_FILE_PATH"`
	DataVolumeSize               string                   `json:"DATA_VOLUME_SIZE"`
	UseSSD                       string                   `json:"USE_SSD"`
	ProvisionEngine              string                   `json:"provisionEngine"`
	MysqlPort                    string                   `json:"MYSQL_PORT"`
	CloudStorageContainer        string                   `json:"CLOUD_STORAGE_CONTAINER"`
	BackupDestination            string                   `json:"BACKUP_DESTINATION"`
	TotalSharedStorage           int                      `json:"totalSharedStorage"`
	ComputeSiteName              string                   `json:"computeSiteName"`
	Patching                     PatchingInfo             `json:"patching"`

	// The reason for the instance going to error state, if available.
	ErrorReason string `json:"error_reason"`
}

type MysqlInfo struct {
	ServiceId                 string                    `json:"serviceId"`
	ComponentId               string                    `json:"componentId"`
	State                     string                    `json:"state"`
	ComponentStateDisplayName string                    `json:"componentStateDisplayName"`
	Version                   string                    `json:"version"`
	ComponentType             string                    `json:"componentType"` // Not in API
	CreationDate              string                    `json:"creationDate"`
	InstanceName              string                    `json:"instanceName"`
	InstanceRole              string                    `json:"instanceRole"`
	IsKeyComponent            bool                      `json:"isKeyComponent"` // Not in API
	Attributes                map[string]AttributeInfo  `json:"attributes"`
	VMInstances               map[string]VMInstanceInfo `json:"vmInstances"`
	AdminHostName             string                    `json:"adminHostName"`
	Hosts                     map[string]HostInfo       `json:"hosts"`       // Not in API
	DisplayName               string                    `json:"displayName"` // Not in API
	// hosts
	// paasServers
}

type VMInstanceInfo struct {
	VmId               string `json:"vmId"`
	Id                 int    `json:"id"`
	Uuid               string `json:"uuid"`
	HostName           string `json:"hostName"`
	Label              string `json:"label"`
	IPAddress          string `json:"ipAddress"`
	PublicIPAddress    string `json:"publicIpAddress"`
	UsageType          string `json:"usageType"`
	Role               string `json:"role"`
	ComponentType      string `json:"componentType"`
	State              string `json:"state"`
	VmStateDisplayName string `json:"vmStateDisplayName"`
	ShapeId            string `json:"shapeId"`
	TotalStorage       int    `json:"totalStorage"`
	CreationDate       string `json:"creationDate"`
	IsAdminNode        bool   `json:"isAdminNode"`
}

type HostInfo struct {
	Vmid               int                   `json:"vmId"`
	Id                 int                   `json:"id"`
	Uuid               string                `json:"uuid"`
	HostName           string                `json:"hostName"`
	Label              string                `json:"label"`
	UsageType          string                `json:"usageType"`
	Role               string                `json:"role"`
	ComponentType      string                `json:"componentType"`
	State              string                `json:"state"`
	VMStateDisplayName string                `json:"vmStateDisplayName"`
	ShapeId            string                `json:"shapeId"`
	TotalStorage       int                   `json:"totalStorage"`
	CreationDate       string                `json:"creationDate"`
	IsAdminNode        bool                  `json:"isAdminNode"`
	Servers            map[string]ServerInfo `json:"servers"`
}

type ServerInfo struct {
	ServerId               string `json:"serverId"`
	ServerName             string `json:"serverName"`
	ServerType             string `json:"serverType"`
	ServerRole             string `json:"serverRole"`
	State                  string `json:"state"`
	ServerStateDisplayName string `json:"serverStateDisplayName"`
	CreationDate           string `json:"creationDate"`
}

type StorageVolumeInfo struct {
	Name       string `json:"name"`
	Partitions string `json:"partitions"`
	Size       string `json:"size"`
}

type ComponentInfo struct {
	Mysql MysqlInfo `json:"mysql"`
}

type PatchingInfo struct {
	CurrentOperation      map[string]string `json:"currentOperation"` //hope this works
	TotalAvailablePatches string            `json:"totalAvailablePatches"`
}

type CreateServiceInstanceInput struct {
	ComponentParameters ComponentParameters `json:"componentParameters"`
	ServiceParameters   ServiceParameters   `json:"serviceParameters"`
}

type ComponentParameters struct {
	Mysql MySQLParameters `json:"mysql"`
}

type MySQLParameters struct {
	DBName                           string `json:"dbName"`
	DBStorage                        string `json:"dbStorage"`
	EnterpriseMonitor                string `json:"enterpriseMonitor,omitempty"`
	EnterpriseMonitorAgentPassword   string `json:"enterpriseMonitorAgentPassword,omitempty"`
	EnterpriseMonitorAgentUser       string `json:"enterpriseMonitorAgentUser,omitempty"`
	EnterpriseMonitorManagerPassword string `json:"enterpriseMonitorManagerPassword,omitempty"`
	EnterpriseMonitorManagerUser     string `json:"enterpriseMonitorManagerUser,omitempty"`
	MysqlCharset                     string `json:"mysqlCharset,omitempty"`
	MysqlCollation                   string `json:"mysqlCollation,omitempty"`
	MysqlEMPort                      string `json:"mysqlEMPort,omitempty"`
	MysqlPort                        string `json:"mysqlPort,omitempty"`
	// This is commented out. Although it shows up in the REST API documentation , attempts to set it have
	// resulted in errors.
	//	MysqlTimezone                    string `json:"mysqlTimezone, omitempty`
	// This is commented out. Although it shows up in the REST API documentation , attempts to set it have
	// resulted in errors.
	//	MysqlOptions                     string `json:"mysqlOptions,omitempty`
	MysqlUserName     string `json:"mysqlUserName"`
	MysqlUserPassword string `json:"mysqlUserPassword,omitempty"`
	Shape             string `json:"shape,omitempty"`
	SnapshotName      string `json:"snapshot,omitempty"`
	SourceServiceName string `json:"sourceServiceName,omitempty"`
}

type ServiceParameters struct {
	BackupDestination                 string `json:"backupDestination,omitempty"`
	CloudStorageContainer             string `json:"cloudStorageContainer,omitempty"`
	CloudStorageContainerAutoGenerate bool   `json:"cloudStorageContainerAutoGenerate,omitempty"`
	CloudStoragePassword              string `json:"cloudStoragePassword,omitempty"`
	CloudStorageUsername              string `json:"cloudStorageUser,omitempty"`
	MeteringFrequency                 string `json:"meteringFrequency,omitempty"`
	Region                            string `json:"region,omitempty"`
	ServiceDescription                string `json:"serviceDescription,omitempty"`
	ServiceName                       string `json:"serviceName"`
	VMPublicKeyText                   string `json:"vmPublicKeyText,omitempty"`
	VMUser                            string `json:"vmUser,omitempty"`
}

func (c *ServiceInstanceClient) CreateServiceInstance(input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	var (
		serviceInstance      *ServiceInstance
		serviceInstanceError error
	)

	if c.PollInterval == 0 {
		c.PollInterval = WaitForServiceInstanceReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceReadyTimeout
	}

	// Since these CloudStorageUsername and CloudStoragePassword are sensitive we'll read them
	// from the client if they haven't specified in the config.
	if input.ServiceParameters.CloudStorageContainer != "" && input.ServiceParameters.CloudStorageUsername == "" && input.ServiceParameters.CloudStoragePassword == "" {
		input.ServiceParameters.CloudStorageUsername = *c.ResourceClient.MySQLClient.client.UserName
		input.ServiceParameters.CloudStoragePassword = *c.ResourceClient.MySQLClient.client.Password
	}

	for i := 0; i < *c.MySQLClient.client.MaxRetries; i++ {
		serviceInstance, serviceInstanceError = c.startServiceInstance(input.ServiceParameters.ServiceName, input)
		if serviceInstanceError == nil {
			return serviceInstance, nil
		}
	}

	return nil, serviceInstanceError
}

func (c *ServiceInstanceClient) startServiceInstance(name string, input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	if err := c.createServiceInstanceResource(*input, nil); err != nil {
		return nil, err
	}

	// Call wait for instance ready now, as creating the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: name,
	}

	serviceInstance, serviceInstanceError := c.WaitForServiceInstanceRunning(getInput, c.PollInterval, c.Timeout)

	if serviceInstanceError != nil {
		c.client.DebugLogString(fmt.Sprintf(": Create Failed %s", serviceInstanceError))
		return nil, serviceInstanceError
	}

	return serviceInstance, nil
}

type GetServiceInstanceInput struct {
	// Name of the MySQL Cloud Service instance.
	// Required.
	Name string `json:"serviceId"`
}

// GetServiceInstance retrieves the ServiceInstance with the given name.
func (c *ServiceInstanceClient) GetServiceInstance(getInput *GetServiceInstanceInput) (*ServiceInstance, error) {
	var serviceInstance ServiceInstance
	if err := c.getServiceInstanceResource(getInput.Name, &serviceInstance); err != nil {
		return nil, err
	}

	return &serviceInstance, nil
}

// WaitForServiceInstanceRunning waits for an instance to be created and completely initialized and available.
func (c *ServiceInstanceClient) WaitForServiceInstanceRunning(input *GetServiceInstanceInput, pollingInterval time.Duration, timeoutSeconds time.Duration) (*ServiceInstance, error) {
	var info *ServiceInstance
	var getErr error

	err := c.client.WaitFor("service instance to be ready", pollingInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetServiceInstance(input)
		if getErr != nil {
			return false, getErr
		}

		c.client.DebugLogString(fmt.Sprintf("ServiceInstance [%s] waiting for ready state. Status : %s", info.ServiceId, info.Status))

		switch s := info.Status; s {

		case ServiceInstanceReady: // Target State
			return true, nil
		case ServiceInstanceInitializing:
			return false, nil
		case ServiceInstanceStarting:
			return false, nil
		default:
			return false, nil
		}
	})
	return info, err
}

type DeleteServiceInput struct {
	//Options string `json:"options,omitempty"`
}

func (c *ServiceInstanceClient) DeleteServiceInstance(serviceName string) error {
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceDeleteTimeout
	}
	if c.PollInterval == 0 {
		c.PollInterval = WaitForServiceInstanceDeletePollInterval
	}

	c.client.DebugLogString(fmt.Sprintf("Deleting Instance : %s", serviceName))

	deleteInput := &DeleteServiceInput{}

	deleteErr := c.deleteServiceInstanceResource(serviceName, deleteInput)
	if deleteErr != nil {
		c.client.DebugLogString(fmt.Sprintf(": Delete Failed %s", deleteErr))
		return deleteErr
	}

	// Call wait for instance deleted now, as deleting the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: serviceName,
	}

	// Wait for instance to be deleted
	return c.WaitForServiceInstanceDeleted(getInput, c.PollInterval, c.Timeout)
}

// WaitForServiceInstanceDeleted waits for a service instance to be fully deleted.
func (c *ServiceInstanceClient) WaitForServiceInstanceDeleted(input *GetServiceInstanceInput, pollingInterval time.Duration, timeoutSeconds time.Duration) error {
	return c.client.WaitFor("service instance to be deleted", pollingInterval, timeoutSeconds, func() (bool, error) {

		c.client.DebugLogString(fmt.Sprintf("Waiting to destroy instance : %s", input.Name))

		info, err := c.GetServiceInstance(input)
		if err != nil {
			if client.WasNotFoundError(err) {
				// Service Instance could not be found, thus deleted
				return true, nil
			}
			// Some other error occurred trying to get instance, exit
			return false, err
		}

		c.client.DebugLogString(fmt.Sprintf("ServiceInstance [%s] waiting for deletion . Status : %s", info.ServiceId, info.Status))

		switch s := info.Status; s {
		case ServiceInstanceError:
			return false, fmt.Errorf("Error stopping instance: %s", info.ErrorReason)
		case ServiceInstanceTerminating:
			return false, nil
		default:
			return false, nil
		}
	})
}
