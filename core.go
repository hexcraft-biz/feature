package feature

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
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

const (
	OwnershipOrganization = "ORGANIZATION"
	OwnershipPrivate      = "PRIVATE"
	OwnershipPublic       = "PUBLIC"
)

const (
	EnumOwnershipUndef int = iota
	EnumOwnershipOrganization
	EnumOwnershipPrivate
	EnumOwnershipPublic
)

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

func (f *Feature) addEndpoint(enumOwnership int, method, relativePath string, scopes []string, ownerParamIndex int) *EndpointHandler {
	segs := strings.Split(path.Join("/", relativePath), "/")
	for i := range segs {
		if strings.HasPrefix(segs[i], ":") {
			segs[i] = "*"
		}
	}

	relativePath = strings.Join(segs, "/")
	u, err := url.Parse(f.Dogmas.AppHost + path.Join("/", f.FeaturePath, relativePath))
	if err != nil {
		panic("Invalid endpoint")
	}

	ownership := ""
	switch enumOwnership {
	case EnumOwnershipOrganization:
		ownership = OwnershipOrganization
		ownerParamIndex = 0

	case EnumOwnershipPrivate:
		ownership = OwnershipPrivate
		if ownerParamIndex > 0 {
			if ownerParamIndex%2 != 0 {
				panic("Invalid endpoint")
			}
		}

	case EnumOwnershipPublic:
		ownership = OwnershipPublic
		ownerParamIndex = 0

	default:
		panic("Invalid endpoint")
	}

	e := &Endpoint{
		EndpointId:      Md5Identifier(fmt.Sprintf("%x", md5.Sum([]byte(method+u.String())))),
		Ownership:       ownership,
		Method:          method,
		UrlHost:         &f.Dogmas.AppHost,
		UrlFeature:      &f.FeaturePath,
		UrlPath:         relativePath,
		OwnerParamIndex: ownerParamIndex,
	}

	if len(scopes) > 0 {
		for _, identifier := range scopes {
			f.Dogmas.endpointsContainer(identifier).addEndpoint(e)
		}
	} else {
		f.Dogmas.endpointsContainer("").addEndpoint(e)
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
	HostUrl *url.URL
	AppHost string
	*scopesHandler
}

func NewDogmas(appHostUrl *url.URL) (*Dogmas, error) {
	u, err := url.ParseRequestURI(os.Getenv("HOST_DOGMAS"))
	if err != nil {
		return nil, err
	}

	return &Dogmas{
		init:          flag.Bool(FlagInit, false, FlagInitDescription),
		HostUrl:       u,
		AppHost:       appHostUrl.String(),
		scopesHandler: newScopesHandler(u),
	}, nil
}

func (d Dogmas) IsInit() bool {
	return *d.init
}

func (d Dogmas) Register() {
	if !*d.init {
		panic("not init mode")
	}
	d.scopesHandler.register()
}

type ResultAccessPermission struct {
	CanAccess bool `json:"canAccess"`
}

// For api-proxy to check
func (d Dogmas) CanAccess(scopes []string, method, endpointUrl string, requesterId *xuuid.UUID) (bool, her.Error) {
	if len(scopes) < 1 {
		return false, her.ErrForbidden
	}

	data := map[string]any{
		"scopes":             scopes,
		"method":             method,
		"requestEndpointUrl": endpointUrl,
	}
	if requesterId != nil {
		data["requesterId"] = requesterId.String()
	}

	jsonbytes, err := json.Marshal(data)
	if err != nil {
		return false, her.NewError(http.StatusInternalServerError, err, nil)
	}

	req, err := http.NewRequest("POST", d.HostUrl.JoinPath("/permissions/v1/proxy").String(), bytes.NewReader(jsonbytes))
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

// ================================================================
//
// ================================================================
func ToEndpointIdWithPath(method, endpointUrl string) (Md5Identifier, SubsetString, error) {
	u, err := url.ParseRequestURI(endpointUrl)
	if err != nil {
		return "", "", errors.New("Invaild enpoint url")
	}

	urlFeature, urlPath, err := extractPathSegments(u.Path)
	if err != nil {
		return "", "", errors.New("Invaild enpoint url")
	}

	pathSegs := strings.Split(urlPath, "/")
	for i := range pathSegs {
		if i > 0 && i%2 == 0 {
			pathSegs[i] = "*"
		}
	}

	urlstring := u.Scheme + "://" + path.Join(u.Host, urlFeature, strings.Join(pathSegs, "/"))

	md5bytes := md5.Sum([]byte(method + urlstring))
	return Md5Identifier(fmt.Sprintf("%x", md5bytes)), SubsetString(urlPath), nil
}

func ToEndpointId(method, endpointUrl string) (Md5Identifier, error) {
	endpointId, _, err := ToEndpointIdWithPath(method, endpointUrl)
	return endpointId, err
}

func extractPathSegments(s string) (string, string, error) {
	re := regexp.MustCompile(`/v[0-9]+`)
	loc := re.FindStringIndex(s)

	if loc == nil {
		return "", "", errors.New("Invalid endpoint")
	}

	return s[0:loc[1]], s[loc[1]:], nil
}
