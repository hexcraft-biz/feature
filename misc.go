package feature

import (
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/publicsuffix"
)

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
func defaultDestHostByString(appRootUrlString string) (string, error) {
	u, err := url.ParseRequestURI(appRootUrlString)
	if err != nil {
		return "", err
	}

	return defaultDestHostByUrl(u), nil
}

func defaultDestHostByUrl(appRootUrl *url.URL) string {
	host := ""
	if appRootUrl.Path != "" {
		segs := strings.Split(appRootUrl.Path, "/")
		host = segs[len(segs)-1]
	} else {

		hostname := appRootUrl.Hostname()

		eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(hostname)
		if err != nil {
			panic(err)
		}

		parts := strings.SplitN(eTLDPlusOne, ".", 2)
		if len(parts) < 2 {
			panic("invalid domain structure")
		}

		hostParts := strings.Split(hostname, ".")
		subParts := hostParts[:len(hostParts)-len(parts)]
		host = subParts[len(subParts)-1]
	}

	return "http://" + host
}
