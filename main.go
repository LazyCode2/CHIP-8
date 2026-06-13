package main

import (
	"github.com/LazyCode2/CHIP-8/cpu"
)

func main() {
	CPU := cpu.Chip8{}
	CPU.LoadROM("./rom/PONG")
}