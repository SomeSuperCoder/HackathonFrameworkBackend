package main

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/internal/bot"
)

func main() {
	b := bot.NewBot()
	b.Start(context.Background())
}
