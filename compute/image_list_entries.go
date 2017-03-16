package compute

const (
	ImageListEntryDescription   = "image list entry"
	ImageListEntryContainerPath = "/imagelist"
	ImageListEntryResourcePath  = "/imagelist"
)

type ImageListEntriesClient struct {
	ResourceClient
}

// ImageListEntries() returns an ImageListEntriesClient that can be used to access the
// necessary CRUD functions for Image List Entrys.
func (c *Client) ImageListEntries(name string, version string) *ImageListEntriesClient {
	var containerPath, resourcePath string
	name = c.getQualifiedName(name)
	containerPath = ImageListEntryContainerPath + name + "/entry/"
	resourcePath = ImageListEntryContainerPath + name + "/entry"
	if version != "" {
		containerPath = containerPath + version
		resourcePath = resourcePath + "/" + version
	}
	return &ImageListEntriesClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: ImageListEntryDescription,
			ContainerPath:       containerPath,
			ResourceRootPath:    resourcePath,
		},
	}
}

// ImageListEntryInfo contains the exported fields necessary to hold all the information about an
// Image List Entry
type ImageListEntryInfo struct {
	// User-defined parameters, in JSON format, that can be passed to an instance of this machine
	// image when it is launched. This field can be used, for example, to specify the location of
	// a database server and login details. Instance metadata, including user-defined data is available
	// at http://192.0.0.192/ within an instance. See Retrieving User-Defined Instance Attributes in Using
	// Oracle Compute Cloud Service (IaaS).
	Attributes map[string]interface{} `json:"attributes"`
	// Name of the imagelist.
	ImageList string `json:"imagelist"`
	// A list of machine images.
	MachineImages []string `json:"machineimages"`
	// Uniform Resource Identifier for the Image List Entry
	Uri string `json:"uri"`
	// Version number of these machineImages in the imagelist.
	Version int `json:"version"`
}

type CreateImageListEntryInput struct {
	// User-defined parameters, in JSON format, that can be passed to an instance of this machine
	// image when it is launched. This field can be used, for example, to specify the location of
	// a database server and login details. Instance metadata, including user-defined data is
	//available at http://192.0.0.192/ within an instance. See Retrieving User-Defined Instance
	//Attributes in Using Oracle Compute Cloud Service (IaaS).
	// Optional
	Attributes map[string]interface{} `json:"attributes"`
	// A list of machine images.
	// Required
	MachineImages []string `json:"machineimages"`
	// The unique version of the entry in the image list.
	// Required
	Version int `json:"version"`
}

// Create a new Image List Entry from an ImageListEntriesClient and an input struct.
// Returns a populated Info struct for the Image List Entry, and any errors
func (c *ImageListEntriesClient) CreateImageListEntry(input *CreateImageListEntryInput) (*ImageListEntryInfo, error) {
	var ipInfo ImageListEntryInfo
	if err := c.createResource(&input, &ipInfo); err != nil {
		return nil, err
	}
	return c.success(&ipInfo)
}

// Returns a populated ImageListEntryInfo struct from an input struct
func (c *ImageListEntriesClient) GetImageListEntry() (*ImageListEntryInfo, error) {
	var ipInfo ImageListEntryInfo
	if err := c.getResource("", &ipInfo); err != nil {
		return nil, err
	}
	return c.success(&ipInfo)
}

func (c *ImageListEntriesClient) DeleteImageListEntry() error {
	return c.deleteResource("")
}

// Unqualifies any qualified fields in the ImageListEntryInfo struct
func (c *ImageListEntriesClient) success(info *ImageListEntryInfo) (*ImageListEntryInfo, error) {
	return info, nil
}