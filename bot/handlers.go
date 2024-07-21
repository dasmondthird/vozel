// bot/handlers.go
package bot

import (
	"strconv"
	"time"

	"vozel/database"

	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

func (b *Bot) SetupHandlers() {
	b.Bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Я ваш VPN бот.", b.Menus.MainMenu)
	})

	b.Bot.Handle("/menu", func(c tele.Context) error {
		return c.Send("Главное меню:", b.Menus.MainMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "myvpn"}, func(c tele.Context) error {
		return c.Send("Выберите VPN протокол:", b.Menus.VpnMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "outline"}, func(c tele.Context) error {
		return c.Send("Выберите сервер:", b.Menus.ServerMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "server1"}, func(c tele.Context) error {
		b.AddTariffButtons()
		return c.Send("Выберите тариф:", b.Menus.TariffMenu)
	})

	for name := range b.Tariffs {
		localName := name
		b.Bot.Handle(&tele.Btn{Unique: "tariff_" + localName}, func(c tele.Context) error {
			return b.handleTariffSelection(c, localName, "Сервер 1")
		})
	}

	b.Bot.Handle(&tele.Btn{Unique: "android"}, func(c tele.Context) error {
		return b.handleDeviceSelection(c, "android")
	})

	b.Bot.Handle(&tele.Btn{Unique: "iphone"}, func(c tele.Context) error {
		return b.handleDeviceSelection(c, "iphone")
	})

	b.Bot.Handle(&tele.Btn{Unique: "backtomain"}, func(c tele.Context) error {
		return c.Send("Главное меню:", b.Menus.MainMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "backtovpn"}, func(c tele.Context) error {
		return c.Send("Выберите VPN протокол:", b.Menus.VpnMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "backtoservers"}, func(c tele.Context) error {
		return c.Send("Выберите сервер:", b.Menus.ServerMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "backtotariffs"}, func(c tele.Context) error {
		b.AddTariffButtons()
		return c.Send("Выберите тариф:", b.Menus.TariffMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "payment"}, func(c tele.Context) error {
		return c.Send("Выберите сумму пополнения:", b.Menus.PaymentMenu)
	})

	b.Bot.Handle(&tele.Btn{Unique: "pay_100"}, func(c tele.Context) error {
		return b.handlePayment(c, 100)
	})

	b.Bot.Handle(&tele.Btn{Unique: "pay_300"}, func(c tele.Context) error {
		return b.handlePayment(c, 300)
	})

	b.Bot.Handle(&tele.Btn{Unique: "pay_1000"}, func(c tele.Context) error {
		return b.handlePayment(c, 1000)
	})

	b.Bot.Handle(&tele.Btn{Unique: "balance"}, func(c tele.Context) error {
		return b.handleBalance(c)
	})

	b.Bot.Handle(&tele.Btn{Unique: "mykey"}, func(c tele.Context) error {
		return b.handleMyKey(c)
	})

	b.Bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Команда не распознана. Используйте /menu для просмотра доступных команд.", b.Menus.MainMenu)
	})
}

func (b *Bot) handleTariffSelection(c tele.Context, tariffName, location string) error {
	return b.DB.Transaction(func(tx *gorm.DB) error {
		user, err := b.getOrCreateUser(c, tx)
		if err != nil {
			return c.Send("Ошибка при получении данных пользователя.", b.Menus.MainMenu)
		}

		selectedTariff, exists := b.Tariffs[tariffName]
		if !exists {
			return c.Send("Некорректный тариф.", b.Menus.MainMenu)
		}

		if user.Balance < selectedTariff.Price {
			return c.Send("Недостаточно средств на балансе. Пополните счет.", b.Menus.MainMenu)
		}

		user.Balance -= selectedTariff.Price
		user.VPNKey = generateVPNKey()
		user.Location = location
		user.Expiry = time.Now().Add(selectedTariff.Duration)
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		return c.Send("Тариф "+tariffName+" выбран. Теперь выберите устройство:", b.Menus.DeviceMenu)
	})
}

func generateVPNKey() string {
	return "UniqueVPNKey123456"
}

func (b *Bot) handleDeviceSelection(c tele.Context, device string) error {
	user, err := b.getOrCreateUser(c, b.DB)
	if err != nil {
		return c.Send("Ошибка при получении данных пользователя.", b.Menus.MainMenu)
	}

	if user.VPNKey == "" {
		return c.Send("У вас нет активного VPN ключа. Пожалуйста, выберите тариф и получите ключ.", b.Menus.MainMenu)
	}

	var instruction string
	if device == "android" {
		instruction = "Установите приложение Outline (https://play.google.com/store/apps/details?id=org.outline.android.client)\n\n" +
			"1. Скопируйте ключ и добавьте его в приложение.\n" +
			"2. Подключитесь к VPN.\n\n" +
			"Ваш ключ: " + user.VPNKey
	} else if device == "iphone" {
		instruction = "Установите приложение Outline (https://apps.apple.com/app/outline-app/id1356177741)\n\n" +
			"1. Скопируйте ключ и добавьте его в приложение.\n" +
			"2. Подключитесь к VPN.\n\n" +
			"Ваш ключ: " + user.VPNKey
	}

	return c.Send(instruction, b.Menus.MainMenu)
}

func (b *Bot) handlePayment(c tele.Context, amount int) error {
	return b.DB.Transaction(func(tx *gorm.DB) error {
		user, err := b.getOrCreateUser(c, tx)
		if err != nil {
			return c.Send("Ошибка при получении данных пользователя.", b.Menus.MainMenu)
		}

		user.Balance += amount
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		return c.Send("Ваш баланс пополнен на "+strconv.Itoa(amount)+" руб.\nТекущий баланс: "+strconv.Itoa(user.Balance)+" руб.", b.Menus.MainMenu)
	})
}

func (b *Bot) handleBalance(c tele.Context) error {
	user, err := b.getOrCreateUser(c, b.DB)
	if err != nil {
		return c.Send("Ошибка при получении данных пользователя.", b.Menus.MainMenu)
	}
	return c.Send("Ваш баланс: "+strconv.Itoa(user.Balance)+" руб.\nПополните счет для продолжения.", b.Menus.MainMenu)
}

func (b *Bot) handleMyKey(c tele.Context) error {
	user, err := b.getOrCreateUser(c, b.DB)
	if err != nil {
		return c.Send("Ошибка при получении данных пользователя.", b.Menus.MainMenu)
	}

	if time.Now().After(user.Expiry) {
		return c.Send("Ваш ключ истек. Пополните баланс для продления доступа.", b.Menus.MainMenu)
	}

	info := "Ваш текущий ключ доступа:\n\n" +
		"Локация: " + user.Location + "\n" +
		"Ключ: " + user.VPNKey + "\n" +
		"Начало действия: " + user.Expiry.Add(-b.Tariffs[user.Location].Duration).Format("02.01.2006 15:04") + "\n" +
		"Истечение срока: " + user.Expiry.Format("02.01.2006 15:04") + "\n\n" +
		"Для продления доступа, пожалуйста, пополните баланс."

	return c.Send(info, b.Menus.MainMenu)
}

func (b *Bot) getOrCreateUser(c tele.Context, tx *gorm.DB) (*database.User, error) {
	var user database.User
	if err := tx.First(&user, c.Sender().ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = database.User{
				ID:           c.Sender().ID,
				Username:     c.Sender().Username,
				Registered:   true,
				RegisteredAt: time.Now(),
			}
			if err := tx.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &user, nil
}
