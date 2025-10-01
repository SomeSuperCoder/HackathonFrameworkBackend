package main

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/SomeSuperCoder/global-chat/internal/bot"
	botregexps "github.com/SomeSuperCoder/global-chat/internal/bot/regexps"
	statemachine "github.com/SomeSuperCoder/global-chat/internal/bot/state_machine"
	botstates "github.com/SomeSuperCoder/global-chat/internal/bot/states"
	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/utils"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	STATE_NONE statemachine.StateUnit = iota
	STATE_ENTER_NAME
	STATE_ENTER_BIRTHDATE
)

func main() {
	b := bot.NewBot()

	// Start command
	b.Handler.Handle(func(ctx *th.Context, update telego.Update) error {
		// Check if user has an account
		user, err := b.UserRepo.GetUserByUsername(ctx, update.Message.From.Username)
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Handle the case where the user does not have an account
			inlineKeyboard := tu.InlineKeyboard(
				tu.InlineKeyboardRow(
					tu.InlineKeyboardButton("Зарегестироваться на хакатон").WithCallbackData("register"),
				),
			)
			message := tu.Messagef(
				tu.ID(update.Message.Chat.ID),
				"Добро пожаловать, %v! Вы пока что не зарегестированы на хакатон 😭 Пожалуйста, пройдите регестрацию по кнопке ниже",
				update.Message.Chat.FirstName,
			).WithReplyMarkup(inlineKeyboard)
			b.Bot.SendMessage(ctx, message)
			return nil
		} else if err != nil {
			b.Bot.SendMessage(ctx, tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Ошибка базы данных",
			))
			return err
		}

		b.Bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"%v, добро пожаловать в личный кабинет!", user.Name,
		))

		return nil
	}, th.CommandEqual("start"))

	// Handle register callback
	b.Handler.Handle(func(ctx *th.Context, update telego.Update) error {
		b.StateMutex.Lock()
		defer b.StateMutex.Unlock()

		b.State.SetState(statemachine.StateKey(update.CallbackQuery.From.ID), STATE_ENTER_NAME, botstates.RegisterState{})
		b.Bot.SendMessage(ctx, tu.Message(
			tu.ID(update.CallbackQuery.From.ID),
			"Как вас зовут? (ФИО)",
		))

		return nil
	}, th.CallbackDataEqual("register"))

	// Handle STATE_ENTER_NAME
	b.Handler.Handle(func(ctx *th.Context, update telego.Update) error {
		// Mutex stuff
		b.StateMutex.Lock()
		defer b.StateMutex.Unlock()

		// Get state
		currentState := b.State.GetState(statemachine.StateKey(update.Message.From.ID))
		// Update state
		data, _ := currentState.Data.(botstates.RegisterState)
		data.Name = update.Message.Text
		// Set state
		b.State.SetState(statemachine.StateKey(update.Message.From.ID), STATE_ENTER_BIRTHDATE, data)

		b.Bot.SendMessage(ctx, tu.Message(
			tu.ID(update.Message.From.ID),
			"Введите вашу дату рождения в формате `31.12.2025`",
		).WithParseMode("MarkdownV2"))

		return nil
	}, func(ctx context.Context, update telego.Update) bool {
		b.StateMutex.RLock()
		defer b.StateMutex.RUnlock()
		stateUnit := b.State.GetState(statemachine.StateKey(update.Message.From.ID)).State
		return stateUnit == STATE_ENTER_NAME
	}, th.TextMatches(regexp.MustCompile(botregexps.NAME_PATTERN)))

	// Handle STATE_ENTER_BIRTHDATE
	b.Handler.Handle(func(ctx *th.Context, update telego.Update) error {
		// Mutex stuff
		b.StateMutex.Lock()
		defer b.StateMutex.Unlock()

		// Get state
		currentState := b.State.GetState(statemachine.StateKey(update.Message.From.ID))
		// Update state
		data, _ := currentState.Data.(botstates.RegisterState)
		parsedBirthDate, _ := time.Parse(update.Message.Text, "02.01.2006")
		data.Birthdate = parsedBirthDate
		// Set state
		b.State.SetState(statemachine.StateKey(update.Message.From.ID), STATE_NONE, nil)

		// Query DB
		newUser := &models.User{
			Username: update.Message.From.Username,
			ChatID:   update.Message.Chat.ID,

			Name:      data.Name,
			Birthdate: data.Birthdate,
			Role:      models.Participant,
			CratedAt:  time.Now(),
		}
		err := b.UserRepo.CreateUser(ctx, newUser)

		if err != nil {

		}

		b.Bot.SendMessage(ctx, tu.Message(
			tu.ID(update.Message.From.ID),
			"Вы были успешно зарегестрированы! 🎉\nНажмите /start чтобы перезапустить бота",
		))

		return nil
	}, func(ctx context.Context, update telego.Update) bool {
		b.StateMutex.RLock()
		defer b.StateMutex.RUnlock()
		stateUnit := b.State.GetState(statemachine.StateKey(update.Message.From.ID)).State
		return stateUnit == STATE_ENTER_BIRTHDATE
	}, th.TextMatches(regexp.MustCompile(botregexps.BIRTHDATE_PATTERN)))

	defer func() { _ = b.Handler.Stop() }()
	err := b.Handler.Start()
	utils.CheckErrorDeadly(err, "Failed to start bot Handler")
}
