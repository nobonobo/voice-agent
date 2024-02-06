package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aethiopicuschan/nanoda"
	"github.com/ebitengine/oto/v3"
	resampler "github.com/hajimehoshi/ebiten/v2/audio/wav"
	"golang.org/x/sys/windows"
)

func init() {
	if err := windows.SetDllDirectory(config.VoiceVoxDir); err != nil {
		log.Fatal(err)
	}
}

func generate(s nanoda.Synthesizer, words string) ([]byte, error) {
	aq, err := s.CreateAudioQuery(words, nanoda.StyleId(config.ActorID))
	if err != nil {
		return nil, err
	}
	aq.SpeedScale = config.TtsSpeed
	for _, p := range aq.AccentPhrases {
		if p.PauseMora != nil {
			p.PauseMora.VowelLength *= config.TtsPause
		}
	}
	w, err := s.Synthesis(aq, nanoda.StyleId(config.ActorID))
	if err != nil {
		return nil, err
	}
	defer w.Close()
	return io.ReadAll(w)
}

func playback(ctxOto *oto.Context, s nanoda.Synthesizer, words string) error {
	b, err := generate(s, words)
	if err != nil {
		return err
	}
	decoded, err := resampler.DecodeWithoutResampling(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	p := ctxOto.NewPlayer(decoded)
	defer p.Close()
	p.Play()
	tick := time.NewTicker(10 * time.Millisecond)
	for range tick.C {
		if !p.IsPlaying() {
			break
		}
	}
	return nil
}

func loadDict(ud *nanoda.UserDict) error {
	fp, err := os.Open(filepath.Join(config.VoiceVoxDir, "user_dict.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		tp, err := strconv.Atoi(fields[2])
		if err != nil {
			tp = 1
		}
		acc, err := strconv.Atoi(fields[3])
		if err != nil {
			acc = 0
		}
		ud.AddWord(nanoda.Word{
			Surface:       fields[0],
			Pronunciation: fields[1],
			WordType:      nanoda.WordType(tp), // PROPER_NOUN=0, COMMON_NOUN=1, VERB=2, ADJECTIVE=3
			AccentType:    uint64(acc),
		})
		log.Println("add:", fields)
	}
	return ud.Use()
}

func TTS(ctx context.Context, input <-chan string) error {
	log.Println("TTS Engine started")
	defer log.Println("TTS Engine stopped")
	ctxOto, _, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   48000,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return err
	}
	vv, err := nanoda.NewVoicevox(
		filepath.Join(config.VoiceVoxDir, "voicevox_core.dll"),
		filepath.Join(config.VoiceVoxDir, "open_jtalk_dic_utf_8-1.11"),
		filepath.Join(config.VoiceVoxDir, "model"))
	if err != nil {
		return err
	}
	dict := vv.NewUserDict()
	if err := loadDict(dict); err != nil {
		return err
	}
	s, err := vv.NewSynthesizer()
	if err != nil {
		return err
	}
	if err := s.LoadModelsFromStyleId(nanoda.StyleId(config.ActorID)); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case text := <-input:
			log.Println("->", text)
			for i := 0; i < 3; i++ {
				if err := playback(ctxOto, s, text); err != nil {
					log.Print(err)
					continue
				}
				break
			}
		}
	}
}
