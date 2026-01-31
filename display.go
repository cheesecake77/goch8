package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const cell_size = 8

func InitDisplay() {
	rl.InitWindow(512, 256, "Chip8")
	rl.SetTargetFPS(60)
}

func DeinitDisplay() {
	rl.CloseWindow()
}

func UpdateDisplay(display *[32]uint64) {
	var x,y int32
	y = 0

	for _,v := range *display {
		x = 0

		for i:= 63; i >= 0; i-- {
			bit := (v >> i) & 1
			if bit == 1 {
				rl.DrawRectangle(x, y, cell_size, cell_size, rl.RayWhite)
			}
			x+=cell_size
		}

		y+=cell_size
	}
}
