package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HanseMerkur/terraform-provider-utils/log"
)

const (
	UserEndpointPrefix = "users"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// The ForemanUser API model represents a user.
type ForemanUser struct {
	// Inherits the base object's attributes
	ForemanObject

	// login name (i.e: username)
	Login string `json:"login"`
	// if user has admin privileges
	Admin bool `json:"admin,omitempty"`

	// "real" firstname of user
	Firstname string `json:"firstname,omitempty"`

	// "real" lastname of user
	Lastname string `json:"lastname,omitempty"`

	// email of user
	Mail string `json:"mail,omitempty"`

	// user description
	Description string `json:"description,omitempty"`

	// user password
	Password string `json:"password,omitempty"`

	// default location for user
	DefaultLocationId int `json:"default_location_id,omitempty"`

	// default organisation for user
	DefaultOrganizationId int `json:"default_organization_id,omitempty"`

	// origin of user authentication
	AuthSourceId int `json:"auth_source_id"`

	// locale setting for user
	Locale string `json:"locale,omitempty"`

	// list of all locations for user
	LocationIds []int `json:"location_ids,omitempty"`

	// list of all organisation for user
	OrganizationIds []int `json:"organization_ids,omitempty"`
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// CreateUser creates a new ForemanUser with the attributes of the
// supplied ForemanUser reference and returns the created ForemanUser
// reference.  The returned reference will have its ID and other API default
// values set by this function.
func (c *Client) CreateUser(u *ForemanUser) (*ForemanUser, error) {
	log.Tracef("foreman/api/user.go#Create")

	reqEndpoint := fmt.Sprintf("/%s", UserEndpointPrefix)

	uJSONBytes, jsonEncErr := c.WrapJSONWithTaxonomy("user", u)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("userJSONBytes: [%s]", uJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPost,
		reqEndpoint,
		bytes.NewBuffer(uJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var createdUser ForemanUser
	sendErr := c.SendAndParse(req, &createdUser)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("createdUser: [%+v]", createdUser)

	return &createdUser, nil
}

// ReadUser reads the attributes of a ForemanUser identified by the
// supplied ID and returns a ForemanUser reference.
func (c *Client) ReadUser(id int) (*ForemanUser, error) {
	log.Tracef("foreman/api/user.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", UserEndpointPrefix, id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readUser ForemanUser
	sendErr := c.SendAndParse(req, &readUser)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readUser: [%+v]", readUser)

	return &readUser, nil
}

// UpdateUser updates a ForemanUser's attributes.  The user with
// the ID of the supplied ForemanUser will be updated. A new
// ForemanUser reference is returned with the attributes from the result
// of the update operation.
func (c *Client) UpdateUser(u *ForemanUser) (*ForemanUser, error) {
	log.Tracef("foreman/api/user.go#Update")

	reqEndpoint := fmt.Sprintf("/%s/%d", UserEndpointPrefix, u.Id)

	uJSONBytes, jsonEncErr := c.WrapJSON("user", u)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("userJSONBytes: [%s]", uJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPut,
		reqEndpoint,
		bytes.NewBuffer(uJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var updatedUser ForemanUser
	sendErr := c.SendAndParse(req, &updatedUser)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("updatedUser: [%+v]", updatedUser)

	return &updatedUser, nil
}

// DeleteUser deletes the ForemanUser identified by the supplied ID
func (c *Client) DeleteUser(id int) error {
	log.Tracef("foreman/api/user.go#Delete")

	reqEndpoint := fmt.Sprintf("/%s/%d", UserEndpointPrefix, id)

	req, reqErr := c.NewRequest(
		http.MethodDelete,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

// QueryUser queries for a ForemanSubnet based on the attributes of the
// supplied ForemanSubnet reference and returns a QueryResponse struct
// containing query/response metadata and the matching subnets
func (c *Client) QueryUser(s *ForemanUser) (QueryResponse, error) {
	log.Tracef("foreman/api/user.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", UserEndpointPrefix)
	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return queryResponse, reqErr
	}

	// dynamically build the query based on the attributes
	// not all api search fields supported
	reqQuery := req.URL.Query()
	if s.Description != "" {
		description := `"` + s.Description + `"`
		reqQuery.Set("search", "description="+description)
	} else if s.Firstname != "" {
		firstname := `"` + s.Firstname + `"`
		reqQuery.Set("search", "firstname="+firstname)
	} else if s.Lastname != "" {
		lastname := `"` + s.Lastname + `"`
		reqQuery.Set("search", "firstname="+lastname)
	} else if s.Mail != "" {
		mail := `"` + s.Mail + `"`
		reqQuery.Set("search", "firstname="+mail)
	} else if s.Login != "" {
		login := `"` + s.Login + `"`
		reqQuery.Set("search", "login="+login)
	}

	req.URL.RawQuery = reqQuery.Encode()
	sendErr := c.SendAndParse(req, &queryResponse)
	if sendErr != nil {
		return queryResponse, sendErr
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	// Results will be Unmarshaled into a []map[string]interface{}
	//
	// Encode back to JSON, then Unmarshal into []ForemanSubnet for
	// the results
	results := []ForemanSubnet{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// convert the search results from []ForemanSubnet to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}
