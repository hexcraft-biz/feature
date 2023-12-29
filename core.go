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
	DstApp      string
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
		DstApp:      defaultDestHostByUrl(appRootUrl),
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
	e := newEndpoint(ownership, method, f.DstApp, f.AppRootUrl.String(), f.FeaturePath, standardizePath(relativePath))

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
	init    *bool
	RootUrl *url.URL
	ScopesHandler
}

func NewDogmas() (*Dogmas, error) {
	uDogmas, err := url.ParseRequestURI(os.Getenv("APP_DOGMAS"))
	if err != nil {
		return nil, err
	}

	return &Dogmas{
		init:          flag.Bool(FlagInit, false, FlagInitDescription),
		RootUrl:       uDogmas,
		ScopesHandler: newScopesHandler(uDogmas),
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

// For proxy to check
func (d Dogmas) CanAccess(scope, method, endpointUrl string, requesterId *xuuid.UUID) (*Route, her.Error) {
	scopeSlice := []string{}
	for _, s := range strings.Split(scope, ScopesDelimiter) {
		if s != "" {
			scopeSlice = append(scopeSlice, s)
		}
	}

	if len(scopeSlice) <= 0 {
		return nil, her.NewErrorWithMessage(http.StatusForbidden, "Invalid scope", nil)
	}

	apiUrl := d.RootUrl.JoinPath("/access/v1/from-proxy")
	return canBeAccessedBy(apiUrl, strings.Split(scope, ScopesDelimiter), method, endpointUrl, requesterId)
}

func canBeAccessedBy(apiUrl *url.URL, scopes []string, method, endpointUrl string, requesterId *xuuid.UUID) (*Route, her.Error) {
	q := apiUrl.Query()
	if scopes != nil {
		q.Set("scopes", strings.Join(scopes, ScopesDelimiter))
	}
	q.Set("method", method)
	q.Set("url", endpointUrl)
	if requesterId != nil {
		q.Set("requester", requesterId.String())
	}
	apiUrl.RawQuery = q.Encode()

	result := new(Route)
	payload := her.NewPayload(result)

	resp, err := http.Get(apiUrl.String())
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
