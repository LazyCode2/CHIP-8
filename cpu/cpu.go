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

func (cpu *Chip8) X() uint8 {
	return uint8((cpu.Opcode & 0x0F00) >> 8)
}

func (cpu *Chip8) Y() uint8 {
	return uint8((cpu.Opcode & 0x00F0) >> 4)
}

func (cpu *Chip8) NN() uint8 {
	return uint8(cpu.Opcode & 0x00FF)
}

func (cpu *Chip8) NNN() uint16 {
	return cpu.Opcode & 0x0FFF
}

func (cpu *Chip8) N() uint8 {
	return uint8(cpu.Opcode & 0x000F)
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
			if cpu.RegisterV[cpu.X()] == cpu.NN() {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break

		case 0x4000:
			if cpu.RegisterV[cpu.X()] != cpu.NN() {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break

		case 0x5000:
			if cpu.RegisterV[cpu.X()] == cpu.RegisterV[cpu.Y()] {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break
		
		case 0x6000:
			cpu.RegisterV[cpu.X()] = cpu.NN()
			cpu.ProgramCounter += 2
			break
		
		case 0x7000:
			result := uint16(cpu.RegisterV[cpu.X()]) + uint16(cpu.NN())
			cpu.RegisterV[cpu.X()] = uint8(result)
			cpu.ProgramCounter += 2
			break
		
		case 0x8000:
			switch cpu.Opcode & 0x000F {
				case 0x0000:
					cpu.RegisterV[cpu.Y()] = cpu.RegisterV[cpu.X()]
					cpu.ProgramCounter += 2
					break

				case 0x0001:

			}

	}
}