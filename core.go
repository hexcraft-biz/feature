package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

// ================================================================
//
// ================================================================
const (
	ByAuthorityOfOrganization = "ORGANIZATION"
	ByAuthorityOfDataOwner    = "DATA_OWNER"
)

type Feature struct {
	*gin.RouterGroup
	*Dogmas
}

func New(e *gin.Engine, startsWith string, d *Dogmas) *Feature {
	return &Feature{
		RouterGroup: e.Group(startsWith),
		Dogmas:      d,
	}
}

type HandlerFunc func(*Endpoint) gin.HandlerFunc

func handlerFuncs(e *Endpoint, handlers []HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = h(e)
	}
	return funcs
}

// ================================================================
//
// ================================================================
type OrganizationHttpMethod struct {
	*Feature
}

type DataOwnerHttpMethod struct {
	*Feature
}

type OrganizationEndpoint Endpoint
type DataOwnerEndpoint Endpoint

type Endpoint struct {
	*Dogmas            `json:"-"`
	EndpointIdentifier string  `json:"endpointIdentifier"`
	ByAuthorityOf      string  `json:"byAuthorityOf"`
	Method             string  `json:"method"`
	UrlHost            *string `json:"urlHost"`
	UrlPath            string  `json:"urlPath"`
}

// ================================================================
func (f *Feature) ByAuthorityOfOrganization() *OrganizationHttpMethod {
	return &OrganizationHttpMethod{
		Feature: f,
	}
}

func (m *OrganizationHttpMethod) GET(path, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "GET", path)
	m.RouterGroup.GET(path, handlerFuncs(e, handlers)...)
	return (*OrganizationEndpoint)(e)
}

func (m *OrganizationHttpMethod) POST(path, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "POST", path)
	m.RouterGroup.POST(path, handlerFuncs(e, handlers)...)
	return (*OrganizationEndpoint)(e)
}

func (m *OrganizationHttpMethod) PUT(path, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "PUT", path)
	m.RouterGroup.PUT(path, handlerFuncs(e, handlers)...)
	return (*OrganizationEndpoint)(e)
}

func (m *OrganizationHttpMethod) PATCH(path, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "PATCH", path)
	m.RouterGroup.PATCH(path, handlerFuncs(e, handlers)...)
	return (*OrganizationEndpoint)(e)
}

func (m *OrganizationHttpMethod) DELETE(path, identifier string, handlers ...HandlerFunc) *OrganizationEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfOrganization, "DELETE", path)
	m.RouterGroup.DELETE(path, handlerFuncs(e, handlers)...)
	return (*OrganizationEndpoint)(e)
}

// ================================================================
func (f *Feature) ByAuthorityOfDataOwner() *DataOwnerHttpMethod {
	return &DataOwnerHttpMethod{
		Feature: f,
	}
}

func (m *DataOwnerHttpMethod) GET(path, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "GET", path)
	m.RouterGroup.GET(path, handlerFuncs(e, handlers)...)
	return (*DataOwnerEndpoint)(e)
}

func (m *DataOwnerHttpMethod) POST(path, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "POST", path)
	m.RouterGroup.POST(path, handlerFuncs(e, handlers)...)
	return (*DataOwnerEndpoint)(e)
}

func (m *DataOwnerHttpMethod) PUT(path, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "PUT", path)
	m.RouterGroup.PUT(path, handlerFuncs(e, handlers)...)
	return (*DataOwnerEndpoint)(e)
}

func (m *DataOwnerHttpMethod) PATCH(path, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "PATCH", path)
	m.RouterGroup.PATCH(path, handlerFuncs(e, handlers)...)
	return (*DataOwnerEndpoint)(e)
}

func (m *DataOwnerHttpMethod) DELETE(path, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "DELETE", path)
	m.RouterGroup.DELETE(path, handlerFuncs(e, handlers)...)
	return (*DataOwnerEndpoint)(e)
}

// ================================================================
//
// ================================================================
type EndpointAuthorizationRule struct {
	EndpointIdentifier         *string `json:"endpointIdentifier"`
	Action                     string  `json:"action"`
	AffectedEndpointIdentifier string  `json:"affectedEndpointIdentifier"`
}

func (e *OrganizationEndpoint) CanAssignUserAccess(identifier string) *OrganizationEndpoint {
	e.Dogmas.EndpointAuthorizationRules = append(e.Dogmas.EndpointAuthorizationRules, &EndpointAuthorizationRule{
		EndpointIdentifier:         &e.EndpointIdentifier,
		Action:                     "ASSIGN",
		AffectedEndpointIdentifier: identifier,
	})
	return e
}

func (e *OrganizationEndpoint) CanGrantUserAccess(identifier string) *OrganizationEndpoint {
	e.Dogmas.EndpointAuthorizationRules = append(e.Dogmas.EndpointAuthorizationRules, &EndpointAuthorizationRule{
		EndpointIdentifier:         &e.EndpointIdentifier,
		Action:                     "GRANT",
		AffectedEndpointIdentifier: identifier,
	})
	return e
}

func (e *OrganizationEndpoint) CanRevokeUserAccess(identifier string) *OrganizationEndpoint {
	e.Dogmas.EndpointAuthorizationRules = append(e.Dogmas.EndpointAuthorizationRules, &EndpointAuthorizationRule{
		EndpointIdentifier:         &e.EndpointIdentifier,
		Action:                     "REVOKE",
		AffectedEndpointIdentifier: identifier,
	})
	return e
}

func (e *OrganizationEndpoint) CanUnassignUserAccess(identifier string) *OrganizationEndpoint {
	e.Dogmas.EndpointAuthorizationRules = append(e.Dogmas.EndpointAuthorizationRules, &EndpointAuthorizationRule{
		EndpointIdentifier:         &e.EndpointIdentifier,
		Action:                     "UNASSIGN",
		AffectedEndpointIdentifier: identifier,
	})
	return e
}

// ================================================================
type UserPrivileges struct {
	EndpointIdentifier          *string  `json:"endpointIdentifier"`
	Action                      string   `json:"action"`
	AffectedEndpointIdentifiers []string `json:"affectedEndpointIdentifiers"`
}

func (e OrganizationEndpoint) AssignUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "ASSIGN", identifiers)
}

func (e OrganizationEndpoint) GrantUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "GRANT", identifiers)
}

func (e OrganizationEndpoint) RevokeUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "REVOKE", identifiers)
}

func (e OrganizationEndpoint) UnassignUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "UNASSIGN", identifiers)
}

func (e OrganizationEndpoint) setUserPrivileges(userId xuuid.UUID, action string, identifiers []string) her.Error {
	if len(identifiers) > 0 {
		jsonbytes, err := json.Marshal(&UserPrivileges{
			EndpointIdentifier:          &e.EndpointIdentifier,
			Action:                      action,
			AffectedEndpointIdentifiers: identifiers,
		})
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		return apiPost(e.Dogmas.HostUrl.JoinPath("/privileges/v1/users", userId.String(), "/endpoints"), jsonbytes)
	}

	return nil
}

func (e OrganizationEndpoint) HasPrivilege(userId xuuid.UUID) (bool, her.Error) {
	apiUrl := e.Dogmas.HostUrl.JoinPath("/privileges/v1/users", userId.String(), "/endpoints", e.EndpointIdentifier)
	req, err := http.NewRequest("GET", apiUrl.String(), nil)
	if err != nil {
		return false, her.NewError(http.StatusInternalServerError, err, nil)
	}

	payload := her.NewPayload(nil)
	client := &http.Client{}

	if resp, err := client.Do(req); err != nil {
		return false, her.NewError(http.StatusInternalServerError, err, nil)
	} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
		return false, err
	} else {
		switch resp.StatusCode {
		case http.StatusOK:
			return true, nil
		case http.StatusNotFound:
			return false, nil
		}
	}

	return false, her.NewErrorWithMessage(http.StatusInternalServerError, "Dogmas: "+payload.Message, nil)
}

// ================================================================
//
// ================================================================
type Dogmas struct {
	HostUrl                    *url.URL
	AppHost                    string
	Endpoints                  []*Endpoint
	EndpointAuthorizationRules []*EndpointAuthorizationRule
}

func NewDogmas(appHostUrl *url.URL) (*Dogmas, error) {
	u, err := url.ParseRequestURI(os.Getenv("HOST_DOGMAS"))
	if err != nil {
		return nil, err
	}

	return &Dogmas{
		HostUrl: u,
		AppHost: appHostUrl.String(),
	}, nil
}

func (d *Dogmas) addEndpoint(identifier, byAuthorityOf, method, path string) *Endpoint {
	e := &Endpoint{
		Dogmas:             d,
		EndpointIdentifier: identifier,
		ByAuthorityOf:      byAuthorityOf,
		Method:             method,
		UrlHost:            &d.AppHost,
		UrlPath:            path,
	}
	d.Endpoints = append(d.Endpoints, e)
	return e
}

func (d Dogmas) RegisterEndpoints() her.Error {
	if len(d.Endpoints) > 0 {
		jsonbytes, err := json.Marshal(d.Endpoints)
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		return apiPost(d.HostUrl.JoinPath("/resources/v1/endpoints"), jsonbytes)
	}

	return nil
}

func (d Dogmas) RegisterEndpointAuthorizationRules() her.Error {
	if len(d.EndpointAuthorizationRules) > 0 {
		jsonbytes, err := json.Marshal(d.EndpointAuthorizationRules)
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		return apiPost(d.HostUrl.JoinPath("/resources/v1/rules"), jsonbytes)
	}

	return nil
}

// ================================================================
func apiPost(apiUrl *url.URL, jsonbytes []byte) her.Error {
	req, err := http.NewRequest("POST", apiUrl.String(), bytes.NewReader(jsonbytes))
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	payload := her.NewPayload(nil)
	client := &http.Client{}

	if resp, err := client.Do(req); err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
		return err
	} else if resp.StatusCode != 201 {
		return her.NewErrorWithMessage(http.StatusInternalServerError, "Dogmas: "+payload.Message, nil)
	}

	return nil
}

// ================================================================
const (
	Delimiter = " "
)

type Identifiers map[string]bool

func NewIdentifiers(input any) Identifiers {
	items := []string{}
	if reflect.TypeOf(input).Kind() == reflect.Slice {
		items = input.([]string)
	} else {
		items = strings.Split(input.(string), Delimiter)
	}

	identifiers := Identifiers{}
	for _, i := range items {
		identifiers.Set(i)
	}

	return identifiers
}

func (i *Identifiers) Set(item string) {
	(*i)[item] = true
}

func (i Identifiers) HasOneOf(sub Identifiers) bool {
	for item, has := range sub {
		if has {
			if val, ok := i[item]; ok && val {
				return true
			}
		}
	}
	return false
}

func (i Identifiers) Contains(sub Identifiers) bool {
	for item, has := range sub {
		if has {
			if val, ok := i[item]; !ok || !val {
				return false
			}
		}
	}
	return true
}

func (i Identifiers) Slice() []string {
	ss := []string{}
	for endpoint, has := range i {
		if has {
			ss = append(ss, endpoint)
		}
	}
	return ss
}
