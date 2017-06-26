package storage

import (
	"net/http"
	"strconv"
	"strings"
)

const CONTAINER_VERSION = "v1"

// Container describes an existing Container.
type Container struct {
	// The name of the Container
	Name string
	// A container access control list (ACL) that grants read access.
	ReadACLs []string
	// A container access control list (ACL) that grants write access
	WriteACLs []string
	// The secret key value for temporary URLs.
	URLKey string
	// The second secret key value for temporary URLs.
	URLKey2 string
	// List of origins to be allowed to make cross-origin Requests.
	AllowedOrigins []string
	// Maximum age in seconds for the origin to hold the preflight results.
	MaxAge int
}

// CreateContainerInput defines an Container to be created.
type CreateContainerInput struct {
	// The unique name for the container. The container name must be from 1 to 256 characters long and can
	// start with any character and contain any pattern. Character set must be UTF-8. The container name
	// cannot contain a slash (/) character because this character delimits the container and object name.
	// For example, /account/container/object.
	// Required
	Name string `json:"name"`
	// Sets a container access control list (ACL) that grants read access.
	ReadACLs []string
	// Sets a container access control list (ACL) that grants read access.
	WriteACLs []string
	// Sets a secret key value for temporary URLs.
	URLKey string
	// Sets a second secret key value for temporary URLs.
	URLKey2 string
	// Sets the list of origins allowed to make cross-origin requests.
	AllowedOrigins []string
	// Sets the maximum age in seconds for the origin to hold the preflight results.
	MaxAge int
}

// CreateContainer creates a new Container with the given name, key and enabled flag.
func (c *StorageClient) CreateContainer(createInput *CreateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	createInput.Name = c.getQualifiedName(CONTAINER_VERSION, createInput.Name)

	// There are default values for these that we don't want to zero out if Read and Write ACLs are not set.
	if len(createInput.ReadACLs) > 0 {
		headers["X-Container-Read"] = strings.Join(createInput.ReadACLs, ",")
	}
	if len(createInput.WriteACLs) > 0 {
		headers["X-Container-Write"] = strings.Join(createInput.WriteACLs, ",")
	}

	headers["X-Container-Meta-Temp-URL-Key"] = createInput.URLKey
	headers["X-Container-Meta-Temp-URL-Key-2"] = createInput.URLKey2
	headers["X-Container-Meta-Access-Control-Expose-Headers"] = strings.Join(createInput.AllowedOrigins, " ")
	headers["X-Container-Meta-Access-Control-Max-Age"] = strconv.Itoa(createInput.MaxAge)

	if err := c.createResource(createInput.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: createInput.Name,
	}

	return c.GetContainer(&getInput)
}

// DeleteKeyInput describes the container to delete
type DeleteContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:name`
}

// DeleteContainer deletes the Container with the given name.
func (c *StorageClient) DeleteContainer(deleteInput *DeleteContainerInput) error {
	deleteInput.Name = c.getQualifiedName(CONTAINER_VERSION, deleteInput.Name)
	return c.deleteResource(deleteInput.Name)
}

// GetContainerInput describes the container to get
type GetContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:name`
}

// GetContainer retrieves the Container with the given name.
func (c *StorageClient) GetContainer(getInput *GetContainerInput) (*Container, error) {
	var (
		container Container
		rsp       *http.Response
		err       error
	)
	getInput.Name = c.getQualifiedName(CONTAINER_VERSION, getInput.Name)

	if rsp, err = c.getResource(getInput.Name, &container); err != nil {
		return nil, err
	}
	return c.success(rsp, &container)
}

// UpdateContainerInput defines an Container to be updated
type UpdateContainerInput struct {
	// The name of the Container
	// Required
	Name string `json:"name"`
	// Updates a container access control list (ACL) that grants read access.
	ReadACLs []string
	// Updates a container access control list (ACL) that grants write access.
	WriteACLs []string
	// Updates the secret key value for temporary URLs.
	URLKey string
	// Update the second secret key value for temporary URLs.
	URLKey2 string
	// Updates the list of origins allowed to make cross-origin requests.
	AllowedOrigins []string
	// Updates the maximum age in seconds for the origin to hold the preflight results.
	MaxAge int
}

// UpdateContainer updates the key and enabled flag of the Container with the given name.
func (c *StorageClient) UpdateContainer(updateInput *UpdateContainerInput) (*Container, error) {
	headers := make(map[string]string)

	// There are default values for these that we don't want to zero out if Read and Write ACLs are not set.
	if len(updateInput.ReadACLs) > 0 {
		headers["X-Container-Read"] = strings.Join(updateInput.ReadACLs, ",")
	}
	if len(updateInput.WriteACLs) > 0 {
		headers["X-Container-Write"] = strings.Join(updateInput.WriteACLs, ",")
	}

	headers["X-Container-Meta-Temp-URL-Key"] = updateInput.URLKey
	headers["X-Container-Meta-Temp-URL-Key-2"] = updateInput.URLKey2
	headers["X-Container-Meta-Access-Control-Expose-Headers"] = strings.Join(updateInput.AllowedOrigins, " ")
	headers["X-Container-Meta-Access-Control-Max-Age"] = strconv.Itoa(updateInput.MaxAge)

	updateInput.Name = c.getQualifiedName(CONTAINER_VERSION, updateInput.Name)
	if err := c.updateResource(updateInput.Name, headers); err != nil {
		return nil, err
	}

	getInput := GetContainerInput{
		Name: updateInput.Name,
	}
	return c.GetContainer(&getInput)
}

func (c *StorageClient) success(rsp *http.Response, container *Container) (*Container, error) {
	var err error
	c.unqualify(&container.Name)
	container.ReadACLs = strings.Split(rsp.Header.Get("X-Container-Read"), ",")
	container.WriteACLs = strings.Split(rsp.Header.Get("X-Container-Write"), ",")
	container.URLKey = rsp.Header.Get("X-Container-Meta-Temp-URL-Key")
	container.URLKey2 = rsp.Header.Get("X-Container-Meta-Temp-URL-Key-2")
	container.AllowedOrigins = strings.Split(rsp.Header.Get("X-Container-Meta-Access-Control-Expose-Headers"), " ")
	container.MaxAge, err = strconv.Atoi(rsp.Header.Get("X-Container-Meta-Access-Control-Max-Age"))
	if err != nil {
		return nil, err
	}

	return container, nil
}
