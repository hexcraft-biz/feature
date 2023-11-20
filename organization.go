package feature

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

// ================================================================
//
// ================================================================
type OrganizationHttpMethods struct {
	*Feature
}

type OrganizationEndpoint struct {
	*Endpoint
}

func newOrganizationEndpoint(e *Endpoint) *OrganizationEndpoint {
	return &OrganizationEndpoint{
		Endpoint: e,
	}
}

func (f *Feature) ByAuthorityOfOrganization() *OrganizationHttpMethods {
	return &OrganizationHttpMethods{
		Feature: f,
	}
}

func (m *OrganizationHttpMethods) GET(relativePath, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "GET", relativePath)
	m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
	return newOrganizationEndpoint(e)
}

func (m *OrganizationHttpMethods) POST(relativePath, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "POST", relativePath)
	m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
	return newOrganizationEndpoint(e)
}

func (m *OrganizationHttpMethods) PUT(relativePath, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "PUT", relativePath)
	m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
	return newOrganizationEndpoint(e)
}

func (m *OrganizationHttpMethods) PATCH(relativePath, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "PATCH", relativePath)
	m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
	return newOrganizationEndpoint(e)
}

func (m *OrganizationHttpMethods) DELETE(relativePath, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "DELETE", relativePath)
	m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
	return newOrganizationEndpoint(e)
}

// ================================================================
//
// ================================================================
type EndpointAccessRules struct {
	Subsets    []string `json:"subsets"`
	Exceptions []string `json:"exceptions,omitempty"`
}

func (r *EndpointAccessRules) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, r)
}

func (j EndpointAccessRules) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (r *EndpointAccessRules) AddSubset(rule string) {
	r.Subsets = append(r.Subsets, rule)
}

func (r *EndpointAccessRules) AddException(rule string) {
	r.Exceptions = append(r.Exceptions, rule)
}

func (r *EndpointAccessRules) RemoveRedundant() {
	r.Subsets = RemoveRedundant(r.Subsets)
	r.Exceptions = RemoveRedundant(r.Exceptions)
}

func (r *EndpointAccessRules) Merge(rules *EndpointAccessRules) {
	r.Subsets = append(r.Subsets, rules.Subsets...)
	r.Exceptions = append(r.Exceptions, rules.Subsets...)
	r.RemoveRedundant()
}

// ================================================================
const (
	ActionAssign = iota
	ActionGrant
	ActionRevoke
)

const (
	WriteBehaviorUndef = iota
	WriteBehaviorIfNotExists
	WriteBehaviorOverwrite
)

type organizationUserAccess struct {
	DogmasApiUrl        *url.URL
	EndpointIdentifier  *string
	accessRulesToCommit map[int]map[string]*EndpointAccessRules
}

func (e *OrganizationEndpoint) ManageAccessFor(userId xuuid.UUID) *organizationUserAccess {
	return &organizationUserAccess{
		DogmasApiUrl:        e.Dogmas.HostUrl.JoinPath("/permissions/v1/users", userId.String()),
		EndpointIdentifier:  &e.EndpointIdentifier,
		accessRulesToCommit: map[int]map[string]*EndpointAccessRules{},
	}
}

func (u *organizationUserAccess) Assign(identifier string, rule string) *organizationUserAccess {
	return u.addAction(ActionAssign, identifier, rule)
}

func (u *organizationUserAccess) Grant(identifier string, rule string) *organizationUserAccess {
	return u.addAction(ActionGrant, identifier, rule)
}

func (u *organizationUserAccess) Revoke(identifier string, rule string) *organizationUserAccess {
	return u.addAction(ActionRevoke, identifier, rule)
}

func (u *organizationUserAccess) addAction(action int, identifier string, rule string) *organizationUserAccess {
	behavior := WriteBehaviorUndef
	switch action {
	case ActionGrant, ActionRevoke:
		behavior = WriteBehaviorOverwrite
	default:
		behavior = WriteBehaviorIfNotExists
	}

	if _, ok := u.accessRulesToCommit[behavior]; !ok {
		u.accessRulesToCommit[behavior] = map[string]*EndpointAccessRules{}
	}

	if _, ok := u.accessRulesToCommit[behavior][identifier]; !ok {
		u.accessRulesToCommit[behavior][identifier] = &EndpointAccessRules{}
	}

	switch action {
	case ActionAssign, ActionGrant:
		u.accessRulesToCommit[behavior][identifier].AddSubset(rule)
	case ActionRevoke:
		u.accessRulesToCommit[behavior][identifier].AddException(rule)
	}

	return u
}

type EndpointAccessRulesWithBehavior struct {
	Behavior                   string               `json:"behavior" db:"-" binding:"required"`
	AffectedEndpointIdentifier string               `json:"affectedEndpointIdentifier" db:"-" binding:"required"`
	AccessRules                *EndpointAccessRules `json:"accessRules" db:"access_rules" binding:"required"`
}

const (
	HeaderEndpointIdentiifer = "X-Endpoint-Identifier"
	HeaderByUserId           = "X-By-User-Id"
)

func (u *organizationUserAccess) Commit(byUserId xuuid.UUID) her.Error {
	rulesWithBehavior := []*EndpointAccessRulesWithBehavior{}
	for behavior, identifiers := range u.accessRulesToCommit {

		behaviorstring := ""
		switch behavior {
		case WriteBehaviorIfNotExists:
			behaviorstring = "IF_NOT_EXISTS"
		case WriteBehaviorOverwrite:
			behaviorstring = "OVERWRITE"
		default:
			return her.NewErrorWithMessage(http.StatusInternalServerError, "Undefined write behavior", nil)
		}

		for identifier, accessRules := range identifiers {
			accessRules.RemoveRedundant()
			rulesWithBehavior = append(rulesWithBehavior, &EndpointAccessRulesWithBehavior{
				Behavior:                   behaviorstring,
				AffectedEndpointIdentifier: identifier,
				AccessRules:                accessRules,
			})
		}
	}

	if len(rulesWithBehavior) > 0 {
		jsonbytes, err := json.Marshal(rulesWithBehavior)
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		req, err := http.NewRequest("POST", u.DogmasApiUrl.String(), bytes.NewReader(jsonbytes))
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		req.Header.Set(HeaderEndpointIdentiifer, *u.EndpointIdentifier)
		req.Header.Set(HeaderByUserId, byUserId.String())

		payload := her.NewPayload(nil)
		client := &http.Client{}

		if resp, err := client.Do(req); err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
			return err
		} else if resp.StatusCode != 201 {
			return her.NewErrorWithMessage(http.StatusInternalServerError, "Dogmas: "+payload.Message, nil)
		}
	}

	return nil
}
