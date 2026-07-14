package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestClassifyAuthTier(t *testing.T) {
	cases := []struct {
		name string
		file authFile
		raw  string
		want string
	}{
		{"free", authFile{AccountType: "free"}, `{}`, tierFree},
		{"super", authFile{}, `{"subscription":{"plan":"SuperGrok"}}`, tierSuper},
		{"heavy", authFile{}, `{"account_tier":"heavy"}`, tierHeavy},
		{"unknown", authFile{}, `{}`, tierUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := classifyAuthTier(tc.file, json.RawMessage(tc.raw)).Tier
			if got != tc.want {
				t.Fatalf("tier=%q want %q", got, tc.want)
			}
		})
	}
}

func TestClassifyAuthTierFromListMetadataWithoutRawJSON(t *testing.T) {
	tests := []struct {
		name string
		file authFile
	}{
		{name: "note", file: authFile{Note: "supergrok"}},
		{name: "label", file: authFile{Label: "Super Grok Account"}},
		{name: "prefix", file: authFile{Prefix: "supergrok"}},
		{name: "tag", file: authFile{Tag: "SuperGrok"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyAuthTier(tt.file, nil)
			if got.Tier != tierSuper {
				t.Fatalf("tier = %q, want %q; sources=%v", got.Tier, tierSuper, got.SourceKeys)
			}
		})
	}
}

func TestClassifyAuthTierOAuthListMetadataDoesNotOverrideSuperSignal(t *testing.T) {
	got := classifyAuthTier(authFile{AccountType: "oauth", Note: "supergrok"}, nil)
	if got.Tier != tierSuper {
		t.Fatalf("tier = %q, want %q; sources=%v", got.Tier, tierSuper, got.SourceKeys)
	}
}

func TestClassifyOfficialSubscriptions(t *testing.T) {
	tests := []struct{ name, body, want string }{
		{"super", `{"subscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK","status":"ACTIVE"}]}`, tierSuper},
		{"heavy", `{"activeSubscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK_HEAVY","status":"ACTIVE"}]}`, tierHeavy},
		{"pro", `{"data":{"subscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK_PRO","status":"ACTIVE"}]}}`, tierHeavy},
		{"inactive", `{"subscriptions":[{"tier":"SUBSCRIPTION_TIER_SUPER_GROK","status":"CANCELED"}]}`, tierFree},
		{"empty", `{"subscriptions":[]}`, tierFree},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := classifyOfficialSubscriptions([]byte(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			if got.Tier != tt.want {
				t.Fatalf("tier=%q want %q", got.Tier, tt.want)
			}
		})
	}
}

func TestExplicitAuthFailureThresholdAndProtection(t *testing.T) {
	old := pluginState
	pluginState = &memoryStore{settings: defaultPluginSettings(), health: map[string]*healthMemory{}}
	defer func() { pluginState = old }()

	file := authFile{AuthIndex: "free-1", Email: "free@example.com", Provider: "xai"}
	eval := healthEvaluation{Health: healthInvalid, ExplicitStatusCode: http.StatusUnauthorized, Reason: "401"}
	for i := 1; i <= 3; i++ {
		rec := updateHealthMemory(file, authClassification{Tier: tierFree}, eval, currentSettings(), testTime(), true, true, nil)
		if i < 3 && rec.DeleteEligible {
			t.Fatalf("eligible at streak %d", i)
		}
		if i == 3 && !rec.DeleteEligible {
			t.Fatal("free account should be eligible at streak 3")
		}
	}

	superFile := authFile{AuthIndex: "super-1", Email: "super@example.com", Provider: "xai"}
	var rec checkRecord
	for i := 0; i < 3; i++ {
		rec = updateHealthMemory(superFile, authClassification{Tier: tierSuper}, eval, currentSettings(), testTime(), true, true, nil)
	}
	if !rec.Protected || rec.DeleteEligible {
		t.Fatal("super account must remain protected")
	}
}

func TestTransientFailureNeverBecomesInvalid(t *testing.T) {
	for _, msg := range []string{"429 rate limited", "503 upstream unavailable", "timeout"} {
		e := evaluateRuntimeHealth(authFile{Status: "error", StatusMessage: msg})
		if e.ExplicitStatusCode == http.StatusUnauthorized || e.ExplicitStatusCode == http.StatusForbidden {
			t.Fatalf("transient %q treated as auth failure", msg)
		}
	}
}

func TestSettingsAlwaysProtectValuableTiers(t *testing.T) {
	s := sanitizeSettings(pluginSettings{InvalidThreshold: 3})
	for _, tier := range []string{tierSuper, tierHeavy, tierUnknown} {
		if !isProtectedTier(tier, s) {
			t.Fatalf("%s must be protected", tier)
		}
	}
}

func TestExtractAuthCredentialsIgnoresUnsafeEndpointAndHeaders(t *testing.T) {
	raw := `{
		"access_token": "test-access-token",
		"base_url": "http://127.0.0.1:8080/v1",
		"headers": {
			"User-Agent": "grok-pager/0.2.93",
			"X-XAI-Token-Auth": "xai-grok-cli",
			"Cookie": "secret-cookie",
			"Proxy-Authorization": "secret-proxy"
		}
	}`
	creds := extractAuthCredentials(json.RawMessage(raw))
	if creds.AccessToken != "test-access-token" {
		t.Fatalf("access_token=%q", creds.AccessToken)
	}
	if creds.Headers["X-XAI-Token-Auth"] != "xai-grok-cli" {
		t.Fatalf("allowed headers=%v", creds.Headers)
	}
	if _, ok := creds.Headers["Cookie"]; ok {
		t.Fatalf("Cookie must not be forwarded: %v", creds.Headers)
	}
	if _, ok := creds.Headers["Proxy-Authorization"]; ok {
		t.Fatalf("Proxy-Authorization must not be forwarded: %v", creds.Headers)
	}
	if officialResponsesEndpoint != "https://cli-chat-proxy.grok.com/v1/responses" {
		t.Fatalf("unexpected fixed endpoint %q", officialResponsesEndpoint)
	}
}

func TestProbeResponsesRejectsNonOfficialEndpoint(t *testing.T) {
	status, _, err := probeResponsesAPI("http://127.0.0.1:8080/responses", authCredentials{AccessToken: "secret"}, "grok-4.5", "ping")
	if err == nil || status != 0 {
		t.Fatalf("unsafe endpoint should be rejected before network: status=%d err=%v", status, err)
	}
}

func TestResponsesProbeDoesNotMutateHealthOrDeletionState(t *testing.T) {
	resetPluginStateForTests()
	file := authFile{AuthIndex: "auth-1", Email: "safe@example.com", Status: "active"}
	settings := defaultPluginSettings()
	now := testTime()
	initial := updateHealthMemory(file, authClassification{Tier: tierFree}, healthEvaluation{Health: healthInvalid, Reason: "explicit_401", ExplicitStatusCode: 401}, settings, now, true, true, nil)
	if initial.InvalidStreak != 1 {
		t.Fatalf("initial streak=%d", initial.InvalidStreak)
	}
	for i := 0; i < 5; i++ {
		_ = evaluateResponsesProbe(file, 403, nil)
	}
	after := snapshotHealthForFile(file, settings)
	if after.InvalidStreak != 1 || after.DeleteEligible {
		t.Fatalf("probe changed deletion state: %+v", after)
	}
}

func TestEvaluateResponsesProbe(t *testing.T) {
	file := authFile{Status: "active"}
	if e := evaluateResponsesProbe(file, 200, nil); e.Health != healthHealthy {
		t.Fatalf("200 => %q", e.Health)
	}
	if e := evaluateResponsesProbe(file, 401, nil); e.Health != healthInvalid || e.ExplicitStatusCode != 401 {
		t.Fatalf("401 => %+v", e)
	}
	if e := evaluateResponsesProbe(file, 429, nil); e.Health != healthUnavailable {
		t.Fatalf("429 => %q", e.Health)
	}
	if e := evaluateResponsesProbe(file, 0, errString("timeout")); e.Health != healthUnavailable {
		t.Fatalf("timeout => %q", e.Health)
	}
}

type errString string

func (e errString) Error() string { return string(e) }

func TestExtractResponsesOutputText(t *testing.T) {
	raw := `{"id":"resp_1","status":"completed","output_text":"Four"}`
	if got := extractResponsesOutputText([]byte(raw)); got != "Four" {
		t.Fatalf("got %q", got)
	}
	raw2 := `{"output":[{"type":"message","content":[{"type":"output_text","text":"4"}]}]}`
	if got := extractResponsesOutputText([]byte(raw2)); got != "4" {
		t.Fatalf("nested got %q", got)
	}
}

func testTime() (v time.Time) { return time.Unix(1700000000, 0).UTC() }
