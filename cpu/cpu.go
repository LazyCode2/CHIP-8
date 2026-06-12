package cpu

type Chip8 struct {
	Stack[16]	 uint16
	StackPointer uint16  

	Memory[4096]  uint8
	RegisterV[16] uint8

	ProgramCounter uint16
	Opcode 		   uint16
	Index		   uint16

	Display[64 * 32] bool
}

func (cpu *Chip8) Init() {
	cpu.ProgramCounter = 0x200
}

func (cpu *Chip8) Emulate() {
	// Fetch opcode - 2 bytes
	cpu.Opcode = uint16(cpu.Memory[cpu.ProgramCounter]) << 8 | uint16(cpu.Memory[cpu.ProgramCounter + 1])

	switch cpu.Opcode & 0xF000 {
		case 0x0000:
			switch cpu.Opcode & 0x00FF {
				case 0x00E0:
					for i := range cpu.Display {
						cpu.Display[i] = false
					}
					cpu.ProgramCounter += 2
					break
			}
		case 0x1000:
			cpu.ProgramCounter = cpu.Opcode & 0x0FFF
			break

		case 0x2000:
			cpu.Stack[cpu.StackPointer] = cpu.ProgramCounter
			cpu.StackPointer++
			cpu.ProgramCounter = cpu.Opcode & 0xFFF
		
		case 0x3000:
			x := (cpu.Opcode & 0x0F00) >> 8
			nn := cpu.Opcode & 0x00FF

			if cpu.RegisterV[x] == uint8(nn) {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break

		case 0x4000:
			x := (cpu.Opcode & 0x0F00) >> 8
			nn := cpu.Opcode & 0x00FF

			if cpu.RegisterV[x] != uint8(nn) {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break

		case 0x5000:
			x := (cpu.Opcode & 0x0F00) >> 8
			y := (cpu.Opcode & 0x00F0) >> 4

			if cpu.RegisterV[x] == cpu.RegisterV[y] {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break
		
		case 0x6000:
			x := (cpu.Opcode & 0x0F00) >> 8
			nn := cpu.Opcode & 0x00FF
			cpu.RegisterV[x] = uint8(nn)
			cpu.ProgramCounter += 2
			break
	}

}



