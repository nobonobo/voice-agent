package main

import (
	"log"

	"github.com/caarlos0/env/v10"
)

var config = struct {
	GcpAPIKey    string  `env:"GCP_API_KEY,required"`
	SttStabirity float64 `env:"STT_STABILITY" envDefault:"0.7"`
	OaiAPIKey    string  `env:"OAI_API_KEY,required"`
	OaiModel     string  `env:"OAI_MODEL" envDefault:"gpt-3.5-turbo"`
	OaiMaxTokens int     `env:"OAI_MAX_TOKENS" envDefault:"1024"`
	VoiceVoxDir  string  `env:"VOICEVOX_DIR" envDefault:"./voicevox"`
	ActorID      int     `env:"ACTOR_ID" envDefault:"3"`
	TtsSpeed     float64 `env:"TTS_SPEED" envDefault:"1.2"`
	TtsPause     float64 `env:"TTS_PAUSE" envDefault:"0.5"`
}{}

func init() {
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}
}
