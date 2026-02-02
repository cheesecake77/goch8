package main

import (
//	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type KeyState uint8

const (
	NotPressed KeyState = iota
	Pressed
	Released
)
const cellSize = 16

func InitDisplay() {
	rl.InitWindow(64*cellSize, 32*cellSize, "Chip8")
	rl.SetTargetFPS(60)
}

func DeinitDisplay() {
	rl.CloseWindow()
}

func UpdateDisplay(display *[32]uint64) {
	var x, y int32
	y = 0

	for _, v := range *display {
		x = 0

		for i := 63; i >= 0; i-- {
			bit := (v >> i) & 1
			if bit == 1 {
				rl.DrawRectangle(x, y, cellSize, cellSize, rl.RayWhite)
			}
			x += cellSize
		}

		y += cellSize
	}
}

func UpdateKeyboard(vm *chip8) {

}

var keyMap = map[int32]uint8{
	rl.KeyOne:   0x1,
	rl.KeyTwo:   0x2,
	rl.KeyThree: 0x3,
	rl.KeyFour:  0xC,
	rl.KeyQ:     0x4,
	rl.KeyW:     0x5,
	rl.KeyE:     0x6,
	rl.KeyR:     0xD,
	rl.KeyA:     0x7,
	rl.KeyS:     0x8,
	rl.KeyD:     0x9,
	rl.KeyF:     0xE,
	rl.KeyZ:     0xA,
	rl.KeyX:     0x0,
	rl.KeyC:     0xB,
	rl.KeyV:     0xF,
}
