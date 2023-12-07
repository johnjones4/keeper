package main

import (
	"main/api"
	"main/core"
	"main/index"
	"main/store"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/johnjones4/errorbot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func startIndex(idx core.Index) {
	for {
		err := idx.ReIndex()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Hour * 6)
	}
}

func main() {
	chatId, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		panic(err)
	}
	bot := errorbot.New(
		"keeper",
		os.Getenv("TELEGRAM_TOKEN"),
		chatId,
	)

	config := zap.NewDevelopmentConfig()
	l, err := config.Build(zap.Hooks(bot.ZapHook([]zapcore.Level{
		zapcore.FatalLevel,
		zapcore.PanicLevel,
		zapcore.DPanicLevel,
		zapcore.ErrorLevel,
		zapcore.WarnLevel,
	})))
	if err != nil {
		panic(err)
	}

	defer l.Sync()
	log := l.Sugar()

	store := store.New(os.Getenv("ROOT_DIR"), strings.Split(os.Getenv("IGNORE_FS"), "|"))
	idx, err := index.New(os.Getenv("INDEX_PATH"), store, log)
	if err != nil {
		panic(err)
	}

	rc := core.RuntimeContext{
		Store:        store,
		Index:        idx,
		PrivateKey:   []byte(os.Getenv("PRIVATE_KEY")),
		PasswordHash: []byte(os.Getenv("PASSWORD")),
		Log:          log,
	}

	go startIndex(idx)

	a := api.New(&rc)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), a)
	panic(err)
}
