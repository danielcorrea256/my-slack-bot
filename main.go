package main

import (
	"encoding/json"
	"context"
	"fmt"
	"log"
	"os"
	"github.com/shomali11/slacker"
	"github.com/joho/godotenv"
	"net/http"
	"io/ioutil"
)

type Quote struct {
	Quote string `json:"content"`
	Author string `json:"author"`
}

func getQuote() Quote{
	response, err := http.Get("https://api.quotable.io/random")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject Quote
	json.Unmarshal(responseData, &responseObject)
	return responseObject
}

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println("***")
	}
}

func main()  {
	godotenv.Load()
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	go printCommandEvents(bot.CommandEvents())

	var last string

	bot.Command("get quote please", &slacker.CommandDefinition{
		Description: "get quote",
		Examples: []string{"get quote"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			timestamp := botCtx.Event().TimeStamp

			// the bot is listening to different events and this can cause repeated messages
			// then to avoid this if a message timestamp is the same than the last one, ignore
			if last == timestamp {
				fmt.Println("ommiting message")
			} else {
				last = timestamp
				q := getQuote()
				r := fmt.Sprintf("%s - %s", q.Quote, q.Author)
				response.Reply(r)	
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}