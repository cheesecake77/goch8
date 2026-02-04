package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	var vm *chip8
	romLoaded := false
	InitDisplay()
	defer DeinitDisplay()
	InitAudio()
	defer DeinitAudio()
	var halt chan bool

	//StartBeeper()
	for !rl.WindowShouldClose() {
		if rl.IsFileDropped() {
			list := rl.LoadDroppedFiles()
			fmt.Println(list[0])
			if halt != nil {
				close(halt)
			}

			vm = newChip8Vm(list[0])
			rl.UnloadDroppedFiles()
			romLoaded = true
			halt = make(chan bool)
			StartBeeper()
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

			go func(localVm *chip8, halt chan bool) {
				delayTicker := time.NewTicker(time.Second / 60)
				defer delayTicker.Stop()
				for {
					select {
					case <-halt:
						return
					case <-delayTicker.C:
						i := localVm.delayTimer.Load()
						if i > 0 {
							localVm.delayTimer.Store(uint32(i - 1))

						}
					}
				}
			}(vm, halt)

			go func(localVm *chip8, halt chan bool) {
				soundTicker := time.NewTicker(time.Second / 60)
				defer soundTicker.Stop()
				for {
					select {
					case <-halt:
						return
					case <-soundTicker.C:
						i := localVm.soundTimer.Load()
						if i > 0 {
							localVm.soundTimer.Store(uint32(i - 1))

						}
					}
				}
			}(vm, halt)
		}
		UpdateBeep()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		if !romLoaded {
			DrawIntro()
		} else {
			vm.displayMU.RLock()
			localDisplay := vm.display
			UpdateDisplay(&localDisplay)
			vm.displayMU.RUnlock()
			UpdateKeyboard(vm)
			if vm.soundTimer.Load() > 0 {
				ResumeBeeper()
			} else {
				PauseBeeper()
			}
		}
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
