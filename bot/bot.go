package bot

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Tariff struct {
	Name     string
	Price    int
	Duration time.Duration
}

type User struct {
	ID           int64 `gorm:"primaryKey"`
	Username     string
	Balance      int
	VPNKey       string
	Location     string
	Expiry       time.Time
	Registered   bool
	RegisteredAt time.Time
}

type Bot struct {
	Bot     *tele.Bot
	DB      *gorm.DB
	Menus   Menus
	Tariffs map[string]Tariff
}

func NewBot() *Bot {
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	db.AutoMigrate(&User{})

	return &Bot{
		DB: db,
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
	tariffs := []Tariff{
		{Name: "1 день", Price: 100, Duration: 24 * time.Hour},
		{Name: "3 дня", Price: 300, Duration: 3 * 24 * time.Hour},
		{Name: "1 неделя", Price: 700, Duration: 7 * 24 * time.Hour},
	}

	b.Tariffs = make(map[string]Tariff)
	for _, tariff := range tariffs {
		b.Tariffs[tariff.Name] = tariff
	}
}

func (b *Bot) Start() {
	b.Bot.Start()
}
