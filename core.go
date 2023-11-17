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

func handlerFuncs(s *Endpoint, handlers []HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = h(s)
	}
	return funcs
}

func (f *Feature) GET(path, identifier, onBehaviorOf string, handlers ...HandlerFunc) *Endpoint {
	s := f.Dogmas.addEndpoint(identifier, onBehaviorOf, "GET", path)
	f.RouterGroup.GET(path, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) POST(path, identifier, onBehaviorOf string, handlers ...HandlerFunc) *Endpoint {
	s := f.Dogmas.addEndpoint(identifier, onBehaviorOf, "POST", path)
	f.RouterGroup.POST(path, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) PUT(path, identifier, onBehaviorOf string, handlers ...HandlerFunc) *Endpoint {
	s := f.Dogmas.addEndpoint(identifier, onBehaviorOf, "PUT", path)
	f.RouterGroup.PUT(path, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) PATCH(path, identifier, onBehaviorOf string, handlers ...HandlerFunc) *Endpoint {
	s := f.Dogmas.addEndpoint(identifier, onBehaviorOf, "PATCH", path)
	f.RouterGroup.PATCH(path, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) DELETE(path, identifier, onBehaviorOf string, handlers ...HandlerFunc) *Endpoint {
	s := f.Dogmas.addEndpoint(identifier, onBehaviorOf, "DELETE", path)
	f.RouterGroup.DELETE(path, handlerFuncs(s, handlers)...)
	return s
}

// ================================================================
//
// ================================================================
const (
	OnBehaviorOfOrganization = "ORGANIZATION"
	OnBehaviorOfDataOwner    = "DATA_OWNER"
)

type Endpoint struct {
	*Dogmas            `json:"-"`
	EndpointIdentifier string  `json:"endpointIdentifier"`
	OnBehaviorOf       string  `json:"onBehaviorOf"`
	Method             string  `json:"method"`
	UrlHost            *string `json:"urlHost"`
	UrlPath            string  `json:"urlPath"`
}

type EndpointAuthorizationRule struct {
	EndpointIdentifier         *string `json:"endpointIdentifier"`
	Action                     string  `json:"action"`
	AffectedEndpointIdentifier string  `json:"affectedEndpointIdentifier"`
}

type UserPrivileges struct {
	EndpointIdentifier          *string  `json:"endpointIdentifier"`
	Action                      string   `json:"action"`
	AffectedEndpointIdentifiers []string `json:"affectedEndpointIdentifiers"`
}

func (e *Endpoint) CanAssignUserAccess(identifier string) *Endpoint {
	e.Dogmas.addEndpointAuthorizationRule(e, "ASSIGN", identifier)
	return e
}

func (e *Endpoint) CanGrantUserAccess(identifier string) *Endpoint {
	e.Dogmas.addEndpointAuthorizationRule(e, "GRANT", identifier)
	return e
}

func (e *Endpoint) CanRevokeUserAccess(identifier string) *Endpoint {
	e.Dogmas.addEndpointAuthorizationRule(e, "REVOKE", identifier)
	return e
}

func (e *Endpoint) CanUnassignUserAccess(identifier string) *Endpoint {
	e.Dogmas.addEndpointAuthorizationRule(e, "UNASSIGN", identifier)
	return e
}

// ================================================================
func (e Endpoint) AssignUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "ASSIGN", identifiers)
}

func (e Endpoint) GrantUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "GRANT", identifiers)
}

func (e Endpoint) RevokeUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "REVOKE", identifiers)
}

func (e Endpoint) UnassignUserPrivileges(userId xuuid.UUID, identifiers []string) her.Error {
	return e.setUserPrivileges(userId, "UNASSIGN", identifiers)
}

func (e Endpoint) setUserPrivileges(userId xuuid.UUID, action string, identifiers []string) her.Error {
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

func (e Endpoint) HasPrivilege(userId xuuid.UUID) (bool, her.Error) {
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

func (d *Dogmas) addEndpoint(identifier, onBehaviorOf, method, path string) *Endpoint {
	e := &Endpoint{
		Dogmas:             d,
		EndpointIdentifier: identifier,
		OnBehaviorOf:       onBehaviorOf,
		Method:             method,
		UrlHost:            &d.AppHost,
		UrlPath:            path,
	}
	d.Endpoints = append(d.Endpoints, e)
	return e
}

func (d *Dogmas) addEndpointAuthorizationRule(e *Endpoint, action, identifier string) {
	if e.OnBehaviorOf == OnBehaviorOfOrganization {
		d.EndpointAuthorizationRules = append(d.EndpointAuthorizationRules, &EndpointAuthorizationRule{
			EndpointIdentifier:         &e.EndpointIdentifier,
			Action:                     action,
			AffectedEndpointIdentifier: identifier,
		})
	}
}

func (d Dogmas) RegisterEndpoints() her.Error {
	jsonbytes, err := json.Marshal(d.Endpoints)
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	return apiPost(d.HostUrl.JoinPath("/resources/v1/endpoints"), jsonbytes)
}

func (d Dogmas) RegisterEndpointAuthorizationRules() her.Error {
	jsonbytes, err := json.Marshal(d.EndpointAuthorizationRules)
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	return apiPost(d.HostUrl.JoinPath("/resources/v1/rules"), jsonbytes)
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
