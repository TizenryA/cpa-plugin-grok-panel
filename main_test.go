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

func testTime() (v time.Time) { return time.Unix(1700000000, 0).UTC() }
