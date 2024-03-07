package main

import (
	"bytes"
	"context"
	"log"
	"time"

	tts "cloud.google.com/go/texttospeech/apiv1"
	ttspb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	oto "github.com/ebitengine/oto/v3"
	mp3 "github.com/hajimehoshi/go-mp3"
)

func TTS(ctx context.Context, input <-chan string) error {
	client, err := tts.NewClient(ctx)
	if err != nil {
		return err
	}
	otoCtx, readyChan, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return err
	}
	<-readyChan

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case text := <-input:
			log.Println("->", text)
			// Set the text input to be synthesized.
			si := ttspb.SynthesisInput{
				InputSource: &ttspb.SynthesisInput_Text{
					Text: text,
				},
			}

			// Set the voice parameters.
			voice := ttspb.VoiceSelectionParams{
				LanguageCode: "ja-JP",
				SsmlGender:   ttspb.SsmlVoiceGender_NEUTRAL,
			}

			// Set the audio configuration.
			audioConfig := ttspb.AudioConfig{
				AudioEncoding: ttspb.AudioEncoding_MP3,
			}

			// Perform the text-to-speech request.
			resp, err := client.SynthesizeSpeech(ctx, &ttspb.SynthesizeSpeechRequest{
				Input:       &si,
				Voice:       &voice,
				AudioConfig: &audioConfig,
			})
			if err != nil {
				return err
			}
			decodedMp3, err := mp3.NewDecoder(bytes.NewBuffer(resp.AudioContent))
			if err != nil {
				return err
			}
			player := otoCtx.NewPlayer(decodedMp3)
			player.Play()
			for player.IsPlaying() {
				time.Sleep(time.Millisecond)
			}
			player.Close()
		}
	}
}
