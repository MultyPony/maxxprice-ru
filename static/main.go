package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// func handler(w http.ResponseWriter, r *http.Request) {
// 	botToken := "1498849323:AAEnx4yglkRZRGx9fbn0_NstCqlX4ATLRJA"
// 	bot, err := tgbotapi.NewBotAPI(botToken)
// 	if err != nil {
// 		log.Panic("Что то не так... ", err)
// 	}
// 	// u := tgbotapi.NewUpdate(0)
// 	// u.Timeout = 60

// 	// updates, err := bot.GetUpdatesChan(u)

// 	b, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var m interface{}

// 	// for update := range updates {
// 	// 	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
// 	err2 := json.Unmarshal(b, &m)
// 	if err2 != nil {
// 		log.Panic("JSON ER", err)
// 	}
// 	str := fmt.Sprintf("%v", m)
// 	msg.Text = str
// 	// msg.Text = update.Message.Text

// 	bot.Send(msg)
// 	// }

// 	fmt.Fprintf(w, "Hi there, I love %s!", b)
// }

type RequestedMessage struct {
	// City   string `json:"city"`
	House  string `json:"house"`
	Square string `json:"square"`
	Stair  string `json:"stair"`
	Rem    string `json:"rem"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
}

func main() {
	const pass = "1234"
	botToken := "1498849323:AAEnx4yglkRZRGx9fbn0_NstCqlX4ATLRJA"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Что то не так... ", err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	// update := <-updates
	var msg tgbotapi.MessageConfig
	// msg.ParseMode = "Markdown"

	go func() {
		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			// tgbotapi.Command(update.Message)
			ent := update.Message.Entities
			entType := (*ent)[0].Type
			log.Printf("%s", entType)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if entType == "bot_command" && update.Message.Text == "/start" {
				msg.Text = "Введите пароль: "
				bot.Send(msg)

			}
			// msg.ReplyToMessageID = update.Message.MessageID

			// bot.Send(msg)
		}
	}()

	handler := func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var m RequestedMessage
		// var m interface{}

		// for update := range updates {
		// update := <-updates
		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		// msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
		// }
		// msg := tgbotapi.NewMessage(updates.Message.Chat.ID, updates.Message.Text)
		err2 := json.Unmarshal(b, &m)
		if err2 != nil {
			log.Panic("JSON ER", err)
		}

		repair := ""
		switch m.Rem {
		case "0":
			repair = "Без отделки"
		case "1":
			repair = "Есть, но требуется обновление"
		case "2":
			repair = "Недавно сделан коcметический ремонт"
		case "3":
			repair = "Евро"
		default:
			repair = ""
		}

		str := fmt.Sprintf("*Aдрес*: %v\n*Площадь*: %v\n*Этаж*: %v\n*Ремонт*: %v\n*Почта*: %v\n*Телефон*: %v", m.House, m.Square, m.Stair, repair, m.Email, m.Phone)
		msg.ParseMode = "Markdown"
		msg.Text = str
		bot.Send(msg)
		fmt.Fprintf(w, "%s", b)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
