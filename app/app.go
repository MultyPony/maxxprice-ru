package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	telegramSender "github.com/newdrtm/message_sender/telegram"
	"github.com/rs/xid"
	"github.com/skeris/flat-grabber/dal"
	"github.com/skeris/flat-grabber/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	ENV     = "env"
	DEFAULT = "default"

	MINE_ID = 504242409
	// MINE_ID = 234597475
)

type App struct {
	err    chan error
	mgo    *mongo.Client
	logger *zap.Logger
}

type Options struct {
	MongoServerURI string `env:"NG_MONGO_SERVER_URI" default:"mongodb://127.0.0.1:27017"`
	MongoDBName    string `env:"NG_MONGO_DB_NAME" default:"nedgrabbing"`

	TCPAddr   string `env:"NG_TCP_ADDR" default:":8002"`
	IsGrabber bool   `env:"NG_IS_GRABBER" default:"true"`
}

var Coefficients []float64 = []float64{1.0, 1.1, 1.2}

func New(ctx context.Context, opts Options) (*App, error) {
	var (
		err     error
		logger  *zap.Logger
		errChan = make(chan error)
	)

	logger, err = zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	logger = logger.With(zap.String("svc_app_version", "1"), zap.String("svc_commit_hash", "test"))

	logger.Info("Application is setting up...")
	logger.Debug("Debug mode enabled.")

	logger.Info("MongoDB driver is setting up...")

	mgoOpts := options.Client().ApplyURI(opts.MongoServerURI)

	mgo, err := mongo.Connect(ctx, mgoOpts)
	if err != nil {
		logger.Error("Error while setting up MongoDB driver", zap.Error(err))
		return nil, err
	}
	logger.Info("MongoDB driver is set up.")

	reqBot, err := telegramSender.New("1060150140:AAF8beCNz0amHulE5aMXm_nxZoBfna6sRwE")

	if err != nil {
		fmt.Println("Telegram error 1", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := reqBot.Bot.GetUpdatesChan(u)

	// create context

	go func() {
		logger.Info("Ping MongoDB...")
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err = mgo.Ping(pingCtx, nil); err != nil {
			logger.Error("Unable to ping MongoDB", zap.Error(err))
			errChan <- err
			return
		}
		logger.Info("Ping MongoDB is successfully done.")
	}()

	db := mgo.Database(opts.MongoDBName)
	neddal := dal.New(db, logger)

	var catsString []byte
	go func() {
		for update := range updates {
			data := strings.Split(update.Message.Text, "-")
			if len(data) == 2 {
				id, err := strconv.ParseInt(strings.ReplaceAll(data[0], " ", ""), 10, 64)
				if err != nil {
					reqBot.Send(ctx, err.Error(), MINE_ID, nil)

				}
				rang := strings.Split(strings.TrimLeft(data[1], " "), " ")
				fmt.Println('p', data[1], 'p', rang)
				if len(rang) != 2 {
					reqBot.Send(ctx, "введи диапазон через пробел", MINE_ID, nil)
				}

				price0, err := strconv.ParseInt(strings.ReplaceAll(rang[0], " ", ""), 10, 64)
				if err != nil {
					reqBot.Send(ctx, err.Error(), MINE_ID, nil)
				}
				price1, err := strconv.ParseInt(strings.ReplaceAll(rang[1], " ", ""), 10, 64)
				if err != nil {
					reqBot.Send(ctx, err.Error(), MINE_ID, nil)
				}
				rec, err := neddal.UpdateRecall(ctx, id, price0, price1)
				if err != nil {
					reqBot.Send(ctx, err.Error(), MINE_ID, nil)
				}

				if err := neddal.PutHouse(ctx, &model.Flat{
					ID:              xid.New().String(),
					Price:           price0,
					LocationName:    rec.Street,
					Square:          rec.Square,
					PerSquare:       int64(price0 / rec.Square),
					SortTimeStamp:   time.Now().Unix(),
					Title:           "",
					Agency:          "",
					SplittedAddress: strings.Split(rec.Street, " "),
				}); err != nil {
					reqBot.Send(ctx, err.Error(), MINE_ID, nil)
				}

				if err := neddal.PutHouse(ctx, &model.Flat{
					ID:              xid.New().String(),
					Price:           price1,
					LocationName:    rec.Street,
					Square:          rec.Square,
					PerSquare:       int64(price1 / rec.Square),
					SortTimeStamp:   time.Now().Unix(),
					Title:           "",
					Agency:          "",
					SplittedAddress: strings.Split(rec.Street, " "),
				}); err != nil {
					reqBot.Send(ctx, err.Error(), MINE_ID, nil)
				}
				fmt.Println("GOTTA ALL", data[1])
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/recall", func(w http.ResponseWriter, r *http.Request) {
		reqText, err := ioutil.ReadAll(r.Body)
		fmt.Println("RECALL ", string(reqText), err)
		type req struct {
			Phone  string `json:"phone"`
			Rem    int64  `json:"rem"`
			Square int64  `json:"square"`
			House  string `json:"house"`
			City   string `json:"city"`
			Stair  string `json:"stair"`
		}
		var re req
		json.Unmarshal(reqText, &re)
		neddal.PutRequest(ctx, re.Phone, re.House, re.Square)

		message := fmt.Sprintf("Поступила заявка в форму обратной связи. Номер телефона для связи %s. Данные из калькулятора: ", re.Phone)

		if re.Stair != "" {
			message += re.Stair
		}

		if re.City != "" {
			message += fmt.Sprintf(", %s", re.City)
		}

		if re.Square != 0 {
			message += fmt.Sprintf(", площадь - %d кв.м. ", re.Square)
		}

		if re.House != "" {
			message += fmt.Sprintf(", по адресу %s", re.House)
		}
		switch re.Rem {
		case 0:
			message += "Без отделки"
		case 1:
			message += "Есть, но требуется обновление"
		case 2:
			message += "Недавно сделан кометический ремонт"
		case 3:
			message += "Евро"
		default:
			message += "кулхацкер"
		}

		reqBot.Send(ctx, message, MINE_ID, nil)
	})
	mux.HandleFunc("/average", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("START AVERAGE")

		type avReq struct {
			Rem    int64  `json:"rem"`
			Square int64  `json:"square"`
			House  string `json:"house"`
			City   string `json:"city"`
		}

		var req avReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {

			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// fmt.Println("GETREQAV", req)

		// if req.Rem < 0 || req.Rem >= int64(len(Coefficients)) {
		// 	w.Write([]byte("Walk to hui"))
		// 	w.WriteHeader(413)
		// 	return
		// }

		// min, err := neddal.GetAverage(r.Context(), req.Square, req.City, req.House)
		// if err != nil {
		// 	fmt.Println(err)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		// squareMultiplier := Coefficients[req.Rem]
		// if req.Square > 0 {
		// 	squareMultiplier *= float64(req.Square)
		// }

		data, _ := neddal.GetRequest(r.Context(), req.Rem)

		w.Write([]byte(fmt.Sprintf(`{"avg": %d, "agv": %d}`, data.Response1, data.Response2)))

		// run task list
		// var res string

		// cctx, cancel := chromedp.NewContext(context.Background() /*, chromedp.WithDebugf(log.Printf)*/)
		// defer cancel()

		// err = chromedp.Run(cctx, submit(`https://www.cian.ru/kalkulator-nedvizhimosti/?address=%D0%A0%D0%BE%D1%81%D1%81%D0%B8%D1%8F%2C%20%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%2C%20%D0%A6%D0%B8%D0%BC%D0%BB%D1%8F%D0%BD%D1%81%D0%BA%D0%B0%D1%8F%20%D1%83%D0%BB%D0%B8%D1%86%D0%B0%2C%2014%20&roomsCount=3&totalArea=74&flatNumber=2`, `div[data-testid="under-price"]`, `chromedp`, &res))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// tags := strings.Split(res, "<")
		// ress := ""

		// for _, val := range tags {
		// 	fmt.Println(val)
		// 	if strings.Contains(val, `data-testid="price"`) {
		// 		ress = strings.Split(strings.Split(val, ">")[1], " ")[0]
		// 		break
		// 	}
		// }

		// res = OnPage(`https://www.cian.ru/kalkulator-nedvizhimosti/?address=%D0%A0%D0%BE%D1%81%D1%81%D0%B8%D1%8F%2C%20%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%2C%20%D0%A6%D0%B8%D0%BC%D0%BB%D1%8F%D0%BD%D1%81%D0%BA%D0%B0%D1%8F%20%D1%83%D0%BB%D0%B8%D1%86%D0%B0%2C%2014%20&roomsCount=3&totalArea=74&flatNumber=2`)

		// // log.Println("got: ", strings.Contains(res, `млн`))
		// fmt.Println(ress, len(tags))
	})

	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		reqText, err := ioutil.ReadAll(r.Body)
		fmt.Println("cat ", string(reqText), err)
		type req struct {
			Phone  string `json:"phone"`
			Rem    int64  `json:"rem"`
			Square int64  `json:"square"`
			House  string `json:"house"`
			City   string `json:"city"`
			Stair  string `json:"stair"`
		}
		var re req
		json.Unmarshal(reqText, &re)
		id, _ := neddal.PutRequest(ctx, re.Phone, re.House, re.Square)

		w.Write([]byte(fmt.Sprintf(`{"avg": %d}`, id)))

		avg, err := neddal.GetAverage(ctx, re.Square, re.City, re.House)
		if err != nil {
			reqBot.Send(ctx, err.Error(), MINE_ID, nil)
		}

		reqBot.Send(ctx, fmt.Sprintf("%d - %s", id, strings.ReplaceAll(fmt.Sprintf("https://www.cian.ru/kalkulator-nedvizhimosti/?address=%s&totalArea=%d&roomsCount=%s&flatNumber=1", url.PathEscape(re.House), re.Square, re.Stair), "+", "%20"))+fmt.Sprintf("  данные из нашей базы %d ", int(avg*re.Square)), MINE_ID, nil)
	})

	grabberServ := &http.Server{
		Addr:    opts.TCPAddr,
		Handler: mux, // FIXME: Should wrap mux in http_middleware.Chain here
	}

	app := &App{
		err:    errChan,
		logger: logger,
	}

	go func() {
		logger.Info(fmt.Sprintf("Grabber server is started on %s", opts.TCPAddr))
		if err := grabberServ.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			logger.Error("Grabber server fatal error", zap.Error(err))
			app.err <- err
		}
	}()

	if opts.IsGrabber {

		// var pos int64

		// i, err := neddal.GetLastPage(ctx)
		// if err != nil {
		// 	fmt.Println(" err in get config ", err)
		// }

		go func() {
			fmt.Println("start grabber")

			//client := &http.Client{}

			//i := 0

			collector := colly.NewCollector()
			collector.OnError(func(r *colly.Response, err error) {
				fmt.Println("ERRROROROROROR ", string(r.Body), err)
				reqBot.Send(ctx, err.Error(), MINE_ID, nil)
			})
			collector.OnHTML("div.snippet-list.js-catalog_serp", func(e *colly.HTMLElement) {
				e.ForEach("div.snippet.snippet-horizontal.snippet-redesign.item-snippet-with-aside.item.item_table.clearfix.js-catalog-item-enum.item-with-contact.js-item-extended", func(i int, e *colly.HTMLElement) {
					fmt.Println("item ", i)
					fmt.Println(e.ChildAttr("a.snippet-link", "title"))
					fmt.Println(e.ChildText("span.snippet-price"))
					fmt.Println(strings.Split(e.ChildText("span.snippet-price"), "₽"))

					fmt.Println(e.ChildAttr("div.snippet-line-row a.snippet-link", "href"))
					fmt.Println(e.ChildAttr("h3.snippet-title a.snippet-link", "href"))
					fmt.Println(e.ChildAttr("div.js-item-contacts-button.item-extended-contacts", "data-props"))
					price := strings.Split(e.ChildText("span.snippet-price"), "₽")

					fmt.Println(len(strings.Split(e.ChildText("span.snippet-price"), "₽")))

					if len(price) > 1 {
						if strings.Trim(price[1], " ") == "" {
							fmt.Println("Нам подходит")

						} else {
							return
						}
					}
					fmt.Println(e.ChildText("span.item-address__string"))

					splitted := strings.Split(e.ChildAttr("a.snippet-link", "title"), " ")

					var square int64

					for i, value := range splitted {
						if value == "м²," {
							flsquare, err := strconv.ParseFloat(splitted[i-1], 64)
							if err != nil {
								fmt.Println(err)
							}
							square = int64(flsquare)
							break
						}
					}

					priceInt, _ := strconv.ParseInt(strings.ReplaceAll(price[0], " ", ""), 10, 64)

					if err := neddal.PutHouse(ctx, &model.Flat{
						ID:              xid.New().String(),
						SplittedAddress: strings.Split(e.ChildText("span.item-address-georeferences-item__content")+" "+e.ChildText("span.item-address__string"), " "),
						Price:           priceInt,
						LocationName:    e.ChildText("span.item-address-georeferences-item__content") + " " + e.ChildText("span.item-address__string"),
						PerSquare:       int64(priceInt / square),
						Square:          square,
						Agency:          e.ChildAttr("div.snippet-line-row a.snippet-link", "href"),
						Link:            e.ChildAttr("h3.snippet-title a.snippet-link", "href"),
						RawPhone:        e.ChildAttr("div.js-item-contacts-button.item-extended-contacts", "data-props"),
					}); err != nil {
						fmt.Println("error in put house ", err)
						panic(err)
					}
				})
			})

			collector.OnRequest(func(r *colly.Request) {
				fmt.Println("Visiting", r.URL)
			})
			// for {
			// 	collector.Visit(fmt.Sprintf("https://www.avito.ru/moskva_i_mo/kvartiry/prodam/vtorichka-ASgBAQICAUSSA8YQAUDmBxSMUg?f=ASgBAQICAUSSA8YQAkDmBxSMUsoIxIpZmqwBmKwBlqwBlKwBiFmGWYRZglmAWfzPMv5Y&p=%d&proprofile=1", i))
			// 	i++

			// 	time.Sleep(1 * time.Minute)
			// }

			// 	req, err := http.NewRequest("GET", "https://www.avito.ru/web/1/main/items", nil)
			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	req.Header.Add("Host", "www.avito.ru")
			// 	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:79.0) Gecko/20100101 Firefox/79.0`)
			// 	req.Header.Add("Accept", "application/json")
			// 	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
			// 	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")
			// 	req.Header.Add("Referer", `https://www.avito.ru/rossiya/nedvizhimost`)
			// 	req.Header.Add("content-type", "application/json")
			// 	req.Header.Add("x-requested-with", "XMLHttpRequest")
			// 	req.Header.Add("x-source", "client-browser")
			// 	req.Header.Add("Connection", "keep-alive")

			// 	q := req.URL.Query()
			// 	q.Add("forceLocation", "false")
			// 	q.Add("locationId", "107620")
			// 	q.Add("categoryId", "4")
			// 	q.Add("limit", "30")
			// 	q.Add("lastStamp", "0")
			// 	q.Add("offset", fmt.Sprintf("%d", i*10))
			// 	req.URL.RawQuery = q.Encode()

			// 	if err := neddal.SaveLastPage(ctx, int64(i)); err != nil {
			// 		fmt.Println(err)
			// 	}

			// 	if err != nil {
			// 		fmt.Println("Errored when sending request to the server")
			// 		return
			// 	}
			// 	resp, err := client.Do(req)

			// 	defer resp.Body.Close()
			// 	resp_body, _ := ioutil.ReadAll(resp.Body)

			// 	var avitoResp model.AvitoFlatResponse

			// 	if err := json.Unmarshal(resp_body, &avitoResp); err != nil {
			// 		panic(err)
			// 	}

			// 	for _, item := range avitoResp.Items {
			// 		func() {
			// 			defer func() {
			// 				recover()
			// 			}()
			// 			fmt.Println("AVIDATA ", item.GeoForItems)

			// 			if err := neddal.PutHouse(ctx, dal.Dto2dao_House(&item)); err != nil {
			// 				fmt.Println("error in put house ", err)
			// 				panic(err)
			// 			}
			// 		}()
			// 	}

			// 	pos++
			// 	i++
			// 	time.Sleep(2 * time.Second)
			// }
		}()
	}

	go func() {
		for {
			var cats []string

			cats, err = neddal.GetCategories(ctx)
			if err != nil {
				fmt.Println(err)
			}

			catsString, err = json.Marshal(cats)
			if err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()

	return app, nil
}

func (a *App) GetLogger() *zap.Logger {
	return a.logger
}

func (a *App) GetErr() chan error {
	return a.err
}

func OnPage(link string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func submit(urlstr, sel, q string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		//chromedp.WaitVisible(sel),
		chromedp.OuterHTML(`html` /*`div[data-testid="price"]`*/, res),
	}
}
