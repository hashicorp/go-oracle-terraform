package mysql

import (
	"fmt"
	"strings"
	"time"
)

var (
	MySQLAccessRuleContainerPath = "/paas/api/v1.1/instancemgmt/%s/services/MySQLCS/instances/%s/accessrules"
	MySQLAccessRuleRootPath      = "/paas/api/v1.1/instancemgmt/%s/services/MySQLCS/instances/%s/accessrules/%s"
)

// AccessRulesClient is a client for the Service functions of the MySQL API.
type AccessRulesClient struct {
	ResourceClient
	Timeout time.Duration
}

const WaitForAccessRuleTimeout = time.Duration(20 * time.Second)
const WaitForAccessRulePollInterval = time.Duration(1 * time.Second)

func (c *MySQLClient) AccessRulesClient() *AccessRulesClient {
	return &AccessRulesClient{
		ResourceClient: ResourceClient{
			MySQLClient:      c,
			ContainerPath:    MySQLAccessRuleContainerPath,
			ResourceRootPath: MySQLAccessRuleRootPath,
		},
	}
}

// Status Constants for an Access Rule
type AccessRuleStatus string

const (
	AccessRuleEnabled  AccessRuleStatus = "enabled"
	AccessRuleDisabled AccessRuleStatus = "disabled"
)

type AccessRuleInfo struct {
	Description string `json:"description"`
	Destination string `json:"destination,omitempty"`
	Ports       string `json:"ports"`
	Protocol    string `json:"protocol"`
	RuleName    string `json:"ruleName"`
	RuleType    string `json:"ruleType,omitempty"`
	Source      string `json:"source"`
	Status      string `json:"status,omitempty"`
}

type AccessRuleOperation string

const (
	AccessRuleUpdate AccessRuleOperation = "update"
	AccessRuleDelete AccessRuleOperation = "delete"
)

type AccessRuleList struct {
	AccessRules []AccessRuleInfo     `json:"accessRules"`
	Activities  []AccessRuleActivity `json:"activities"`
}

type AccessRuleActivity struct {
	AccessRuleActivityInfo AccessRuleActivityInfo `json:"activity"`
}

type AccessRuleActivityInfo struct {
	RuleName string `json:"ruleName"`
	Message  string `json:"message"`
	Errors   string `json:"errors"`
	Status   string `json:"status"`
}

type UpdateAccessRuleInput struct {
	ServiceInstanceID string              `json:"-"`
	Name              string              `json:"-"`
	Operation         AccessRuleOperation `json:"operation"`
	Status            AccessRuleStatus    `json:"status"`
	PollInterval      time.Duration       `json:"-"`
	Timeout           time.Duration       `json:"-"`
}

type CreateAccessRuleInput struct {
	ServiceInstanceID string        `json:"-"`
	Description       string        `json:"description"`
	Destination       string        `json:"destination,omitempty"`
	Ports             string        `json:"ports"`
	Protocol          string        `json:"protocol"`
	RuleName          string        `json:"ruleName"`
	Source            string        `json:"source"`
	Status            string        `json:"status,omitempty"`
	PollInterval      time.Duration `json:"-"`
	Timeout           time.Duration `json:"-"`
}

type GetAccessRuleInput struct {
	ServiceInstanceID string `json:"-"`
	Name              string `json:"-"`
}

type DeleteAccessRuleInput struct {
	Name              string              `json:"-"`
	ServiceInstanceID string              `json:"-"`
	Operation         AccessRuleOperation `json:"operation"`
	Status            AccessRuleStatus    `json:"status"`
	PollInterval      time.Duration       `json:"-"`
	Timeout           time.Duration       `json:"-"`
}

func (c *AccessRulesClient) CreateAccessRule(input *CreateAccessRuleInput) error {

	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	if err := c.createAccessRuleResource(input, nil); err != nil {
		return err
	}

	pollInterval := input.PollInterval
	if pollInterval == 0 {
		pollInterval = WaitForAccessRulePollInterval
	}

	timeout := input.Timeout
	if timeout == 0 {
		timeout = WaitForAccessRuleTimeout
	}

	getRuleInput := &GetAccessRuleInput{
		ServiceInstanceID: input.ServiceInstanceID,
		Name:              input.RuleName,
	}

	return c.WaitForAccessRuleReady(getRuleInput, pollInterval, timeout)
}

func (c *AccessRulesClient) GetAllAccessRules(input *GetAccessRuleInput) (*AccessRuleList, error) {

	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRules AccessRuleList
	if err := c.getAccessRulesResource(&accessRules); err != nil {
		return nil, err
	}

	// Iterated through entire slice, rule was not found.
	// No error occurred though, return a nil struct, and allow the Provdier to handle
	// a Nil response case.
	return &accessRules, nil
}

func (c *AccessRulesClient) GetAccessRule(input *GetAccessRuleInput) (*AccessRuleInfo, error) {

	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	var accessRules AccessRuleList
	if err := c.getAccessRulesResource(&accessRules); err != nil {
		return nil, err
	}

	// This is likely not the most optimal path for this, however, the upper bound on
	// performance here is the actual API request, not the iteration.
	for _, rule := range accessRules.AccessRules {
		if rule.RuleName == input.Name {
			return &rule, nil
		}
	}

	// Iterated through entire slice, rule was not found.
	// No error occurred though, return a nil struct, and allow the Provdier to handle
	// a Nil response case.
	return nil, nil
}

func (c *AccessRulesClient) WaitForAccessRuleReady(input *GetAccessRuleInput, pollInterval time.Duration, timeoutSeconds time.Duration) error {

	err := c.client.WaitFor("access rule to be created.", pollInterval, timeoutSeconds, func() (bool, error) {

		var info AccessRuleList
		if err := c.getAccessRulesResource(&info); err != nil {
			return false, err
		}

		c.client.DebugLogString(fmt.Sprintf("[DEBUG] Checking Activities : %v", info))
		for _, accessRule := range info.AccessRules {
			if accessRule.RuleName == input.Name {
				return true, nil
			}
		}

		for _, activity := range info.Activities {
			if activity.AccessRuleActivityInfo.RuleName == input.Name {

				switch s := strings.ToUpper(activity.AccessRuleActivityInfo.Status); s {
				case "FAILED":
					return false, fmt.Errorf("Error creating Access Rule : %s", activity.AccessRuleActivityInfo.Message)
				case "SUCCESS":
					return true, nil
				case "RUNNING":
					return false, nil
				default:
					return false, nil
				}
			}
		}
		return false, nil
	})
	return err
}

// Updates an AccessRule with the provided input struct. Returns a fully populated Info struct
// and any errors encountered

func (c *AccessRulesClient) UpdateAccessRule(input *UpdateAccessRuleInput,
) (*AccessRuleInfo, error) {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	// Since this is strictly an Update call, set the Operation constant
	input.Operation = AccessRuleUpdate
	// Initialize the response struct
	var accessRule AccessRuleInfo
	if err := c.updateAccessRulesResource(input.Name, input, &accessRule); err != nil {
		return nil, err
	}
	return &accessRule, nil
}

// Deletes an AccessRule with the provided input struct. Returns any errors that occurred.
func (c *AccessRulesClient) DeleteAccessRule(input *DeleteAccessRuleInput) error {
	if input.ServiceInstanceID != "" {
		c.ServiceInstanceID = input.ServiceInstanceID
	}

	// Since this is strictly an Update call, set the Operation constant
	input.Operation = AccessRuleDelete
	// The Update API call with a `DELETE` operation actually returns the same access rule info
	// in a response body. As we are deleting the AccessRule, we don't actually need to parse that.
	// However, the Update API call requires a pointer to parse, or else we throw an error during the
	// json unmarshal
	var result AccessRuleInfo
	if err := c.updateAccessRulesResource(input.Name, input, &result); err != nil {
		c.client.DebugLogString(fmt.Sprintf("[DEBUG] Failed to delete access rule : %v", err))
		return err
	}

	pollInterval := input.PollInterval
	if pollInterval == 0 {
		pollInterval = WaitForAccessRulePollInterval
	}

	timeout := input.Timeout
	if timeout == 0 {
		timeout = WaitForAccessRuleTimeout
	}

	getInput := &GetAccessRuleInput{
		Name: input.Name,
	}

	_, err := c.WaitForAccessRuleDeleted(getInput, pollInterval, timeout)
	if err != nil {
		c.client.DebugLogString(fmt.Sprintf("[DEBUG] Failed to delete access rule : %v", err))
		return err
	}

	return nil
}

func (c *AccessRulesClient) WaitForAccessRuleDeleted(input *GetAccessRuleInput, pollInterval time.Duration, timeout time.Duration) (*AccessRuleInfo, error) {
	var info *AccessRuleInfo
	var getErr error
	err := c.client.WaitFor("access rule to be deleted", pollInterval, timeout, func() (bool, error) {
		info, getErr = c.GetAccessRule(input)
		if getErr != nil {
			c.client.DebugLogString(fmt.Sprintf("[DEBUG] Failed waiting for access rule delete : %v",getErr))
			return true, nil
		}
		if info != nil {
			// Rule found, continue
			return false, nil
		}
		// Rule not found, return. Desired case
		return true, nil
	})
	return info, err
}
