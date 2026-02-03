package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var stream rl.AudioStream

func InitAudio() {
	// Maybe load sound here
	rl.InitAudioDevice()
	stream = rl.LoadAudioStream(44000, 16, 2)
	if !rl.IsAudioStreamValid(stream) {
		panic("Failed to init sound")
	}

}

func DeinitAudio() {
	rl.UnloadAudioStream(stream)
	rl.CloseAudioDevice()
}

func StartBeep() {
	if rl.IsAudioStreamValid(stream) {
		rl.PlayAudioStream(stream)
	}

}
func StopBeep() {
	if rl.IsAudioStreamValid(stream) {
		rl.StopAudioStream(stream)
	}

}
