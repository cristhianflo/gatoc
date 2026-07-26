package embedfixer

import "testing"

func TestPlatformByHostAliases(t *testing.T) {
	tests := []struct {
		host           string
		expectedKey    string
		expectedDomain string
	}{
		{host: "twitter.com", expectedKey: "twitter", expectedDomain: "fxtwitter.com"},
		{host: "x.com", expectedKey: "twitter", expectedDomain: "fxtwitter.com"},
		{host: "reddit.com", expectedKey: "reddit", expectedDomain: "vxreddit.com"},
	}

	for _, tc := range tests {
		platform, ok := platformByHost(tc.host)
		if !ok {
			t.Fatalf("expected host %s to resolve a platform", tc.host)
		}

		if platform.Key != tc.expectedKey {
			t.Fatalf("expected key %s, got %s", tc.expectedKey, platform.Key)
		}

		if platform.DefaultDomain != tc.expectedDomain {
			t.Fatalf("expected domain %s, got %s", tc.expectedDomain, platform.DefaultDomain)
		}
	}
}

func TestActiveDomainResolution(t *testing.T) {
	platform, ok := platformByKey("twitter")
	if !ok {
		t.Fatal("expected twitter platform")
	}

	custom, customMode := activeDomain(platform, "customfx.example", true)
	if !customMode {
		t.Fatal("expected custom mode")
	}
	if custom != "customfx.example" {
		t.Fatalf("expected custom domain, got %s", custom)
	}

	defaultDomain, defaultMode := activeDomain(platform, "", false)
	if defaultMode {
		t.Fatal("expected default mode")
	}
	if defaultDomain != platform.DefaultDomain {
		t.Fatalf("expected default domain %s, got %s", platform.DefaultDomain, defaultDomain)
	}
}

func TestNormalizeDomainHost(t *testing.T) {
	domain, err := normalizeDomainHost("HTTPS://WWW.FxTwitter.com")
	if err != nil {
		t.Fatalf("expected valid domain, got error %v", err)
	}

	if domain != "fxtwitter.com" {
		t.Fatalf("expected fxtwitter.com, got %s", domain)
	}
}

func TestNormalizeDomainHostInvalidInput(t *testing.T) {
	invalid := []string{
		"",
		"fxtwitter.com/path",
		"fxtwitter.com?query=1",
		"no-dot-domain",
	}

	for _, candidate := range invalid {
		if _, err := normalizeDomainHost(candidate); err == nil {
			t.Fatalf("expected invalid domain error for %q", candidate)
		}
	}
}

func TestPlatformByKeyInvalid(t *testing.T) {
	if _, ok := platformByKey("tiktok"); ok {
		t.Fatal("expected unsupported platform")
	}
}
