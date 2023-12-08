package feature

import (
	"errors"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

func newEndpoint(ownership, method, urlHost, urlFeature, urlPath string) *Endpoint {
	return &Endpoint{
		Ownership:  ownership,
		Method:     method,
		DestHost:   defaultDestHost(urlHost),
		UrlHost:    urlHost,
		UrlFeature: urlFeature,
		UrlPath:    urlPath,
	}
}

type Endpoint struct {
	EndpointId xuuid.UUID `json:"endpointId" db:"endpoint_id" binding:"-"`
	Actived    bool       `json:"actived" db:"actived" binding:"-"`
	Ownership  string     `json:"ownership" db:"ownership" binding:"required"`
	Method     string     `json:"method" db:"method" binding:"required"`
	DestHost   string     `json:"destHost" db:"dest_host" binding:"required"`
	UrlHost    string     `json:"urlHost" db:"url_host" binding:"required"`
	UrlFeature string     `json:"urlFeature" db:"url_feature" binding:"required"`
	UrlPath    string     `json:"urlPath" db:"url_path" binding:"required"`
}

type EndpointHandler struct {
	*Dogmas `json:"-"`
	*Endpoint
}

func (e *EndpointHandler) SetAccessRulesFor(custodianId xuuid.UUID) *Authorizer {
	if e.Ownership != "ORGANIZATION" && e.Ownership != "PRIVATE" {
		return nil
	}

	return newAuthorizer(e.Dogmas.RootUrl, custodianId, e.EndpointId)
}

// For resource to check
func (e EndpointHandler) CanBeAccessedBy(requesterId xuuid.UUID, requestUrlPath string) her.Error {
	_, err := e.Dogmas.canBeAccessedBy(nil, e.Method, e.UrlHost+path.Join("/", e.UrlFeature, requestUrlPath), &requesterId)
	return err
}

// ================================================================
type RequestedUrlString string

func (s RequestedUrlString) Parse(method string) (*RequestedEndpointHandler, error) {
	u, err := url.ParseRequestURI(string(s))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`/v[0-9]+`)
	loc := re.FindStringIndex(u.Path)

	if loc == nil {
		return nil, errors.New("Invalid endpoint")
	}

	urlHost, appRoot := u.Scheme+"://", u.Host
	urlFeature, requestedPath := u.Path[0:loc[1]], u.Path[loc[1]:]

	segs := strings.Split(urlFeature, "/")
	if len(segs) > 3 {
		appRoot = path.Join(appRoot, strings.Join(segs[0:len(segs)-2], "/"))
		urlFeature = path.Join("/", strings.Join(segs[len(segs)-2:], "/"))
	}
	urlHost += appRoot

	segs = strings.Split(requestedPath, "/")
	subsetSegs := []string{""}
	possibleOwnerIdString := ""
	for i := range segs {
		if i > 0 && i%2 == 0 {
			subsetSegs = append(subsetSegs, segs[i])
			segs[i] = "*"
		}

		if i == 2 {
			possibleOwnerIdString = segs[i]
		}
	}

	if u.RawQuery != "" {
		requestedPath += "?" + u.RawQuery
	}

	return &RequestedEndpointHandler{
		Endpoint: &Endpoint{
			Method:     method,
			DestHost:   defaultDestHost(urlHost),
			UrlHost:    urlHost,
			UrlFeature: urlFeature,
			UrlPath:    strings.Join(segs, "/"),
		},
		RequestedPath:         requestedPath,
		SubsetToCheck:         strings.Join(subsetSegs, "/"),
		possibleOwnerIdString: possibleOwnerIdString,
	}, nil
}

type RequestedEndpointHandler struct {
	*Endpoint
	RequestedPath         string
	SubsetToCheck         string
	possibleOwnerIdString string
}

func (h RequestedEndpointHandler) GetOwnerId() (xuuid.UUID, error) {
	if h.Ownership == OwnershipPrivate {
		return xuuid.Parse(h.possibleOwnerIdString)
	} else {
		return xuuid.UUID(uuid.Nil), errors.New("invalid ownership endpoint")
	}
}

func (h RequestedEndpointHandler) Route() *Route {
	return &Route{
		Method:  h.Method,
		RootUrl: h.DestHost,
		Feature: h.UrlFeature,
		Path:    h.RequestedPath,
	}
}

type Route struct {
	Method  string `json:"method"`
	RootUrl string `json:"rootUrl"`
	Feature string `json:"feature"`
	Path    string `json:"path"`
}
