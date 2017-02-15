package compute

// IPReservationsClient is a client for the IP Reservations functions of the Compute API.
type IPReservationsClient struct {
	*ResourceClient
}

// IPReservations obtains an IPReservationsClient which can be used to access to the
// IP Reservations functions of the Compute API
func (c *AuthenticatedClient) IPReservations() *IPReservationsClient {
	return &IPReservationsClient{
		ResourceClient: &ResourceClient{
			AuthenticatedClient: c,
			ResourceDescription: "ip reservation",
			ContainerPath:       "/ip/reservation/",
			ResourceRootPath:    "/ip/reservation",
		}}
}

// CreateIPReservationInfo defines an IP reservation to be created.
type CreateIPReservationInfo struct {
	Name       string   `json:"name"`
	ParentPool string   `json:"parentpool"`
	Permanent  bool     `json:"permanent"`
	Tags       []string `json:"tags"`
}

// IPReservationInfo describes an existing IP reservation.
type IPReservation struct {
	Account    string   `json:account`
	IP         string   `json:"ip"`
	Name       string   `json:"name"`
	ParentPool string   `json:"parentpool"`
	Permanent  bool     `json:"permanent"`
	Tags       []string `json:"tags"`
	Uri        string   `json:uri`
	Used       bool     `json:used`
}

// UpdateIPReservationInfo defines an IP Reservation to be updated
type UpdateIPReservationInfo struct {
	Name       string   `json:"name"`
	ParentPool string   `json:"parentpool"`
	Permanent  bool     `json:"permanent"`
	Tags       []string `json:"tags"`
}

// DeleteIPReservationINfo defines an IP Reservation to delete
type DeleteIPReservationInfo struct {
	Name string
}

// GetIPReservationInfo defines an IP Reservation to get
type GetIPReservationInfo struct {
	Name string
}

func (c *IPReservationsClient) success(result *IPReservation) (*IPReservation, error) {
	c.unqualify(&result.Name)
	return result, nil
}

// CreateIPReservation creates a new IP reservation with the given parentpool, tags and permanent flag.
func (c *IPReservationsClient) CreateIPReservation(createInfo CreateIPReservationInfo) (*IPReservation, error) {
	var ipInfo IPReservation
	if err := c.createResource(&createInfo, &ipInfo); err != nil {
		return nil, err
	}

	return c.success(&ipInfo)
}

// GetIPReservation retrieves the IP reservation with the given name.
func (c *IPReservationsClient) GetIPReservation(getInfo GetIPReservationInfo) (*IPReservation, error) {
	var ipInfo IPReservation
	if err := c.getResource(getInfo.Name, &ipInfo); err != nil {
		return nil, err
	}

	return c.success(&ipInfo)
}

// DeleteIPReservation deletes the IP reservation with the given name.
func (c *IPReservationsClient) DeleteIPReservation(deleteInfo DeleteIPReservationInfo) error {
	return c.deleteResource(deleteInfo.Name)
}

// UpdateIPReservation updates the IP reservation.
func (c *IPReservationsClient) UpdateIPReservation(updateInfo UpdateIPReservationInfo) (*IPReservation, error) {
	var ipInfo IPReservation
	if err := c.updateResource(updateInfo.Name, updateInfo, &ipInfo); err != nil {
		return nil, err
	}
	return c.success(&ipInfo)
}
