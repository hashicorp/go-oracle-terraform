package lbaas

type SSLCertificateInfo struct {
	Name        string `json:"name"`
	Certificate string `json:"certificate"`
	State       string `json:"state"`
	Trusted     bool   `json:"trusted"`
	URI         string `json:"uri"`
}

type CreateSSLCertificateInput struct {
	Name        string `json:"name"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
	Trusted     bool   `json:"trusted"`
}

// CreateSSLCertificate creates a new SSL certificate
func (c *SSLCertificateClient) CreateSSLCertificate(input *CreateSSLCertificateInput) (*SSLCertificateInfo, error) {
	var info SSLCertificateInfo
	if err := c.createResource(&input, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// DeleteSSLCertificate deletes the SSL certificate with the specified name
func (c *SSLCertificateClient) DeleteSSLCertificate(name string) (*SSLCertificateInfo, error) {
	var info SSLCertificateInfo
	if err := c.deleteResource(name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetSSLCertificate fetch the SSL Certificate details
func (c *SSLCertificateClient) GetSSLCertificate(name string) (*SSLCertificateInfo, error) {
	var info SSLCertificateInfo
	if err := c.getResource(name, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
