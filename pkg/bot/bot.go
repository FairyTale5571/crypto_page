package bot

import (
	"github.com/fairytale5571/crypto_page/pkg/database"
	"github.com/fairytale5571/crypto_page/pkg/logger"
	"github.com/fairytale5571/crypto_page/pkg/models"
	"github.com/fairytale5571/crypto_page/pkg/storage/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"time"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	updates  tgbotapi.UpdatesChannel
	cfg      *models.Config
	redis    *redis.Redis
	database *database.DB
	logger   *logger.LoggerWrapper
}

const (
	updateOffset  = 0
	updateTimeout = 64
)

var (
	location, _ = time.LoadLocation("Europe/Kiev")
)

func NewTelegram(cfg *models.Config, redis *redis.Redis, database *database.DB) (*Bot, error) {

	b := &Bot{
		cfg:      cfg,
		redis:    redis,
		database: database,
		logger:   logger.New("telegram"),
	}
	bot, err := tgbotapi.NewBotAPIWithClient(cfg.TelegramToken, tgbotapi.APIEndpoint, &http.Client{
		Timeout: updateTimeout * time.Second,
	})
	bot.Debug = cfg.Debug
	if err != nil {
		b.logger.Errorf("error start telegram api with client: %v", err)
		return nil, err
	}
	b.bot = bot
	return b, nil
}

func (b *Bot) Start() {

	botUpdate := tgbotapi.NewUpdate(updateOffset)
	botUpdate.Timeout = updateTimeout
	b.updates = b.bot.GetUpdatesChan(botUpdate)

	for update := range b.updates {
		go b.eventUpdates(update)
	}
}
