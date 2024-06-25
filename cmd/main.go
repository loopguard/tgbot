package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	menu             = &tele.ReplyMarkup{}
	categorySelector = &tele.ReplyMarkup{}

	// Reply buttons.
	btnSetCategory = menu.Text("Выбрать категорию")
	btnHelp        = menu.Text("Помощь")
	btnUnsubscribe = menu.Text("Отписаться")
	btnSendMessage = menu.Text("Напиши мне")

	btnJobCategory   = menu.Text("Работа")
	btnAutoCategory  = menu.Text("Авто")
	btnRealtCategory = menu.Text("Недвижимость")
	btnCategoryBack  = menu.Text("Назад")
)

type User struct {
	userID string
}

func (u *User) Recipient() string {
	return u.userID
}

func main() {
	pref := tele.Settings{
		// Token:  os.Getenv("TOKEN"),
		Token:  "6596525393:AAG3CxadbZQPGmcciUBOEc5rOtpQDFtdfSc",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	menu.Reply(
		menu.Row(btnSetCategory),
		menu.Row(btnHelp),
		menu.Row(btnUnsubscribe),
		menu.Row(btnSendMessage),
	)

	categorySelector.Reply(
		categorySelector.Row(btnJobCategory),
		categorySelector.Row(btnAutoCategory),
		categorySelector.Row(btnRealtCategory),
		categorySelector.Row(btnCategoryBack),
	)

	var channelID int64

	b.Handle("/start", func(c tele.Context) error {
		channelID = c.Chat().ID

		return c.Send("Привет!", menu)
	})

	// основное меню
	b.Handle(&btnSetCategory, func(c tele.Context) error {
		return c.Send("Выберете категорию", categorySelector)
	})
	b.Handle(&btnHelp, func(c tele.Context) error {
		return c.Send("По всем вопросам обращайтесь в Службу поддержки", menu)
	})
	b.Handle(&btnUnsubscribe, func(c tele.Context) error {
		println(fmt.Sprintf("UserID: %d отписался", channelID))
		return c.Send("Вы отписались от уведомлений...", menu)
	})

	// отправка сообщения юзеру
	b.Handle(&btnSendMessage, func(c tele.Context) error {
		sendMessage(b, channelID)
		return nil
	})

	// селектор категорий
	b.Handle(&btnJobCategory, func(c tele.Context) error {
		println(fmt.Sprintf("UserID: %d подписался на Работа", channelID))
		return c.Send("Оформлена подписка на обновления в категории Работа", menu)
	})
	b.Handle(&btnRealtCategory, func(c tele.Context) error {
		println(fmt.Sprintf("UserID: %d подписался на Недвижимость", channelID))
		return c.Send("Оформлена подписка на обновления в категории Недвижимость", menu)
	})
	b.Handle(&btnAutoCategory, func(c tele.Context) error {
		println(fmt.Sprintf("UserID: %d подписался на Авто", channelID))
		return c.Send("Оформлена подписка на обновления в категории Авто", menu)
	})
	b.Handle(&btnCategoryBack, func(c tele.Context) error {
		return c.Send("Возврат в основное меню", menu)
	})

	// sendMessage(b)

	b.Start()
}

func sendMessage(b *tele.Bot, userID int64) {
	sendMsg := "```mermaid\nsequenceDiagram\nAlice ->> Bob: Hello Bob, how are you?\nBob-->>John: How about you John?\nBob--x Alice: I am good thanks!\nBob-x John: I am good thanks!\nNote right of John: Bob thinks a long<br/>long time, so long<br/>that the text does<br/>not fit on a row.\n\nBob-->Alice: Checking with John...\nAlice->John: Yes... John, how are you?\n```"
	msg, err := b.Send(
		&User{userID: strconv.FormatInt(userID, 10)},
		sendMsg,
		&tele.SendOptions{
			ParseMode: tele.ModeMarkdownV2,
		})
	if err != nil {
		println(err)
		return
	}
	println(msg)
}
