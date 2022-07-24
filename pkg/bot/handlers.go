package bot

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) eventUpdates(update tgbotapi.Update) {

	switch {
	case update.Poll != nil:
		b.onPoll(update.Poll)
	case update.MyChatMember != nil:
		b.onMyChatMember(update.MyChatMember)
	case update.CallbackQuery != nil:
		b.onCallback(update.CallbackQuery)
	case update.Message != nil:
		b.onMessageCreate(update.Message)
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
	switch poll.Question {
	case "Как давно интересуетесь криптовалютой?":
		b.createPolls(b.getIdFromPoll(poll.ID), "Дополнительные интересы")
	case "Дополнительные интересы":
		for _, v := range poll.Options {
			if v.VoterCount > 0 {
				b.createPolls(b.getIdFromPoll(poll.ID), v.Text)
				return
			}
		}
	default:
		b.lastVerify(b.getIdFromPoll(poll.ID))
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
		b.cleanWaiting(message.Chat.ID)
		b.updateUserNames(message.Chat.ID, message.From.UserName, message.From.FirstName, message.From.LastName, referralID)
	}

	switch message.Text {

	case buttonsAbout:
		b.about(message)
	case buttonsReferral:
		b.referral(message)
	case "export":
		if b.getUserStatus(message.Chat.ID) == "admin" {
			b.export(message)
		}
	default:
		b.onHandleWait(message)
	}
}

func (b *Bot) onCallback(callback *tgbotapi.CallbackQuery) {
	switch callback.Data {
	case "start_register":
		b.startRegister(callback)
	case "check_subscriptions":
		fmt.Println("check_subscriptions")
		b.verifyTelegram(callback)
	case "twitter_old":
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.verifyTwitter(callback.From.ID)
	case "want_yes":
		b.wantYes(callback)
	case "want_no":
		b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.finishRegistration(callback.Message.Chat.ID)

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

func (b *Bot) onHandleWait(message *tgbotapi.Message) {

	if _, ok := waitInstagram[message.Chat.ID]; ok {
		b.checkInstagram(message)
		return
	}

	if _, ok := waitWhyYouCanHelp[message.Chat.ID]; ok {
		_, err := b.database.Exec("UPDATE users SET want_help = ? WHERE telegram_id = ?", message.Text, message.Chat.ID)
		if err != nil {
			b.logger.Errorf("Error update want help: %v", err)
		}
		b.finishRegistration(message.Chat.ID)
		delete(waitWhyYouCanHelp, message.Chat.ID)
		return
	}
}
