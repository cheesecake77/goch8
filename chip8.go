package main

import "fmt"

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
}

// Methods
func (vm *chip8) loadFont() {
	pointer := 0x50
	for _, value := range fontSet {
		vm.memory[pointer] = value
		pointer++
	}
}

func (vm *chip8) loadROM(path string) {

}

func (vm *chip8) cycle() {
	opcode := (uint16)(vm.memory[vm.pc])<<8 | (uint16)(vm.memory[vm.pc+1])
	vm.pc += 2

	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := (opcode & 0x000F)
	nn := (opcode & 0x00FF)
	nnn := (opcode & 0x0FFF)

	switch opcode & 0xF000 {
	case 0x0000:
		fmt.Println(x, y, n, nn, nnn)

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
	var sum uint16 = (uint16)(vm.v[X]) + (uint16)(vm.v[Y])
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

// pointer
// Set I to NNN
func opANNN(vm *chip8, NNN uint16) {
	vm.i = NNN
}

// Update display
func opDXYN(vm *chip8, X uint8, Y uint16) {
	// TODO
	// draw
}
