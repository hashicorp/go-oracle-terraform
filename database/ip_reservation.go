package database

// API URI Paths for Container and Root objects
const (
	ipReservationContainerPath = "/paas/api/v1.1/network/%s/services/dbaas/ipreservations"
	ipReservationResourcePath  = "/paas/api/v1.1/network/%s/services/dbaas/ipreservations/%s"
)

type IPReservationClient struct {
	IPReservationResourceClient
}

// IPReservationClient obtains an new ResourceClient which can be used to access the
// Database Cloud IP Reservation API
func (c *Client) IPReservationClient() *IPReservationClient {
	return &IPReservationClient{
		IPReservationResourceClient{
			Client:           c,
			ContainerPath:    ipReservationContainerPath,
			ResourceRootPath: ipReservationResourcePath,
		}}
}

type CreateIPReservationInput struct {
	// Identity domain ID for the Database Cloud Service account
	// For a Cloud account with Identity Cloud Service: the identity service ID, which has the form idcs-letters-and-numbers.
	// For a traditional cloud account: the name of the identity domain.
	// Required
	IdentityDomainID string
	// Name of the IP reservation to create.
	// Required
	Name string `json:"ipResName"`
	// Indicates whether the IP reservation is for instances attached to IP networks or the shared network
	// set to `IPNetwork` for IP Network, or omit for shared network
	NetworkType string `json:"networkType,omitempty"`
	// Name of the region to create the IP reservation in
	// Required
	Region string `json:"region"`
}

type IPReservationInfo struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	NetworkType          string `json:"networkType"`
	Status               string `json:"status"`
	IdentityDomain       string `json:"identityDomain"`
	ServiceType          string `json:"serviceType"`
	ComputeSiteName      string `json:"computeSiteName"`
	ServiceEntitlementID string `json:"serviceEntitlementId"`
}

// IPReservations - used for the GET request that returns all reservations
type IPReservations struct {
	IPReservations []IPReservationInfo `json:"ipReservations"`
}

// CreateServiceInstance creates a new ServiceInstace.
func (c *IPReservationClient) CreateIPReservation(input *CreateIPReservationInput) (*IPReservationInfo, error) {
	var ipReservation *IPReservationInfo
	if err := c.createResource(input, ipReservation); err != nil {
		return nil, err
	}
	return ipReservation, nil
}

func (c *IPReservationClient) GetIPReservation(name string) (*IPReservationInfo, error) {
	var ipReservations *IPReservations
	if err := c.getResource(name, ipReservations); err != nil {
		return nil, err
	}

	// API returns all IP Reservations, iterate to find the one we want
	for _, ipres := range ipReservations.IPReservations {
		if ipres.Name == name {
			return &ipres, nil
		}
	}
	return nil, nil
}

func (c *IPReservationClient) DeleteIPReservation(name string) error {
	if err := c.deleteResource(name); err != nil {
		return err
	}
	return nil
}
