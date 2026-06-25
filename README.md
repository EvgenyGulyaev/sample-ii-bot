# sample-ii-bot

Минимальный Telegram-бот на Go. Он принимает текстовые сообщения в Telegram,
отправляет их в OpenAI-compatible LLM endpoint и возвращает ответ обратно в чат.

## Что умеет

- Отвечает на обычные текстовые сообщения.
- Работает через Telegram long polling, без открытого HTTP-порта.
- Использует Kimi через OpenAI-compatible API.
- Держит короткую историю диалога в памяти процесса.
- Поддерживает `/start` и `/reset`.
- Блокирует параллельный запрос в одном чате, чтобы не плодить дорогие ответы.
- Может ограничивать доступ по Telegram user id.

## Как устроен

- `cmd/bot-ii` - точка входа.
- `internal/app` - основной цикл бота и обработка сообщений.
- `internal/telegram` - минимальный клиент Telegram Bot API.
- `internal/llm` - клиент `/chat/completions`.
- `internal/config` - чтение env-настроек.
- `internal/envfile` - загрузка локального `.env` при запуске руками.

База данных не используется. История хранится только в оперативной памяти и
сбрасывается при рестарте сервиса.

## Какая ИИ используется

По умолчанию бот ходит в роутер из Zed:

```env
LLM_BASE_URL=http://31.56.177.191:8317/v1
LLM_MODEL=kimi-k2.7-code
```

Endpoint совместим с OpenAI Chat Completions:

```text
POST /chat/completions
```

Для Kimi выставлена `temperature=1`, потому что эта модель не принимает другие
значения температуры.

## Команды в Telegram

- `/start` - короткое приветствие.
- `/reset` - очистить историю текущего чата.

Все остальные текстовые сообщения отправляются в нейронку.

## Env

- `TELEGRAM_BOT_TOKEN` - токен Telegram-бота.
- `LLM_BASE_URL` - OpenAI-compatible base URL.
- `LLM_MODEL` - модель, например `kimi-k2.7-code`.
- `LLM_API_KEY` - API key.
- `BOT_ALLOWED_USERS` - список Telegram user id через запятую. Если пусто, бот доступен всем.
- `BOT_SYSTEM_PROMPT` - системный промпт.
- `BOT_HISTORY_MESSAGES` - сколько последних сообщений хранить в памяти на чат.
- `BOT_POLL_TIMEOUT_SECONDS` - timeout long polling к Telegram.
- `BOT_REQUEST_TIMEOUT_SECONDS` - timeout запроса к нейронке.

`.env` нельзя пушить в репозиторий. В git лежит только `.env.example`.

## Локальный запуск

```bash
cp .env.example .env
go run ./cmd/bot-ii
```
