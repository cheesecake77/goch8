package main

import (
	"math/rand/v2"
	"os"
	"sync"
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type chip8 struct {
	memory      [4096]uint8
	display     [32]uint64
	pc          uint16
	i           uint16
	stack       []uint16
	delay_timer uint8
	sound_timer uint8
	v           [16]uint8
	mu			sync.RWMutex
}

// Methods
func (vm *chip8) loadFont() {
	pointer := 0x50
	for _, value := range fontSet {
		vm.memory[pointer] = value
		pointer++
	}
}

func (vm *chip8) loadROM(path string) error {
data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    
    // ROMs start at 0x200
    for i, b := range data {
        vm.memory[0x200+i] = b
    }
    
    return nil
}

func (vm *chip8) togglePixel(x, y uint8) bool {
    mask := uint64(1) << (63 - x)
    collision := vm.display[y]&mask != 0
    vm.display[y] ^= mask
    return collision
}



func (vm *chip8) Cycle() {
	if vm.pc >= 0x0FFF {
    return
	}
	opcode := (uint16)(vm.memory[vm.pc])<<8 | (uint16)(vm.memory[vm.pc+1])
	vm.pc += 2

	x := uint8((opcode & 0x0F00) >> 8)
	y := uint8((opcode & 0x00F0) >> 4)
	n := (opcode & 0x000F)
	nn := uint8(opcode & 0x00FF)
	nnn := (opcode & 0x0FFF)

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			op00E0(vm)
		case 0x00EE:
			op00EE(vm)
		}
	case 0x1000:
		op1NNN(vm, nnn)
	case 0x2000:
		op2NNN(vm, nnn)
	case 0x3000:
		op3XNN(vm, x, nn)
	case 0x4000:
		op4XNN(vm, x, nn)
	case 0x5000:
		op5XY0(vm, x, y)
	case 0x6000:
		op6XNN(vm, x, nn)
	case 0x7000:
		op7XNN(vm, x, nn)
	case 0x8000:
		switch n {
		case 0x0:
			op8XY0(vm, x, y)
		case 0x1:
			op8XY1(vm, x, y)
		case 0x2:
			op8XY2(vm, x, y)
		case 0x3:
			op8XY3(vm, x, y)
		case 0x4:
			op8XY4(vm, x, y)
		case 0x5:
			op8XY5(vm, x, y)
		case 0x6:
			op8XY6(vm, x)
		case 0x7:
			op8XY7(vm, x, y)
		case 0xE:
			op8XYE(vm, x)
		}
	case 0x9000:
		op9XY0(vm, x, y)
	case 0xA000:
		opANNN(vm, nnn)
	case 0xB000:
		opBNNN(vm, nnn)
	case 0xC000:
		opCXNN(vm, x, nn)
	case 0xD000:
		vm.mu.Lock()
		opDXYN(vm, x, y, uint8(n))
		vm.mu.Unlock()	
	case 0xE000:
		switch nn {
		case 0x9E:
			opEX9E(vm, x)
		case 0xA1:
			opEXA1(vm, x)
		}
	case 0xF000:
		switch nn {
		case 0x07:
			opFX07(vm, x)
		case 0x0A:
			opFX0A(vm, x)
		case 0x15:
			opFX15(vm, x)
		case 0x18:
			opFX18(vm, x)
		case 0x1E:
			opFX1E(vm, x)
		case 0x29:
			opFX29(vm, x)
		case 0x33:
			opFX33(vm, x)
		case 0x55:
			opFX55(vm, x)
		case 0x65:
			opFX65(vm, x)
		}
	}
}

func (vm *chip8) push(val uint16) {
	vm.stack = append(vm.stack, val)
}

func (vm *chip8) pop() uint16 {
	if len(vm.stack) == 0 {
		// or panic
		return 0
	}
	i := len(vm.stack) - 1
	val := vm.stack[i]
	vm.stack = vm.stack[:i]
	return val
}

// ### Op codes ###

// Clear display
func op00E0(vm *chip8) {
	vm.display = [32]uint64{}
}

// Pop stack and set PC equal to the value
func op00EE(vm *chip8) {
	res := vm.pop()
	if res != 0 {
		vm.pc = res
	}
}

// Set PC to NNN
func op1NNN(vm *chip8, NNN uint16) {
	vm.pc = NNN
}

// Go to subroutine in NNN address. Save current PC to stack
func op2NNN(vm *chip8, NNN uint16) {
	vm.push(vm.pc)
	vm.pc = NNN
}

// Skip instruction if Vx == NN
func op3XNN(vm *chip8, X uint8, NN uint8) {
	if vm.v[X] == NN {
		vm.pc += 2
	}
}

// Skip instruction if Vx != NN
func op4XNN(vm *chip8, X, NN uint8) {
	if vm.v[X] != NN {
		vm.pc += 2
	}
}

// Skip instruction if Vx = vy
func op5XY0(vm *chip8, X, Y uint8) {
	if vm.v[X] == vm.v[Y] {
		vm.pc += 2
	}
}

// Set Vx to NN
func op6XNN(vm *chip8, X uint8, NN uint8) {
	vm.v[X] = NN
}

// Increment Vx by NN
func op7XNN(vm *chip8, X uint8, NN uint8) {
	vm.v[X] += NN
}

// Copy value from Vy to Vx
func op8XY0(vm *chip8, X, Y uint8) {
	vm.v[X] = vm.v[Y]
}

// Set Vx to Vx OR Vy
func op8XY1(vm *chip8, X, Y uint8) {
	vm.v[X] = vm.v[X] | vm.v[Y]
}

// Set Vx to Vx AND Vy
func op8XY2(vm *chip8, X, Y uint8) {
	vm.v[X] = vm.v[X] & vm.v[Y]
}

// Set Vx to Vx XOR Vy
func op8XY3(vm *chip8, X, Y uint8) {
	vm.v[X] = vm.v[X] ^ vm.v[Y]
}

// Set Vx to Vx + Vy and set VF to 1 if result is greater than 8 bits.
// Only lowest 8 bits is stored in Vx
func op8XY4(vm *chip8, X, Y uint8) {
	sum := (uint16)(vm.v[X]) + (uint16)(vm.v[Y])
	if sum > 255 {
		vm.v[0xf] = 1
	} else {
		vm.v[0xf] = 0
	}
	vm.v[X] = (uint8)(sum & 0x00FF)
}

// Sub Vx - Vy. If Vy > Vx set Vf to 1, else 0
func op8XY5(vm *chip8, X, Y uint8) {
	if vm.v[X] > vm.v[Y] {
		vm.v[0xf] = 1
	} else {
		vm.v[0xf] = 0
	}
	vm.v[X] = vm.v[X] - vm.v[Y]
}

// If most significant bit is 1 set VF to 1, else 0. Vx devided by 2
func op8XY6(vm *chip8, X uint8) {
	if (vm.v[X] & 0x01) == 1 {
		vm.v[0xf] = 1
	} else {
		vm.v[0xf] = 0
	}
	vm.v[X] = vm.v[X] >> 1
}

// Sub Vx - Vy. If Vy < Vx set Vf to 1, else 0
func op8XY7(vm *chip8, X, Y uint8) {
	if vm.v[X] < vm.v[Y] {
		vm.v[0xf] = 1
	} else {
		vm.v[0xf] = 0
	}
	vm.v[X] = vm.v[Y] - vm.v[X]
}

// If most significant bit is 1 set VF to 1, else 0. Vx devided by 2
func op8XYE(vm *chip8, X uint8) {
	if (vm.v[X] & 0x80) == 0x80 {
		vm.v[0xf] = 1
	} else {
		vm.v[0xf] = 0
	}
	vm.v[X] = vm.v[X] << 1
}

// Skip next instruction if Vx != Vy
func op9XY0(vm *chip8, X, Y uint8) {
	if vm.v[X] != vm.v[Y] {
		vm.pc += 2
	}
}

// pointer
// Set I to NNN
func opANNN(vm *chip8, NNN uint16) {
	vm.i = NNN
}

// Set PC equal to V0 + NNN
func opBNNN(vm *chip8, NNN uint16) {
	vm.pc = (uint16)(vm.v[0]) + NNN
}

// Set Vx to random byte AND NN
func opCXNN(vm *chip8, X uint8, NN uint8) {
	vm.v[X] = (uint8)(rand.IntN(256)) & NN
}

// Draw sprite
func opDXYN(vm *chip8, X, Y, N uint8) {
    x0 := vm.v[X] & 63 // 0..63
    y0 := vm.v[Y] & 31 // 0..31

    vm.v[0xF] = 0

    for row := uint8(0); row < N; row++ {
        y := (y0 + row) & 31
        sprite := vm.memory[vm.i+uint16(row)]

        for bit := uint8(0); bit < 8; bit++ {
            x := (x0 + bit) & 63

            if (sprite>>(7-bit))&1 == 1 {
                if vm.togglePixel(x, y) {
                    vm.v[0xF] = 1
                }
            }
        }
    }
}

// Skip instruction if key numbered Vx is pressed
func opEX9E(vm *chip8, X uint8) {
	// TODO
}

// Skip instruction if key numbered Vx is not pressed
func opEXA1(vm *chip8, X uint8) {
	// TODO
}

// Set Vx = delay timer
func opFX07(vm *chip8, X uint8) {
	vm.v[X] = vm.delay_timer
}

// Wait for key press and store it in Vx
func opFX0A(vm *chip8, X uint8) {
	//TODO
}

// Set delay timer to Vx
func opFX15(vm *chip8, X uint8) {
	vm.delay_timer = vm.v[X]
}

// Set sound timer to Vx
func opFX18(vm *chip8, X uint8) {
	vm.sound_timer = vm.v[X]
}

// Increment I by Vx
func opFX1E(vm *chip8, X uint8) {
	vm.i += uint16(vm.v[X])
}

// Set I to the address of sprite for X
func opFX29(vm *chip8, X uint8) {
	// Start of font + size of one sprite multiplied by X
	vm.i = uint16(0x50 + X*5)
}

// Store BCD of Vx in I, I+1 and I+2
func opFX33(vm *chip8, X uint8) {
	vm.memory[vm.i] = vm.v[X] / 100
	vm.memory[vm.i+1] = (vm.v[X] / 10) % 10
	vm.memory[vm.i+2] = vm.v[X] % 10
}

// Copy V0 to Vx to memory starting at I
func opFX55(vm *chip8, X uint8) {
	for i := 0; i <= int(X); i++ {
		vm.memory[vm.i+uint16(i)] = vm.v[i]
	}
}

// Write memory from I to registers Vx
func opFX65(vm *chip8, X uint8) {
	for i := 0; i <= int(X); i++ {
		vm.v[i] = vm.memory[vm.i+uint16(i)]
	}

}
