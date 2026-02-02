package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

func main() {
	var vm *chip8
	romLoaded := false
	InitDisplay()
	defer DeinitDisplay()
	InitAudio()
	defer DeinitAudio()
	// start chip8 cycles 600Hz
	var halt chan bool

	// start sound timer 60Hz
	// start delay timer 60Hz

	// 60Hz display loop
	for !rl.WindowShouldClose() {
		if rl.IsFileDropped() {
			list := rl.LoadDroppedFiles()
			fmt.Println(list[0])
			if halt != nil {
				close(halt)
				halt = nil
			}

			vm = newChip8Vm(list[0])
			rl.UnloadDroppedFiles()
			romLoaded = true
			halt = make(chan bool)
			go func(localVm *chip8, halt chan bool) {
				cpuTicker := time.NewTicker(time.Second / 500)
				defer cpuTicker.Stop()
				for {
					select {
					case <-halt:
						return
					case <-cpuTicker.C:
						localVm.Cycle()
					}
				}
			}(vm, halt)
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		if !romLoaded {
			// TODO move to display and scale to fit different cell_size
			rl.DrawText("Drag&Drop .c8 or .ch8 ROM here!", 130, 120, 20, rl.RayWhite)
		} else {
			vm.displayMU.RLock()
			localDisplay := vm.display
			vm.displayMU.RUnlock()
			UpdateDisplay(&localDisplay)
			UpdateKeyboard(vm)
		}
		rl.EndDrawing()
	}
}

func newChip8Vm(romPath string) *chip8 {
	vm := chip8{
		pc:        0x200,
		stack:     make([]uint16, 16),
		keyEvents: make(chan uint8, 1),
	}
	vm.loadFont()
	err := vm.loadROM(romPath)
	if err != nil {
		panic(err)
	}
	return &vm
}
