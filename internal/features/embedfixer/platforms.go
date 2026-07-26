package embedfixer

import (
	"slices"
	"strings"
)

type platformConfig struct {
	Key           string
	DisplayName   string
	SourceHosts   []string
	DefaultDomain string
	LinkLabel     string
}

var supportedPlatforms = []platformConfig{
	{
		Key:           "twitter",
		DisplayName:   "Twitter/X",
		SourceHosts:   []string{"twitter.com", "x.com"},
		DefaultDomain: "fxtwitter.com",
		LinkLabel:     "Tweet",
	},
	{
		Key:           "reddit",
		DisplayName:   "Reddit",
		SourceHosts:   []string{"reddit.com"},
		DefaultDomain: "vxreddit.com",
		LinkLabel:     "Reddit",
	},
	{
		Key:           "instagram",
		DisplayName:   "Instagram",
		SourceHosts:   []string{"instagram.com"},
		DefaultDomain: "kkinstagram.com",
		LinkLabel:     "Instagram",
	},
}

func platformByKey(key string) (platformConfig, bool) {
	normalizedKey := strings.ToLower(strings.TrimSpace(key))
	for _, platform := range supportedPlatforms {
		if platform.Key == normalizedKey {
			return platform, true
		}
	}

	return platformConfig{}, false
}

func platformByHost(host string) (platformConfig, bool) {
	normalizedHost := strings.ToLower(strings.TrimSpace(host))
	normalizedHost = strings.TrimPrefix(normalizedHost, "www.")

	for _, platform := range supportedPlatforms {
		if slices.Contains(platform.SourceHosts, normalizedHost) {
			return platform, true
		}
	}

	return platformConfig{}, false
}
