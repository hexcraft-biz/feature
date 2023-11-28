package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"

	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

type Endpoint struct {
	EndpointId      Md5Identifier `json:"endpointId" db:"endpoint_id" binding:"required"`
	Ownership       string        `json:"ownership" db:"ownership" binding:"required"`
	Method          string        `json:"method" db:"method" binding:"required"`
	UrlHost         *string       `json:"urlHost" db:"url_host" binding:"required"`
	UrlFeature      *string       `json:"urlFeature" db:"url_feature" binding:"required"`
	UrlPath         string        `json:"urlPath" db:"url_path" binding:"required"`
	OwnerParamIndex int           `json:"ownerParamIndex" db:"owner_param_index" binding:"required"`
}

type EndpointHandler struct {
	*Dogmas `json:"-"`
	*Endpoint
}

func (e *EndpointHandler) SetAccessRulesFor(custodianId xuuid.UUID) *Authorizer {
	return &Authorizer{
		dogmasApiUrl:        e.Dogmas.HostUrl.JoinPath("/permissions/v1/custodians", custodianId.String()),
		EndpointId:          &e.EndpointId,
		accessRulesToCommit: map[int]map[Md5Identifier]*AccessRules{},
	}
}

// For resource to check
func (e EndpointHandler) CanBeAccessedBy(requesterId xuuid.UUID, subset string) (bool, her.Error) {
	jsonbytes, err := json.Marshal(map[string]string{
		"method":             e.Method,
		"requestEndpointUrl": *e.UrlHost + path.Join("/", *e.UrlFeature, subset),
		"requesterId":        requesterId.String(),
	})
	if err != nil {
		return false, her.NewError(http.StatusInternalServerError, err, nil)
	}

	req, err := http.NewRequest("POST", e.Dogmas.HostUrl.JoinPath("/permissions/v1/internal").String(), bytes.NewReader(jsonbytes))
	if err != nil {
		return false, her.NewError(http.StatusInternalServerError, err, nil)
	}

	result := new(ResultAccessPermission)
	payload := her.NewPayload(result)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return false, her.NewError(http.StatusInternalServerError, err, nil)
	} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return result.CanAccess, nil
	default:
		return false, her.NewErrorWithMessage(http.StatusInternalServerError, "Dogmas: "+payload.Message, nil)
	}
}
