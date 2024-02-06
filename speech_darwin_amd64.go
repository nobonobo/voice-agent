//go:build darwin && !ios && amd64

package main

const (
	downloadFile = "download-osx-x64"
	IsWindows    = false
)

var voiceVoxFiles = []string{
	"libvoicevox_core.dylib",
	"libonnxruntime.1.14.0.dylib",
	"open_jtalk_dic_utf_8-1.11",
	"model",
}

func voiceVoxPreSetup() error { return nil }
