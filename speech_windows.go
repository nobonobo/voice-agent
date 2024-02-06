//go:build windows

package main

import "golang.org/x/sys/windows"

const (
	downloadFile = "download-windows-x64.exe"
	IsWindows    = true
)

var voiceVoxFiles = []string{
	"voicevox_core.dll",
	"onnxruntime_providers_shared.dll",
	"onnxruntime.dll",
	"open_jtalk_dic_utf_8-1.11",
	"model",
}

func voiceVoxPreSetup() error {
	if err := windows.SetDllDirectory(config.VoiceVoxDir); err != nil {
		return err
	}
	return nil
}
