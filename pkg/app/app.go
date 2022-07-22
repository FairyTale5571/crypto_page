package app

import (
	"github.com/caarlos0/env/v6"
	"github.com/fairytale5571/crypto_page/pkg/bot"
	"github.com/fairytale5571/crypto_page/pkg/database"
	"github.com/fairytale5571/crypto_page/pkg/logger"
	"github.com/fairytale5571/crypto_page/pkg/models"
	"github.com/fairytale5571/crypto_page/pkg/server"
	"github.com/fairytale5571/crypto_page/pkg/storage/redis"
)

type App struct {
	DB     *database.DB
	Config *models.Config
	Logger *logger.LoggerWrapper
	Server *server.Router
}

func New() (*App, error) {
	log := logger.New("application")

	cfg := &models.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Errorf("error parse config: %v", err)
		return nil, err
	}

	db, err := database.New(cfg.MysqlUri)
	if err != nil {
		log.Errorf("error start database: %v", err)
		return nil, err
	}

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Errorf("error start redis: %v", err)
		return nil, err
	}

	telegram, err := bot.NewTelegram(cfg, rdb, db)
	if err != nil {
		log.Errorf("error start telegram: %v", err)
		return nil, err
	}

	srv := server.New(cfg, telegram, rdb)
	if err != nil {
		log.Errorf("error start server: %v", err)
		return nil, err
	}

	go srv.Start()
	go telegram.Start()

	log.Info("application started")
	return &App{
		DB:     db,
		Config: cfg,
		Logger: log,
		Server: srv,
	}, nil

}
