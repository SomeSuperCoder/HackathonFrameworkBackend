package main

import (
	"context"

	"github.com/SomeSuperCoder/global-chat/internal/bot"
)

func main() {
	b := &bot.Bot{}
	b.Start(context.Background())
}
