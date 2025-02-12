package main

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/iarsham/price-tracker/configs"
	"github.com/iarsham/price-tracker/internal/services"
	"github.com/iarsham/price-tracker/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
	"net/http"
	"time"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}
	logs, err := logger.NewZapLog(cfg.App.Debug)
	if err != nil {
		panic(err)
	}
	defer logs.Sync()
	pref := telebot.Settings{
		Token:  cfg.App.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		logs.Fatal("Failed to create bot", zap.Error(err))
	}
	backGround, err := gocron.NewScheduler()
	if err != nil {
		logs.Fatal("Failed to create scheduler", zap.Error(err))
	}
	defer backGround.Shutdown()
	tracker := &services.TrackerService{
		Client:     &http.Client{},
		Bot:        bot,
		BackGround: backGround,
		Logger:     logs,
		Cfg:        cfg,
	}
	logs.Info("Tracker started...")
	go tracker.Run()
	logs.Info("Background task started...")
	backGround.Start()
	logs.Info("Bot started...")
	bot.Start()
}
