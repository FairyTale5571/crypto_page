package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) start(message *tgbotapi.Message) {

	msg := b.photoConfigUrl(message.Chat.ID, b.cfg.URL+"/assets/images/main_date.jpg", "Будь вовлечён в проект Crypto.Page! \n\n"+
		"Случайным образом, разыграем 500 USDT☑️\n\n"+
		"💥 Заполни анкету, выполнив все условия и стань участником розыгрыша!\n"+
		"⚡️ Имей план Б, пригласи друга и получи гарантированный приз за его победу.")

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
