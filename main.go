package main

import (
	"fmt"

	"github.com/LazyCode2/CHIP-8/cpu"
)

func main() {
	CPU := cpu.Chip8{}
	CPU.LoadROM("./README.md")
	CPU.Emulate()
	
	// Dummy Instruction
	CPU.Memory[0x200] = 0x62
	CPU.Memory[0x201] = 0x34
	CPU.Memory[0x202] = 0x73
	CPU.Memory[0x203] = 0x14

	CPU.Emulate()
	CPU.Emulate()

	// Register Debug
	// cpu.Register[1] = 10
	for i, value := range CPU.RegisterV {
		fmt.Printf("V%X: %d\n", i, value)
	}

}