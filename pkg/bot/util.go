package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
		_, err = b.database.Exec("UPDATE users SET user_name = ?, user_first_name = ?, user_last_name = ?, referaled_by = ? WHERE telegram_id = ?", username, firstName, lastName, referral, userId)
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
