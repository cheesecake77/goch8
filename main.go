package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var vm *chip8
	romLoaded := false

	InitDisplay()
	defer DeinitDisplay()

	InitAudio()
	defer DeinitAudio()


	for !rl.WindowShouldClose() {
		if rl.IsFileDropped() {
			list := rl.LoadDroppedFiles()
			fmt.Println(list[0])
			vm = newChip8Vm("test")
			romLoaded = true
		}
		rl.BeginDrawing()
		if !romLoaded {
			rl.ClearBackground(rl.Black)
			rl.DrawText("Drag&Drop ROM here!", 130, 120, 20, rl.RayWhite)
		} else {
			// основная работа
			vm.cycle()
			if vm.sound_timer > 0{
				// beep
			}
		}
		rl.EndDrawing()
	}
}

func newChip8Vm(romPath string) *chip8 {
	vm := chip8{pc: 0x200, stack: make([]uint16, 16)}
	vm.loadFont()
	vm.loadROM(romPath)
	return &vm
}
