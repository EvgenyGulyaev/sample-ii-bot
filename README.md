# sample-ii-bot

Минимальный Telegram-бот на Go, который ходит в OpenAI-compatible LLM endpoint.

## Локальный запуск

```bash
cp .env.example .env
go run ./cmd/bot-ii
```

## Env

- `TELEGRAM_BOT_TOKEN` - токен Telegram-бота.
- `LLM_BASE_URL` - OpenAI-compatible base URL, например `http://31.56.177.191:8317/v1`.
- `LLM_MODEL` - модель, например `kimi-k2.7-code`.
- `LLM_API_KEY` - API key, можно оставить пустым, если endpoint не требует авторизацию.
- `BOT_ALLOWED_USERS` - опциональный список Telegram user id через запятую. Если пусто, бот доступен всем.
- `BOT_SYSTEM_PROMPT` - системный промпт.
- `BOT_HISTORY_MESSAGES` - сколько последних сообщений хранить в памяти на чат.

## VPS

Рабочая директория: `/var/go/bot-ii`.

```bash
systemctl restart bot-ii
journalctl -u bot-ii -f
```
