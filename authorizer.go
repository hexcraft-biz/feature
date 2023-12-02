package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

func newAuthorizer(dogmasHost *url.URL, custodianId, viaEndpointId xuuid.UUID) *Authorizer {
	return &Authorizer{
		dogmasApiUrl:        dogmasHost.JoinPath("/permissions/v1/custodians", custodianId.String()),
		viaEndpointId:       viaEndpointId,
		accessRulesToCommit: accessRulesToCommit{},
	}
}

type Authorizer struct {
	dogmasApiUrl  *url.URL
	viaEndpointId xuuid.UUID
	accessRulesToCommit
}

func (u *Authorizer) AffectedEndpoint(affectedEndpointId xuuid.UUID) *affectedEndpointAccessRules {
	return &affectedEndpointAccessRules{
		Authorizer:         u,
		affectedEndpointId: affectedEndpointId,
	}
}

func (u Authorizer) Commit(byCustodianId xuuid.UUID) her.Error {
	rulesWithBehavior := u.toAccessRulesWithBehavior()

	if len(rulesWithBehavior) > 0 {
		jsonbytes, err := json.Marshal(rulesWithBehavior)
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		req, err := http.NewRequest("POST", u.dogmasApiUrl.String(), bytes.NewReader(jsonbytes))
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		req.Header.Set(HeaderViaEndpointId, u.viaEndpointId.String())
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

type affectedEndpointAccessRules struct {
	*Authorizer
	affectedEndpointId xuuid.UUID
}

func (r *affectedEndpointAccessRules) Assign(rule string) {
	r.accessRulesToCommit.add(ActionAssign, rule, r.affectedEndpointId)
}

func (r *affectedEndpointAccessRules) Grant(rule string) {
	r.accessRulesToCommit.add(ActionGrant, rule, r.affectedEndpointId)
}

func (r *affectedEndpointAccessRules) Revoke(rule string) {
	r.accessRulesToCommit.add(ActionRevoke, rule, r.affectedEndpointId)
}

type accessRulesToCommit map[string]map[xuuid.UUID]*AccessRules

func (r *accessRulesToCommit) add(action int, rule string, affectedEndpointId xuuid.UUID) {
	behavior := ""
	switch action {
	case ActionGrant, ActionRevoke:
		behavior = WriteBehaviorOverwrite
	default:
		behavior = WriteBehaviorIdempotent
	}

	if _, ok := (*r)[behavior]; !ok {
		(*r)[behavior] = map[xuuid.UUID]*AccessRules{}
	}

	if _, ok := (*r)[behavior][affectedEndpointId]; !ok {
		(*r)[behavior][affectedEndpointId] = &AccessRules{}
	}

	switch action {
	case ActionAssign, ActionGrant:
		(*r)[behavior][affectedEndpointId].AddSubset(rule)
	case ActionRevoke:
		(*r)[behavior][affectedEndpointId].AddException(rule)
	}
}

func (r accessRulesToCommit) toAccessRulesWithBehavior() []*AccessRulesWithBehavior {
	rulesWithBehavior := []*AccessRulesWithBehavior{}
	for behavior, idAccessRules := range r {
		for affectedEndpointId, accessRules := range idAccessRules {
			accessRules.RemoveRedundant()
			rulesWithBehavior = append(rulesWithBehavior, &AccessRulesWithBehavior{
				Behavior: behavior,
				AuthorizationAction: &AuthorizationAction{
					AffectedEndpointId: affectedEndpointId,
					AccessRules:        accessRules,
				},
			})
		}
	}

	return rulesWithBehavior
}
