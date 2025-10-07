package main

import (
	"context"
	"log"

	"github.com/SomeSuperCoder/global-chat/application"
)

func main() {
	app := application.New()

	err := app.Start(context.Background())
	if err != nil {
		log.Fatalf("failed to start app: %v", err.Error())
	}
}
