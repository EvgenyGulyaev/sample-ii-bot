package config

import "testing"

func TestLoadRequiresLLMSettings(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_MODEL", "")

	if _, err := Load(); err == nil {
		t.Fatalf("expected error without LLM settings")
	}
}

func TestAllowedUsers(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("LLM_BASE_URL", "http://example.test/v1")
	t.Setenv("LLM_MODEL", "test-model")
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
