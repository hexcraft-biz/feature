package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

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

type HandlerFunc func(*Scope) gin.HandlerFunc

func handlerFuncs(s *Scope, handlers []HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = h(s)
	}
	return funcs
}

func (f *Feature) GET(relativePath, identifier, description string, handlers ...HandlerFunc) *Scope {
	s := f.Dogmas.addScope(identifier, description)
	f.RouterGroup.GET(relativePath, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) POST(relativePath, identifier, description string, handlers ...HandlerFunc) *Scope {
	s := f.Dogmas.addScope(identifier, description)
	f.RouterGroup.POST(relativePath, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) PUT(relativePath, identifier, description string, handlers ...HandlerFunc) *Scope {
	s := f.Dogmas.addScope(identifier, description)
	f.RouterGroup.PUT(relativePath, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) PATCH(relativePath, identifier, description string, handlers ...HandlerFunc) *Scope {
	s := f.Dogmas.addScope(identifier, description)
	f.RouterGroup.PATCH(relativePath, handlerFuncs(s, handlers)...)
	return s
}

func (f *Feature) DELETE(relativePath, identifier, description string, handlers ...HandlerFunc) *Scope {
	s := f.Dogmas.addScope(identifier, description)
	f.RouterGroup.DELETE(relativePath, handlerFuncs(s, handlers)...)
	return s
}

// ================================================================
//
// ================================================================
type Scope struct {
	*Dogmas     `json:"-"`
	Identifier  string `json:"identifier"`
	Description string `json:"description"`
}

type ScopeAuthorizationRule struct {
	Scope         *string `json:"scope"`
	Action        string  `json:"action"`
	AffectedScope string  `json:"affectedScope"`
}

type UserPermissions struct {
	Scope         *string  `json:"scope"`
	AffectedScope []string `json:"affectedScopes"`
}

func (s *Scope) CanAssignUserAccess(affectedScope string) *Scope {
	s.Dogmas.addAssignScopeAuthorizationRule(s, affectedScope)
	return s
}

func (s *Scope) CanGrantUserAccess(affectedScope string) *Scope {
	s.Dogmas.addGrantScopeAuthorizationRule(s, affectedScope)
	return s
}

func (s *Scope) CanRevokeUserAccess(affectedScope string) *Scope {
	s.Dogmas.addRevokeScopeAuthorizationRule(s, affectedScope)
	return s
}

func (s Scope) SetUserPermissions(userId xuuid.UUID, scopes []string) her.Error {
	jsonbytes, err := json.Marshal(&UserPermissions{
		Scope:         &s.Identifier,
		AffectedScope: scopes,
	})
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	return apiPost(s.Dogmas.Endpoint.JoinPath("/permissions/v1/users", userId.String(), "/scopes"), jsonbytes)
}

// ================================================================
//
// ================================================================
type Dogmas struct {
	Host                    string
	Endpoint                *url.URL
	Scopes                  []*Scope
	ScopeAuthorizationRules []*ScopeAuthorizationRule
}

func NewDogmas() (*Dogmas, error) {
	u, err := url.ParseRequestURI(os.Getenv("HOST_DOGMAS"))
	if err != nil {
		return nil, err
	}

	return &Dogmas{
		Host:     u.String(),
		Endpoint: u,
	}, nil
}

func (d *Dogmas) addScope(identifier, description string) *Scope {
	s := &Scope{
		Dogmas:      d,
		Identifier:  identifier,
		Description: description,
	}
	d.Scopes = append(d.Scopes, s)
	return s
}

func (d *Dogmas) addAssignScopeAuthorizationRule(scope *Scope, affectedScope string) {
	d.ScopeAuthorizationRules = append(d.ScopeAuthorizationRules, &ScopeAuthorizationRule{
		Scope:         &scope.Identifier,
		Action:        "ASSIGN",
		AffectedScope: affectedScope,
	})
}

func (d *Dogmas) addGrantScopeAuthorizationRule(scope *Scope, affectedScope string) {
	d.ScopeAuthorizationRules = append(d.ScopeAuthorizationRules, &ScopeAuthorizationRule{
		Scope:         &scope.Identifier,
		Action:        "GRANT",
		AffectedScope: affectedScope,
	})
}

func (d *Dogmas) addRevokeScopeAuthorizationRule(scope *Scope, affectedScope string) {
	d.ScopeAuthorizationRules = append(d.ScopeAuthorizationRules, &ScopeAuthorizationRule{
		Scope:         &scope.Identifier,
		Action:        "REVOKE",
		AffectedScope: affectedScope,
	})
}

func (d Dogmas) ScopesRegister() her.Error {
	jsonbytes, err := json.Marshal(d.Scopes)
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	return apiPost(d.Endpoint.JoinPath("/resources/v1/scopes"), jsonbytes)
}

func (d Dogmas) AuthorizationRulesRegister() her.Error {
	jsonbytes, err := json.Marshal(d.ScopeAuthorizationRules)
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	return apiPost(d.Endpoint.JoinPath("/resources/v1/rules"), jsonbytes)
}

// ================================================================
func apiPost(endpoint *url.URL, jsonbytes []byte) her.Error {
	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(jsonbytes))
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
