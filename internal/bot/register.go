package bot

import (
	statemachine "github.com/SomeSuperCoder/global-chat/internal/bot/state_machine"
	stateunits "github.com/SomeSuperCoder/global-chat/internal/bot/state_units"
	botstates "github.com/SomeSuperCoder/global-chat/internal/bot/states"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (b *Bot) Register(ctx *th.Context, update telego.Update) error {
	b.StateMutex.Lock()
	defer b.StateMutex.Unlock()

	b.State.SetState(statemachine.StateKey(update.CallbackQuery.From.ID), stateunits.STATE_ENTER_NAME, botstates.RegisterState{})
	b.Bot.SendMessage(ctx, tu.Message(
		tu.ID(update.CallbackQuery.From.ID),
		"Как вас зовут? (ФИО)",
	))

	return nil
}
