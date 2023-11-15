package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/her"
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

func (f *Feature) GET(relativePath, identifier, description string, handlers ...gin.HandlerFunc) *Scope {
	f.RouterGroup.GET(relativePath, handlers...)
	return f.Dogmas.addScope(identifier, description)
}

func (f *Feature) POST(relativePath, identifier, description string, handlers ...gin.HandlerFunc) *Scope {
	f.RouterGroup.POST(relativePath, handlers...)
	return f.Dogmas.addScope(identifier, description)
}

func (f *Feature) PUT(relativePath, identifier, description string, handlers ...gin.HandlerFunc) *Scope {
	f.RouterGroup.PUT(relativePath, handlers...)
	return f.Dogmas.addScope(identifier, description)
}

func (f *Feature) PATCH(relativePath, identifier, description string, handlers ...gin.HandlerFunc) *Scope {
	f.RouterGroup.PATCH(relativePath, handlers...)
	return f.Dogmas.addScope(identifier, description)
}

func (f *Feature) DELETE(relativePath, identifier, description string, handlers ...gin.HandlerFunc) *Scope {
	f.RouterGroup.DELETE(relativePath, handlers...)
	return f.Dogmas.addScope(identifier, description)
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
		Endpoint: u.JoinPath("/resources/v1"),
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

	req, err := http.NewRequest("POST", d.Endpoint.JoinPath("scopes").String(), bytes.NewReader(jsonbytes))
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

func (d Dogmas) AuthorizationRulesRegister() her.Error {
	jsonbytes, err := json.Marshal(d.ScopeAuthorizationRules)
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	req, err := http.NewRequest("POST", d.Endpoint.JoinPath("rules").String(), bytes.NewReader(jsonbytes))
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
