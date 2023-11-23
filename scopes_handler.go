package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hexcraft-biz/her"
)

type endpoints map[Md5Identifier]*Endpoint

type scopesHandler struct {
	apiUrl string
	scopes map[string]*scopeWithEndpoints
}

type scopeWithEndpoints struct {
	*scope
	endpoints
}

type scope struct {
	identifier  string `json:"identifier"`
	description string `json:"description"`
}

func newScopesHandler(dogmasHost *url.URL) *scopesHandler {
	return &scopesHandler{
		apiUrl: dogmasHost.JoinPath("/resources/v1/scopes").String(),
		scopes: map[string]*scopeWithEndpoints{
			"": &scopeWithEndpoints{
				scope:     nil,
				endpoints: endpoints{},
			},
		},
	}
}

func (h *scopesHandler) AddScope(identifier, description string) {
	h.scopes[identifier] = &scopeWithEndpoints{
		scope: &scope{
			identifier:  identifier,
			description: description,
		},
		endpoints: endpoints{},
	}
}

func (h scopesHandler) endpointsContainer(identifier string) *scopeWithEndpoints {
	se, ok := h.scopes[identifier]
	if !ok {
		panic("No such scope(s) to add endpoint")
	}
	return se
}

type registerScope struct {
	*scope
	endpoints []*Endpoint `json:"endpoints"`
}

func (h scopesHandler) Register() {
	scopes := []*registerScope{}

	for _, se := range h.scopes {
		scope := &registerScope{
			scope:     se.scope,
			endpoints: []*Endpoint{},
		}
		scopes = append(scopes, scope)
		for _, e := range se.endpoints {
			scope.endpoints = append(scope.endpoints, e)
		}
	}

	if len(scopes) > 0 {
		jsonbytes, err := json.Marshal(scopes)
		if err != nil {
			panic(err.Error())
		}

		req, err := http.NewRequest("POST", h.apiUrl, bytes.NewReader(jsonbytes))
		if err != nil {
			panic(err.Error())
		}

		payload := her.NewPayload(nil)
		client := &http.Client{}

		if resp, err := client.Do(req); err != nil {
			panic(err.Error())
		} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
			panic(err.Error())
		} else if resp.StatusCode != http.StatusCreated {
			panic("Dogmas: " + payload.Message)
		}
	}
}

func (se *scopeWithEndpoints) addEndpoint(e *Endpoint) *scopeWithEndpoints {
	se.endpoints[e.EndpointId] = e
	return se
}
