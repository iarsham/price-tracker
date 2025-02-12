package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/iarsham/price-tracker/configs"
	"github.com/iarsham/price-tracker/internal/entities"
	"github.com/iarsham/price-tracker/internal/helpers"
	ptime "github.com/yaa110/go-persian-calendar"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
	"io"
	"net/http"
	"time"
)

const ApiURL string = "https://brsapi.ir/FreeTsetmcBourseApi/Api_Free_Gold_Currency_v2.json"

type TrackerService struct {
	Client     *http.Client
	Bot        *telebot.Bot
	BackGround gocron.Scheduler
	Logger     *zap.Logger
	Cfg        *configs.Config
}

func (t *TrackerService) Run() {
	_, err := t.BackGround.NewJob(
		gocron.DurationJob(
			time.Hour,
		),
		gocron.NewTask(func() {
			t.scrapeTgju()
		}),
		gocron.WithName(
			"Scrape Prices",
		),
	)
	if err != nil {
		t.Logger.Fatal("Failed to create background price for scrapers.", zap.Error(err))
	}
}

func (t *TrackerService) scrapeTgju() {
	defer func() {
		if err := recover(); err != nil {
			t.Logger.Error("Failed to scrape tgju", zap.Any("error", err))
		}
	}()
	var data entities.PriceData
	req, err := http.NewRequest(http.MethodGet, ApiURL, nil)
	if err != nil {
		t.Logger.Error("Failed to read api data", zap.Error(err))
	}
	resp, err := t.Client.Do(req)
	if err != nil {
		t.Logger.Error("Failed to get api response", zap.Error(err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Logger.Error("Failed to read response body", zap.Error(err))
	}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Logger.Error("Failed to unmarshal api data", zap.Error(err))
	}
	extractedItems := make(map[string]string)
	for _, item := range data.Gold {
		extractedItems[item.Symbol] = helpers.ThousandSeparator(fmt.Sprintf("%.2f", item.Price))
	}
	for _, item := range data.Currency {
		extractedItems[item.Symbol] = helpers.ThousandSeparator(fmt.Sprintf("%.2f", item.Price))
	}
	for _, item := range data.Cryptocurrency {
		extractedItems[item.Symbol] = helpers.ThousandSeparator(fmt.Sprintf("%.2f", item.Price))
	}
	fmt.Println(extractedItems)
	fmt.Println(extractedItems["IMCOIN"])
	pt := ptime.Now()
	outPut := fmt.Sprintf(`📊 لیست قیمت‌ها:
🕒 زمان به‌روزرسانی: %s 

🌟 قیمت‌های طلا:
🔸 سکه امامی : %s تومان
🔸 سکه بهار آزادی : %s تومان
🔸 طلای ۱۸ عیار : %s تومان
🔸 ربع سکه : %s تومن
🔸 نیم سکه : %s تومن

💵 قیمت‌های ارز:
🔸 دلار (USD): %s تومان
🔸 یورو (EUR): %s تومان
🔸 پوند انگلیس (GBP): %s تومان
🔸 لیر ترکیه (TRY): %s تومان
🔸 یوان چین (CNY): %s تومان
🔸 دینار عراق  (IQD): %s تومان
🔸 روبل روسیه (RUB): %s تومان
🔸 روپیه هند (INR): %s تومان

💹 قیمت‌های ارز دیجیتال:
🔸 بیت کوین (BTC): %s دلار
🔸 اتریوم (ETH): %s دلار
🔸 لایت کوین (LTC): %s دلار
🔸 ریپل (XRP): %s دلار


📢 کانال ما: @price_hourly`,
		pt.Format("yyyy/MM/dd E hh:mm:ss a"),
		extractedItems["IMCOIN"],
		extractedItems["BACOIN"],
		extractedItems["Gold18"],
		extractedItems["QUCOIN"],
		extractedItems["HACOIN"],
		extractedItems["USD"],
		extractedItems["EUR"],
		extractedItems["GBP"],
		extractedItems["TRY"],
		extractedItems["CNY"],
		extractedItems["IQD"],
		extractedItems["RUB"],
		extractedItems["INR"],
		extractedItems["BTC"],
		extractedItems["ETH"],
		extractedItems["LTC"],
		extractedItems["XRP"],
	)
	if _, err := t.Bot.Send(&telebot.Chat{ID: t.Cfg.App.ChannelID}, outPut); err != nil {
		t.Logger.Error("Failed to send info to channel", zap.Error(err))
	}
}
