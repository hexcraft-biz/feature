package feature

import (
	"flag"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

// ================================================================
//
// ================================================================
type HandlerFunc func(*EndpointHandler) gin.HandlerFunc

func handlerFuncs(e *EndpointHandler, handlers []HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = h(e)
	}
	return funcs
}

// ================================================================
//
// ================================================================
type Feature struct {
	*gin.RouterGroup
	AppRootUrl  *url.URL
	FeaturePath string
	*Dogmas
}

func New(e *gin.Engine, appRootUrl *url.URL, featurePath string, d *Dogmas) *Feature {
	var group *gin.RouterGroup
	if e != nil {
		group = e.Group(featurePath)
	}

	return &Feature{
		RouterGroup: group,
		AppRootUrl:  appRootUrl,
		FeaturePath: featurePath,
		Dogmas:      d,
	}
}

func (f *Feature) OrganizationAssets() *OrganizationAssets {
	return &OrganizationAssets{
		Feature: f,
	}
}

func (f *Feature) PrivateAssets() *PrivateAssets {
	return &PrivateAssets{
		Feature: f,
	}
}

func (f *Feature) PublicAssets() *PublicAssets {
	return &PublicAssets{
		Feature: f,
	}
}

// ================================================================
func (f *Feature) addEndpoint(ownership, method, relativePath string, scopes []string) *EndpointHandler {
	e := newEndpoint(ownership, method, f.AppRootUrl.String(), f.FeaturePath, standardizePath(relativePath))

	if len(scopes) > 0 {
		for _, identifier := range scopes {
			f.Dogmas.Scope(identifier).AddEndpoint(e)
		}
	} else {
		f.Dogmas.Scope("").AddEndpoint(e)
	}

	return &EndpointHandler{
		Dogmas:   f.Dogmas,
		Endpoint: e,
	}
}

// ================================================================
//
// ================================================================
const (
	FlagInit            = "initdogmas"
	FlagInitDescription = "To register scopes and endpoints on dogmas"
)

type Dogmas struct {
	init             *bool
	DogmasRootUrl    *url.URL
	DoctrinesRootUrl *url.URL
	CreedsRootUrl    *url.URL
	ScopesHandler
}

func NewDogmas() (*Dogmas, error) {
	uDogmas, err := url.ParseRequestURI(os.Getenv("APP_DOGMAS"))
	if err != nil {
		return nil, err
	}

	uDoctrines, err := url.ParseRequestURI(os.Getenv("APP_DOCTRINES"))
	if err != nil {
		return nil, err
	}

	uCreeds, err := url.ParseRequestURI(os.Getenv("APP_CREEDS"))
	if err != nil {
		return nil, err
	}

	return &Dogmas{
		init:             flag.Bool(FlagInit, false, FlagInitDescription),
		DogmasRootUrl:    uDogmas,
		DoctrinesRootUrl: uDoctrines,
		CreedsRootUrl:    uCreeds,
		ScopesHandler:    newScopesHandler(uDogmas),
	}, nil
}

func (d Dogmas) IsInit() bool {
	return *d.init
}

func (d Dogmas) Register() {
	if !*d.init {
		panic("not init mode")
	}

	if err := d.ScopesHandler.register(); err != nil {
		panic(err)
	}
}

// For api-proxy to check
func (d Dogmas) CanAccess(scope, method, endpointUrl string, requesterId *xuuid.UUID) (*Route, her.Error) {
	if scope == "" {
		return nil, her.ErrForbidden
	}
	return canBeAccessedBy(d.DoctrinesRootUrl, strings.Split(scope, " "), method, endpointUrl, requesterId)
}

func canBeAccessedBy(rootUrl *url.URL, scopes []string, method, endpointUrl string, requesterId *xuuid.UUID) (*Route, her.Error) {
	u := rootUrl.JoinPath("/routes/v1/endpoints")
	q := u.Query()
	if scopes != nil {
		q.Set("scopes", strings.Join(scopes, " "))
	}
	q.Set("method", method)
	q.Set("url", endpointUrl)
	if requesterId != nil {
		q.Set("requester", requesterId.String())
	}
	u.RawQuery = q.Encode()

	result := new(Route)
	payload := her.NewPayload(result)

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, her.NewError(http.StatusInternalServerError, err, nil)
	} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return result, nil
	default:
		return nil, her.NewErrorWithMessage(resp.StatusCode, payload.Message, nil)
	}
}
