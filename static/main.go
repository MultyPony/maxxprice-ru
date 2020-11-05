package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RequestedMessage struct {
	House               string `json:"house"`
	Square              string `json:"square"`
	Rooms               string `json:"rooms"`
	Rem                 string `json:"rem"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	CommunicationMethod string `json:"communicationMethod"`
}

func main() {
	port := os.Getenv("PORT")

	const pass = "1234"
	botToken := "1498849323:AAEnx4yglkRZRGx9fbn0_NstCqlX4ATLRJA"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Что то не так... ", err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	var msg tgbotapi.MessageConfig
	var chatIDMap = make(map[int64]struct{})

	go func() {
		for update := range updates {

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			ent := update.Message.Entities
			entType := "No entity"
			if ent != nil {
				entType = (*ent)[0].Type
			}
			log.Printf("%s", entType)

			if _, ok := chatIDMap[update.Message.Chat.ID]; ok {
				continue
			} else {
				chatIDMap[update.Message.Chat.ID] = struct{}{}
			}
		}
	}()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var m RequestedMessage

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

		str := fmt.Sprintf("*Aдрес*: %v\n*Площадь*: %v\n*Комнат*: %v\n*Ремонт*: %v\n*Предпочитаемый вид связи*: %v\n*Почта*: %v\n*Телефон*: %v", m.House, m.Square, m.Rooms, repair, m.CommunicationMethod, m.Email, m.Phone)

		for key := range chatIDMap {
			msg = tgbotapi.NewMessage(key, "")
			msg.ParseMode = "Markdown"
			msg.Text = str
			bot.Send(msg)
		}
		fmt.Fprintf(w, "%s", b)
	}
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/", fs)
	http.HandleFunc("/submit-contacts", handler)
	log.Print("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
