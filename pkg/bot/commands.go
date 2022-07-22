package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) start(message *tgbotapi.Message) {

	msg := b.photoConfigUrl(message.Chat.ID, b.cfg.URL+"/assets/images/crypto_page_main.jpg", "Crypto.Page - decentralized cross-chain social network and nft marketplace")

	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Начать регистрацию!", "start_register"),
		),
	)

	msg.ReplyMarkup = buttons
	_, err := b.bot.Send(msg)
	if err != nil {
		b.logger.Errorf("Error send message: %v", err)
	}
}
