package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

var (
	stream     rl.AudioStream
	isPlaying  bool
	phase      float64
	data       []int16   = make([]int16, maxSamples)
	writeBuf   []float32 = make([]float32, maxSamplesPerUpdate) // float32 для UpdateAudioStrea
	waveLength int       = 128 
)

const (
	frequency           = 440.0
	sampleRate          = 44100
	maxSamplesPerUpdate = 4096
	maxSamples          = 512
)

func InitAudio() {
	rl.InitAudioDevice()

	rl.SetAudioStreamBufferSizeDefault(maxSamplesPerUpdate)
	stream = rl.LoadAudioStream(sampleRate, 16, 1)
	rl.SetAudioStreamVolume(stream, 0.1)
	isPlaying = true

	for i := 0; i < waveLength*2; i++ {
		data[i] = int16(math.Sin(2*math.Pi*float64(i)/float64(waveLength)) * 32000)
	}
	for j := waveLength * 2; j < maxSamples; j++ {
		data[j] = 0
	}
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
		readCursor := 0
		writeCursor := 0

		for writeCursor < maxSamplesPerUpdate {
			writeLength := maxSamplesPerUpdate - writeCursor
			readLength := waveLength - readCursor

			if writeLength > readLength {
				writeLength = readLength
			}

			// Конвертируем int16 в float32 (normalized -1.0 to 1.0)
			for i := 0; i < writeLength; i++ {
				writeBuf[writeCursor+i] = float32(data[readCursor+i]) / 32768.0
			}

			readCursor = (readCursor + writeLength) % waveLength
			writeCursor += writeLength
		}

		rl.UpdateAudioStream(stream, writeBuf)
	}
}
