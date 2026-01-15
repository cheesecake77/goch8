package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func InitDisplay() {
	rl.InitWindow(512, 256, "Chip8")
	rl.SetTargetFPS(60)
}

func DeinitDisplay() {
	rl.CloseWindow()
}
