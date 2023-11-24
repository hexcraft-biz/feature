package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

type Endpoint struct {
	*Dogmas       `json:"-"`
	EndpointId    Md5Identifier `json:"endpointId"`
	ByAuthorityOf string        `json:"byAuthorityOf"`
	Method        string        `json:"method"`
	UrlHost       *string       `json:"urlHost"`
	UrlFeature    *string       `json:"urlFeature"`
	UrlPath       string        `json:"urlPath"`
}

func (e Endpoint) CanBeAccessedBy(userId xuuid.UUID, subset *url.URL) her.Error {
	jsonbytes, err := json.Marshal(map[string]string{
		"method":             e.Method,
		"requestEndpointUrl": *e.UrlHost + path.Join("/", *e.UrlFeature, subset.String()),
		"userId":             userId.String(),
	})
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	req, err := http.NewRequest("POST", e.Dogmas.HostUrl.JoinPath("/permissions/v1/internal").String(), bytes.NewReader(jsonbytes))
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	payload := her.NewPayload(nil)
	client := &http.Client{}

	if resp, err := client.Do(req); err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return her.NewErrorWithMessage(http.StatusForbidden, "Dogmas: "+payload.Message, nil)
	}

	return nil
}
