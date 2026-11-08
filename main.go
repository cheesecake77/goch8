// Chip8 interpreter/emulator in Go using Raylib
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

	romLoaded := false

	for !rl.WindowShouldClose() {
		if rl.IsFileDropped() {
			list := rl.LoadDroppedFiles()
			fmt.Println(list[0])

			vm = NewChip8Vm(list[0])
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


