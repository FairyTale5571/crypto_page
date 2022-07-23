package bot

import (
	"fmt"
	"github.com/fairytale5571/crypto_page/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
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

	b.verifyTwitter(message.From.ID)
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

func (b *Bot) verifyTwitter(id int64) {
	_msg := b.photoConfigUrl(id, b.cfg.URL+"/assets/images/crypto_page_main.jpg", "Подпишитесь на наш Twitter и нажмите  \"✅ Проверить\"")
	_msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Crypto.page", "https://twitter.com/cryptopage_web3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("✅ Проверить", b.getTwitterUrl(id)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ссылка устарела", "twitter_old"),
		),
	)
	_, err := b.bot.Send(_msg)
	if err != nil {
		b.logger.Errorf("Error send message: %v", err)
	}
}

func (b *Bot) verifyInstagram(id int64) {
	msg := tgbotapi.NewMessage(id, "Подпишитесь на наш Instagram и нажмите введите свой логин ниже")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Подписаться", "https://www.instagram.com/cryptopage_web3"),
		),
	)
	waitInstagram[id] = struct{}{}
	_, err := b.bot.Send(msg)
	if err != nil {
		b.logger.Errorf("Error send message: %v", err)
	}
}

func (b *Bot) checkInstagram(message *tgbotapi.Message) {
	if ok, _ := b.validInstagram(strings.TrimLeft(message.Text, "@")); !ok {
		b.sendMessage(message.Chat.ID, "Не удалось проверить введенные данные, попробуйте еще раз")
		return
	}
	if b.isAlreadyRegistered("instagram", message.Text) {
		b.sendMessage(message.Chat.ID, "Этот аккаунт уже зарегистрирован")
		return
	}
	b.redis.Set(fmt.Sprintf("instagram:%d", message.Chat.ID), message.Text, storage.UserInstagram)
	b.sendMessage(message.Chat.ID, "Ваша инстаграм подтвержден")
	delete(waitInstagram, message.Chat.ID)
	b.createPolls(message.Chat.ID, "Как давно интересуетесь криптовалютой?")
	return
}

func (b *Bot) createPolls(id int64, i string) {
	var err error
	switch i {
	case "Как давно интересуетесь криптовалютой?":
		err = b.createPoll(&tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: id,
			},
			Question: "Как давно интересуетесь криптовалютой?",
			Options: []string{
				"Менее 1 года",
				"Более 1 года",
				"2 и более",
				"5 и более",
			},
			IsAnonymous: true,
		})
	case "Дополнительные интересы":
		err = b.createPoll(&tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: id,
			},
			Question: "Дополнительные интересы",
			Options: []string{
				"Криптовалюты",
				"Маркетинг, реклама, PR",
				"Программирование, IT",
				"Создатель контента",
			},
			IsAnonymous: true,
		})
	case "Программирование, IT":
		err = b.createPoll(&tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: id,
			},
			Question: "Укажите ваше направление?",
			Options: []string{
				"Web Developer",
				"Software developer",
				"Android/iOS developer",
				"Game development",
				"Web3/Crypto development",
				"DevOps",
				"Team manager",
				"QA",
				"Другое",
			},
			AllowsMultipleAnswers: true,
			IsAnonymous:           true,
		})
	case "Криптовалюты":
		err = b.createPoll(&tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: id,
			},
			Question: "Чем интересуетесь в этом направлении?",
			Options: []string{
				"Трейдинг",
				"Инвестирование",
				"Блокчейн-разработка",
				"NFT",
				"Play-to-earn",
				"Ноды, тестнеты",
			},
			AllowsMultipleAnswers: true,
			IsAnonymous:           true,
		})
	case "Маркетинг, реклама, PR":
		err = b.createPoll(&tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: id,
			},
			Question: "Чем занимаетесь в сфере маркетинга?",
			Options: []string{
				"Медийная реклама",
				"Контекст",
				"Таргет",
				"SMM",
				"Продажи",
			},
			AllowsMultipleAnswers: true,
			IsAnonymous:           true,
		})
	case "Создатель контента":
		err = b.createPoll(&tgbotapi.SendPollConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: id,
			},
			Question: "Вы?",
			Options: []string{
				"Блогер",
				"Влогер",
				"Подкастер",
				"Художник",
				"Другое",
			},
			IsAnonymous: true,
		})
	}

	if err != nil {
		b.logger.Errorf("Error send message: %v", err)
	}
}
