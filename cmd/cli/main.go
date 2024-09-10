package main

import (
	"luganodes/internal/database"
	"luganodes/internal/logger"
	"luganodes/internal/services"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.Init()
	database.Init()

	notifier := services.NewNotifier()
	tracker := services.NewTracker("0x00000000219ab540356cBB839Cbe05303d7705Fa", "0x10", time.Minute, notifier)

	go notifier.InitWebhook()
	go tracker.Start()

	select {}
}
