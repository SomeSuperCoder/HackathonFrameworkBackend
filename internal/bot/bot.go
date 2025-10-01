package bot

import (
	"context"
	"os"
	"sync"

	statemachine "github.com/SomeSuperCoder/global-chat/internal/bot/state_machine"
	"github.com/SomeSuperCoder/global-chat/repository"
	"github.com/SomeSuperCoder/global-chat/utils"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Bot struct {
	Bot        *telego.Bot
	Handler    *th.BotHandler
	State      *statemachine.BotState
	StateMutex *sync.RWMutex
	UserRepo   *repository.UserRepo
}

func (b *Bot) Start() {
	defer func() { _ = b.Handler.Stop() }()
	err := b.Handler.Start()
	utils.CheckErrorDeadly(err, "Failed to start bot Handler")
}

func NewBot() *Bot {
	newBot := &Bot{}
	ctx := context.Background()

	// Load .env
	err := godotenv.Load()
	utils.CheckErrorDeadly(err, "Failedd to load .env")
	botToken := os.Getenv("TELEGRAM_TOKEN")

	// Create newBot.Bot
	newBot.Bot, err = telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	utils.CheckErrorDeadly(err, "Faileed to create newBot.Bot")

	// Create handler
	updates, _ := newBot.Bot.UpdatesViaLongPolling(ctx, nil)
	defer newBot.Bot.StopPoll(ctx, nil)
	newBot.Handler, _ = th.NewBotHandler(newBot.Bot, updates)

	// Connect to MongoDB
	connectionString := "mongodb://localhost:27017"
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString))
	utils.CheckErrorDeadly(err, "Failed to conneect to MongoDB")
	defer client.Disconnect(ctx)
	database := client.Database("hackathonframework")

	// Init database repos
	newBot.UserRepo = &repository.UserRepo{
		Database: database,
	}

	// Init state manager
	newBot.State = statemachine.NewBotState()
	newBot.StateMutex = &sync.RWMutex{}

	return newBot
}
