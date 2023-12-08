package feature

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	paging "github.com/hexcraft-biz/envmod-mysql"
	"github.com/hexcraft-biz/her"
)

func newScopesHandler(dogmasRootUrl *url.URL) ScopesHandler {
	return ScopesHandler{
		dogmasRootUrl: dogmasRootUrl,
		Maps: Maps{
			"": newScopeWithEndpoints("", ""),
		},
	}
}

type Maps map[string]*scopeWithEndpoints

type ScopesHandler struct {
	dogmasRootUrl *url.URL
	Maps
}

func (h *ScopesHandler) AddScope(identifier, description string) {
	if _, ok := h.Maps[identifier]; !ok {
		h.Maps[identifier] = newScopeWithEndpoints(identifier, description)
	}
}

func (h ScopesHandler) Scope(identifier string) *scopeWithEndpoints {
	se, ok := h.Maps[identifier]
	if !ok {
		panic("No such scope(s) to add endpoint")
	}

	return se
}

func (h ScopesHandler) register() error {
	scopes := []*scopeWithEndpoints{}

	for _, se := range h.Maps {
		scopes = append(scopes, se)
	}

	if len(scopes) > 0 {
		jsonbytes, err := json.Marshal(scopes)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", h.dogmasRootUrl.JoinPath("/resources/v1/scopes").String(), bytes.NewReader(jsonbytes))
		if err != nil {
			return err
		}

		payload := her.NewPayload(nil)
		client := &http.Client{}

		if resp, err := client.Do(req); err != nil {
			return err
		} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
			return errors.New("Dogmas: " + err.Error())
		} else if resp.StatusCode != http.StatusCreated {
			return errors.New("Dogmas: " + payload.Message)
		}
	}

	return nil
}

func (h *ScopesHandler) SyncEndpoints(appRootUrl *url.URL) error {
	endpoints := map[*Endpoint]struct{}{}
	for _, se := range h.Maps {
		for _, e := range se.Endpoints {
			endpoints[e] = struct{}{}
		}
	}

	u := h.dogmasRootUrl.JoinPath("/resources/v1/endpoints")
	q := u.Query()
	q.Set("host", appRootUrl.String())
	u.RawQuery = q.Encode()

	urlstring := u.String()
	next := &urlstring
	for next != nil {
		result := new(resultSyncEndpoints)
		payload := her.NewPayload(result)
		if resp, err := http.Get(*next); err != nil {
			return err
		} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
			return err
		} else if resp.StatusCode != http.StatusOK {
			return errors.New("Dogmas: " + payload.Message)
		}

		for _, r := range result.Endpoints {
			for e := range endpoints {
				if e.Method == r.Method &&
					e.UrlHost == r.UrlHost &&
					e.UrlFeature == r.UrlFeature &&
					e.UrlPath == r.UrlPath {
					e.EndpointId = r.EndpointId
					e.DestHost = r.DestHost
					e.Actived = r.Actived
					delete(endpoints, e)
				}
			}
		}

		next = result.Paging.Next
	}

	return nil
}

func (h ScopesHandler) HasUnsynchronized() bool {
	for _, se := range h.Maps {
		for _, e := range se.Endpoints {
			if !e.Actived || e.EndpointId.IsZero() {
				return true
			}
		}
	}

	return false
}

type resultSyncEndpoints struct {
	Endpoints      []*Endpoint `json:"endpoints"`
	*paging.Paging `json:"paging"`
}

// ================================================================
func newScopeWithEndpoints(identifier, description string) *scopeWithEndpoints {
	if identifier == "" {
		return &scopeWithEndpoints{
			Scope:     nil,
			Endpoints: []*Endpoint{},
		}
	} else {
		if description == "" {
			panic("Empty scope description")
		}

		return &scopeWithEndpoints{
			Scope: &Scope{
				Identifier:  identifier,
				Description: description,
			},
			Endpoints: []*Endpoint{},
		}
	}
}

type scopeWithEndpoints struct {
	*Scope
	Endpoints []*Endpoint `json:"endpoints"`
}

func (se *scopeWithEndpoints) AddEndpoint(e *Endpoint) {
	se.Endpoints = append(se.Endpoints, e)
}

// ================================================================
type Scope struct {
	Identifier  string `json:"identifier"`
	Description string `json:"description"`
}
