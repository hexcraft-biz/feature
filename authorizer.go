package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

const (
	HeaderEndpointId    = "X-Endpoint-Id"
	HeaderByCustodianId = "X-By-Custodian-Id"
)

const (
	ActionAssign int = iota
	ActionGrant
	ActionRevoke
)

const (
	writeBehaviorUndef int = iota
	writeBehaviorIfNotExists
	writeBehaviorOverwrite
)

type Authorizer struct {
	dogmasApiUrl        *url.URL
	EndpointId          *Md5Identifier
	accessRulesToCommit map[int]map[Md5Identifier]*EndpointAccessRules
}

func (u *Authorizer) AffectedEndpoint(affectedEndpointId Md5Identifier) *affectedEndpointAccessRules {
	return &affectedEndpointAccessRules{
		Authorizer:         u,
		affectedEndpointId: affectedEndpointId,
	}
}

func (u Authorizer) Commit(byCustodianId xuuid.UUID) her.Error {
	rulesWithBehavior := []*AccessRulesWithBehavior{}
	for behavior, idAccessRules := range u.accessRulesToCommit {

		behaviorstring := ""
		switch behavior {
		case writeBehaviorIfNotExists:
			behaviorstring = "IF_NOT_EXISTS"
		case writeBehaviorOverwrite:
			behaviorstring = "OVERWRITE"
		default:
			return her.NewErrorWithMessage(http.StatusInternalServerError, "Undefined write behavior", nil)
		}

		for id, accessRules := range idAccessRules {
			accessRules.RemoveRedundant()
			rulesWithBehavior = append(rulesWithBehavior, &AccessRulesWithBehavior{
				AffectedEndpointId: id,
				AccessRulesSettingBehavior: &AccessRulesSettingBehavior{
					Behavior:    behaviorstring,
					AccessRules: accessRules,
				},
			})
		}
	}

	if len(rulesWithBehavior) > 0 {
		jsonbytes, err := json.Marshal(rulesWithBehavior)
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		req, err := http.NewRequest("POST", u.dogmasApiUrl.String(), bytes.NewReader(jsonbytes))
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		req.Header.Set(HeaderEndpointId, string(*u.EndpointId))
		req.Header.Set(HeaderByCustodianId, byCustodianId.String())

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

type AccessRulesWithBehavior struct {
	AffectedEndpointId Md5Identifier `json:"affectedEndpointId" db:"endpoint_id" binding:"required"`
	*AccessRulesSettingBehavior
}

type affectedEndpointAccessRules struct {
	*Authorizer
	affectedEndpointId Md5Identifier
}

func (u *affectedEndpointAccessRules) Assign(rule string) *affectedEndpointAccessRules {
	return u.addAction(ActionAssign, rule)
}

func (u *affectedEndpointAccessRules) Grant(rule string) *affectedEndpointAccessRules {
	return u.addAction(ActionGrant, rule)
}

func (u *affectedEndpointAccessRules) Revoke(rule string) *affectedEndpointAccessRules {
	return u.addAction(ActionRevoke, rule)
}

func (u *affectedEndpointAccessRules) addAction(action int, rule string) *affectedEndpointAccessRules {
	behavior := writeBehaviorUndef
	switch action {
	case ActionGrant, ActionRevoke:
		behavior = writeBehaviorOverwrite
	default:
		behavior = writeBehaviorIfNotExists
	}

	if _, ok := u.accessRulesToCommit[behavior]; !ok {
		u.accessRulesToCommit[behavior] = map[Md5Identifier]*EndpointAccessRules{}
	}

	if _, ok := u.accessRulesToCommit[behavior][u.affectedEndpointId]; !ok {
		u.accessRulesToCommit[behavior][u.affectedEndpointId] = &EndpointAccessRules{}
	}

	switch action {
	case ActionAssign, ActionGrant:
		u.accessRulesToCommit[behavior][u.affectedEndpointId].AddSubset(rule)
	case ActionRevoke:
		u.accessRulesToCommit[behavior][u.affectedEndpointId].AddException(rule)
	}

	return u
}
