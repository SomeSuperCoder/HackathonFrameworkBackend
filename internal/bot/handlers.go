package bot

import (
	"regexp"

	botregexps "github.com/SomeSuperCoder/global-chat/internal/bot/regexps"
	th "github.com/mymmrac/telego/telegohandler"
)

func (b *Bot) registerHandlers() {
	b.Handler.Handle(b.StartCommand, th.CommandEqual("start"))

	// Register callback
	b.Handler.Handle(b.Register, th.CallbackDataEqual("register"))

	// Handle STATE_ENTER_NAME
	b.Handler.Handle(b.EnterName, b.EnterNamePredicate, th.TextMatches(regexp.MustCompile(botregexps.NAME_PATTERN)))
	b.Handler.Handle(b.InvalidEnterName, b.EnterNamePredicate, th.AnyMessageWithText())

	// Handle STATE_ENTER_BIRTHDATE
	b.Handler.Handle(b.EnterBirthdate, b.EnterBirthdatePredicate, th.TextMatches(regexp.MustCompile(botregexps.BIRTHDATE_PATTERN)))
	b.Handler.Handle(b.InvalidEnterBirthdate, b.EnterBirthdatePredicate, th.AnyMessageWithText())
}
