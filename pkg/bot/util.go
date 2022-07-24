package bot

import (
	"fmt"
	"strconv"

	"github.com/fairytale5571/crypto_page/pkg/models"
	"github.com/fairytale5571/crypto_page/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) isAlreadyRegistered(field, data string) bool {
	var id string
	query := fmt.Sprintf("select id from users where %s = ?", field)
	_ = b.database.QueryRow(query, data).Scan(&id)
	return id != ""
}

func (b *Bot) TwitterValid(id string, username string) {
	n, _ := strconv.ParseInt(id, 10, 64)
	if b.isAlreadyRegistered("twitter", username) {
		b.TwitterNotValid(id)
		return
	}
	err := b.redis.Set(fmt.Sprintf("twitter_id:%d", n), username, storage.UserTwitter)
	if err != nil {
		b.logger.Errorf("error TwitterValid: %v", err)
		return
	}
	b.sendMessage(n, "Проверка подписки на Twitter прошла успешно")
	b.verifyInstagram(n)
}

func (b *Bot) TwitterNotValid(id string) {
	n, _ := strconv.ParseInt(id, 10, 64)
	_ = b.sendMessage(n, "Подписка на Twitter не прошла проверку")
	b.verifyTwitter(n)
}

func (b *Bot) getAllChats() map[string]string {
	rows, err := b.database.Query("SELECT name, username FROM chats")
	defer rows.Close() // nolint: not needed

	if err != nil {
		b.logger.Errorf("error getAllChats: %v", err)
		return nil
	}

	chats := make(map[string]string)
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			b.logger.Errorf("error getAllChats: %v", err)
			return nil
		}
		if name == "" {
			continue
		}
		chats[id] = name
	}

	return chats

}

func (b *Bot) startRegister(callback *tgbotapi.CallbackQuery) {
	b.deleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	msg := b.photoConfigUrl(callback.From.ID, b.cfg.URL+"/assets/images/telegram.jpg", "Подпишитесь на каналы по ссылкам ниже и нажмите \n\"✅ Проверить подписки\"")

	channels := b.getAllChats()
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for k, v := range channels {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(k, "t.me/"+v),
		))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Проверить подписки", "check_subscriptions"),
	))

	msg.ReplyMarkup = keyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		b.logger.Errorf("Error send message: %v", err)
	}
}

func (b *Bot) createPoll(poll *tgbotapi.SendPollConfig) error {
	_msg, err := b.bot.Send(poll)
	if err != nil {
		b.logger.Errorf("error createPoll: %v", err)
		return err
	}
	_, err = b.database.Exec("INSERT INTO polls_result (id, telegram_id, poll, insert_time) VALUES (?,?,?, NOW())", _msg.Poll.ID, _msg.Chat.ID, _msg.Poll.Question)
	if err != nil {
		b.logger.Errorf("error createPoll: %v", err)
		return err
	}
	return nil
}

func (b *Bot) photoConfigUrl(id int64, url, caption string) *tgbotapi.PhotoConfig {
	return &tgbotapi.PhotoConfig{
		Caption: caption,
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: id,
			},

			File: tgbotapi.FileURL(url),
		},
	}
}

func (b *Bot) isInDatabase(userId int64) bool {
	var id string
	err := b.database.QueryRow("select id from users where telegram_id = ?", userId).Scan(&id)
	if err != nil {
		b.logger.Errorf("error isInDatabase: %v", err)
		return false
	}
	return id != ""
}

func (b *Bot) insertNewUser(userId int64, username, firstName, lastName string) {
	stmt, err := b.database.Prepare("INSERT INTO users (telegram_id, user_name, user_first_name, user_last_name, status, registered_at) VALUES (?,?,?,?, 'member', now())")
	if err != nil {
		b.logger.Errorf("error insertNewUser: %v", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, username, firstName, lastName)
	if err != nil {
		b.logger.Errorf("error insertNewUser: %v", err)
		return
	}
}

func (b *Bot) updateUserNames(userId int64, username, firstName, lastName, referral string) {
	var err error
	if referral != "" {
		_, err = b.database.Exec("UPDATE users SET user_name = ?, user_first_name = ?, user_last_name = ?, referred_by = ? WHERE telegram_id = ?", username, firstName, lastName, referral, userId)
	} else {
		_, err = b.database.Exec("UPDATE users SET user_name = ?, user_first_name = ?, user_last_name = ? WHERE telegram_id = ?", username, firstName, lastName, userId)
	}

	if err != nil {
		b.logger.Errorf("error updateUserNames: %v", err)
	}
}

func (b *Bot) validInstagram(id string) (bool, string) {
	b.logger.Infof("%s is subscriber on instagram", id)
	return checkString(id), ""
}

func (b *Bot) getTwitterUrl(id int64) string {
	consumer := models.TwitterConsumer
	tokenUrl := fmt.Sprintf("%s/auth/twitter/callback", b.cfg.URL)
	token, requestUrl, err := consumer.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		b.logger.Errorf("error getTwitterTokens: %v", err)
		return ""
	}
	models.Tokens[token.Token] = token

	_ = b.redis.Set(fmt.Sprintf("twitter:ts_%s_id", token.Token), fmt.Sprintf("%d", id), storage.UserTwitter)
	return requestUrl
}

func (b *Bot) sendMessage(id int64, s string) tgbotapi.Message {
	msg := tgbotapi.NewMessage(id, s)
	send, _ := b.bot.Send(msg)
	return send
}

func (b *Bot) deleteMessage(chatID int64, messageID int) {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, _ = b.bot.Send(msg)
}

func (b *Bot) getIdFromPoll(id string) int64 {
	var telegramId int64
	err := b.database.QueryRow("select telegram_id from polls_result where id = ?", id).Scan(&telegramId)
	if err != nil {
		b.logger.Errorf("error getIdFromPoll: %v", err)
		return 0
	}
	return telegramId
}

func (b *Bot) SendMessage(id int64, s string) tgbotapi.Message {
	return b.sendMessage(id, s)
}

func (b *Bot) lastVerify(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Хотите ли быть вовлеченным в процесс развития децентрализованной социальной сети Crypto.Page?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да", "want_yes"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет", "want_no"),
		),
	)
	_, _ = b.bot.Send(msg)
}

func (b *Bot) wantYes(callback *tgbotapi.CallbackQuery) {
	waitWhyYouCanHelp[callback.From.ID] = struct{}{}
	msg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "Почему вы хотите участвовать?\nОпишите одним сообщением")
	_, _ = b.bot.Send(msg)
}

func (b *Bot) finishRegistration(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Поздравляем 🎉 \nРегистрация завершена!\nОсталось только пригласить друга и увеличить свои шансы!\nЖми👇Реферальная программа.")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttonsAbout),
			tgbotapi.NewKeyboardButton(buttonsReferral),
		),
	)
	_, _ = b.bot.Send(msg)
}

func (b *Bot) about(message *tgbotapi.Message) {
	msg := b.photoConfigUrl(message.Chat.ID, b.cfg.URL+"/assets/images/crypto_page.jpg", "Crypto.Page — это децентрализованная cross-chain социальная сеть , которая поддерживает Binance Smart Chain, Ethereum, Polygon, Tron, Solana.\nПользователи платят комиссию, размер которой зависит от выбранного блокчейна, в обмен на токены PAGE за каждое действие на платформе Crypto.Page (публикация контента, оценка контента, подтверждение имени или получение статуса амбасадора).\n\n"+
		"Весь опубликованный контент создается и остается в виде невзаимозаменяемых токенов (NFT), которые можно продавать с аукциона, передавать или сжигать при фиксации токенов PAGE. "+
		"Создатель NFT получает процент с каждой перепродажи усилий в токенах PAGE, а также с других действий, например, по рейтингу контента. "+
		"Из-за децентрализованного характера социальной сети только владелец контента может удалить (сжечь) опубликованный контент.")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Web", "https://crypto.page"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Приложение", "https://app.crypto.page"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Официальный Твиттер", "https://twitter.com/cryptopage_web3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Инстаграм", "https://instagram.com/cryptopage_web3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Англоязычное комьюнити", "https://t.me/cryptopage_web3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Русскоязычное комьюнити", "https://t.me/cryptopage_web_3"),
		),
	)

	_, _ = b.bot.Send(msg)
}

func (b *Bot) countReferrals(id int64) uint {
	var count uint
	err := b.database.QueryRow("select count(*) from users where referred_by = ? and telegram_id != ? "+
		"and users.instagram is not null "+
		"and users.twitter is not null "+
		"and users.status = 'member'", id, id).Scan(&count)
	if err != nil {
		b.logger.Errorf("error isInDatabase: %v", err)
		return count
	}
	return count
}

func (b *Bot) referral(message *tgbotapi.Message) {

	var involvedText string
	referrals := b.countReferrals(message.Chat.ID)
	if referrals == 0 {
		involvedText = "не участвуете ❌ - пригласите хотя бы одного реферала"
	} else {
		involvedText = "вы учавствуете"
	}
	msg := b.photoConfigUrl(message.Chat.ID, b.cfg.URL+"/assets/images/about.jpg", fmt.Sprintf(
		"Будь вовлечён в проект Crypto.Page!\n\n"+
			"⚡️Имей план Б, пригласи друга и получи гарантированный приз за его победу.\n\n"+
			"Ваш статус: "+involvedText+"\n"+
			"Ваши друзья: %d 👥\n"+
			"Для участия в розыгрыше необходимо выполнить все условия и пригласить как минимум одного друга. \n"+
			"Ваша личная ссылка для приглашений 🔗:\n"+
			"https://t.me/crypto_page_bot?start=%d", referrals, message.Chat.ID))
	_, _ = b.bot.Send(msg)
}

func (b *Bot) getUserStatus(id int64) string {
	var status string
	err := b.database.QueryRow("select status from users where telegram_id = ?", id).Scan(&status)
	if err != nil {
		b.logger.Errorf("error getUserStatus: %v", err)
		return ""
	}
	return status
}

func checkString(str string) bool {
	allowedCharacters := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		".", "_", "-",
	}
	for _, v := range str {
		for _, v2 := range allowedCharacters {
			if string(v) == v2 {
				return true
			}
		}
	}
	return false
}
