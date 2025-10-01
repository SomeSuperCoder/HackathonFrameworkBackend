package main

import (
	"context"
	"errors"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/SomeSuperCoder/global-chat/models"
	"github.com/SomeSuperCoder/global-chat/repository"
	botstates "github.com/SomeSuperCoder/global-chat/services/bot/bot_states"
	statemachine "github.com/SomeSuperCoder/global-chat/state_machine"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	STATE_NONE statemachine.StateUnit = iota
	STATE_ENTER_NAME
	STATE_ENTER_BIRTHDATE
)

// Regexp patterns
const NAME_PATTERN = `^[А-ЯЁ][а-яё]+(\s+[А-ЯЁ][а-яё]+){1,2}$`
const BIRTHDATE_PATTERN = `^(0[1-9]|[12][0-9]|3[01])\.(0[1-9]|1[0-2])\.(19|20)\d{2}$`

func CheckErrorDeadly(err error, message string) {
	if err != nil {
		logrus.Fatalf("%v: %v", message, err.Error())
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	CheckErrorDeadly(err, "Failedd to load .env")

	botToken := os.Getenv("TELEGRAM_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	CheckErrorDeadly(err, "Faileed to create bot")

	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)
	defer bot.StopPoll(ctx, nil)

	bh, _ := th.NewBotHandler(bot, updates)

	// Connect to MongoDB
	connectionString := "mongodb://localhost:27017"
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString))
	CheckErrorDeadly(err, "Failed to conneect to MongoDB")
	defer client.Disconnect(ctx)
	database := client.Database("hackathonframework")

	// Init database repos
	userRepo := repository.UserRepo{
		Database: database,
	}

	// Init state manager
	botState := statemachine.NewBotState()
	botStateMutex := sync.RWMutex{}

	// Start command
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Check if user has an account
		user, err := userRepo.GetUserByUsername(ctx, update.Message.From.Username)
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
			bot.SendMessage(ctx, message)
			return nil
		} else if err != nil {
			bot.SendMessage(ctx, tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Ошибка базы данных",
			))
			return err
		}

		bot.SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID),
			"%v, добро пожаловать в личный кабинет!", user.Name,
		))

		return nil
	}, th.CommandEqual("start"))

	// Handle register callback
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		botStateMutex.Lock()
		defer botStateMutex.Unlock()

		botState.SetState(statemachine.StateKey(update.CallbackQuery.From.ID), STATE_ENTER_NAME, botstates.RegisterState{})
		bot.SendMessage(ctx, tu.Message(
			tu.ID(update.CallbackQuery.From.ID),
			"Как вас зовут? (ФИО)",
		))

		return nil
	}, th.CallbackDataEqual("register"))

	// Handle STATE_ENTER_NAME
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Mutex stuff
		botStateMutex.Lock()
		defer botStateMutex.Unlock()

		// Get state
		currentState := botState.GetState(statemachine.StateKey(update.Message.From.ID))
		// Update state
		data, _ := currentState.Data.(botstates.RegisterState)
		data.Name = update.Message.Text
		// Set state
		botState.SetState(statemachine.StateKey(update.Message.From.ID), STATE_ENTER_BIRTHDATE, data)

		bot.SendMessage(ctx, tu.Message(
			tu.ID(update.Message.From.ID),
			"Введите вашу дату рождения в формате `31.12.2025`",
		).WithParseMode("MarkdownV2"))

		return nil
	}, func(ctx context.Context, update telego.Update) bool {
		botStateMutex.RLock()
		defer botStateMutex.RUnlock()
		stateUnit := botState.GetState(statemachine.StateKey(update.Message.From.ID)).State
		return stateUnit == STATE_ENTER_NAME
	}, th.TextMatches(regexp.MustCompile(NAME_PATTERN)))

	// Handle STATE_ENTER_BIRTHDATE
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Mutex stuff
		botStateMutex.Lock()
		defer botStateMutex.Unlock()

		// Get state
		currentState := botState.GetState(statemachine.StateKey(update.Message.From.ID))
		// Update state
		data, _ := currentState.Data.(botstates.RegisterState)
		parsedBirthDate, _ := time.Parse(update.Message.Text, "02.01.2006")
		data.Birthdate = parsedBirthDate
		// Set state
		botState.SetState(statemachine.StateKey(update.Message.From.ID), STATE_NONE, nil)

		// Query DB
		newUser := &models.User{
			Username: update.Message.From.Username,
			ChatID:   update.Message.Chat.ID,

			Name:      data.Name,
			Birthdate: data.Birthdate,
			Role:      models.Participant,
			CratedAt:  time.Now(),
		}
		err := userRepo.CreateUser(ctx, newUser)

		if err != nil {

		}

		bot.SendMessage(ctx, tu.Message(
			tu.ID(update.Message.From.ID),
			"Вы были успешно зарегестрированы! 🎉\nНажмите /start чтобы перезапустить бота",
		))

		return nil
	}, func(ctx context.Context, update telego.Update) bool {
		botStateMutex.RLock()
		defer botStateMutex.RUnlock()
		stateUnit := botState.GetState(statemachine.StateKey(update.Message.From.ID)).State
		return stateUnit == STATE_ENTER_BIRTHDATE
	}, th.TextMatches(regexp.MustCompile(BIRTHDATE_PATTERN)))

	defer func() { _ = bh.Stop() }()
	err = bh.Start()
	CheckErrorDeadly(err, "Failed to start bot handler")
}
