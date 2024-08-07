package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	ctx context.Context
	bot *tele.Bot

	menu             *tele.ReplyMarkup
	categorySelector *tele.ReplyMarkup
}

func NewHandler(ctx context.Context, bot *tele.Bot, menu *tele.ReplyMarkup, categorySelector *tele.ReplyMarkup) *Handler {
	return &Handler{
		ctx:              ctx,
		bot:              bot,
		menu:             menu,
		categorySelector: categorySelector,
	}
}

type User struct {
	userID string
}

func (u *User) Recipient() string {
	return u.userID
}

func main() {
	pref := tele.Settings{
		Token: os.Getenv("TOKEN"),
		// Token:  "6596525393:AAG3CxadbZQPGmcciUBOEc5rOtpQDFtdfSc",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Reply menu buttons
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnSubscribe := menu.Text("Записаться")
	btnSendMessage := menu.Text("Как устроены тренировки")
	btnUnsubscribe := menu.Text("Отписаться")

	// Category selector
	categorySelector := &tele.ReplyMarkup{ResizeKeyboard: true, Selective: true}
	btnOnlineCategory := categorySelector.Text("Онлайн")
	btnOfflineCategory := categorySelector.Text("Оффлайн")
	btnAnketaCategory := categorySelector.Text("Заполнить анкету")
	btnCategoryBack := categorySelector.Text("Назад")

	h := NewHandler(context.Background(), bot, menu, categorySelector)

	h.menu.Reply(
		menu.Row(btnSubscribe),
		menu.Row(btnSendMessage),
		menu.Row(btnUnsubscribe),
	)

	h.categorySelector.Reply(
		categorySelector.Row(btnAnketaCategory),
		categorySelector.Row(btnOnlineCategory),
		categorySelector.Row(btnOfflineCategory),
		categorySelector.Row(btnCategoryBack),
	)

	var channelID int64

	// start
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle(&btnUnsubscribe, h.handleRestart)

	// основное меню
	h.bot.Handle(&btnSubscribe, func(c tele.Context) error {
		return c.Send("Выбрать опцию", categorySelector)
	})
	h.bot.Handle(&btnSendMessage, func(c tele.Context) error {
		channelID = c.Chat().ID
		sendMessage(h.bot, channelID)
		return nil
	})

	// селектор категорий
	h.bot.Handle(&btnOnlineCategory, func(c tele.Context) error {
		return c.Send("Оформлена онлайн подписка", menu)
	})
	h.bot.Handle(&btnOfflineCategory, func(c tele.Context) error {
		return c.Send("Оформлена онлайн подписка", menu)
	})
	h.bot.Handle(&btnAnketaCategory, func(c tele.Context) error {
		return c.Send("https://ссылочка на гуглоформу.com", menu)
	})
	h.bot.Handle(&btnCategoryBack, func(c tele.Context) error {
		return c.Send("Возврат в основное меню", menu)
	})

	h.bot.Start()
}

func (h *Handler) handleStart(c tele.Context) error {
	return c.Send("Привет!", h.menu)
}

func (h *Handler) handleRestart(c tele.Context) error {
	return c.Send("Вы отписались!", h.menu)
}

func sendMessage(b *tele.Bot, userID int64) {
	sendMsg1 := "\n**Как устроены тренировки**\n\n  \n\n1️⃣ Я спрашиваю о целях и твоем образе жизни.\n\n  \n\n▪️Цели. Похудеть, набрать вес. Подготовится к забегу, заплыву, сезону в игровом спорте.\n\n▪️Тренировочный опыт.\n\n▪️Перенесенные травмы.\n\n▪️Текущие боли и ограничения в движении.\n\n▪️Желаемое и реально доступное количество тренировок в неделю.\n\n▪️Продолжительность тренировочных сессий.\n\n▪️Прочая активность и тренировки, которые у тебя есть на неделе (футбол, волейбол, плавание…).\n\n▪️Формат питания.\n\n▪️Качество сна.\n\n▪️Уровень стресса.\n\n▪️Доступные условия для тренировок (дома или в зале, в каком зале, какое там есть оборудование).\n\n▪️ Как мы будем работать (онлайн или офлайн, под контролем или без). Подробнее об этом НИЖЕ ⬇️\n\n  \n\n2️⃣ Делаем набор двигательных тестов (смотрим какие есть ограничения подвижности).\n\n  \n\n3️⃣ Я анализирую всю собранную информацию и:\n\n  \n\n— Для себя отмечаю ограничения и красные флаги (что делать точно нельзя).\n\n— Прописываю нужные адаптации. Например, сила+гипертрофия мышц ног, увеличение мобильности в грудном отделе позвоночника, локальные мышечная выносливость — только то, что поможет достичь оговоренных целей и задач.\n\n— Раскладываю все адаптации в рамках недельного микроцикла. Он зависит от доступного времени и количества тренировочных сессий.\n\n— Под полученную недельную конструкцию подбираю упражнения.\n\n— К каждому упражнению пишу количество повторений, весá, темп и комментарии к выполнению.\n\n— Периодизирую нагрузку в рамках 4 недель. От недели к неделю интенсивность и объем будут меняться.\n\n  \n\n4️⃣ Как мы можем работать:\n\n  \n\n**1. Оффлайн.**\n\nТут все просто: приходишь в зал в Калининграде, и мы тренируемся.\n\n💳 2 500 рублей в час (в эту стоимость заложена аренда зала).\n\n  \n\n**2. Онлайн — есть 3 варианта.**\n\n**1) По видео**\n\nЗвонишь из зала каждую тренировку, и мы работаем по видео. Сначала непривычно, но многим это подходит.\n\n💳 2 000 рублей в час.\n\n  \n\n**2) С индивидуальной программой в приложении TRAINHEROIC**\n\nДля тебя оно бесплатное, а я оплачиваю тренерскую подписку. Ты приходишь в зал, открываешь приложение, ставишь галочки в опроснике по сну и качеству восстановления и запускаешь тренировочную сессию. Там написаны мои комментарии и видео упражнений — как, сколько и что. По окончании сессии мне приходит уведомление, и я анализирую результаты, корректирую нагрузку, меняю предписания, даю обратную связь. В таком формате все детально и адаптировано именно для тебя. Обратная связь в мессенджерах и приложении 24/7.\n\n💳 8 000 рублей за месяц. Без разницы, сколько тренировок — 8 или 12.\n\n  \n\n**3) По общей программе в TRAINHEROIC**\n\nТоже занимаемся с приложением TRAINHEROIC, но программа общая, без детальных предписаний и адаптаций. Только самое необходимое с минимальным количеством тестов и минимальной обратной связью. Все на твоем самоконтроле. Обратная связь в приложении.\n\n💳 2 000 рублей в месяц.\n\n  \n\n  \n\nЕрофеев Константин\n\nhttps://t.me/konstantinSRT"
	// sendMsg2 := "2️⃣ Делаем набор двигательных тестов (смотрим какие есть ограничения подвижности).\n\n  \n\n3️⃣ Я анализирую всю собранную информацию и:\n\n  \n\n— Для себя отмечаю ограничения и красные флаги (что делать точно нельзя).\n\n— Прописываю нужные адаптации. Например, сила+гипертрофия мышц ног, увеличение мобильности в грудном отделе позвоночника, локальные мышечная выносливость — только то, что поможет достичь оговоренных целей и задач.\n\n— Раскладываю все адаптации в рамках недельного микроцикла. Он зависит от доступного времени и количества тренировочных сессий.\n\n— Под полученную недельную конструкцию подбираю упражнения.\n\n— К каждому упражнению пишу количество повторений, весá, темп и комментарии к выполнению.\n\n— Периодизирую нагрузку в рамках 4 недель. От недели к неделю интенсивность и объем будут меняться."
	// sendMsg3 := "4️⃣ Как мы можем работать:\n\n  \n\n**1. Оффлайн.**\n\nТут все просто: приходишь в зал в Калининграде, и мы тренируемся.\n\n💳 2 500 рублей в час (в эту стоимость заложена аренда зала).\n\n  \n\n**2. Онлайн — есть 3 варианта.**\n\n**1) По видео**\n\nЗвонишь из зала каждую тренировку, и мы работаем по видео. Сначала непривычно, но многим это подходит.\n\n💳 2 000 рублей в час.\n\n  \n\n**2) С индивидуальной программой в приложении TRAINHEROIC**\n\nДля тебя оно бесплатное, а я оплачиваю тренерскую подписку. Ты приходишь в зал, открываешь приложение, ставишь галочки в опроснике по сну и качеству восстановления и запускаешь тренировочную сессию. Там написаны мои комментарии и видео упражнений — как, сколько и что. По окончании сессии мне приходит уведомление, и я анализирую результаты, корректирую нагрузку, меняю предписания, даю обратную связь. В таком формате все детально и адаптировано именно для тебя. Обратная связь в мессенджерах и приложении 24/7.\n\n💳 8 000 рублей за месяц. Без разницы, сколько тренировок — 8 или 12.\n\n  \n\n**3) По общей программе в TRAINHEROIC**\n\nТоже занимаемся с приложением TRAINHEROIC, но программа общая, без детальных предписаний и адаптаций. Только самое необходимое с минимальным количеством тестов и минимальной обратной связью. Все на твоем самоконтроле. Обратная связь в приложении.\n\n💳 2 000 рублей в месяц.\n\n  \n\n  \n\nЕрофеев Константин\n\nhttps://t.me/konstantinSRT"
	msg, err := b.Send(
		&User{userID: strconv.FormatInt(userID, 10)},
		sendMsg1,
	)
	if err != nil {
		println(err)
		return
	}
	println(msg)
}
