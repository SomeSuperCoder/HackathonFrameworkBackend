package bot

import (
	"errors"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (b *Bot) StartCommand(ctx *th.Context, update telego.Update) error {
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
	botUser, _ := b.Bot.GetMe(ctx)
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Открыть мини-приложение").WithURL(fmt.Sprintf("t.me/%v/controlcenter", botUser.Username)),
		),
	)

	b.Bot.SendMessage(ctx, tu.Messagef(
		tu.ID(update.Message.Chat.ID),
		"%v, добро пожаловать в личный кабинет!", user.Name,
	).WithReplyMarkup(inlineKeyboard))

	return nil
}
