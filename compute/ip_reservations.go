package compute

// IPReservationsClient is a client for the IP Reservations functions of the Compute API.
type IPReservationsClient struct {
	*ResourceClient
}

// IPReservations obtains an IPReservationsClient which can be used to access to the
// IP Reservations functions of the Compute API
func (c *Client) IPReservations() *IPReservationsClient {
	return &IPReservationsClient{
		ResourceClient: &ResourceClient{
			Client:              c,
			ResourceDescription: "ip reservation",
			ContainerPath:       "/ip/reservation/",
			ResourceRootPath:    "/ip/reservation",
		}}
}

// IPReservationInput describes an existing IP reservation.
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

// CreateIPReservationInput defines an IP reservation to be created.
type CreateIPReservationInput struct {
	Name       string   `json:"name"`
	ParentPool string   `json:"parentpool"`
	Permanent  bool     `json:"permanent"`
	Tags       []string `json:"tags"`
}

// CreateIPReservation creates a new IP reservation with the given parentpool, tags and permanent flag.
func (c *IPReservationsClient) CreateIPReservation(createInput CreateIPReservationInput) (*IPReservation, error) {
	var ipInput IPReservation
	if err := c.createResource(&createInput, &ipInput); err != nil {
		return nil, err
	}

	return c.success(&ipInput)
}

// GetIPReservationInput defines an IP Reservation to get
type GetIPReservationInput struct {
	Name string
}

// GetIPReservation retrieves the IP reservation with the given name.
func (c *IPReservationsClient) GetIPReservation(getInput GetIPReservationInput) (*IPReservation, error) {
	var ipInput IPReservation
	if err := c.getResource(getInput.Name, &ipInput); err != nil {
		return nil, err
	}

	return c.success(&ipInput)
}

// DeleteIPReservationInput defines an IP Reservation to delete
type DeleteIPReservationInput struct {
	Name string
}

// DeleteIPReservation deletes the IP reservation with the given name.
func (c *IPReservationsClient) DeleteIPReservation(deleteInput DeleteIPReservationInput) error {
	return c.deleteResource(deleteInput.Name)
}

// UpdateIPReservationInput defines an IP Reservation to be updated
type UpdateIPReservationInput struct {
	Name       string   `json:"name"`
	ParentPool string   `json:"parentpool"`
	Permanent  bool     `json:"permanent"`
	Tags       []string `json:"tags"`
}

// UpdateIPReservation updates the IP reservation.
func (c *IPReservationsClient) UpdateIPReservation(updateInput UpdateIPReservationInput) (*IPReservation, error) {
	var ipInput IPReservation
	if err := c.updateResource(updateInput.Name, updateInput, &ipInput); err != nil {
		return nil, err
	}
	return c.success(&ipInput)
}

func (c *IPReservationsClient) success(result *IPReservation) (*IPReservation, error) {
	c.unqualify(&result.Name)
	return result, nil
}
