package bot

import (
	"context"

	statemachine "github.com/SomeSuperCoder/global-chat/internal/bot/state_machine"
	stateunits "github.com/SomeSuperCoder/global-chat/internal/bot/state_units"
	botstates "github.com/SomeSuperCoder/global-chat/internal/bot/states"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
)

func (b *Bot) EnterName(ctx *th.Context, update telego.Update) error {
	// Mutex stuff
	b.StateMutex.Lock()
	defer b.StateMutex.Unlock()

	// Get state
	currentState := b.State.GetState(statemachine.StateKey(update.Message.From.ID))
	// Update state
	data, _ := currentState.Data.(botstates.RegisterState)
	data.Name = update.Message.Text
	// Set state
	b.State.SetState(statemachine.StateKey(update.Message.From.ID), stateunits.STATE_ENTER_BIRTHDATE, data)

	b.Bot.SendMessage(ctx, tu.Message(
		tu.ID(update.Message.From.ID),
		"Введите вашу дату рождения в формате `31.12.2025`",
	).WithParseMode("MarkdownV2"))

	return nil
}

func (b *Bot) InvalidEnterName(ctx *th.Context, update telego.Update) error {
	logrus.Info("Invalid enter name!!!!!!!")
	_, err := b.Bot.SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Некорректное ФИО\\!\nПример ФИО: `Иванов Иван Иванович`\nЕсли у вас нет отчества, вводите фамилию и имя в формате:\n`Иванов Иван`",
	).WithParseMode("MarkdownV2"))
	return err
}

func (b *Bot) EnterNamePredicate(ctx context.Context, update telego.Update) bool {
	b.StateMutex.RLock()
	defer b.StateMutex.RUnlock()
	stateUnit := b.State.GetState(statemachine.StateKey(update.Message.From.ID)).State
	return stateUnit == stateunits.STATE_ENTER_NAME
}
