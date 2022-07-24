package bot

import (
	"database/sql"
	"time"
)

var (
	waitInstagram     = map[int64]struct{}{}
	waitWhyYouCanHelp = map[int64]struct{}{}
	nextPoll          = map[int64]string{}
)

// Buttons
const (
	buttonsAbout    = "О проекте"
	buttonsReferral = "Реферальная программа"
)

func (b *Bot) cleanWaiting(id int64) {
	delete(waitInstagram, id)
	delete(nextPoll, id)
}

type User struct {
	RegisteredAt time.Time
	TelegramID   string
	Username     sql.NullString
	LastName     sql.NullString
	Instagram    sql.NullString
	Twitter      sql.NullString
	WantHelp     sql.NullString
	ReferredBy   sql.NullString
	FirstName    sql.NullString
	TotalInvites sql.NullString
}

type Polls struct {
	QuestionOne   string
	QuestionTwo   string
	QuestionThree string
	Answers       string
	User
}
