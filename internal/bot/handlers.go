package bot

import (
	"context"
	"regexp"

	botregexps "github.com/SomeSuperCoder/global-chat/internal/bot/regexps"
	statemachine "github.com/SomeSuperCoder/global-chat/internal/bot/state_machine"
	stateunits "github.com/SomeSuperCoder/global-chat/internal/bot/state_units"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func (b *Bot) registerHandlers() {
	b.Handler.Handle(b.StartCommand, th.CommandEqual("start"))
	// Register callback
	b.Handler.Handle(b.Register, th.CallbackDataEqual("register"))

	// Handle STATE_ENTER_NAME
	b.Handler.Handle(b.EnterName, func(ctx context.Context, update telego.Update) bool {
		b.StateMutex.RLock()
		defer b.StateMutex.RUnlock()
		stateUnit := b.State.GetState(statemachine.StateKey(update.Message.From.ID)).State
		return stateUnit == stateunits.STATE_ENTER_NAME
	}, th.TextMatches(regexp.MustCompile(botregexps.NAME_PATTERN)))

	// Handle STATE_ENTER_BIRTHDATE
	b.Handler.Handle(b.EnterName, func(ctx context.Context, update telego.Update) bool {
		b.StateMutex.RLock()
		defer b.StateMutex.RUnlock()
		stateUnit := b.State.GetState(statemachine.StateKey(update.Message.From.ID)).State
		return stateUnit == stateunits.STATE_ENTER_BIRTHDATE
	}, th.TextMatches(regexp.MustCompile(botregexps.BIRTHDATE_PATTERN)))

}
