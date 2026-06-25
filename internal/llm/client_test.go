package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCompleteCallsChatCompletions(t *testing.T) {
	var auth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		auth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"ok"}}]}`))
	}))
	defer server.Close()

	client := New(server.URL+"/v1", "secret", "test-model", time.Second)
	answer, err := client.Complete(context.Background(), []Message{{Role: "user", Content: "hi"}})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if answer != "ok" {
		t.Fatalf("unexpected answer: %q", answer)
	}
	if auth != "Bearer secret" {
		t.Fatalf("unexpected auth header: %q", auth)
	}
}
