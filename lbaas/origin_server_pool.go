package lbaas

var (
	originserverpoolContainerPath = "/vlbrs/%s/%s/originserverpools"
	originserverpoolResourcePath  = "/vlbrs/%s/%s/originserverpools/%s"
)

// OriginServerPoolClient is a client for the Load Balancer Origin Server Pool resources.
type OriginServerPoolClient struct {
	LBaaSResourceClient
}

// OriginServerPoolClient returns an Client which is used to access the
// Load Balancer Origin Server Pool API
func (c *Client) OriginServerPoolClient() *OriginServerPoolClient {
	c.ContentType = CONTENT_TYPE_ORIGIN_SERVER_POOL_JSON
	c.Accept = CONTENT_TYPE_ORIGIN_SERVER_POOL_JSON
	return &OriginServerPoolClient{
		LBaaSResourceClient: LBaaSResourceClient{
			Client:           c,
			ContainerPath:    originserverpoolContainerPath,
			ResourceRootPath: originserverpoolResourcePath,
		},
	}
}

type OriginServerInfo struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Status   string `json:"status"`
}

type OriginServerPoolInfo struct {
	Consumers          string             `json:"consumers"`
	HealthCheck        HealthCheckInfo    `json:"health_check"`
	Name               string             `json:"name"`
	OperationDetails   string             `json:"operation_details"`
	OriginServers      []OriginServerInfo `json:"origin_servers"`
	ReasonForDisabling string             `json:"reason_for_disabling"`
	State              string             `json:"state"`
	Status             string             `json:"status"`
	Tags               []string           `json:"tags"`
	URL                string             `json:"uri"`
	VnicSetName        string             `json:"vnic_set_name"`
}

type CreateOriginServerPoolInput struct {
	Name          string             `json:"name"`
	OriginServers []OriginServerInfo `json:"origin_servers,omitempty"`
	Status        string             `json:"status,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
}

type UpdateOriginServerPoolInput struct {
	OriginServers []OriginServerInfo `json:"origin_servers,omitempty"`
	Status        string             `json:"status,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
}

// CreateOriginServerPool creates a new listener
func (c *OriginServerPoolClient) CreateOriginServerPool(lbRegion, lbName string, input *CreateOriginServerPoolInput) (*OriginServerPoolInfo, error) {
	var info OriginServerPoolInfo
	if err := c.createResource(lbRegion, lbName, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// DeleteOriginServerPool deletes the listener with the specified input
func (c *OriginServerPoolClient) DeleteOriginServerPool(lbRegion, lbName, name string) (*OriginServerPoolInfo, error) {
	var info OriginServerPoolInfo
	if err := c.deleteResource(lbRegion, lbName, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetOriginServerPool fetchs the listener details
func (c *OriginServerPoolClient) GetOriginServerPool(lbRegion, lbName, name string) (*OriginServerPoolInfo, error) {
	var info OriginServerPoolInfo
	if err := c.getResource(lbRegion, lbName, name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// UpdateOriginServerPool fetchs the listener details
func (c *OriginServerPoolClient) UpdateOriginServerPool(lbRegion, lbName, name string, input *UpdateOriginServerPoolInput) (*OriginServerPoolInfo, error) {
	var info OriginServerPoolInfo
	if err := c.updateResource(lbRegion, lbName, name, &input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
