package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var vm *chip8

	InitDisplay()
	defer DeinitDisplay()

	InitAudio()
	defer DeinitAudio()

	const cyclesPerFrame = 8 //500Hz/60Hz

	//StartBeeper()
	var romLoaded bool = false

	for !rl.WindowShouldClose() {
		if rl.IsFileDropped() {
			list := rl.LoadDroppedFiles()
			fmt.Println(list[0])

			vm = newChip8Vm(list[0])
			rl.UnloadDroppedFiles()
			romLoaded = true
			StartBeeper()
		}

		if !romLoaded {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)
			DrawIntro()
			rl.EndDrawing()
			continue
		}

		UpdateKeyboard(vm)
		UpdateBeep()
		vm.tickTimers()

		for range cyclesPerFrame {
			vm.Cycle()
		}

   		if vm.soundTimer > 0 {
       		ResumeBeeper()
   		} else {
       		PauseBeeper()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		UpdateDisplay(&vm.display)
		rl.EndDrawing()
	}
}

func newChip8Vm(romPath string) *chip8 {
	vm := chip8{
		pc:       0x200,
		stack:    make([]uint16, 0, 16),
		keyboard: make([]bool, 16),
	}
	vm.loadFont()
	err := vm.loadROM(romPath)
	if err != nil {
		panic(err)
	}
	return &vm
}
