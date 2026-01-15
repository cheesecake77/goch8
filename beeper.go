package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func InitAudio(){
	// Maybe load sound here
	rl.InitAudioDevice()
}

func DeinitAudio(){
	rl.CloseAudioDevice()
}
