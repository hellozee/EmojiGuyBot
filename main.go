package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"local.proj/config"
)

type Emojis struct {
	Results []struct {
		Text string `json:"text"`
	} `json:"results"`
}

func getEmojis(text string) string {
	url := "https://api.getdango.com/api/emoji?q=" + text

	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		panic(err)
	}

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var parsed string

	test := Emojis{}
	err = json.Unmarshal(body, &test)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 4; i++ {
		parsed += test.Results[i].Text
	}

	return parsed

}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  config.TelegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnQuery, func(q *tb.Query) {

		results := make(tb.Results, 1)

		emojis := getEmojis(q.Text)

		results[0] = &tb.ArticleResult{Text: emojis, Title: "Emojis", Description: emojis}

		err := b.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60,
		})

		if err != nil {
			log.Fatal(err)
		}

	})

	b.Start()
}
