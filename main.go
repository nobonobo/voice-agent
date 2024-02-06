package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
)

func main() {
	log.SetFlags(log.Lshortfile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// start signal handler
	go func() {
		<-signalChan
		cancel()
	}()

	// start capture
	wg.Add(1)
	wave := make(chan []byte, 128)
	go func() {
		defer wg.Done()
		if err := Capture(ctx, wave); err != nil {
			log.Print(err)
		}
	}()

	// start transcription
	wg.Add(1)
	trans := make(chan *speechpb.StreamingRecognitionResult, 128)
	go func() {
		log.Println("transcription start")
		defer log.Println("transcription end")
		defer wg.Done()
		defer close(trans)
		for {
			if err := Transcription(ctx, wave, trans); err != nil {
				log.Print(err)
				return
			}
		}
	}()

	// start choise proc
	wg.Add(1)
	prompt := make(chan string, 128)
	go func() {
		log.Println("choise start")
		defer log.Println("choise end")
		defer wg.Done()
		for result := range trans {
			log.Println("<-", result.Alternatives[0].Transcript)
			prompt <- result.Alternatives[0].Transcript
		}
	}()

	// start chatgtp completion
	wg.Add(1)
	output := make(chan string, 128)
	go func() {
		log.Println("completion start")
		defer log.Println("completion end")
		defer wg.Done()
		Completion(ctx, prompt, output)
	}()

	// start tts engine
	func() {
		log.Println("text-to-speech start")
		defer log.Println("text-to-speech end")
		if err := TTS(ctx, output); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
