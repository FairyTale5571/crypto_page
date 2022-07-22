package server

import (
	"github.com/fairytale5571/crypto_page/pkg/bot"
	"github.com/fairytale5571/crypto_page/pkg/logger"
	"github.com/fairytale5571/crypto_page/pkg/models"
	"github.com/fairytale5571/crypto_page/pkg/storage/redis"
	"github.com/gin-gonic/gin"
	"github.com/mrjones/oauth"
)

type Router struct {
	router   *gin.Engine
	bot      *bot.Bot
	cfg      *models.Config
	settings AuthConfig
	redis    *redis.Redis
	logger   *logger.LoggerWrapper
}

type AuthConfig struct {
	TwitterConfig *oauth.Consumer
}

const DSID = "981627353857937509"

func New(cfg *models.Config, bot *bot.Bot, rdb *redis.Redis) *Router {
	r := Router{
		bot:    bot,
		cfg:    cfg,
		redis:  rdb,
		logger: logger.New("server"),
		router: gin.Default(),
	}

	r.settings = AuthConfig{
		TwitterConfig: models.TwitterConsumer,
	}
	return &r
}

func (r *Router) Start() {

	r.router.Static("/assets/", "webApp/assets/")
	r.logger.Info("gin opened")

	r.mainRouter()
	err := r.router.Run(":" + r.cfg.PORT)
	if err != nil {
		r.logger.Errorf("cant open gin engine: %v", err)
		return
	}
}

func (r *Router) mainRouter() {
	r.router.GET("/auth/twitter/callback", r.twitterHandle)
}

func (r *Router) Stop() {
	r.logger.Info("gin closed")
}
