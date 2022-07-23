package bot

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
