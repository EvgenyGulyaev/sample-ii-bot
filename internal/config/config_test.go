package config

import "testing"

func TestLoadDefaultsToKimi(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_MODEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.LLMBaseURL != "http://31.56.177.191:8317/v1" {
		t.Fatalf("unexpected base url: %q", cfg.LLMBaseURL)
	}
	if cfg.LLMModel != "kimi-k2.7-code" {
		t.Fatalf("unexpected model: %q", cfg.LLMModel)
	}
}

func TestAllowedUsers(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("BOT_ALLOWED_USERS", "42, 100")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.UserAllowed(42) {
		t.Fatalf("expected user 42 to be allowed")
	}
	if cfg.UserAllowed(7) {
		t.Fatalf("expected user 7 to be denied")
	}
}
