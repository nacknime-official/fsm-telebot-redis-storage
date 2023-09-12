package redisfsm_test

import (
	"log"
	"os"

	redisfsm "github.com/nacknime-official/fsm-telebot-redis-storage"
	"github.com/redis/go-redis/v9"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"golang.org/x/net/context"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func ExampleNewDefaultStorage() {
	bot, err := tele.NewBot(tele.Settings{
		Token:   os.Getenv("BOT_TOKEN"),
		Offline: true,
	})

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	fsmStorage := redisfsm.NewDefaultStorage(client)

	defer func(fsmStorage *redisfsm.Storage) {
		err := fsmStorage.Close()
		if err != nil {
			log.Printf("redis storage Close: %v", err)
		}
	}(fsmStorage)

	g := bot.Group()
	g.Use(middleware.AutoRespond())

	manager := fsm.NewManager(bot, g, fsmStorage, nil)
	manager.Use(middleware.Recover(func(err error) {
		log.Printf("panic recovered: %v", err)
	}))

}
