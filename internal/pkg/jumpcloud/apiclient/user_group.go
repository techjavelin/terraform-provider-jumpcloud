package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	apiVersion  = "v2"
	apiEndpoint = "usergroups"
)

type (
	UserGroup struct {
		Attributes              *UserGroupAttributes             `json:"attributes,omitempty"`
		Description             string                           `json:"description,omitempty"`
		Email                   string                           `json:"email,omitempty"`
		Id                      string                           `json:"id"`
		MemberQuery             *UserGroupMemberQuery            `json:"memberQuery,omitempty"`
		MemberQueryExceptions   []UserGroupMemberQueryExceptions `json:"memberQueryExceptions,omitempty"`
		MemberSuggestionsNotify bool                             `json:"memberSuggestionsNotify,omitempty"`
		MembershipAutomated     bool                             `json:"membershipAutomated,omitempty"`
		Name                    string                           `json:"name"`
		SuggestionCounts        *UserGroupMemberSuggestionCounts `json:"suggestionCounts,omitempty"`
		Type                    string                           `json:"type,omitempty"`
	}

	UserGroupAttributes struct {
		Sudo         *UserGroupSudoConfig   `json:"sudo,omitempty"`
		LdapGroups   []LdapGroup            `json:"ldapGroups,omitempty"`
		PosixGroups  []PosixGroup           `json:"posixGroups,omitempty"`
		Radius       *UserGroupRadiusConfig `json:"radius,omitempty"`
		SambaEnabled bool                   `json:"sambaEnabled,omitempty"`
	}

	UserGroupSudoConfig struct {
		Enabled         bool `json:"enabled,omitempty"`
		WithoutPassword bool `json:"withoutPassword,omitempty"`
	}

	LdapGroup struct {
		Name string `json:"name,omitempty"`
	}

	PosixGroup struct {
		Id   int64  `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	UserGroupRadiusConfig struct {
		Reply []RadiusReply `json:"reply,omitempty"`
	}

	RadiusReply struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}

	UserGroupMemberQuery struct {
		QueryType string        `json:"queryType,omitempty"`
		Filters   []QueryFilter `json:"filters,omitempty"`
	}

	QueryFilter struct {
		Field    string `json:"field,omitempty"`
		Operator string `json:"operator,omitempty"`
		Value    string `json:"value,omitempty"`
	}

	UserGroupMemberQueryExceptions struct {
		Attributes []interface{} `json:"attributes,omitempty"`
		Id         string        `json:"id,omitempty"`
		Type       string        `json:"type,omitempty"`
	}

	UserGroupMemberSuggestionCounts struct {
		Add    int64 `json:"add,omitempty"`
		Remove int64 `json:"remove,omitempty"`
		Total  int64 `json:"total,omitempty"`
	}
)

func (c *Client) CallApiWithBody(method string, group *UserGroup) (payload UserGroup, response *http.Response, err error) {
	tflog.SubsystemInfo(c.Context, SUBSYSTEM_NAME, "Upsert UserGroup named "+group.Name, map[string]interface{}{
		"client":  "UserGroup",
		"func":    "CallApiWithBody",
		"method":  method,
		"payload": c.getPayloadAsString(payload),
	})

	endpoint := apiEndpoint
	if group.Id != "" {
		endpoint = fmt.Sprintf("%s/%s", apiEndpoint, group.Id)
	}

	request, err := c.prepareRequest(method, apiVersion, endpoint, group, nil, nil)
	if err != nil {
		return payload, nil, err
	}

	request_reader := c.ReusableReader(request.Body)
	request.Body = io.NopCloser(request_reader)
	tflog.SubsystemInfo(c.Context, SUBSYSTEM_NAME, "Sending Request", map[string]interface{}{
		"client": "UserGroup",
		"func":   "CallApiWithBody",
		"method": request.Method,
		"url":    request.URL.String(),
		"body":   c.ReadBody(request_reader),
	})

	response, err = c.httpClient.Do(request)

	response_reader := c.ReusableReader(response.Body)
	response.Body = io.NopCloser(response_reader)

	tflog.SubsystemInfo(c.Context, SUBSYSTEM_NAME, "Got Response from API", map[string]interface{}{
		"client":   "UserGroup",
		"func":     "CallApiWithBody",
		"method":   request.Method,
		"url":      request.URL.String(),
		"request":  c.ReadBody(request_reader),
		"status":   response.Status,
		"response": c.ReadBody(response_reader),
		"err":      err,
	})

	if err != nil || response == nil {
		return payload, response, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		tflog.SubsystemInfo(c.Context, SUBSYSTEM_NAME, "Got Response from API", map[string]interface{}{
			"client":        "UserGroup",
			"func":          "CallApiWithBody",
			"method":        request.Method,
			"url":           request.URL.String(),
			"request_body":  c.getPayloadAsString(payload),
			"status":        response.Status,
			"response_body": body,
			"err":           err,
		})
		return payload, response, fmt.Errorf("status: %v, body: %s", response.StatusCode, body)
	}

	// if err = json.NewDecoder(response.Body).Decode(&payload); err != nil {
	// 	return payload, response, err
	// }

	response_buf, _ := io.ReadAll(response_reader)
	err = json.Unmarshal(response_buf, &payload)
	if err != nil {
		tflog.SubsystemError(c.Context, SUBSYSTEM_NAME, "Error while Unmarshalling Response", map[string]interface{}{
			"client":   "UserGroup",
			"func":     "CallApiWithBody",
			"response": c.ReadBody(response_reader),
			"err":      err,
		})
	}

	jsonPayload, _ := json.Marshal(payload)
	tflog.SubsystemInfo(c.Context, SUBSYSTEM_NAME, "Unmarshalled Response into UserGroup", map[string]interface{}{
		"client":    "UserGroup",
		"func":      "CallApiWithBody",
		"response":  c.ReadBody(response_reader),
		"usergroup": string(jsonPayload),
	})

	return payload, response, err
}

func (c *Client) CallApiWithID(method string, id string) (payload UserGroup, response *http.Response, err error) {
	endpoint := fmt.Sprintf("%s/%s", apiEndpoint, id)
	request, err := c.prepareRequest(method, apiVersion, endpoint, nil, nil, nil)
	if err != nil {
		return payload, nil, err
	}

	response, err = c.httpClient.Do(request)

	if err != nil || response == nil {
		return payload, response, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(response.Body)
		return payload, response, errors.New(fmt.Sprintf("Status: %v, Body: %s", response.StatusCode, body))
	}

	if err = json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return payload, response, err
	}

	return payload, response, err
}

func (c *Client) CreateUserGroup(create *UserGroup) (UserGroup, *http.Response, error) {
	return c.CallApiWithBody(http.MethodPost, create)
}

func (c *Client) GetUserGroupDetails(id string) (UserGroup, *http.Response, error) {
	return c.CallApiWithID(http.MethodGet, id)
}

func (c *Client) DeleteUserGroup(id string) (*http.Response, error) {
	_, response, err := c.CallApiWithID(http.MethodDelete, id)
	return response, err
}

func (c *Client) UpdateUserGroup(update *UserGroup) (UserGroup, *http.Response, error) {
	return c.CallApiWithBody(http.MethodPut, update)
}
