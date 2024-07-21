// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User struct to map database columns
type User struct {
	ID         int64 `gorm:"primaryKey"`
	Username   string
	Registered bool
	VPNKey     string
	Location   string
	Expiry     time.Time
	Balance    int
}

// Tariff describes the tariff plans
type Tariff struct {
	Name     string
	Price    int
	Duration time.Duration
}

// Global variables
var (
	mainMenu, vpnMenu, serverMenu, tariffMenu, deviceMenu, paymentMenu *tele.ReplyMarkup
	DB                                                                 *gorm.DB
	tariffs                                                            = map[string]Tariff{
		"1 день":   {Name: "1 день", Price: 100, Duration: 24 * time.Hour},
		"1 неделя": {Name: "1 неделя", Price: 300, Duration: 7 * 24 * time.Hour},
		"1 месяц":  {Name: "1 месяц", Price: 1000, Duration: 30 * 24 * time.Hour},
	}
)

// InitDB initializes the database connection
func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Error migrating the database:", err)
	}
}

// GetUserByID fetches a user by their ID from the database
func GetUserByID(userID int64) (*User, error) {
	var user User
	if err := DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// UpdateUser updates a user's information in the database
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	InitDB()

	// Get the bot token from the environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is missing")
	}

	// Bot settings
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	// Create a new bot
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Setup main menu
	mainMenu = &tele.ReplyMarkup{}
	btnMyVPN := mainMenu.Data("Мой VPN", "myvpn")
	btnMyKey := mainMenu.Data("Мой ключ", "mykey")
	btnBalance := mainMenu.Data("Баланс", "balance")
	btnPayment := mainMenu.Data("Пополнить баланс", "payment")
	mainMenu.Inline(
		mainMenu.Row(btnMyVPN),
		mainMenu.Row(btnMyKey),
		mainMenu.Row(btnBalance),
		mainMenu.Row(btnPayment),
	)

	// Setup VPN protocol menu
	vpnMenu = &tele.ReplyMarkup{}
	btnOutline := vpnMenu.Data("Outline", "outline")
	btnBackToMain := vpnMenu.Data("Назад", "backtomain")
	vpnMenu.Inline(
		vpnMenu.Row(btnOutline),
		vpnMenu.Row(btnBackToMain),
	)

	// Setup server menu
	serverMenu = &tele.ReplyMarkup{}
	btnServer1 := serverMenu.Data("Сервер 1", "server1")
	btnBackToVPN := serverMenu.Data("Назад", "backtovpn")
	serverMenu.Inline(
		serverMenu.Row(btnServer1),
		serverMenu.Row(btnBackToVPN),
	)

	// Setup tariff menu
	tariffMenu = &tele.ReplyMarkup{}
	btnDay := tariffMenu.Data("1 день - 100 руб", "tariff_1day")
	btnWeek := tariffMenu.Data("1 неделя - 300 руб", "tariff_1week")
	btnMonth := tariffMenu.Data("1 месяц - 1000 руб", "tariff_1month")
	btnBackToServers := tariffMenu.Data("Назад", "backtoservers")
	tariffMenu.Inline(
		tariffMenu.Row(btnDay),
		tariffMenu.Row(btnWeek),
		tariffMenu.Row(btnMonth),
		tariffMenu.Row(btnBackToServers),
	)

	// Setup device menu
	deviceMenu = &tele.ReplyMarkup{}
	btnAndroid := deviceMenu.Data("Android", "android")
	btnIphone := deviceMenu.Data("iPhone", "iphone")
	btnBackToTariffs := deviceMenu.Data("Назад", "backtotariffs")
	deviceMenu.Inline(
		deviceMenu.Row(btnAndroid, btnIphone),
		deviceMenu.Row(btnBackToTariffs),
	)

	// Setup payment menu
	paymentMenu = &tele.ReplyMarkup{}
	btnPay100 := paymentMenu.Data("Пополнить на 100 руб", "pay_100")
	btnPay300 := paymentMenu.Data("Пополнить на 300 руб", "pay_300")
	btnPay1000 := paymentMenu.Data("Пополнить на 1000 руб", "pay_1000")
	btnBackToMainFromPayment := paymentMenu.Data("Назад", "backtomain")
	paymentMenu.Inline(
		paymentMenu.Row(btnPay100, btnPay300, btnPay1000),
		paymentMenu.Row(btnBackToMainFromPayment),
	)

	// Handler for /start command
	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Я ваш VPN бот.", mainMenu)
	})

	// Handler for /menu command
	b.Handle("/menu", func(c tele.Context) error {
		return c.Send("Главное меню:", mainMenu)
	})

	// Handlers for inline buttons
	b.Handle(&btnMyVPN, func(c tele.Context) error {
		return c.Send("Выберите VPN протокол:", vpnMenu)
	})

	b.Handle(&btnOutline, func(c tele.Context) error {
		return c.Send("Выберите сервер:", serverMenu)
	})

	b.Handle(&btnServer1, func(c tele.Context) error {
		return c.Send("Выберите тариф:", tariffMenu)
	})

	b.Handle(&btnDay, func(c tele.Context) error {
		return handleTariffSelection(c, "1 день", "Сервер 1")
	})

	b.Handle(&btnWeek, func(c tele.Context) error {
		return handleTariffSelection(c, "1 неделя", "Сервер 1")
	})

	b.Handle(&btnMonth, func(c tele.Context) error {
		return handleTariffSelection(c, "1 месяц", "Сервер 1")
	})

	b.Handle(&btnAndroid, func(c tele.Context) error {
		return handleDeviceSelection(c, "android")
	})

	b.Handle(&btnIphone, func(c tele.Context) error {
		return handleDeviceSelection(c, "iphone")
	})

	b.Handle(&btnBackToMain, func(c tele.Context) error {
		return c.Send("Главное меню:", mainMenu)
	})

	b.Handle(&btnBackToVPN, func(c tele.Context) error {
		return c.Send("Выберите VPN протокол:", vpnMenu)
	})

	b.Handle(&btnBackToServers, func(c tele.Context) error {
		return c.Send("Выберите сервер:", serverMenu)
	})

	b.Handle(&btnBackToTariffs, func(c tele.Context) error {
		return c.Send("Выберите тариф:", tariffMenu)
	})

	b.Handle(&btnPayment, func(c tele.Context) error {
		return c.Send("Выберите сумму пополнения:", paymentMenu)
	})

	b.Handle(&btnBackToMainFromPayment, func(c tele.Context) error {
		return c.Send("Главное меню:", mainMenu)
	})

	b.Handle(&btnPay100, func(c tele.Context) error {
		return handlePayment(c, 100)
	})

	b.Handle(&btnPay300, func(c tele.Context) error {
		return handlePayment(c, 300)
	})

	b.Handle(&btnPay1000, func(c tele.Context) error {
		return handlePayment(c, 1000)
	})

	b.Handle(&btnBalance, func(c tele.Context) error {
		return handleBalance(c)
	})

	b.Handle(&btnMyKey, func(c tele.Context) error {
		return handleMyKey(c)
	})

	// Start the bot
	b.Start()
}

// handleTariffSelection handles the tariff selection
func handleTariffSelection(c tele.Context, tariffName, location string) error {
	user, err := GetUserByID(c.Sender().ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = &User{
				ID:         c.Sender().ID,
				Username:   c.Sender().Username,
				Registered: true,
				Balance:    0,
			}
			CreateUser(user)
		} else {
			return c.Send("Ошибка при получении данных пользователя.", mainMenu)
		}
	}

	selectedTariff := tariffs[tariffName]
	if user.Balance < selectedTariff.Price {
		return c.Send("Недостаточно средств на балансе. Пополните счет.", mainMenu)
	}

	user.Balance -= selectedTariff.Price
	user.VPNKey = "ТестовыйVPNКлюч" // Here should be the logic to generate a real VPN key
	user.Location = location
	user.Expiry = time.Now().Add(selectedTariff.Duration)
	UpdateUser(user)

	return c.Send("Тариф "+tariffName+" выбран. Теперь выберите устройство:", deviceMenu)
}

// handleDeviceSelection handles the device selection and provides instructions
func handleDeviceSelection(c tele.Context, device string) error {
	user, err := GetUserByID(c.Sender().ID)
	if err != nil {
		return c.Send("Ошибка при получении данных пользователя.", mainMenu)
	}

	instruction := "Установите приложение Outline (https://play.google.com/store/apps/details?id=org.outline.android.client)\n\n" +
		"1. Скопируйте ключ и добавьте его в приложение.\n" +
		"2. Подключитесь к VPN.\n\n" +
		"Ваш ключ: " + user.VPNKey

	if device == "iphone" {
		instruction = "Установите приложение Outline (https://apps.apple.com/app/outline-app/id1356177741)\n\n" +
			"1. Скопируйте ключ и добавьте его в приложение.\n" +
			"2. Подключитесь к VPN.\n\n" +
			"Ваш ключ: " + user.VPNKey
	}

	return c.Send(instruction, mainMenu)
}

// handlePayment handles the balance top-up
func handlePayment(c tele.Context, amount int) error {
	user, err := GetUserByID(c.Sender().ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = &User{
				ID:         c.Sender().ID,
				Username:   c.Sender().Username,
				Registered: true,
				Balance:    0,
			}
			CreateUser(user)
		} else {
			return c.Send("Ошибка при получении данных пользователя.", mainMenu)
		}
	}

	user.Balance += amount
	UpdateUser(user)

	return c.Send("Ваш баланс пополнен на "+strconv.Itoa(amount)+" руб.\nТекущий баланс: "+strconv.Itoa(user.Balance)+" руб.", mainMenu)
}

// handleBalance handles the balance query
func handleBalance(c tele.Context) error {
	user, err := GetUserByID(c.Sender().ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user = &User{
				ID:         c.Sender().ID,
				Username:   c.Sender().Username,
				Registered: true,
				Balance:    0,
			}
			CreateUser(user)
		} else {
			return c.Send("Ошибка при получении данных пользователя.", mainMenu)
		}
	}
	return c.Send("Ваш баланс: "+strconv.Itoa(user.Balance)+" руб.\nПополните счет для продолжения.", mainMenu)
}

// handleMyKey handles the request to view current key information
func handleMyKey(c tele.Context) error {
	user, err := GetUserByID(c.Sender().ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Send("У вас еще нет активного ключа. Пожалуйста, выберите тариф и получите ключ.", mainMenu)
		}
		return c.Send("Ошибка при получении данных пользователя.", mainMenu)
	}

	if time.Now().After(user.Expiry) {
		return c.Send("Ваш ключ истек. Пополните баланс для продления доступа.", mainMenu)
	}

	info := "Ваш текущий ключ доступа:\n\n" +
		"Локация: " + user.Location + "\n" +
		"Ключ: " + user.VPNKey + "\n" +
		"Начало действия: " + user.Expiry.Add(-tariffs[user.Location].Duration).Format("02.01.2006 15:04") + "\n" +
		"Истечение срока: " + user.Expiry.Format("02.01.2006 15:04") + "\n\n" +
		"Для продления доступа, пожалуйста, пополните баланс."

	return c.Send(info, mainMenu)
}
