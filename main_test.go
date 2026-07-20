package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestManagementRegistrationKeepsDynamicDataBehindManagementAuth(t *testing.T) {
	body, err := handleMethod("management.register", nil)
	if err != nil {
		t.Fatal(err)
	}
	var env struct {
		OK     bool                   `json:"ok"`
		Result managementRegistration `json:"result"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatal(err)
	}
	if !env.OK {
		t.Fatalf("registration envelope is not ok: %s", body)
	}
	dataRoute := managementBasePath + "/data"
	foundDataRoute := false
	for _, route := range env.Result.Routes {
		if route.Path == dataRoute && route.Method == http.MethodGet {
			foundDataRoute = true
		}
	}
	if !foundDataRoute {
		t.Fatalf("authenticated data route %q is missing", dataRoute)
	}
	for _, resource := range env.Result.Resources {
		if resource.Path == resourcePanelDataPath {
			t.Fatalf("dynamic data must not be registered as public resource: %+v", resource)
		}
	}
}

func TestPanelFetchesDataThroughManagementRoute(t *testing.T) {
	if !strings.Contains(htmlPage, "managementPluginGet('data')") {
		t.Fatal("panel must load account data through the authenticated management route")
	}
	if strings.Contains(htmlPage, "fetch(fixedApiUrl('data')") {
		t.Fatal("panel must not load dynamic account data from a public resource path")
	}
}

func TestLegacyPublicDataPathCannotDispatchDynamicData(t *testing.T) {
	oldHostCaller := hostCaller
	hostCaller = func(method string, payload any) (json.RawMessage, error) {
		t.Fatalf("legacy public data path reached privileged host callback %q", method)
		return nil, nil
	}
	defer func() { hostCaller = oldHostCaller }()

	body, err := handleManagement([]byte(`{"method":"GET","path":"/panel/data"}`))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), `"files"`) {
		t.Fatalf("legacy public path returned dynamic account data: %s", body)
	}
}

func TestManagementDataRouteAcceptsRegisteredPathForms(t *testing.T) {
	for _, path := range []string{managementBasePath + "/data", "/v0/management" + managementBasePath + "/data"} {
		if !managementRouteMatches(path, managementBasePath+"/data") {
			t.Fatalf("registered management data path did not match: %q", path)
		}
	}
	if managementRouteMatches(resourcePanelDataPath, managementBasePath+"/data") {
		t.Fatal("legacy public resource path must not match the management data route")
	}
}

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
	status, _, _, err := probeResponsesAPI("http://127.0.0.1:8080/responses", authCredentials{AccessToken: "secret"}, "grok-4.5", "ping")
	if err == nil || status != 0 {
		t.Fatalf("unsafe endpoint should be rejected before network: status=%d err=%v", status, err)
	}
}

func TestParseResponsesRateLimits(t *testing.T) {
	header := http.Header{
		"X-Ratelimit-Limit-Tokens":       []string{"2000000"},
		"X-Ratelimit-Remaining-Tokens":   []string{"1999900"},
		"X-Ratelimit-Limit-Requests":     []string{"21"},
		"X-Ratelimit-Remaining-Requests": []string{"20"},
	}
	got := parseResponsesRateLimits(header)
	if got.TokenLimit != 2000000 || got.TokenRemaining != 1999900 || got.TokenUsed != 100 {
		t.Fatalf("token limits=%+v", got)
	}
	if got.RequestLimit != 21 || got.RequestRemaining != 20 || got.RequestUsed != 1 {
		t.Fatalf("request limits=%+v", got)
	}
	if got.Source != "responses_headers" || !got.Available || !got.TokenAvailable || !got.RequestAvailable {
		t.Fatalf("source/available=%+v", got)
	}
}

func TestParseResponsesRateLimitsMissingAndInvalid(t *testing.T) {
	got := parseResponsesRateLimits(http.Header{"X-Ratelimit-Limit-Tokens": []string{"not-a-number"}})
	if got.Available || got.TokenLimit != 0 || got.Source != "" {
		t.Fatalf("invalid headers must remain unavailable: %+v", got)
	}

	requestOnly := parseResponsesRateLimits(http.Header{
		"X-Ratelimit-Limit-Requests":     []string{"21"},
		"X-Ratelimit-Remaining-Requests": []string{"20"},
	})
	if !requestOnly.Available || requestOnly.TokenAvailable || !requestOnly.RequestAvailable {
		t.Fatalf("request-only headers=%+v", requestOnly)
	}

	partialToken := parseResponsesRateLimits(http.Header{"X-Ratelimit-Limit-Tokens": []string{"2000000"}})
	if partialToken.Available || partialToken.TokenAvailable {
		t.Fatalf("partial token pair must not be advertised: %+v", partialToken)
	}
}

func TestCacheResponsesRateLimitsKeepsDeletionStateIsolated(t *testing.T) {
	resetPluginStateForTests()
	file := authFile{AuthIndex: "auth-limit", Email: "limit@example.com", Status: "active"}
	settings := defaultPluginSettings()
	now := testTime()
	initial := updateHealthMemory(file, authClassification{Tier: tierFree}, healthEvaluation{Health: healthInvalid, Reason: "explicit_401", ExplicitStatusCode: 401}, settings, now, true, true, nil)
	cacheResponsesRateLimits(file, responsesRateLimits{Available: true, Source: "responses_headers", TokenLimit: 2000000, TokenRemaining: 1500000, TokenUsed: 500000, RequestLimit: 21, RequestRemaining: 20, RequestUsed: 1, MeasuredAt: now})
	after := snapshotHealthForFile(file, settings)
	if after.InvalidStreak != initial.InvalidStreak || after.DeleteEligible != initial.DeleteEligible {
		t.Fatalf("rate limit cache changed deletion state: before=%+v after=%+v", initial, after)
	}
	limits := snapshotResponsesRateLimits(file)
	if !limits.Available || limits.TokenRemaining != 1500000 || limits.RequestRemaining != 20 {
		t.Fatalf("cached limits=%+v", limits)
	}
}

func TestCacheResponsesRateLimitsDoesNotCreateHealthMemory(t *testing.T) {
	resetPluginStateForTests()
	file := authFile{AuthIndex: "auth-limit-only", Email: "limit-only@example.com", Status: "active"}
	cacheResponsesRateLimits(file, responsesRateLimits{
		Available: true, TokenAvailable: true, Source: "responses_headers",
		TokenLimit: 2000000, TokenRemaining: 1999900, TokenUsed: 100,
	})
	pluginState.mu.Lock()
	_, hasHealth := pluginState.health[authMemoryKey(file)]
	_, hasLimits := pluginState.limits[authMemoryKey(file)]
	pluginState.mu.Unlock()
	if hasHealth || !hasLimits {
		t.Fatalf("limit cache must be isolated: hasHealth=%v hasLimits=%v", hasHealth, hasLimits)
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
