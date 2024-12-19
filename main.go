package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/reiver/go-atproto/com/atproto/repo"
	"github.com/reiver/go-atproto/com/atproto/server"
)

type WthrResp struct {
	Current struct {
		Properties struct {
			Periods []struct {
				DetailedForecast string `json:"detailedForecast"`
			} `json:"periods"`
		} `json:"properties"`
	} `json:"current"`
}

type LlamaResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LlamaReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

var Prompt = `

You are a helpful weather posting bot. You live in Chicago. You only use weather reporting terms like: wind speed in mph, temperature in fahrenheit. You never make up the weather, because the actual forecast from the National Weather service is included in your context. It's labeled "Current Forecast"

Reformat the following Chicago weather report as a tweet less than 240 characters.

Please use only one emoji that best represents the current weather report. It should be at the beginning of the post.

Do not make up anything.
Do not editorialize.
Do not add any hashtags.
Please include the current date and time. I've provided it as "Current Date".
Always include one emoji that best describes the current weather conditions.

If it's raining or going to rain soon, let people know they should pack a rainjacket or umbrella.

If it's snowing or going to snow soon, include an emoji of a snowman.

If it's sunny, include an emoji of a bright sun and sunglasses.

Never make up a time. Remove the time if unsure.

Be as creative and descriptive as the 240 characters allow.

Current Date: %s
Current Forecast:

%s
`

func main() {
	resp, err := http.Get("https://wthr.cloud/api/forecast?latitude=41.975844&longitude=-87.6633969")
	if err != nil {
		log.Panic(err)
	}

	var wthr WthrResp
	err = json.NewDecoder(resp.Body).Decode(&wthr)
	if err != nil {
		log.Panicf("could not decode wthr %s", err)
	}

	now := time.Now().Format(time.RubyDate)

	req := &LlamaReq{
		Model: "mistral-nemo",
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf(Prompt, now, wthr.Current.Properties.Periods[0].DetailedForecast),
			},
		},
	}

	log.Printf("sending current forecast %+v", req)

	b, _ := json.Marshal(req)
	llmReq, err := http.NewRequest(
		"POST",
		"https://ollama.home.benoist.dev/v1/chat/completions",
		bytes.NewBuffer(b),
	)

	r, err := http.DefaultClient.Do(llmReq)
	if err != nil {
		log.Panicf("could not get completion %s %+v", err, r)
	}

	var llm LlamaResp
	err = json.NewDecoder(r.Body).Decode(&llm)
	if err != nil {
		log.Panicf("could not decode llm %s", err)
	}

	log.Printf("Going to post: %+v", llm.Choices[0].Message.Content)

	password := os.Getenv("CHICAGO_WTHR_BSKY_PASS")
	identifier := "chicago.wthr.cloud"

	var dst server.CreateSessionResponse
	err = server.CreateSession(&dst, identifier, password)
	if nil != err {
		log.Panicf("CREATE SESSION: %s", err)
	}

	bearerToken := dst.AccessJWT

	when := time.Now().Format("2006-01-02T15:04:05.999Z")

	text := strings.Trim(llm.Choices[0].Message.Content, "\"")
	// log.Printf("GOT TOKEN %s", bearerToken)
	post := map[string]any{
		"$type":     "app.bsky.feed.post",
		"text":      text,
		"createdAt": when,
	}

	err = repo.CreateRecord(&dst, bearerToken, identifier, "app.bsky.feed.post", post)
	if nil != err {
		log.Panicf("CREATE RECORD: %s", err)
	}
}
