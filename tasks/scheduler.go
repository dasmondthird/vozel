package tasks

import (
	"log"
	"time"

	"telegram-vpn-bot/bot"

	"github.com/go-co-op/gocron"
	tele "gopkg.in/telebot.v3"
)

type Scheduler struct {
	bot *bot.Bot
}

func NewScheduler(bot *bot.Bot) *Scheduler {
	return &Scheduler{bot: bot}
}

func (s *Scheduler) Start() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().Do(s.checkAndNotifyExpiry)
	scheduler.StartAsync()
}

func (s *Scheduler) checkAndNotifyExpiry() {
	var users []bot.User
	expiryThreshold := time.Now().Add(3 * 24 * time.Hour)
	if err := s.bot.DB.Where("expiry < ?", expiryThreshold).Find(&users).Error; err != nil {
		log.Printf("Error fetching users for expiry check: %v", err)
		return
	}

	for _, user := range users {
		if time.Now().Before(user.Expiry) {
			msg := "Ваш VPN ключ истекает " + user.Expiry.Format("02.01.2006 15:04") + ". Пожалуйста, продлите подписку, чтобы избежать перерыва в услуге."
			_, err := s.bot.Bot.Send(tele.ChatID(user.ID), msg)
			if err != nil {
				log.Printf("Error sending expiry notification to user %d: %v", user.ID, err)
			}
		}
	}
}
