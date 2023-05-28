package telegram

import (
	"context"
	"fmt"

	"github.com/c0rwin/jira/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot represents a Telegram bot
type Bot struct {
	tg   *tgbotapi.BotAPI
	jira service.Jira
}

// NewBot creates a new Bot instance
func NewBot(token string, jira service.Jira) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{tg: bot, jira: jira}, nil
}

// Negotiate negotiates with Telegram API
func (b *Bot) Negotiate(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	fmt.Println("Negotiating with Telegram API...")
	updates, err := b.tg.GetUpdatesChan(updateConfig)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Negotiation with Telegram API is done")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			chatID := update.Message.Chat.ID
			text := update.Message.Text

			msg := tgbotapi.NewMessage(chatID, "")

			var issues service.Issues
			var err error

			switch text {
			case "/start":
				msg.Text = "Hello, I'm Jira Bot. You can type /recent or /all to get list of recent or all open tasks"
				_, err = b.tg.Send(msg)
				if err != nil {
					panic(err)
				}
			case "/ping":
				msg.Text = "pong"
				_, err = b.tg.Send(msg)
				if err != nil {
					panic(err)
				}
			case "/recent":
				issues, err = b.jira.GetRecentOpenTasks()
				if err != nil {
					msg.Text = "Error: " + err.Error()
				}
			case "/all":
				issues, err = b.jira.GetAllOpenTasks()
				if err != nil {
					msg.Text = "Error: " + err.Error()
				}
			default:
				msg.Text = "I don't know that command"
				_, err = b.tg.Send(msg)
				if err != nil {
					panic(err)
				}
			}

			for _, issue := range issues {
				msg.Text = issue.String()
				reply := tgbotapi.NewMessage(chatID, issue.String())
				_, err = b.tg.Send(reply)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
