// main.go
package main

import (
	"log"

	"vozel/bot"
	"vozel/database"
	"vozel/tasks"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	database.InitDB()
	botInstance := bot.NewBot()
	botInstance.LoadTariffs()
	botInstance.InitTelegramBot()
	botInstance.SetupHandlers()
	botInstance.Start()

	scheduler := tasks.NewScheduler(botInstance)
	scheduler.Start()
}
