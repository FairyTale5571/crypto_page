package bot

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/fairytale5571/crypto_page/pkg/database"
	"github.com/fairytale5571/crypto_page/pkg/logger"
	"github.com/fairytale5571/crypto_page/pkg/models"
	"github.com/fairytale5571/crypto_page/pkg/storage/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func NewTelegram(cfg *models.Config, rdb *redis.Redis, db *database.DB) (*Bot, error) {

	b := &Bot{
		cfg:      cfg,
		redis:    rdb,
		database: db,
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

func (b *Bot) getUsers() ([]*Polls, error) {
	rows, err := b.database.Query(`
		SELECT u.telegram_id,
			   u.user_name,
			   u.user_first_name,
			   u.user_last_name,
			   u.instagram,
			   u.twitter,
			   u.want_help,
			   u.referred_by,
			   (select count(*) from users us where us.referred_by = u.telegram_id and us.telegram_id != u.referred_by
					and us.instagram is not null
					and us.twitter is not null
				   ) as TotalInvites,
			   u.registered_at FROM users u
		`)
	defer rows.Close()

	if err != nil {
		b.logger.Errorf("error get users: %v", err)
		return nil, err
	}

	type poll struct {
		Text       string `json:"text"`
		VoterCount int    `json:"voter_count"`
	}

	var users []*Polls
	for rows.Next() {
		var user Polls
		err = rows.Scan(&user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.Instagram, &user.Twitter, &user.WantHelp, &user.ReferredBy, &user.TotalInvites, &user.RegisteredAt)
		if err != nil {
			b.logger.Errorf("error scan users: %v", err)
			return nil, err
		}
		polls, err := b.database.Query("SELECT poll, result FROM polls_result WHERE telegram_id = ?", user.TelegramID)
		if err != nil {
			b.logger.Errorf("error get polls: %v", err)
		}
		for polls.Next() {
			var question string
			var result string

			err = polls.Scan(&question, &result)
			if err != nil {
				b.logger.Errorf("error scan polls: %v", err)
			}
			var p []poll
			err = json.Unmarshal([]byte(result), &p)
			if err != nil {
				b.logger.Errorf("error unmarshal polls: %v", err)
			}
			for _, v := range p {
				if v.VoterCount > 0 {
					user.Answers += question + "\n" + v.Text + "\n\n"
				}
			}
		}
		users = append(users, &user)
	}
	return users, nil
}

func (b *Bot) export(message *tgbotapi.Message) {
	err := os.Remove("users.csv")
	if err != nil {
		b.logger.Errorf("error remove file: %v", err)
	}

	f, err := os.Create("users.csv")
	defer f.Close() // nolint: not needed
	if err != nil {
		b.logger.Errorf("error create file: %v", err)
		return
	}
	defer b.sendUserFile(message.Chat.ID)

	w := csv.NewWriter(f)
	_ = w.Write([]string{
		"Telegram ID",
		"Юзернейм",
		"Имя",
		"Фамилия",
		"Instagram",
		"Twitter",
		"Хочет помочь",
		"Приглашен пользователем",
		"Пригласил пользователей",
		"Дата регистрации",
		"Опросник",
	})

	users, err := b.getUsers()
	if err != nil {
		b.logger.Errorf("error get users: %v", err)
		return
	}
	for _, v := range users {
		_ = w.Write([]string{
			v.TelegramID,
			v.Username.String,
			v.FirstName.String,
			v.LastName.String,
			v.Instagram.String,
			v.Twitter.String,
			v.WantHelp.String,
			v.ReferredBy.String,
			v.TotalInvites.String,
			v.RegisteredAt.In(location).Format("15:04:05 02.01.2006"),
			v.Answers,
		})
	}
	w.Flush()
}

func (b *Bot) sendUserFile(chatID int64) {
	f, err := os.OpenFile("users.csv", os.O_RDONLY, 0o666)
	if err != nil {
		b.logger.Errorf("error open file: %v", err)
		return
	}
	msg := tgbotapi.NewDocument(chatID, tgbotapi.FileReader{
		Name:   "users.csv",
		Reader: f,
	})
	_, err = b.bot.Send(msg)
}
