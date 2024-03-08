package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"

	tts "cloud.google.com/go/texttospeech/apiv1"
	ttspb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/youpy/go-wav"
	"google.golang.org/api/option"
)

// AUDIODEV=plughw:2,0 play -e signed-integer -b 16 --endian little -r 24000 -
const Options = `-q -t raw -c 1 -e signed-integer -b 16 --endian little -r 24000 -`

func TTS(ctx context.Context, input <-chan string) error {
	client, err := tts.NewClient(ctx, option.WithAPIKey(config.GcpAPIKey))
	if err != nil {
		return fmt.Errorf("NewClient failed: %w", err)
	}
	cmd := exec.CommandContext(ctx, "play", strings.Split(Options, " ")...)
	cmd.Env = append(cmd.Env, "AUDIODEV="+config.TtsDevice)
	if runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, "AUDIODRIVER=waveaudio")
	}
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("StdinPipe failed: %w", err)
	}
	cmd.Stderr = log.Default().Writer()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd start failed: %w", err)
	}
	defer in.Close()
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
				AudioEncoding:   ttspb.AudioEncoding_LINEAR16,
				SpeakingRate:    config.TtsSpeed,
				SampleRateHertz: 24000,
			}

			// Perform the text-to-speech request.
			resp, err := client.SynthesizeSpeech(ctx, &ttspb.SynthesizeSpeechRequest{
				Input:       &si,
				Voice:       &voice,
				AudioConfig: &audioConfig,
			})
			if err != nil {
				return fmt.Errorf("SynthesizeSpeech failed: %w", err)
			}
			if len(resp.AudioContent) > 0 {
				reader := wav.NewReader(bytes.NewReader(resp.AudioContent))
				if _, err := io.Copy(in, reader); err != nil {
					return fmt.Errorf("io.Copy failed: %w", err)
				}
			}
		}
	}
}
