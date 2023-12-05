package feature

import (
	"net/url"
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func removeRedundant(rules []string) []string {
	var patterns []*regexp.Regexp
	result := []string{}
	resultMap := map[string]bool{}

	for _, key := range rules {
		if strings.Contains(key, "*") {
			pattern := strings.ReplaceAll(key, "*", ".*")
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				continue
			}
			patterns = append(patterns, re)
		}
	}

	for _, key := range rules {
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
func isCovered(rule string, rules []string) bool {
	if rules == nil {
		return false
	}

	var patterns []*regexp.Regexp

	for _, key := range rules {
		if strings.Contains(key, "*") {
			pattern := strings.ReplaceAll(key, "*", ".*")
			if re, err := regexp.Compile("^" + pattern + "$"); err == nil {
				patterns = append(patterns, re)
			}
		} else if rule == key {
			return true
		}
	}

	return isCoveredBy(rule, patterns)
}

func isCoveredBy(key string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(key) {
			return true
		}
	}
	return false
}

// ================================================================
func standardizePath(relativePath string) string {
	segs := strings.Split(path.Join("/", relativePath), "/")
	for i := range segs {
		if strings.HasPrefix(segs[i], ":") {
			segs[i] = "*"
		}
	}

	return strings.Join(segs, "/")
}

// ================================================================
func defaultDestHost(inputURL string) string {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		panic(err)
	}

	hostname := parsedURL.Hostname()

	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(eTLDPlusOne, ".", 2)
	if len(parts) < 2 {
		panic("invalid domain structure")
	}

	subParts := strings.Split(parts[0], ".")
	secondLevelDomain := subParts[len(subParts)-1]

	return "http://" + secondLevelDomain
}
