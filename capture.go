package main

import (
	"context"
	"log"
	"os/exec"
	"strings"
)

type writer chan<- []byte

func (w writer) Write(b []byte) (int, error) {
	buff := make([]byte, len(b))
	copy(buff, b)
	w <- buff
	return len(b), nil
}

func Capture(ctx context.Context, ch chan<- []byte) error {
	log.Println("capture start")
	defer log.Println("capture end")
	args := strings.Fields("sox -e signed-integer -b 16 -c 1 -r 16000 -t waveaudio 0 -e signed-integer -b 16 -t raw -")
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	//cmd.Stderr = os.Stderr
	cmd.Stdout = writer(ch)
	return cmd.Run()
}
