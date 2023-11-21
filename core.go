package feature

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/her"
)

// ================================================================
//
// ================================================================
const (
	ByAuthorityOfOrganization = "ORGANIZATION"
	ByAuthorityOfDataOwner    = "DATA_OWNER"
)

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

func (f *Feature) addEndpoint(byAuthorityOf, method, relativePath string) *Endpoint {
	urlPath := GetAuthorizedEndpointPath(relativePath)
	u, err := url.Parse(path.Join(f.Dogmas.AppHost, f.FeaturePath, urlPath))
	if err != nil {
		panic("Invalid endpoint")
	}
	e := &Endpoint{
		Dogmas:             f.Dogmas,
		EndpointIdentifier: fmt.Sprintf("%x", md5.Sum([]byte(method+u.String()))),
		ByAuthorityOf:      byAuthorityOf,
		Method:             method,
		UrlHost:            &f.Dogmas.AppHost,
		UrlFeature:         &f.FeaturePath,
		UrlPath:            urlPath,
	}

	f.Dogmas.Endpoints = append(f.Dogmas.Endpoints, e)
	return e
}

func EndpointIdentifier(method, appHost, urlFeature, urlPath string) string {
	urlPath = GetAuthorizedEndpointPath(urlPath)
	u, err := url.Parse(path.Join(appHost, urlFeature, urlPath))
	if err != nil {
		panic("Invalid endpoint")
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(method+u.String())))
}

type HandlerFunc func(*Endpoint) gin.HandlerFunc

func handlerFuncs(e *Endpoint, handlers []HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		funcs[i] = h(e)
	}
	return funcs
}

type Endpoint struct {
	*Dogmas            `json:"-"`
	EndpointIdentifier string  `json:"endpointIdentifier"`
	ByAuthorityOf      string  `json:"byAuthorityOf"`
	Method             string  `json:"method"`
	UrlHost            *string `json:"urlHost"`
	UrlFeature         *string `json:"urlFeature"`
	UrlPath            string  `json:"urlPath"`
}

// ================================================================
//
// ================================================================
type Dogmas struct {
	HostUrl   *url.URL
	AppHost   string
	Endpoints []*Endpoint
}

func NewDogmas(appHostUrl *url.URL) (*Dogmas, error) {
	u, err := url.ParseRequestURI(os.Getenv("HOST_DOGMAS"))
	if err != nil {
		return nil, err
	}

	return &Dogmas{
		HostUrl: u,
		AppHost: appHostUrl.String(),
	}, nil
}

func GetAuthorizedEndpointPath(relativePath string) string {
	u, _ := url.Parse(path.Join("/", relativePath))
	segs := strings.Split(u.Path, "/")
	for i := range segs {
		if strings.HasPrefix(segs[i], ":") {
			segs[i] = "*"
		}
	}

	return strings.Join(segs, "/")
}

func GetAuthorizedURIPath(relativePath string, params ...string) string {
	u, _ := url.Parse(path.Join("/", relativePath))
	segs := strings.Split(u.Path, "/")
	for i := range segs {
		if len(params) == 0 {
			panic("args not matched")
		}

		if strings.HasPrefix(segs[i], ":") {
			segs[i] = params[0]
			params = params[1:]
		}
	}

	return strings.Join(segs, "/")
}

func (d Dogmas) RegisterEndpoints() her.Error {
	if len(d.Endpoints) > 0 {
		jsonbytes, err := json.Marshal(d.Endpoints)
		if err != nil {
			return her.NewError(http.StatusInternalServerError, err, nil)
		}

		return apiPost(d.HostUrl.JoinPath("/resources/v1/endpoints"), jsonbytes)
	}

	return nil
}

// ================================================================
func apiPost(apiUrl *url.URL, jsonbytes []byte) her.Error {
	req, err := http.NewRequest("POST", apiUrl.String(), bytes.NewReader(jsonbytes))
	if err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	}

	payload := her.NewPayload(nil)
	client := &http.Client{}

	if resp, err := client.Do(req); err != nil {
		return her.NewError(http.StatusInternalServerError, err, nil)
	} else if err := her.FetchHexcApiResult(resp, payload); err != nil {
		return err
	} else if resp.StatusCode != 201 {
		return her.NewErrorWithMessage(http.StatusInternalServerError, "Dogmas: "+payload.Message, nil)
	}

	return nil
}

// ================================================================
func RemoveRedundant(scopes []string) []string {
	var patterns []*regexp.Regexp
	result := []string{}
	resultMap := map[string]bool{}

	for _, key := range scopes {
		if strings.Contains(key, "*") {
			pattern := strings.ReplaceAll(key, "*", ".*")
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				fmt.Println("Regex compile error:", err)
				continue
			}
			patterns = append(patterns, re)
		}
	}

	for _, key := range scopes {
		if !isCoveredByMoreGeneralPattern(key, patterns) && !resultMap[key] {
			result = append(result, key)
			resultMap[key] = true
		}
	}

	return result
}

func isCoveredByMoreGeneralPattern(key string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(key) && pattern.String() != "^"+strings.ReplaceAll(key, "*", ".*")+"$" {
			return true
		}
	}
	return false
}

// ================================================================
func IsCovered(scope string, scopes []string) bool {
	var patterns []*regexp.Regexp

	for _, key := range scopes {
		if strings.Contains(key, "*") {
			pattern := strings.ReplaceAll(key, "*", ".*")
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				fmt.Println("Regex compile error:", err)
				continue
			}
			patterns = append(patterns, re)
		}
	}

	return isCoveredBy(scope, patterns)
}

func isCoveredBy(key string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(key) {
			return true
		}
	}
	return false
}
