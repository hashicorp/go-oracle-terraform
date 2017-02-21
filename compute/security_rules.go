package compute

import (
	"fmt"
	"strings"
)

// SecurityRulesClient is a client for the Security Rules functions of the Compute API.
type SecurityRulesClient struct {
	ResourceClient
}

// SecurityRules obtains a SecurityRulesClient which can be used to access to the
// Security Rules functions of the Compute API
func (c *Client) SecurityRules() *SecurityRulesClient {
	return &SecurityRulesClient{
		ResourceClient: ResourceClient{
			Client:              c,
			ResourceDescription: "security ip list",
			ContainerPath:       "/secrule/",
			ResourceRootPath:    "/secrule",
		}}
}

// SecurityRuleInfo describes an existing security rule.
type SecurityRuleInfo struct {
	Action          string `json:"action"`
	Application     string `json:"application"`
	Description     string `json:"description"`
	Disabled        bool   `json:"disabled"`
	DestinationList string `json:"dst_list"`
	Name            string `json:"name"`
	SourceList      string `json:"src_list"`
	URI             string `json:"uri"`
}

// CreateSecurityRuleInput defines a security rule to be created.
type CreateSecurityRuleInput struct {
	Action          string `json:"action"`
	Application     string `json:"application"`
	Description     string `json:"description"`
	Disabled        bool   `json:"disabled"`
	DestinationList string `json:"dst_list"`
	Name            string `json:"name"`
	SourceList      string `json:"src_list"`
}

// CreateSecurityRule creates a new security rule.
func (c *SecurityRulesClient) CreateSecurityRule(createInput *CreateSecurityRuleInput) (*SecurityRuleInfo, error) {
	createInput.Name = c.getQualifiedName(createInput.Name)
	createInput.SourceList = c.getQualifiedListName(createInput.SourceList)
	createInput.DestinationList = c.getQualifiedListName(createInput.DestinationList)
	createInput.Application = c.getQualifiedName(createInput.Application)

	var ruleInfo SecurityRuleInfo
	if err := c.createResource(createInput, &ruleInfo); err != nil {
		return nil, err
	}

	return c.success(&ruleInfo)
}

// GetSecurityRuleInput describes the Security Rule to get
type GetSecurityRuleInput struct {
	Name string `json:"name"`
}

// GetSecurityRule retrieves the security rule with the given name.
func (c *SecurityRulesClient) GetSecurityRule(getInput *GetSecurityRuleInput) (*SecurityRuleInfo, error) {
	var ruleInfo SecurityRuleInfo
	if err := c.getResource(getInput.Name, &ruleInfo); err != nil {
		return nil, err
	}

	return c.success(&ruleInfo)
}

// UpdateSecurityRuleInput describes a secruity rule to update
type UpdateSecurityRuleInput struct {
	Action          string `json:"action"`
	Application     string `json:"application"`
	Description     string `json:"description"`
	Disabled        bool   `json:"disabled"`
	DestinationList string `json:"dst_list"`
	Name            string `json:"name"`
	SourceList      string `json:"src_list"`
}

// UpdateSecurityRule modifies the properties of the security rule with the given name.
func (c *SecurityRulesClient) UpdateSecurityRule(udpateInput *UpdateSecurityRuleInput) (*SecurityRuleInfo, error) {
  updateInput.Name = c.getQualifiedName(updateInput.Name)

	var ruleInfo SecurityRuleInfo
	if err := c.updateResource(updateInput.Name), updateInput, &ruleInfo); err != nil {
		return nil, err
	}

	return c.success(&ruleInfo)
}

// DeleteSecurityRuleInput describes the security rule to delete
type DeleteSecurityRuleInput struct {
	Name string `json:"name"`
}

// DeleteSecurityRule deletes the security rule with the given name.
func (c *SecurityRulesClient) DeleteSecurityRule(deleteInput *DeleteSecurityRuleInput) error {
	return c.deleteResource(deleteInput.Name)
}

func (c *SecurityRulesClient) getQualifiedListName(name string) string {
	nameParts := strings.Split(name, ":")
	listType := nameParts[0]
	listName := nameParts[1]
	return fmt.Sprintf("%s:%s", listType, c.getQualifiedName(listName))
}

func (c *SecurityRulesClient) unqualifyListName(qualifiedName string) string {
	nameParts := strings.Split(qualifiedName, ":")
	listType := nameParts[0]
	listName := nameParts[1]
	return fmt.Sprintf("%s:%s", listType, c.getUnqualifiedName(listName))
}

func (c *SecurityRulesClient) success(ruleInfo *SecurityRuleInfo) (*SecurityRuleInfo, error) {
	ruleInfo.Name = c.getUnqualifiedName(ruleInfo.Name)
	ruleInfo.SourceList = c.unqualifyListName(ruleInfo.SourceList)
	ruleInfo.DestinationList = c.unqualifyListName(ruleInfo.DestinationList)
	ruleInfo.Application = c.getUnqualifiedName(ruleInfo.Application)
	return ruleInfo, nil
}
