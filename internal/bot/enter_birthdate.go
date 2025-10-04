package bot

import (
	"context"
	"time"

	"github.com/SomeSuperCoder/global-chat/internal"
	statemachine "github.com/SomeSuperCoder/global-chat/internal/bot/state_machine"
	stateunits "github.com/SomeSuperCoder/global-chat/internal/bot/state_units"
	botstates "github.com/SomeSuperCoder/global-chat/internal/bot/states"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func (b *Bot) EnterBirthdate(ctx *th.Context, update telego.Update) error {
	// Mutex stuff
	b.StateMutex.Lock()
	defer b.StateMutex.Unlock()

	// Get state
	currentState := b.State.GetState(statemachine.StateKey(update.Message.From.ID))
	// Update state
	data, _ := currentState.Data.(botstates.RegisterState)
	parsedBirthDate, _ := time.Parse("02.01.2006", update.Message.Text)
	data.Birthdate = parsedBirthDate
	// Set state
	b.State.SetState(statemachine.StateKey(update.Message.From.ID), stateunits.STATE_NONE, nil)

	// Query DB
	newUser := &models.User{
		Username: update.Message.From.Username,
		ChatID:   update.Message.Chat.ID,

		Name:      data.Name,
		Birthdate: data.Birthdate,
		Role:      models.Participant,
		Team:      internal.UndefinedObjectID,
	}
	err := b.UserRepo.Create(ctx, newUser)

	if err != nil {

	}

	b.Bot.SendMessage(ctx, tu.Message(
		tu.ID(update.Message.From.ID),
		"Вы были успешно зарегестрированы! 🎉\nНажмите /start чтобы перезапустить бота",
	))

	return nil
}

func (b *Bot) InvalidEnterBirthdate(ctx *th.Context, update telego.Update) error {
	b.Bot.SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Некорректная дата рождения\\! Пример даты рождения: `31.12.2025`",
	).WithParseMode("MarkdownV2"))
	return nil
}

func (b *Bot) EnterBirthdatePredicate(ctx context.Context, update telego.Update) bool {
	b.StateMutex.RLock()
	defer b.StateMutex.RUnlock()
	stateUnit := b.State.GetState(statemachine.StateKey(update.Message.From.ID)).State
	return stateUnit == stateunits.STATE_ENTER_BIRTHDATE
}
