package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) eventUpdates(update tgbotapi.Update) {
	switch {
	case update.CallbackQuery != nil:
		b.onCallback(update.CallbackQuery)
	case update.Message != nil:
		b.onMessageCreate(update.Message)
	case update.MyChatMember != nil:
		b.onMyChatMember(update.MyChatMember)
	}
}

func (b *Bot) onMessageCreate(message *tgbotapi.Message) {

}

func (b *Bot) onCallback(callback *tgbotapi.CallbackQuery) {
	switch callback.Data {

	}
}

func (b *Bot) onMyChatMember(member *tgbotapi.ChatMemberUpdated) {
	switch {
	case member.Chat.ID < 0:
		if member.NewChatMember.Status == "left" || member.NewChatMember.Status == "kicked" {
			_, err := b.database.Exec("DELETE FROM chats WHERE id = ?", member.Chat.ID)
			if err != nil {
				b.logger.Errorf("Error delete status: %v", err)
				return
			}
			return
		}
		_, err := b.database.Exec("INSERT INTO chats (id, name) VALUES (?,?)", member.Chat.ID, member.Chat.Title)
		if err != nil {
			b.logger.Errorf("Error insert chat: %v", err)
			return
		}
	case member.OldChatMember.User.IsBot:
		_, err := b.database.Exec("UPDATE users SET status = ? WHERE telegram_id = ?", member.NewChatMember.Status, member.Chat.ID)
		if err != nil {
			b.logger.Errorf("Error update user status: %v", err)
			return
		}
	}
}
