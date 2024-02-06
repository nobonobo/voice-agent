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

func proc(ctx context.Context, msg string, output chan<- string) func() {
	//log.Println("proc:", config.OaiModel, config.OaiMaxTokens)
	req := openai.ChatCompletionRequest{
		Model:     config.OaiModel,
		MaxTokens: config.OaiMaxTokens,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "チャットボットとしてロールプレイします。ずんだもんという名前のボットとして振る舞ってください。性格はポジティブで元気です。"},
			{Role: openai.ChatMessageRoleUser, Content: msg},
		},
		Stream: true,
	}
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

func Completion(ctx context.Context, prompt <-chan string, output chan<- string) {
	if oaiClient == nil {
		oaiClient = openai.NewClient(config.OaiAPIKey)
	}
	cancel := func() {}
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-prompt:
			cancel()
			cancel = sync.OnceFunc(proc(ctx, msg, output))
		}
	}
}
