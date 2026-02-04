package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	stream    rl.AudioStream
	isPlaying bool
	phase     float64
)

const (
	frequency  = 440.0
	sampleRate = 44100
)

func InitAudio() {
	rl.InitAudioDevice()

	stream = rl.LoadAudioStream(sampleRate, 32, 2)
	rl.SetAudioStreamVolume(stream, 0.1)
	rl.SetAudioStreamBufferSizeDefault(4096)
	rl.PauseAudioStream(stream)

	isPlaying = false
}

func DeinitAudio() {
	if rl.IsAudioStreamPlaying(stream) {
		rl.StopAudioStream(stream)
	}
	rl.UnloadAudioStream(stream)
	rl.CloseAudioDevice()
}

func StartBeeper() {
	if !rl.IsAudioStreamPlaying(stream) {
		rl.PlayAudioStream(stream)
		isPlaying = true
	}
}

func StopBeeper() {
	if rl.IsAudioStreamPlaying(stream) {
		rl.StopAudioStream(stream)
		isPlaying = false
	}
}

func PauseBeeper() {
	if !isPlaying {
		return
	}
	rl.PauseAudioStream(stream)
	isPlaying = false
}

func ResumeBeeper() {
	if isPlaying {
		return
	}
	rl.ResumeAudioStream(stream)
	isPlaying = true
}

func UpdateBeep() {
	if rl.IsAudioStreamProcessed(stream) {
		numFrames := 512
		samples := make([]float32, numFrames*2)

		for i := range numFrames {
			sample := float32(math.Sin(phase))
			samples[i*2] = sample
			samples[i*2+1] = sample

			phase += 2 * math.Pi * frequency / sampleRate
		}
		rl.UpdateAudioStream(stream, samples)
	}
}
