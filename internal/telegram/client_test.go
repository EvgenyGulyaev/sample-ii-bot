package telegram

import (
	"errors"
	"strings"
	"testing"
)

func TestSanitizeErrorHidesTelegramToken(t *testing.T) {
	err := sanitizeError(errors.New(`Get "https://api.telegram.org/bot123456789:AASecret_token/getUpdates": timeout`))
	if strings.Contains(err.Error(), "AASecret_token") {
		t.Fatalf("expected token to be hidden, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "bot[TOKEN]") {
		t.Fatalf("expected sanitized marker, got %q", err.Error())
	}
}
