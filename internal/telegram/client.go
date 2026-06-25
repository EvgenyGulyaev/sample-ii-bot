package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Text      string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

func New(token string) *Client {
	return &Client{
		baseURL: "https://api.telegram.org/bot" + strings.TrimSpace(token),
		http:    &http.Client{Timeout: 90 * time.Second},
	}
}

func (c *Client) GetUpdates(ctx context.Context, offset int64, timeoutSeconds int) ([]Update, error) {
	values := url.Values{}
	values.Set("offset", strconv.FormatInt(offset, 10))
	values.Set("timeout", strconv.Itoa(timeoutSeconds))
	values.Set("allowed_updates", `["message"]`)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/getUpdates?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	if err := c.doJSON(req, &parsed); err != nil {
		return nil, err
	}
	if !parsed.OK {
		return nil, fmt.Errorf("telegram getUpdates failed")
	}
	return parsed.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	for _, part := range splitTelegramText(text) {
		body := map[string]any{
			"chat_id": chatID,
			"text":    part,
		}
		if err := c.postJSON(ctx, "/sendMessage", body, nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) SendTyping(ctx context.Context, chatID int64) error {
	body := map[string]any{
		"chat_id": chatID,
		"action":  "typing",
	}
	return c.postJSON(ctx, "/sendChatAction", body, nil)
}

func (c *Client) postJSON(ctx context.Context, path string, body any, target any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doJSON(req, target)
}

func (c *Client) doJSON(req *http.Request, target any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram status %d: %s", resp.StatusCode, string(data))
	}
	if target == nil {
		return nil
	}
	return json.Unmarshal(data, target)
}

func splitTelegramText(text string) []string {
	const limit = 3900
	runes := []rune(text)
	if len(runes) <= limit {
		return []string{text}
	}

	parts := make([]string, 0, len(runes)/limit+1)
	for len(runes) > 0 {
		end := limit
		if len(runes) < end {
			end = len(runes)
		}
		parts = append(parts, string(runes[:end]))
		runes = runes[end:]
	}
	return parts
}
