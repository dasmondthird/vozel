package bot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

type Menus struct {
	MainMenu, VpnMenu, ServerMenu, TariffMenu, DeviceMenu, PaymentMenu *tele.ReplyMarkup
}

func (b *Bot) SetupMenus() {
	b.Menus.MainMenu = &tele.ReplyMarkup{}
	btnMyVPN := b.Menus.MainMenu.Data("Мой VPN", "myvpn")
	btnMyKey := b.Menus.MainMenu.Data("Мой ключ", "mykey")
	btnBalance := b.Menus.MainMenu.Data("Баланс", "balance")
	btnPayment := b.Menus.MainMenu.Data("Пополнить баланс", "payment")
	b.Menus.MainMenu.Inline(
		b.Menus.MainMenu.Row(btnMyVPN),
		b.Menus.MainMenu.Row(btnMyKey),
		b.Menus.MainMenu.Row(btnBalance),
		b.Menus.MainMenu.Row(btnPayment),
	)

	b.Menus.VpnMenu = &tele.ReplyMarkup{}
	btnOutline := b.Menus.VpnMenu.Data("Outline", "outline")
	btnBackToMain := b.Menus.VpnMenu.Data("Назад", "backtomain")
	b.Menus.VpnMenu.Inline(
		b.Menus.VpnMenu.Row(btnOutline),
		b.Menus.VpnMenu.Row(btnBackToMain),
	)

	b.Menus.ServerMenu = &tele.ReplyMarkup{}
	btnServer1 := b.Menus.ServerMenu.Data("Сервер 1", "server1")
	btnBackToVPN := b.Menus.ServerMenu.Data("Назад", "backtovpn")
	b.Menus.ServerMenu.Inline(
		b.Menus.ServerMenu.Row(btnServer1),
		b.Menus.ServerMenu.Row(btnBackToVPN),
	)

	b.Menus.TariffMenu = &tele.ReplyMarkup{}
	btnBackToServers := b.Menus.TariffMenu.Data("Назад", "backtoservers")
	b.Menus.TariffMenu.Inline(b.Menus.TariffMenu.Row(btnBackToServers))

	b.AddTariffButtons()
}

func (b *Bot) AddTariffButtons() {
	for name, tariff := range b.Tariffs {
		btn := b.Menus.TariffMenu.Data(fmt.Sprintf("%s - %d руб", name, tariff.Price), "tariff_"+name)
		b.Menus.TariffMenu.Inline(b.Menus.TariffMenu.Row(btn))
	}
}
