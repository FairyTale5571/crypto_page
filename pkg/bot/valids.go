package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) verifyTelegram(message *tgbotapi.CallbackQuery) {
	isSub, channel := b.checkTelegram(message)
	if !isSub {
		msg := tgbotapi.NewMessage(message.Message.Chat.ID, "Вы не подписаны на телеграм канал "+channel+". Пожалуйста, подпишитесь на него и попробуйте снова.")
		_, err := b.bot.Send(msg)
		if err != nil {
			b.logger.Errorf("Error send message: %v", err)
		}
		b.logger.Infof("%d is not subscriber on %s telegram", message.From.ID, err)
		return
	}
	msg := tgbotapi.NewMessage(message.Message.Chat.ID, "Проверка пройдена")
	_, err := b.bot.Send(msg)
	if err != nil {
		b.logger.Errorf("Error send message: %v", err)
	}

}

func (b *Bot) checkTelegram(message *tgbotapi.CallbackQuery) (bool, string) {
	chats := func() map[int64]string {
		chats := map[int64]string{}
		rows, err := b.database.Query(`select id, name from chats`)
		defer rows.Close()
		if err != nil {
			b.logger.Errorf("error get chats: %v", err)
			return chats
		}
		for rows.Next() {
			var id int64
			var chat string
			err = rows.Scan(&id, &chat)
			if err != nil {
				b.logger.Errorf("error get chats: %v", err)
				return chats
			}
			chats[id] = chat
		}
		return chats
	}
	for k, v := range chats() {
		user, err := b.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: k,
				UserID: message.Message.Chat.ID,
			},
		})
		if err != nil {
			b.logger.Errorf("error get chat member: %v", err)
			return false, ""
		}
		if user.Status == "left" || user.Status == "kicked" {
			return false, v
		}
	}
	b.logger.Infof("%d is subscriber on telegram", message.From.ID)
	return true, ""

}
