// bot/bot.go
package bot

import (
	"log"
	"os"
	"time"

	"vozel/database"

	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type Bot struct {
	Bot     *tele.Bot
	DB      *gorm.DB
	Menus   Menus
	Tariffs map[string]database.Tariff
}

func NewBot() *Bot {
	return &Bot{
		DB: database.DB,
	}
}

func (b *Bot) InitTelegramBot() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is missing")
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	var err error
	b.Bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	b.SetupMenus()
}

func (b *Bot) LoadTariffs() {
	var tariffs []database.Tariff
	if err := b.DB.Find(&tariffs).Error; err != nil {
		log.Printf("Error loading tariffs: %v", err)
		return
	}
	b.Tariffs = make(map[string]database.Tariff)
	for _, tariff := range tariffs {
		b.Tariffs[tariff.Name] = tariff
	}
}

func (b *Bot) Start() {
	b.Bot.Start()
}
