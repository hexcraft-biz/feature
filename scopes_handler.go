package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hexcraft-biz/her"
)

type endpointsMap map[Md5Identifier]*Endpoint

type Scope struct {
	Identifier  string `json:"identifier"`
	Description string `json:"description"`
}

type scopesHandler struct {
	apiUrl string
	scopes map[string]*scopeWithEndpoints
}

type scopeWithEndpoints struct {
	*Scope
	endpointsMap
}

func newScopesHandler(dogmasHost *url.URL) *scopesHandler {
	return &scopesHandler{
		apiUrl: dogmasHost.JoinPath("/resources/v1/scopes").String(),
		scopes: map[string]*scopeWithEndpoints{
			"": &scopeWithEndpoints{
				Scope:        nil,
				endpointsMap: endpointsMap{},
			},
		},
	}
}

func (h *scopesHandler) AddScope(identifier, description string) {
	if _, ok := h.scopes[identifier]; !ok {
		h.scopes[identifier] = &scopeWithEndpoints{
			Scope: &Scope{
				Identifier:  identifier,
				Description: description,
			},
			endpointsMap: endpointsMap{},
		}
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
	*Scope
	Endpoints []*Endpoint `json:"endpoints"`
}

func (h scopesHandler) register() {
	scopes := []*registerScope{}

	for _, se := range h.scopes {
		scope := &registerScope{
			Scope:     se.Scope,
			Endpoints: []*Endpoint{},
		}
		scopes = append(scopes, scope)
		for _, e := range se.endpointsMap {
			scope.Endpoints = append(scope.Endpoints, e)
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
	se.endpointsMap[e.EndpointId] = e
	return se
}
