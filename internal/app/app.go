package app

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/EvgenyGulyaev/sample-ii-bot/internal/config"
	"github.com/EvgenyGulyaev/sample-ii-bot/internal/llm"
	"github.com/EvgenyGulyaev/sample-ii-bot/internal/telegram"
)

type App struct {
	cfg      config.Config
	tg       *telegram.Client
	llm      *llm.Client
	mu       sync.Mutex
	history  map[int64][]llm.Message
	inFlight map[int64]bool
}

func New(cfg config.Config) *App {
	return &App{
		cfg:      cfg,
		tg:       telegram.New(cfg.TelegramToken),
		llm:      llm.New(cfg.LLMBaseURL, cfg.LLMAPIKey, cfg.LLMModel, cfg.RequestTimeout),
		history:  make(map[int64][]llm.Message),
		inFlight: make(map[int64]bool),
	}
}

func (a *App) Run(ctx context.Context) error {
	log.Printf("bot started with model=%s base_url=%s", a.cfg.LLMModel, a.cfg.LLMBaseURL)
	var offset int64

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		updates, err := a.tg.GetUpdates(ctx, offset, a.cfg.PollTimeoutSeconds)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("telegram getUpdates: %v", err)
			continue
		}

		for _, update := range updates {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			if update.Message == nil {
				continue
			}
			a.handleMessage(ctx, *update.Message)
		}
	}
}

func (a *App) handleMessage(ctx context.Context, message telegram.Message) {
	text := strings.TrimSpace(message.Text)
	if text == "" {
		return
	}

	if !a.cfg.UserAllowed(message.From.ID) {
		_ = a.tg.SendMessage(ctx, message.Chat.ID, "У тебя пока нет доступа к этому боту.")
		return
	}

	if strings.EqualFold(text, "/start") {
		_ = a.tg.SendMessage(ctx, message.Chat.ID, "Привет. Напиши вопрос, я отвечу через нейронку.")
		return
	}
	if strings.EqualFold(text, "/reset") {
		a.resetHistory(message.Chat.ID)
		_ = a.tg.SendMessage(ctx, message.Chat.ID, "Историю диалога очистил.")
		return
	}

	if !a.tryLockChat(message.Chat.ID) {
		_ = a.tg.SendMessage(ctx, message.Chat.ID, "Я еще думаю над прошлым сообщением. Подожди пару секунд.")
		return
	}
	defer a.unlockChat(message.Chat.ID)

	_ = a.tg.SendTyping(ctx, message.Chat.ID)
	messages := a.messagesForChat(message.Chat.ID, text)
	answer, err := a.llm.Complete(ctx, messages)
	if err != nil {
		log.Printf("llm complete: %v", err)
		_ = a.tg.SendMessage(ctx, message.Chat.ID, "Не получилось получить ответ от нейронки. Попробуй еще раз чуть позже.")
		return
	}

	a.appendHistory(message.Chat.ID, text, answer)
	if err := a.tg.SendMessage(ctx, message.Chat.ID, answer); err != nil {
		log.Printf("telegram sendMessage: %v", err)
	}
}

func (a *App) messagesForChat(chatID int64, text string) []llm.Message {
	a.mu.Lock()
	defer a.mu.Unlock()

	messages := []llm.Message{{Role: "system", Content: a.cfg.SystemPrompt}}
	messages = append(messages, a.history[chatID]...)
	messages = append(messages, llm.Message{Role: "user", Content: text})
	return messages
}

func (a *App) appendHistory(chatID int64, question, answer string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	history := append(a.history[chatID],
		llm.Message{Role: "user", Content: question},
		llm.Message{Role: "assistant", Content: answer},
	)
	if limit := a.cfg.HistoryMessages; limit > 0 && len(history) > limit {
		history = history[len(history)-limit:]
	}
	a.history[chatID] = history
}

func (a *App) resetHistory(chatID int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.history, chatID)
}

func (a *App) tryLockChat(chatID int64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.inFlight[chatID] {
		return false
	}
	a.inFlight[chatID] = true
	return true
}

func (a *App) unlockChat(chatID int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.inFlight, chatID)
}
