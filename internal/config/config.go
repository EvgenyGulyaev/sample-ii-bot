package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	TelegramToken string
	LLMBaseURL    string
	LLMModel      string
	LLMAPIKey     string

	AllowedUsers       map[int64]bool
	SystemPrompt       string
	HistoryMessages    int
	PollTimeoutSeconds int
	RequestTimeout     time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		TelegramToken:      strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		LLMBaseURL:         trimRightSlash(env("LLM_BASE_URL", "http://31.56.177.191:8317/v1")),
		LLMModel:           env("LLM_MODEL", "kimi-k2.7-code"),
		LLMAPIKey:          strings.TrimSpace(os.Getenv("LLM_API_KEY")),
		AllowedUsers:       parseAllowedUsers(os.Getenv("BOT_ALLOWED_USERS")),
		SystemPrompt:       env("BOT_SYSTEM_PROMPT", "Ты дружелюбный и полезный ассистент. Отвечай кратко и по делу на языке пользователя."),
		HistoryMessages:    intEnv("BOT_HISTORY_MESSAGES", 12),
		PollTimeoutSeconds: intEnv("BOT_POLL_TIMEOUT_SECONDS", 30),
		RequestTimeout:     time.Duration(intEnv("BOT_REQUEST_TIMEOUT_SECONDS", 120)) * time.Second,
	}

	if cfg.TelegramToken == "" {
		return Config{}, errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.LLMBaseURL == "" {
		return Config{}, errors.New("LLM_BASE_URL is required")
	}
	if cfg.LLMModel == "" {
		return Config{}, errors.New("LLM_MODEL is required")
	}
	if cfg.HistoryMessages < 0 {
		cfg.HistoryMessages = 0
	}
	if cfg.PollTimeoutSeconds < 5 {
		cfg.PollTimeoutSeconds = 5
	}
	if cfg.RequestTimeout < 5*time.Second {
		cfg.RequestTimeout = 5 * time.Second
	}

	return cfg, nil
}

func (c Config) UserAllowed(id int64) bool {
	if len(c.AllowedUsers) == 0 {
		return true
	}
	return c.AllowedUsers[id]
}

func env(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func intEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseAllowedUsers(value string) map[int64]bool {
	users := make(map[int64]bool)
	for _, raw := range strings.Split(value, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		id, err := strconv.ParseInt(raw, 10, 64)
		if err == nil {
			users[id] = true
		}
	}
	return users
}

func trimRightSlash(value string) string {
	return strings.TrimRight(strings.TrimSpace(value), "/")
}
