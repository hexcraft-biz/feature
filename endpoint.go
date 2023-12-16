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

func newEndpoint(ownership, method, srcApp, appFeature, appPath string) *Endpoint {
	return &Endpoint{
		Activated:   true,
		FullProxied: true,
		Ownership:   ownership,
		Method:      method,
		DstApp:      defaultDestHost(srcApp),
		SrcApp:      srcApp,
		AppFeature:  appFeature,
		AppPath:     appPath,
	}
}

type Endpoint struct {
	EndpointId  xuuid.UUID `json:"endpointId" db:"endpoint_id" binding:"-"`
	Activated   bool       `json:"activated" db:"activated" binding:"-"`
	FullProxied bool       `json:"fullProxied" db:"full_proxied" binding:"-"`
	Ownership   string     `json:"ownership" db:"ownership" binding:"required"`
	Method      string     `json:"method" db:"method" binding:"required"`
	DstApp      string     `json:"dstApp" db:"dst_app" binding:"required"`
	SrcApp      string     `json:"srcApp" db:"src_app" binding:"required"`
	AppFeature  string     `json:"appFeature" db:"app_feature" binding:"required"`
	AppPath     string     `json:"appPath" db:"app_path" binding:"required"`
}

type EndpointHandler struct {
	*Dogmas `json:"-"`
	*Endpoint
}

func (e *EndpointHandler) SetAccessRulesFor(custodianId xuuid.UUID) *Authorizer {
	if e.Ownership != "ORGANIZATION" && e.Ownership != "PRIVATE" {
		return nil
	}

	return newAuthorizer(e.Dogmas.CreedsRootUrl, custodianId)
}

// For resource to check
func (e EndpointHandler) CanBeAccessedBy(requesterId xuuid.UUID, requestUrlPath string) her.Error {
	_, err := e.Dogmas.canBeAccessedBy(nil, e.Method, e.SrcApp+path.Join("/", e.AppFeature, requestUrlPath), &requesterId)
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

	srcApp, appHost := u.Scheme+"://", u.Host
	appFeature, requestedPath := u.Path[0:loc[1]], u.Path[loc[1]:]

	segs := strings.Split(appFeature, "/")
	if len(segs) > 3 {
		appHost = path.Join(appHost, strings.Join(segs[0:len(segs)-2], "/"))
		appFeature = path.Join("/", strings.Join(segs[len(segs)-2:], "/"))
	}
	srcApp += appHost

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
			DstApp:     defaultDestHost(srcApp),
			SrcApp:     srcApp,
			AppFeature: appFeature,
			AppPath:    strings.Join(segs, "/"),
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
		RootUrl: h.DstApp,
		Feature: h.AppFeature,
		Path:    h.RequestedPath,
	}
}

type Route struct {
	Method  string `json:"method"`
	RootUrl string `json:"rootUrl"`
	Feature string `json:"feature"`
	Path    string `json:"path"`
}
