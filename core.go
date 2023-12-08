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
	FeaturePath string
	*gin.RouterGroup
	*Dogmas
}

func New(e *gin.Engine, featurePath string, d *Dogmas) *Feature {
	return &Feature{
		FeaturePath: featurePath,
		RouterGroup: e.Group(featurePath),
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

func (f *Feature) addEndpoint(ownership, method, relativePath string, scopes []string) *EndpointHandler {
	e := newEndpoint(ownership, method, f.Dogmas.AppRootUrl.String(), f.FeaturePath, standardizePath(relativePath))

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
	init       *bool
	RootUrl    *url.URL
	AppRootUrl *url.URL
	ScopesHandler
}

func NewDogmas(appRootUrl *url.URL) (*Dogmas, error) {
	u, err := url.ParseRequestURI(os.Getenv("APP_DOGMAS"))
	if err != nil {
		return nil, err
	}

	return &Dogmas{
		init:          flag.Bool(FlagInit, false, FlagInitDescription),
		RootUrl:       u,
		AppRootUrl:    appRootUrl,
		ScopesHandler: newScopesHandler(u),
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
	return d.canBeAccessedBy(strings.Split(scope, " "), method, endpointUrl, requesterId)
}

func (d Dogmas) canBeAccessedBy(scopes []string, method, endpointUrl string, requesterId *xuuid.UUID) (*Route, her.Error) {
	u := d.RootUrl.JoinPath("/routes/v1/endpoints")
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
