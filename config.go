package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

var config = struct {
	GcpAPIKey    string  `env:"GCP_API_KEY,required"`
	SttStabirity float64 `env:"STT_STABILITY" envDefault:"0.7"`
	OaiAPIKey    string  `env:"OAI_API_KEY,required"`
	OaiModel     string  `env:"OAI_MODEL" envDefault:"gpt-3.5-turbo"`
	OaiMaxTokens int     `env:"OAI_MAX_TOKENS" envDefault:"1024"`
	TtsSpeed     float64 `env:"TTS_SPEED" envDefault:"1.2"`
	TtsPitch     float64 `env:"TTS_PITCH" envDefault:"1.0"`
	TtsDevice    string  `env:"TTS_DEVICE" envDefault:"default"`
}{}

func init() {
	var dotenv string
	flag.StringVar(&dotenv, "env", ".env", "load .env file")
	flag.Parse()
	if err := godotenv.Load(dotenv); err != nil {
		log.Print(err)
	}
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}
}
