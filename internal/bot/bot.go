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
	client     *mongo.Client
	database   *mongo.Database
	Bot        *telego.Bot
	Handler    *th.BotHandler
	State      *statemachine.BotState
	StateMutex *sync.RWMutex
	UserRepo   *repository.UserRepo
}

func NewBot() *Bot {
	bot := &Bot{}

	return bot
}

func (b *Bot) Start(ctx context.Context) {
	// Load .env
	err := godotenv.Load()
	utils.CheckErrorDeadly(err, "Failed to load .env")
	botToken := os.Getenv("TELEGRAM_TOKEN")

	// Create a bot
	b.Bot, err = telego.NewBot(botToken, telego.WithDefaultLogger(false, true))
	utils.CheckErrorDeadly(err, "Faileed to create b.Bot")

	// Create handler
	updates, _ := b.Bot.UpdatesViaLongPolling(ctx, nil)
	defer b.Bot.StopPoll(ctx, nil)
	b.Handler, _ = th.NewBotHandler(b.Bot, updates)

	// Register handlers
	b.registerHandlers()

	// Connect to MongoDB
	connectionString := "mongodb://localhost:27017"
	b.client, err = mongo.Connect(options.Client().ApplyURI(connectionString))
	utils.CheckErrorDeadly(err, "Failed to conneect to MongoDB")
	defer b.client.Disconnect(ctx)
	b.database = b.client.Database("hackathonframework")

	// Init database repos
	b.UserRepo = repository.NewUserRepo(b.database)

	// Init state manager
	b.State = statemachine.NewBotState()
	b.StateMutex = &sync.RWMutex{}

	defer func() { _ = b.Handler.Stop() }()
	err = b.Handler.Start()
	utils.CheckErrorDeadly(err, "Failed to start bot Handler")
}
