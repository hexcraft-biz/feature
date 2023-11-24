package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

// ================================================================
//
// ================================================================
type OrganizationHttpMethods struct {
	*Feature
}

func (m *OrganizationHttpMethods) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfOrganization, "GET", relativePath, scopes)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationHttpMethods) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfOrganization, "POST", relativePath, scopes)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationHttpMethods) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfOrganization, "PUT", relativePath, scopes)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationHttpMethods) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfOrganization, "PATCH", relativePath, scopes)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationHttpMethods) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfOrganization, "DELETE", relativePath, scopes)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}

// ================================================================
const (
	ActionAssign = iota
	ActionGrant
	ActionRevoke
)

const (
	writeBehaviorUndef = iota
	writeBehaviorIfNotExists
	writeBehaviorOverwrite
)

func (e *Endpoint) SetEndpointAccessRulesFor(userId xuuid.UUID) *organizationEndpointPermission {
	return &organizationEndpointPermission{
		dogmasApiUrl:        e.Dogmas.HostUrl.JoinPath("/permissions/v1/users", userId.String()),
		EndpointId:          &e.EndpointId,
		accessRulesToCommit: map[int]map[Md5Identifier]*EndpointAccessRules{},
	}
}

type organizationEndpointPermission struct {
	dogmasApiUrl        *url.URL
	EndpointId          *Md5Identifier
	accessRulesToCommit map[int]map[Md5Identifier]*EndpointAccessRules
}

func (u *organizationEndpointPermission) AffectedEndpoint(affectedEndpointId Md5Identifier) *affectedEndpointAccessRules {
	return &affectedEndpointAccessRules{
		organizationEndpointPermission: u,
		affectedEndpointId:             affectedEndpointId,
	}
}

type affectedEndpointAccessRules struct {
	*organizationEndpointPermission
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

type EndpointAccessRulesWithBehavior struct {
	Behavior           string               `json:"behavior" db:"-" binding:"required"`
	AffectedEndpointId Md5Identifier        `json:"affectedEndpointId" db:"endpoint_id" binding:"required"`
	AccessRules        *EndpointAccessRules `json:"accessRules" db:"access_rules" binding:"required"`
}

const (
	HeaderEndpointId = "X-Endpoint-Id"
	HeaderByUserId   = "X-By-User-Id"
)

func (u organizationEndpointPermission) Commit(byUserId xuuid.UUID) her.Error {
	rulesWithBehavior := []*EndpointAccessRulesWithBehavior{}
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
			rulesWithBehavior = append(rulesWithBehavior, &EndpointAccessRulesWithBehavior{
				Behavior:           behaviorstring,
				AffectedEndpointId: id,
				AccessRules:        accessRules,
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
		req.Header.Set(HeaderByUserId, byUserId.String())

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
