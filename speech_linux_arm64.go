//go:build linux && arm64

package main

const (
	downloadFile = "download-linux-arm64"
	IsWindows    = false
)

var voiceVoxFiles = []string{
	"libvoicevox_core.so",
	"open_jtalk_dic_utf_8-1.11",
	"model",
	"libonnxruntime.so.1.14.0",
}

func voiceVoxPreSetup() error { return nil }
