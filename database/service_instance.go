package database

import (
	"fmt"
	"time"

	"log"

	"github.com/hashicorp/go-oracle-terraform/client"
)

const WaitForServiceInstanceReadyPollInterval = time.Duration(60 * time.Second)
const WaitForServiceInstanceReadyTimeout = time.Duration(3600 * time.Second)
const WaitForServiceInstanceDeletePollInterval = time.Duration(60 * time.Second)
const WaitForServiceInstanceDeleteTimeout = time.Duration(3600 * time.Second)
const ServiceInstanceDeleteRetry = 30

var (
	ServiceInstanceContainerPath = "/paas/service/dbcs/api/v1.1/instances/%s"
	ServiceInstanceResourcePath  = "/paas/service/dbcs/api/v1.1/instances/%s/%s"
)

// ServiceInstanceClient is a client for the Service functions of the Database API.
type ServiceInstanceClient struct {
	ResourceClient
	PollInterval time.Duration
	Timeout      time.Duration
}

// ServiceInstanceClient obtains an ServiceInstanceClient which can be used to access to the
// Service Instance functions of the Database Cloud API
func (c *DatabaseClient) ServiceInstanceClient() *ServiceInstanceClient {
	return &ServiceInstanceClient{
		ResourceClient: ResourceClient{
			DatabaseClient:   c,
			ContainerPath:    ServiceInstanceContainerPath,
			ResourceRootPath: ServiceInstanceResourcePath,
		}}
}

type ServiceInstanceEdition string

const (
	// SE: Standard Edition
	ServiceInstanceStandardEdition ServiceInstanceEdition = "SE"
	// EE: Enterprise Edition
	ServiceInstanceEnterpriseEdition ServiceInstanceEdition = "EE"
	// EE_HP: Enterprise Edition - High Performance
	ServiceInstanceEnterpriseEditionHighPerformance ServiceInstanceEdition = "EE_HP"
	// EE_EP: Enterprise Edition - Extreme Performance
	ServiceInstanceEnterpriseEditionExtremePerformance ServiceInstanceEdition = "EE_EP"
)

type ServiceInstanceLevel string

const (
	// PAAS: The Oracle Database Cloud Service service level
	ServiceInstanceLevelPAAS ServiceInstanceLevel = "PAAS"
	// PAAS_EXADATA: The Oracle Exadata Cloud Service service level
	ServiceInstanceLevelEXADATA ServiceInstanceLevel = "PAAS_EXADATA"
	// BASIC: The Oracle Database Cloud Service - Virtual Image service level
	ServiceInstanceLevelBasic ServiceInstanceLevel = "BASIC"
)

type ServiceInstanceBackupDestination string

const (
	// BOTH - Both Cloud Storage and Local Storage
	ServiceInstanceBackupDestinationBoth ServiceInstanceBackupDestination = "BOTH"
	// OSS - Cloud Storage only
	ServiceInstanceBackupDestinationOSS ServiceInstanceBackupDestination = "OSS"
	// NONE - None
	ServiceInstanceBackupDestinationNone ServiceInstanceBackupDestination = "NONE"
)

type ServiceInstanceNCharSet string

const (
	ServiceInstanceNCharSetUTF16 ServiceInstanceNCharSet = "AL16UTF16"
	ServiceInstanceNCharSetUTF8  ServiceInstanceNCharSet = "UTF8"
)

type ServiceInstanceType string

const (
	ServiceInstanceTypeDB ServiceInstanceType = "db"
)

type ServiceInstanceShape string

const (
	// Suportted OCI Classic Shapes
	// oc3: 1 OCPU, 7.5 GB memory
	ServiceInstanceShapeOC3 ServiceInstanceShape = "oc3"
	// oc4: 2 OCPUs, 15 GB memory
	ServiceInstanceShapeOC4 ServiceInstanceShape = "oc4"
	// oc5: 4 OCPUs, 30 GB memory
	ServiceInstanceShapeOC5 ServiceInstanceShape = "oc5"
	// oc6: 8 OCPUs, 60 GB memory
	ServiceInstanceShapeOC6 ServiceInstanceShape = "oc6"
	// oc7: 16 OCPUS, 120 GB memory
	ServiceInstanceShapeOC7 ServiceInstanceShape = "oc7"
	// oc1m: 1 OCPU, 15 GB memory
	ServiceInstanceShapeOC1M ServiceInstanceShape = "oc1m"
	// oc2m: 2 OCPUs, 30 GB memory
	ServiceInstanceShapeOC2M ServiceInstanceShape = "oc2m"
	// oc3m: 4 OCPUs, 60 GB memory
	ServiceInstanceShapeOC3M ServiceInstanceShape = "oc3m"
	// oc4m: 8 OCPUs, 120 GB memory
	ServiceInstanceShapeOC4M ServiceInstanceShape = "oc4m"
	// oc5m: 16 OCPUS, 240 GB memory
	ServiceInstanceShapeOC5M ServiceInstanceShape = "oc5m"

	// Supported OCI VM shapes
	// VM.Standard1.1: 1 OCPU, 7 GB memory
	ServiceInstanceShapeVMStandard1_1 ServiceInstanceShape = "VM.Standard1.1"
	// VM.Standard1.2: 2 OCPU, 14 GB memory
	ServiceInstanceShapeVMStandard1_2 ServiceInstanceShape = "VM.Standard1.2"
	// VM.Standard1.4: 4 OCPU, 28 GB memory
	ServiceInstanceShapeVMStandard1_4 ServiceInstanceShape = "VM.Standard1.4"
	// VM.Standard1.8: 8 OCPU, 56 GB memory
	ServiceInstanceShapeVMStandard1_8 ServiceInstanceShape = "VM.Standard1.8"
	// VM.Standard1.16: 16 OCPU, 112 GB memory
	ServiceInstanceShapeVMStandard1_16 ServiceInstanceShape = "VM.Standard1.16"
	// VM.Standard2.1: 1 OCPU, 15 GB memory
	ServiceInstanceShapeVMStandard2_1 ServiceInstanceShape = "VM.Standard2.1"
	// VM.Standard2.2: 2 OCPU, 30 GB memory
	ServiceInstanceShapeVMStandard2_2 ServiceInstanceShape = "VM.Standard2.2"
	// VM.Standard2.4: 4 OCPU, 60 GB memory
	ServiceInstanceShapeVMStandard2_4 ServiceInstanceShape = "VM.Standard2.4"
	// VM.Standard2.8: 8 OCPU, 120 GB memory
	ServiceInstanceShapeVMStandard2_8 ServiceInstanceShape = "VM.Standard2.8"
	// VM.Standard2.16: 16 OCPU, 240 GB memory
	ServiceInstanceShapeVMStandard2_16 ServiceInstanceShape = "VM.Standard2.16"
	// VM.Standard2.24: 24 OCPU, 320 GB memory
	ServiceInstanceShapeVMStandard2_24 ServiceInstanceShape = "VM.Standard2.24"

	// Supported OCI Bare Metal shapes
	// BM.Standard1.36: 36 OCPU, 256 GB memory
	ServiceInstanceShapeBMStandard1_36 ServiceInstanceShape = "BM.Standard1.36"
	// BM.Standard2.52: 52 OCPU, 768 GB memory
	ServiceInstanceShapeBMStandard2_52 ServiceInstanceShape = "BM.Standard2.52"
)

type ServiceInstanceSubscriptionType string

const (
	ServiceInstanceSubscriptionTypeHourly  ServiceInstanceSubscriptionType = "HOURLY"
	ServiceInstanceSubscriptionTypeMonthly ServiceInstanceSubscriptionType = "MONTHLY"
)

type ServiceInstanceVersion string

const (
	// 18.0.0.0
	ServiceInstanceVersion18000 ServiceInstanceVersion = "18.0.0.0"
	// 12.2.0.1
	ServiceInstanceVersion12201 ServiceInstanceVersion = "12.2.0.1"
	// 12.1.0.2
	ServiceInstanceVersion12102 ServiceInstanceVersion = "12.1.0.2"
	// 11.2.0.4
	ServiceInstanceVersion11204 ServiceInstanceVersion = "11.2.0.4"
)

type ServiceInstanceState string

const (
	// 	Configured: the service instance has been configured.
	ServiceInstanceConfigured ServiceInstanceState = "Configured"
	//	In Progress: the service instance is being created.
	ServiceInstanceInProgress ServiceInstanceState = "In Progress"
	//	Maintenance: the service instance is being stopped, started, restarted or scaled.
	ServiceInstanceMaintenance ServiceInstanceState = "Maintenance"
	//	Running: the service instance is running.
	ServiceInstanceRunning ServiceInstanceState = "Running"
	//	Stopped: the service instance is stopped.
	ServiceInstanceStopped ServiceInstanceState = "Stopped"
	//	Terminating: the service instance is being deleted.
	ServiceInstanceTerminating ServiceInstanceState = "Terminating"
)

type ServiceInstanceDatabaseTemplate string

const (
	// oltp: configures the starter database for a transactional workload
	ServiceInstanceTemplateOLPT ServiceInstanceDatabaseTemplate = "oltp"
	// dw: configures the starter database for a decision support or data warehouse workload
	ServiceInstanceTemplateDW ServiceInstanceDatabaseTemplate = "dw"
)

type ServiceInstance struct {
	// The URL to use to connect to Oracle Application Express on the service instance.
	ApexURL string `json:"apex_url"`
	// Applicable only in Oracle Cloud Infrastructure regions.
	// Name of the availability domain within the region where the Oracle Database Cloud Service instance is provisioned.
	AvailabilityDomain string `json:"availability_domain"`
	// The backup configuration of the service instance.
	BackupDestination string `json:"backup_destination"`
	// The version of cloud tooling for backup and recovery supported by the service instance.
	BackupSupportedVersion string `json:"backup_supported_version"`
	// The database character set of the database.
	CharSet string `json:"charset"`
	// The Oracle Storage Cloud container for backups.
	CloudStorageContainer string `json:"cloud_storage_container"`
	// Name of the cluster supporting the Exadata Cloud Service database deployment.
	ClusterName string `json:"cluster_names"`
	// The Oracle Cloud location housing the service instance.
	ComputeSiteName string `json:"compute_site_name"`
	// The connection descriptor for Oracle Net Services (SQL*Net).
	ConnectDescriptor string `json:"connect_descriptor"`
	// The connection descriptor for Oracle Net Services (SQL*Net) with IP addresses instead of host names.
	ConnectorDescriptorWithPublicIP string `json:"connect_descriptor_with_public_ip"`
	// The user name of the Oracle Cloud user who created the service instance.
	CreatedBy string `json:"created_by"`
	// The job id of the job that created the service instance.
	CreationJobID string `json:"creation_job_id"`
	// The date-and-time stamp when the service instance was created.
	CreationTime string `json:"creation_time"`
	// The Oracle Database version on the service instance, including the patch level.
	CurrentVersion string `json:"current_version"`
	// The URL to use to connect to Oracle DBaaS Monitor on the service instance.
	DBAASMonitorURL string `json:"dbaasmonitor_url"`
	// The description of the service instance, if one was provided when the instance was created.
	Description string `json:"description"`
	// The software edition of the service instance.
	Edition ServiceInstanceEdition `json:"edition"`
	// The URL to use to connect to Enterprise Manager on the service instance.
	EMURL string `json:"em_url"`
	// Indicates whether the service instance hosts an Oracle Data Guard configuration.
	FailoverDatabase bool `json:"failover_database"`
	// The URL to use to connect to the Oracle GlassFish Server Administration Console on the service instance.
	GlassFishURL string `json:"glassfish_url"`
	// Data Guard Role of the on-premise instance in Oracle Hybrid Disaster Recovery configuration.
	HDGPremIP string `json:"hdgPremIP"`
	// Indicates whether the service instance hosts an Oracle Hybrid Disaster Recovery configuration.
	HybridDG string `json:"hybrid_db"`
	// The identity domain housing the service instance.
	IdentityDomain string `json:"identity_domain"`
	// This attribute is only applicable to accounts where regions are supported.
	// The three-part name of an IP network to which the service instance is added.
	// For example: /Compute-identity_domain/user/object
	IPNetwork string `json:"ipNetwork"`
	// Groups one or more IP reservations in use on this service instance.
	// This attribute is only applicable to accounts where regions are supported.
	// This attribute is returned when you set the ?outputLevel=verbose query parameter.
	IPReservations string `json:"ipReservations"`
	// Indicates whether service instance was provisioned with the 'Bring Your Own License' (BYOL) option.
	IsBYOL bool `json:"isBYOL"`
	// The Oracle Java Cloud Service instances using this Database Cloud Service instance.
	JAASInstancesUsingService string `json:"jaas_instances_using_service"`
	// The date-and-time stamp when the service instance was last modified.
	LastModifiedTime string `json:"last_modified_time"`
	// The service level of the service instance.
	Level ServiceInstanceLevel `json:"level"`
	// The listener port for Oracle Net Services (SQL*Net) connections.
	ListenerPort int `json:"listener_port"`
	// The national character set of the database.
	NCharSet string `json:"ncharset"`
	// List of compute nodes that host database instances for the database deployment.
	NodeList string `json:"nodelist"`
	// The number of Oracle Compute Service IP reservations assigned to the service instance.
	NumIPReservations int `json:"num_ip_reservations"`
	// The number of compute nodes in the service instance.
	NumNodes string `json:"num_nodes"`
	// The name of the default PDB (pluggable database) created when the service instance was created.
	PDBName string `json:"pdbName"`
	// Indicates whether the service instance hosts an Oracle RAC database.
	RACDatabase bool `json:"rac_database"`
	// This attribute is only applicable to accounts where regions are supported.
	// Location where the service instance is provisioned (only for accounts where regions are supported).
	Region string `json:"region"`
	// The name of the service instance.
	Name string `json:"service_name"`
	// The REST endpoint URI of the service instance.
	URI string `json:"service_uri"`
	// The Oracle Compute Cloud shape of the service instance.
	Shape string `json:"shape"`
	// The SID of the database.
	SID string `json:"sid"`
	// The version of the cloud tooling service manager plugin supported by the service instance.
	SMPluginVersion string `json:"sm_plugin_version"`
	// The status of the service instance
	Status ServiceInstanceState `json:"status"`
	// Applicable only in Oracle Cloud Infrastructure regions.
	// Name of the subnet within the region where the Oracle Database Cloud Service instance is provisioned.
	Subnet string `json:"subnet"`
	// The subscription name
	SubscriptionName string `json:"subscription_name"`
	// The billing frequency of the service instance; either MONTHLY or HOURLY.
	SubscriptionType ServiceInstanceSubscriptionType `json:"subscriptionType"`
	// The time zone of the operating system.
	Timezone string `json:"timezone"`
	// For service instances hosting an Oracle RAC database, the size in GB of the storage shared
	// and accessed by the nodes of the RAC database.
	TotalSharedStorage int `json:"total_shared_storage"`
	// Indicates whether the service instance was provisioned with high performance storage.
	UseHighPerformanceStorage bool `json:"useHighPerformanceStorage"`
	// The Oracle Database version on the service instance.
	Version string `json:"version"`
}

type CreateServiceInstanceInput struct {
	// Name of the availability domain within the region where the Oracle Database Cloud Service instance is to be provisioned.
	// Optional
	AvailabilityDomain string `json:"availabilityDomain,omitempty"`
	// Name of the cluster supporting the Exadata Cloud Service database deployment.
	// Optional. Exadata Cloud Service only.
	ClusterName string `json:"clusterName,omitempty"`
	// The template to use for the starter database. `dw` or `oltp`.
	// Include this parameter only when creating a starter database. The starter database is the first Exadata Cloud Service database deployment that you create on your Exadata system. Subsequent databases are created with a standardized database configuration.
	// Optional. Exadata Cloud Service only.
	DatabaseTemplate string `json:"dbTemplate,omitempty"`
	// Free-form text that provides additional information about the service instance.
	// Optional.
	Description string `json:"description,omitempty"`
	// Database edition for the service instance:
	// If you specify SE, a Standard Edition 2 database is created if you specify 12.2.0.1 or 12.1.0.2
	// for version and a Standard Edition One database is created if you specify 11.2.0.4 for version.
	// Edition must be Enterprise Edition - Extreme Performance to configure the Database
	// Cloud Service instance as Cluster Database.
	// Required.
	Edition ServiceInstanceEdition `json:"edition"`
	// Specify if you want an email notification sent upon successful or unsuccessful completion of the instance-creation operation.
	// When true, you must also specify notificationEmail. Valid values are true and false. Default value is false.
	// Optional
	EnableNotification bool `json:"enableNotification,omitempty"`
	// Required if level is `PAAS_EXADATA`
	// Name of the Exadata system on which to create the Exadata Cloud Service database deployment.
	// Optional. Exadata Cloud Service only.
	ExadataSystemName string `json:"exadataSystemName,omitempty"`
	// Specify if you want to use an existing perpetual license to Oracle Database to establish the right to use Oracle Database on the new instance.
	// When true, your Oracle Cloud account will be charged a lesser amount for the new instance because the right to use Oracle Database is covered by your perpetual license agreement.
	// Valid values are true and false. Default value is false.
	// Optional
	IsBYOL bool `json:"isBYOL,omitempty"`
	// This parameter is not available on Oracle Cloud at Customer.
	// Applicable only if region is an Oracle Cloud Infrastructure Classic region.
	// The three-part name of a custom IP network to use. For example: /Compute-identity_domain/user/object.
	// A region must be specified in order to use ipNetwork. Only IP networks created in the specified region can be used.
	// ipNetwork cannot be used with ipReservations.
	// Optional
	IPNetwork string `json:"ipNetwork,omitempty"`
	// Applicable only if region is an Oracle Cloud Infrastructure Classic region.
	// A single IP reservation name or multiple IP reservation names separated by commas. Only IP reservations created in the specified region can be used.
	// When IP reservations are used, all compute nodes of an instance must be provisioned with IP reservations, so the number of names in ipReservations must match the number of compute nodes in the service instance.
	// Optional
	IPReservations []string `json:"ipReservations,omitempty"`
	// Service level for the service instance
	// Required.
	Level ServiceInstanceLevel `json:"level"`
	// Array of one JSON object that specifies configuration details of the services instance.
	// This array is not required if the level value is BASIC.
	// Required if level value is PAAS.
	Parameter ParameterInput `json:"-"`
	// Name of Database Cloud Service instance. The service name:
	// Must not exceed 50 characters.
	// Must start with a letter.
	// Must contain only letters, numbers, or hyphens.
	// Must not contain any other special characters.
	// Must be unique within the identity domain.
	// Required.
	Name string `json:"serviceName"`
	// Specifies the list of compute nodes that host database instances for the database deployment.
	// Separate compute node names with a comma. If nodelist is not specified the database is deployed across all compute nodes
	// Optional. Exadata Cloud Service only.
	NodeList string `json:"nodelist,omitempty"`
	// Required if enableNotification is set to true.
	// The email address to send completion notifications to.
	// This parameter is not available on Oracle Cloud at Customer.
	// Optional
	NotificationEmail string `json:"notificationEmail,omitempty"`
	// Applicable only to accounts that support regions.
	// Name of the Oracle Cloud Infrastructure or Oracle Cloud Infrastructure Classic region where the Oracle Database Cloud Service instance is to be provisioned.
	// Optional
	Region string `json:"region,omitempty"`
	// Desired compute shape. A shape defines the number of Oracle Compute Units (OCPUs) and amount
	// of memory (RAM).
	// Required.
	Shape ServiceInstanceShape `json:"shape,omitempty"`
	// Required if region is an Oracle Cloud Infrastructure region.
	// Name of the subnet within the region where the Oracle Database Cloud Service instance is to be provisioned.
	// Optional
	Subnet string `json:"subnet,omitempty"`
	// Billing unit. Valid values are:
	// HOURLY: Pay only for the number of hours used during your billing period. This is the default.
	// MONTHLY: Pay one price for the full month irrespective of the number of hours used.
	// Required.
	SubscriptionType ServiceInstanceSubscriptionType `json:"subscriptionType"`
	// Applicable only in Oracle Cloud Infrastructure regions.
	// Array of one JSON object that specifies configuration details of the standby database when failoverDatabase is set to true. disasterRecovery must be set to true.
	Standbys []StandBy `json:"standbys,omitempty"`
	// Specify if high performance storage should be used for the Database Cloud Service instance. Default data storage will allocate your database
	// block storage on spinning devices. By checking this box, your block storage will be allocated on solid state devices. Valid values are true and false.
	// Default value is false.
	// Optional
	UseHighPerformanceStorage bool `json:"useHighPerformanceStorage,omitempty"`
	// Oracle Database software version
	// Required.
	Version ServiceInstanceVersion `json:"version"`
	// Public key for the secure shell (SSH). This key will be used for authentication when
	// connecting to the Database Cloud Service instance using an SSH client. You generate an
	// SSH public-private key pair using a standard SSH key generation tool. See Connecting to
	// a Compute Node Through Secure Shell (SSH) in Using Oracle Database Cloud Service.
	// Required.
	VMPublicKey string `json:"vmPublicKeyText"`
}

type StandBy struct {
	// Name of the availability domain within the region where the standby database of the Oracle Database
	// Cloud Service instance is to be provisioned.
	// Required.
	AvailabilityDomain string `json:"availabilityDomain"`
	// Name of the subnet within the region where the standby database of the Oracle Database Cloud Service
	// instance is to be provisioned.
	// Required.
	Subnet string `json:"subnet"`
}

type CreateServiceInstanceRequest struct {
	CreateServiceInstanceInput
	ParameterRequest []ParameterRequest `json:"parameters"`
}

type ParameterInput struct {
	AdditionalParameters AdditionalParameters `json:"additionalParams,omitempty"`
	// Password for Oracle Database administrator users sys and system. The password must meet the following requirements:
	// Starts with a letter
	// Is between 8 and 30 characters long
	// Contains letters, at least one number, and optionally, any number of these special characters: dollar sign ($), pound sign (#), and underscore (_).
	// Required
	AdminPassword string `json:"adminPassword"`
	//Backup destination.
	// Required
	BackupDestination ServiceInstanceBackupDestination `json:"backupDestination"`
	// Character Set for the Database Cloud Service Instance.
	// Valid values are AL32UTF8, AR8ADOS710, AR8ADOS720, AR8APTEC715, AR8ARABICMACS,
	// AR8ASMO8X, AR8ISO8859P6, AR8MSWIN1256, AR8MUSSAD768, AR8NAFITHA711, AR8NAFITHA721,
	// AR8SAKHR706, AR8SAKHR707, AZ8ISO8859P9E, BG8MSWIN, BG8PC437S, BLT8CP921, BLT8ISO8859P13,
	// BLT8MSWIN1257, BLT8PC775, BN8BSCII, CDN8PC863, CEL8ISO8859P14, CL8ISO8859P5, CL8ISOIR111,
	// CL8KOI8R, CL8KOI8U, CL8MACCYRILLICS, CL8MSWIN1251, EE8ISO8859P2, EE8MACCES, EE8MACCROATIANS,
	// EE8MSWIN1250, EE8PC852, EL8DEC, EL8ISO8859P7, EL8MACGREEKS, EL8MSWIN1253, EL8PC437S, EL8PC851,
	// EL8PC869, ET8MSWIN923, HU8ABMOD, HU8CWI2, IN8ISCII, IS8PC861, IW8ISO8859P8, IW8MACHEBREWS,
	// IW8MSWIN1255, IW8PC1507, JA16EUC, JA16EUCTILDE, JA16SJIS, JA16SJISTILDE, JA16VMS, KO16KSC5601,
	// KO16KSCCS, KO16MSWIN949, LA8ISO6937, LA8PASSPORT, LT8MSWIN921, LT8PC772, LT8PC774, LV8PC1117,
	// LV8PC8LR, LV8RST104090, N8PC865, NE8ISO8859P10, NEE8ISO8859P4, RU8BESTA, RU8PC855, RU8PC866,
	// SE8ISO8859P3, TH8MACTHAIS, TH8TISASCII, TR8DEC, TR8MACTURKISHS, TR8MSWIN1254, TR8PC857, US7ASCII,
	// US8PC437, UTF8, VN8MSWIN1258, VN8VN3, WE8DEC, WE8DG, WE8ISO8859P1, WE8ISO8859P15, WE8ISO8859P9,
	// WE8MACROMAN8S, WE8MSWIN1252, WE8NCR4970, WE8NEXTSTEP, WE8PC850, WE8PC858, WE8PC860, WE8ROMAN8,
	// ZHS16CGB231280, ZHS16GBK, ZHT16BIG5, ZHT16CCDC, ZHT16DBT, ZHT16HKSCS, ZHT16MSWIN950, ZHT32EUC,
	// ZHT32SOPS, ZHT32TRIS.
	// Default value is AL32UTF8.
	// Optional.
	CharSet string `json:"charset,omitempty"`
	// Name of the Oracle Storage Cloud Service container used to provide storage for your service
	// instance backups. Use the following format to specify the container name:
	// <storageservicename>-<storageidentitydomain>/<containername>
	// Notes:
	// An Oracle Storage Cloud Service container is not required when provisioning a Database
	// Cloud Service - Virtual Image instance.
	// Do not use an Oracle Storage Cloud container that you use to back up Database Cloud Service
	// instances for any other purpose. For example, do not also use it to back up Oracle Java Cloud
	// Service instances. Using the container for multiple purposes can result in billing errors.
	// Optional.
	CloudStorageContainer string `json:"cloudStorageContainer,omitempty"`
	// Password for the Oracle Storage Cloud Service administrator.
	// Optional.
	CloudStoragePassword string `json:"cloudStoragePwd,omitempty"`
	// Username for the Oracle Storage Cloud Service administrator.
	// Optional.
	CloudStorageUsername string `json:"cloudStorageUser,omitempty"`
	// Specify if the given cloudStorageContainer is to be created if it does not already exist.
	// Default value is false.
	// Optional.
	CreateStorageContainerIfMissing bool `json:"createStorageContainerIfMissing,omitempty"`
	// Specify if an Oracle Data Guard configuration is created using the Disaster Recovery option
	// or the High Availability option.
	// true - The Disaster Recovery option is used, which places the compute node hosting the primary
	// database and the compute node hosting the standby database in compute zones of different data centers.
	// false - The High Availability option is used, which places the compute node hosting the primary
	// database and the compute node hosting the standby database in different compute zones of the same
	// data center.
	// Default value is false.
	// This option is applicable only when failoverDatabase is set to true.
	// Optional
	DisasterRecovery bool `json:"-"`
	// Specify if an Oracle Data Guard configuration comprising a primary database and a
	// standby database is created. Default value is false.
	// You cannot set both failoverDatabase and isRac to false.
	// Optional
	FailoverDatabase bool `json:"-"`
	// Specify if the database should be configured for use as the replication database of an
	// Oracle GoldenGate Cloud Service instance. Default value is false.
	// You cannot set goldenGate to true if either isRac or failoverDatabase is set to true.
	// Optional
	GoldenGate bool `json:"-"`
	// Specify if an Oracle Hybrid Disaster Recovery configuration comprising a primary database on customer premisesand a standby database in Oracle Public Cloud should be configured.
	// Valid values are yes and no. Default value is no.
	// You cannot set failoverDatabase or isRac to yes if Hybrid Disaster Recovery options is chosen.
	// Optional
	HDG bool `json:"-"`
	// Name of the Oracle Storage Cloud Service container where the backup from on-premise instance is stored. This parameter is required if hdg is set to yes.
	// Optional
	HDGCloudStorageContainer string `json:"hdgCloudStorageContainer,omitempty"`
	// Password of the Oracle Cloud user specified in hdgCloudStorageUser. This parameter is required if hdg is set to yes.
	// Optional
	HDGCloudStoragePassword string `json:"hdgCloudStoragePassword,omitempty"`
	// User name of an Oracle Cloud user who has read access to the container specified in hdgCloudStorageContainer. This parameter is required if hdg is set to yes.
	// Optional
	HDGCloudStorageUser string `json:"hdgCloudStorageUser,omitempty"`
	// Specify if the service instance's database should, after the instance is created, be replaced
	// by a database stored in an existing cloud backup that was created using Oracle Database Backup
	// Cloud Service. Default value is false.
	// Optional
	IBKUP bool `json:"-"`
	// Name of the Oracle Storage Cloud Service container where the existing cloud backup is stored. This parameter is required if ibkup is set to yes and ibkupOnPremise is set to yes.
	// Optional
	IBKUPCloudStorageContainer string `json:"ibkupCloudStorageContainer,omitempty"`
	// Name of the Oracle Storage Cloud Service container where the existing cloud backup is stored.
	// This parameter is required if ibkup is set to yes.
	// Optional
	IBKUPCloudStoragePassword string `json:"ibkupCloudStoragePassword,omitempty"`
	// User name of an Oracle Cloud user who has read access to the container specified in
	// ibkupCloudStorageContainer.
	// This parameter is required if ibkup is set to yes.
	// Optional
	IBKUPCloudStorageUser string `json:"ibkupCloudStorageUser,omitempty"`
	// Database id of the database from which the existing cloud backup was created.
	// This parameter is required if ibkup is set to yes.
	// Optional
	IBKUPDatabaseID string `json:"ibkupDatabaseID,omitempty"`
	// Password used to create the existing, password-encrypted cloud backup.
	// This password is used to decrypt the backup.
	// This parameter is required if ibkup is set to yes.
	// Optional
	IBKUPDecryptionKey string `json:"ibkupDecryptionKey,omitempty"`
	// Specify if the existing cloud backup being used to replace the database is from an on-premises database or another Database Cloud Service instance.
	// Valid values are true for an on-premises database and false for a Database Cloud Service instance. Default value is true.
	// Optional
	IBKUPOnPremise bool `json:"ibkupOnPremise,omitempty"`
	// Oracle Databsae Cloud Service instance name from which the database of new Oracle Database Cloud Service instance should be created.
	// This parameter is required if ibkup is set to yes and ibkupOnPremise is set to no.
	// Optional
	IBKUPServiceID string `json:"ibkupServiceID,omitempty"`
	// String containing the xsd:base64Binary representation of the cloud backup's wallet archive file.
	// Optional
	IBKUPWalletFileContent string `json:"ibkupWalletFileContent,omitempty"`
	// Specify if a cluster database using Oracle Real Application Clusters should be configured.
	// Valid values are yes and no. Default value is no.
	// Optional
	IsRAC bool `json:"-"`
	// National Character Set for the Database Cloud Service instance.
	// Default value is AL16UTF16.
	// Optional.
	NCharSet ServiceInstanceNCharSet `json:"ncharset,omitempty"`
	// Note: This attribute is valid when Database Cloud Service instance is configured with version 12c.
	// Pluggable Database Name for the Database Cloud Service instance.
	// Default value is pdb1.
	// Optional.
	PDBName string `json:"pdbName,omitempty"`
	// Database Name (sid) for the Database Cloud Service instance.
	// Default value is ORCL.
	// Required.
	SID string `json:"sid"`
	// The name of the snapshot of the service instance specified by sourceServiceName
	// that is to be used to create a "snapshot clone".
	// This parameter is valid only if sourceServiceName is specified.
	// Optional.
	SnapshotName string `json:"snapshotName,omitempty"`
	// When present, indicates that the service instance should be created as a
	// "snapshot clone" of another service instance. Provide the name of the existing service
	// instance whose snapshot is to be used.
	// Optional.
	SourceServiceName string `json:"sourceServiceName,omitempty"`
	// Time Zone for the Database Cloud Service instance.
	// Valid values are Africa/Cairo, Africa/Casablanca, Africa/Harare, Africa/Monrovia,
	// Africa/Nairobi, Africa/Tripoli, Africa/Windhoek, America/Araguaina, America/Asuncion,
	// America/Bogota, America/Caracas, America/Chihuahua, America/Cuiaba, America/Denver,
	// America/Fortaleza, America/Guatemala, America/Halifax, America/Manaus, America/Matamoros,
	// America/Monterrey, America/Montevideo, America/Phoenix, America/Santiago, America/Tijuana,
	// Asia/Amman, Asia/Ashgabat, Asia/Baghdad, Asia/Baku, Asia/Bangkok, Asia/Beirut, Asia/Calcutta,
	// Asia/Damascus, Asia/Dhaka, Asia/Irkutsk, Asia/Jerusalem, Asia/Kabul, Asia/Karachi,
	// Asia/Kathmandu, Asia/Krasnoyarsk, Asia/Magadan, Asia/Muscat, Asia/Novosibirsk, Asia/Riyadh,
	// Asia/Seoul, Asia/Shanghai, Asia/Singapore, Asia/Taipei, Asia/Tehran, Asia/Tokyo, Asia/Ulaanbaatar,
	// Asia/Vladivostok, Asia/Yakutsk, Asia/Yerevan, Atlantic/Azores, Australia/Adelaide,
	// Australia/Brisbane, Australia/Darwin, Australia/Hobart, Australia/Perth, Australia/Sydney,
	// Brazil/East, Canada/Newfoundland, Canada/Saskatchewan, Europe/Amsterdam, Europe/Athens,
	// Europe/Dublin, Europe/Helsinki, Europe/Istanbul, Europe/Kaliningrad, Europe/Moscow,
	// Europe/Paris, Europe/Prague, Europe/Sarajevo, Pacific/Auckland, Pacific/Fiji, Pacific/Guam,
	// Pacific/Honolulu, Pacific/Samoa, US/Alaska, US/Central, US/Eastern, US/East-Indiana,
	// US/Pacific, UTC.
	// Default value is UTC.
	// Optional.
	Timezone string `json:"timezone,omitempty"`
	// Component type to which the set of parameters applies.
	// Valid values are: db - Oracle Database
	// Required.
	Type ServiceInstanceType `json:"type"`
	// Storage size for data (in GB). Minimum value is 15. Maximum value depends on the backup
	// destination: if BOTH is specified, the maximum value is 1200; if OSS or NONE is specified,
	// the maximum value is 2048.
	// Required.
	UsableStorage string `json:"usableStorage"`
}

type ParameterRequest struct {
	ParameterInput
	DisasterRecoveryString string `json:"disasterRecovery,omitempty"`
	FailoverDatabaseString string `json:"failoverDatabase,omitempty"`
	GoldenGateString       string `json:"goldenGate,omitempty"`
	HDGString              string `json:"hdg,omitempty"`
	IsRACString            string `json:"isRac,omitempty"`
	IBKUPString            string `json:"ibkup,omitempty"`
}

type AdditionalParameters struct {
	// Indicates whether to include the Demos PDB
	// Optional
	DBDemo string `json:"db_demo,omitempty"`
}

// CreateServiceInstance creates a new ServiceInstace.
func (c *ServiceInstanceClient) CreateServiceInstance(input *CreateServiceInstanceInput) (*ServiceInstance, error) {
	if c.PollInterval == 0 {
		c.PollInterval = WaitForServiceInstanceReadyPollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceReadyTimeout
	}
	// Since these CloudStorageUsername and CloudStoragePassword are sensitive we'll read them
	// from the client if they haven't specified in the config.
	if input.Parameter.CloudStorageContainer != "" {
		if input.Parameter.CloudStorageUsername == "" {
			input.Parameter.CloudStorageUsername = *c.ResourceClient.DatabaseClient.client.UserName
		}
		if input.Parameter.CloudStoragePassword == "" {
			input.Parameter.CloudStoragePassword = *c.ResourceClient.DatabaseClient.client.Password
		}
	}
	if input.Parameter.IBKUPCloudStorageContainer != "" {
		if input.Parameter.IBKUPCloudStorageUser == "" {
			input.Parameter.IBKUPCloudStorageUser = *c.ResourceClient.DatabaseClient.client.UserName
		}
		if input.Parameter.IBKUPCloudStoragePassword == "" {
			input.Parameter.IBKUPCloudStoragePassword = *c.ResourceClient.DatabaseClient.client.Password
		}
	}
	if input.Parameter.HDGCloudStorageContainer != "" {
		if input.Parameter.HDGCloudStorageUser == "" {
			input.Parameter.HDGCloudStorageUser = *c.ResourceClient.DatabaseClient.client.UserName
		}
		if input.Parameter.HDGCloudStoragePassword == "" {
			input.Parameter.HDGCloudStoragePassword = *c.ResourceClient.DatabaseClient.client.Password
		}
	}

	// Create request where bools(true/false) are switched to strings(yes/no).
	request := createRequest(input)

	serviceInstance, err := c.startServiceInstance(request.Name, request)
	if err != nil {
		return serviceInstance, fmt.Errorf("unable to create Database Service Instance %q: %+v", request.Name, err)
	}
	return serviceInstance, nil
}

func createRequest(input *CreateServiceInstanceInput) *CreateServiceInstanceRequest {
	parameterRequest := ParameterRequest{
		ParameterInput:         input.Parameter,
		DisasterRecoveryString: convertOracleBool(input.Parameter.DisasterRecovery),
		FailoverDatabaseString: convertOracleBool(input.Parameter.FailoverDatabase),
		GoldenGateString:       convertOracleBool(input.Parameter.GoldenGate),
		HDGString:              convertOracleBool(input.Parameter.HDG),
		IsRACString:            convertOracleBool(input.Parameter.IsRAC),
		IBKUPString:            convertOracleBool(input.Parameter.IBKUP),
	}
	request := &CreateServiceInstanceRequest{
		CreateServiceInstanceInput: *input,
		ParameterRequest:           []ParameterRequest{parameterRequest},
	}

	return request
}

func (c *ServiceInstanceClient) startServiceInstance(name string, input *CreateServiceInstanceRequest) (*ServiceInstance, error) {
	if err := c.createResource(*input, nil); err != nil {
		return nil, err
	}

	// Call wait for instance ready now, as creating the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: name,
	}

	// Wait for the service instance to be running and return the result
	// Don't have to unqualify any objects, as the GetServiceInstance method will handle that
	serviceInstance, serviceInstanceError := c.WaitForServiceInstanceRunning(getInput, c.PollInterval, c.Timeout)
	// If the service instance enters an error state we need to delete the instance and retry
	if serviceInstanceError != nil {
		deleteInput := &DeleteServiceInstanceInput{
			Name: name,
		}
		err := c.DeleteServiceInstance(deleteInput)
		if err != nil {
			return nil, fmt.Errorf("Error deleting service instance %s: %s", name, err)
		}
		return nil, serviceInstanceError
	}
	return serviceInstance, nil
}

// WaitForServiceInstanceRunning waits for a service instance to be completely initialized and available.
func (c *ServiceInstanceClient) WaitForServiceInstanceRunning(input *GetServiceInstanceInput, pollInterval, timeoutSeconds time.Duration) (*ServiceInstance, error) {
	var info *ServiceInstance
	var getErr error
	err := c.client.WaitFor("service instance to be ready", pollInterval, timeoutSeconds, func() (bool, error) {
		info, getErr = c.GetServiceInstance(input)
		if getErr != nil {
			return false, getErr
		}
		c.client.DebugLogString(fmt.Sprintf("Service instance name is %v, Service instance info is %+v", info.Name, info))
		switch s := info.Status; s {
		case ServiceInstanceRunning: // Target State
			c.client.DebugLogString("Service Instance Running")
			return true, nil
		case ServiceInstanceConfigured:
			c.client.DebugLogString("Service Instance Configured")
			return false, nil
		case ServiceInstanceInProgress:
			c.client.DebugLogString("Service Instance is being created")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
	return info, err
}

type GetServiceInstanceInput struct {
	// Name of the Database Cloud Service instance.
	// Required.
	Name string `json:"serviceId"`
}

// GetServiceInstance retrieves the SeriveInstance with the given name.
func (c *ServiceInstanceClient) GetServiceInstance(getInput *GetServiceInstanceInput) (*ServiceInstance, error) {
	var serviceInstance ServiceInstance
	if err := c.getResource(getInput.Name, &serviceInstance); err != nil {
		return nil, err
	}

	return &serviceInstance, nil
}

type DeleteServiceInstanceInput struct {
	// Name of the Database Cloud Service instance.
	// Required.
	Name string
	// Flag that when set to true deletes all backups of the service instance from Oracle Cloud Storage container.
	// Use caution in specifying this option. If this option is specified, instance can not be recovered as all backups
	// will be deleted. This option is not currently supported for Cluster Databases.
	// Default value is false.
	// Optional
	DeleteBackup bool
}

func (c *ServiceInstanceClient) DeleteServiceInstance(input *DeleteServiceInstanceInput) error {
	if c.PollInterval == 0 {
		c.PollInterval = WaitForServiceInstanceDeletePollInterval
	}
	if c.Timeout == 0 {
		c.Timeout = WaitForServiceInstanceDeleteTimeout
	}
	// Call to delete the service instance
	// If delete fails, rerun it in case the instance still has not been setup correctly.
	// An instance takes additional time to setup after it's configured.
	var deleteErr error
	for i := 0; i < ServiceInstanceDeleteRetry; i++ {
		if deleteErr = c.deleteResource(input.Name, input.DeleteBackup); deleteErr != nil {
			log.Printf("Error during delete, waiting 30s: %+v", deleteErr)
			time.Sleep(30 * time.Second)
			continue
		}
		break
	}
	if deleteErr != nil {
		return deleteErr
	}

	// Call wait for instance deleted now, as deleting the instance is an eventually consistent operation
	getInput := &GetServiceInstanceInput{
		Name: input.Name,
	}

	// Wait for instance to be deleted
	return c.WaitForServiceInstanceDeleted(getInput, c.PollInterval, c.Timeout)
}

// WaitForServiceInstanceDeleted waits for a service instance to be fully deleted.
func (c *ServiceInstanceClient) WaitForServiceInstanceDeleted(input *GetServiceInstanceInput, pollInterval, timeoutSeconds time.Duration) error {
	return c.client.WaitFor("service instance to be deleted", pollInterval, timeoutSeconds, func() (bool, error) {
		info, err := c.GetServiceInstance(input)
		if err != nil {
			if client.WasNotFoundError(err) {
				// Service Instance could not be found, thus deleted
				return true, nil
			}
			// Some other error occurred trying to get instance, exit
			return false, err
		}
		switch s := info.Status; s {
		case ServiceInstanceTerminating:
			c.client.DebugLogString("Service Instance terminating")
			return false, nil
		default:
			c.client.DebugLogString(fmt.Sprintf("Unknown instance state: %s, waiting", s))
			return false, nil
		}
	})
}

func convertOracleBool(val bool) string {
	if val {
		return "yes"
	}
	// set false as blank rather than "no" so omitempty is honored
	return ""
}
