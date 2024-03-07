package main

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

var oaiClient *openai.Client

func checkSentence(res string) bool {
	r := []rune(res)
	if len(r) > 5 {
		for _, ch := range []rune{'、', '。', '.', '！', '？', '\n'} {
			if r[len(r)-1] == ch {
				return true
			}
		}
	}
	return false
}

func proc(ctx context.Context, req openai.ChatCompletionRequest, output chan<- string) func() {
	//log.Println("proc:", config.OaiModel, config.OaiMaxTokens)
	stream, err := oaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("CompletionStream error: %v", err)
		return nil
	}
	go func() {
		buff := strings.Builder{}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Printf("Stream error: %v", err)
				return
			}
			for _, v := range response.Choices[0].Delta.ToolCalls {
				log.Println(v.Function.Name, v.Function.Arguments)
				switch v.Function.Name {
				case "chat_reset":
					req.Messages = ResetMessages()
					if _, err := buff.WriteString("チャットをリセットしました。"); err != nil {
						log.Printf("Output writing error: %v", err)
						return
					}
				}
			}
			if _, err := buff.WriteString(response.Choices[0].Delta.Content); err != nil {
				log.Printf("Output writing error: %v", err)
				return
			}
			sentence := buff.String()
			if checkSentence(sentence) {
				output <- sentence
				buff.Reset()
			}
			if response.Choices[0].FinishReason == "stop" {
				return
			}
		}
	}()
	return func() {
		stream.Close()
	}
}

const SYSTEM_PROMPT = `You are a helpfule assistant.
You need to follow the following rules:
- lang:ja
- please be polite (e.g. use です, ます)
- short reply (less than 100 japanese charactors)
`

func ResetMessages() []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: SYSTEM_PROMPT},
	}
}

func Tools() []openai.Tool {
	t := openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionDefinition{
			Name:        "chat_reset",
			Description: "チャットをリセット",
			Parameters:  nil,
		},
	}
	return []openai.Tool{t}
}

func Completion(ctx context.Context, prompt <-chan string, output chan<- string) {
	if oaiClient == nil {
		oaiClient = openai.NewClient(config.OaiAPIKey)
	}
	msgs := ResetMessages()
	req := openai.ChatCompletionRequest{
		Model:     config.OaiModel,
		MaxTokens: config.OaiMaxTokens,
		Messages:  msgs,
		Stream:    true,
		Tools:     Tools(),
	}
	cancel := func() {}
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-prompt:
			cancel()
			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleUser, Content: msg,
			})
			cancel = sync.OnceFunc(proc(ctx, req, output))
		}
	}
}
