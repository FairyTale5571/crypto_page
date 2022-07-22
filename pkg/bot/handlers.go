package bot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) eventUpdates(update tgbotapi.Update) {

	switch {
	case update.MyChatMember != nil:
		b.onMyChatMember(update.MyChatMember)
	case update.CallbackQuery != nil:
		b.onCallback(update.CallbackQuery)
	case update.Message != nil:
		b.onMessageCreate(update.Message)
	case update.Poll != nil:
		b.onPoll(update.Poll)
	}
}

func (b *Bot) onPoll(poll *tgbotapi.Poll) {
	js, err := json.Marshal(poll.Options)
	if err != nil {
		b.logger.Errorf("Error marshal poll: %v", err)
		return
	}
	_, err = b.database.Exec("UPDATE polls_result SET result = ? WHERE id = ?", string(js), poll.ID)
	if err != nil {
		b.logger.Errorf("Error update poll result: %v", err)
		return
	}
}

func (b *Bot) onMessageCreate(message *tgbotapi.Message) {
	if !b.isInDatabase(message.Chat.ID) {
		b.insertNewUser(message.Chat.ID, message.From.UserName, message.From.FirstName, message.From.LastName)
	}

	switch message.Command() {
	case "start":
		b.start(message)
		referralID := message.CommandArguments()
		b.updateUserNames(message.Chat.ID, message.From.UserName, message.From.FirstName, message.From.LastName, referralID)
	}

	switch message.Text {

	case "x":
		poll := tgbotapi.NewPoll(message.Chat.ID, "Как давно интересуетесь криптовалютой?", "Менее 1 года", "Более 1 года", "2 и более", "5+")
		poll.AllowsMultipleAnswers = true
		err := b.createPoll(&poll)
		if err != nil {
			b.logger.Errorf("Error create poll: %v", err)
			return
		}
	case "y":
		poll := tgbotapi.NewPoll(message.Chat.ID, "Как давно интересуетесь криптовалютой?", "Менее 1 года", "Более 1 года", "2 и более", "5+")
		poll.AllowsMultipleAnswers = false
		err := b.createPoll(&poll)
		if err != nil {
			b.logger.Errorf("Error create poll: %v", err)
			return
		}
	}
}

func (b *Bot) onCallback(callback *tgbotapi.CallbackQuery) {
	switch callback.Data {
	case "start_register":
		b.startRegister(callback)
	case "check_subscriptions":
		fmt.Println("check_subscriptions")
		b.verifyTelegram(callback)
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
		_, err := b.database.Exec("INSERT INTO chats (id, name, username) VALUES (?,?,?)", member.Chat.ID, member.Chat.Title, member.Chat.UserName)
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
