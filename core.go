package feature

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
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
type HandlerFunc func(*Endpoint) gin.HandlerFunc

func handlerFuncs(e *Endpoint, handlers []HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = h(e)
	}
	return funcs
}

const (
	ByAuthorityOfOrganization = "ORGANIZATION"
	ByAuthorityOfDataOwner    = "DATA_OWNER"
	ByAuthorityOfNone         = "NONE"
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

func (f *Feature) addEndpoint(byAuthorityOf, method, relativePath string, scopes []string) *Endpoint {
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

	e := &Endpoint{
		Dogmas:        f.Dogmas,
		EndpointId:    Md5Identifier(fmt.Sprintf("%x", md5.Sum([]byte(method+u.String())))),
		ByAuthorityOf: byAuthorityOf,
		Method:        method,
		UrlHost:       &f.Dogmas.AppHost,
		UrlFeature:    &f.FeaturePath,
		UrlPath:       relativePath,
	}

	for _, identifier := range scopes {
		f.Dogmas.endpointsContainer(identifier).addEndpoint(e)
	}

	return e
}

// ================================================================
//
// ================================================================
type Dogmas struct {
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
		HostUrl:       u,
		AppHost:       appHostUrl.String(),
		scopesHandler: newScopesHandler(u),
	}, nil
}

func (d Dogmas) CanAccess(scopes []string, method, endpointUrl string, userId *xuuid.UUID) her.Error {
	if len(scopes) < 1 {
		return her.ErrForbidden
	}

	data := map[string]any{
		"scopes":             scopes,
		"method":             method,
		"requestEndpointUrl": endpointUrl,
	}
	if userId != nil {
		data["userId"] = userId.String()
	}

	jsonbytes, err := json.Marshal(data)
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	req, err := http.NewRequest("POST", d.HostUrl.JoinPath("/permissions/v1/proxy").String(), bytes.NewReader(jsonbytes))
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	payload := her.NewPayload(nil)
	client := &http.Client{}

	if resp, err := client.Do(req); err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return her.NewErrorWithMessage(http.StatusForbidden, "Dogmas: "+payload.Message, nil)
	}

	return nil
}

// ================================================================
//
// ================================================================
func ToEndpointIdWithPath(method, endpointUrl string) (Md5Identifier, string, error) {
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
	return Md5Identifier(fmt.Sprintf("%x", md5bytes)), urlPath, nil
}

func extractPathSegments(s string) (string, string, error) {
	re := regexp.MustCompile(`/v[0-9]+`)
	loc := re.FindStringIndex(s)

	if loc == nil {
		return "", "", fmt.Errorf("Invalid endpoint")
	}

	return s[0:loc[1]], s[loc[1]:], nil
}
