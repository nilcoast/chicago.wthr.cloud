package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
Reformat the following Chicago weather report from %s as a tweet less than 240 characters.

Please use emoji.
Do not make up anything.
Do not editorialize.
Do not add any hashtags.

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
		log.Panic(err)
	}

	req := &LlamaReq{
		Model: "llama3.2",
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf(Prompt, "current-time", wthr.Current.Properties.Periods[0].DetailedForecast),
			},
		},
	}

	b, _ := json.Marshal(req)
	http.Post(
		"https://ollama.home.benoist.dev/v1/chat/completions",
		"application/json",
		bytes.NewBuffer(b),
	)

	// do the post
}
