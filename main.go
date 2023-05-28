package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/c0rwin/jira/service"
	"github.com/c0rwin/jira/telegram"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SDBOT")
	viper.SetConfigFile(".env") // Optional: If you have a .env file

	// Read environment variables
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file or environment variables: ", err)
	}

	// Access environment variables using Viper
	jiraURL := viper.GetString("JIRA_URL")
	jiraToken := viper.GetString("JIRA_TOKEN")
	username := viper.GetString("JIRA_USERNAME")
	botAPIKey := viper.GetString("BOT_API_KEY")
	projectKey := viper.GetString("PROJECT_KEY")

	bot, err := telegram.NewBot(botAPIKey,
		service.NewJira(jiraURL+"/rest/api/2", username, jiraToken, projectKey))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
		<-interruptChannel
		cancel()
	}()

	bot.Negotiate(ctx)
}
