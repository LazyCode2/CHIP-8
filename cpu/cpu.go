package cpu

import "math/rand"

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

// Return Vx00 field
func (cpu *Chip8) X() uint8 {
	return uint8((cpu.Opcode & 0x0F00) >> 8)
}

// Return Vy00 field
func (cpu *Chip8) Y() uint8 {
	return uint8((cpu.Opcode & 0x00F0) >> 4)
}

// Return 00nn field
func (cpu *Chip8) NN() uint8 {
	return uint8(cpu.Opcode & 0x00FF)
}

// Return 0nnn field
func (cpu *Chip8) NNN() uint16 {
	return cpu.Opcode & 0x0FFF
}

// Return 000n field
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
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break

				case 0x0001:
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] | cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break
				case 0x0002:
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] & cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break

				case 0x0003:
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] ^ cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break

				case 0x0004:
					sum := uint16(cpu.RegisterV[cpu.X()]) + uint16(cpu.RegisterV[cpu.Y()])
					
					if sum > 255 {
						cpu.RegisterV[0xF] = 1
					} else {
						cpu.RegisterV[0xF] = 0
					}

					cpu.ProgramCounter += 2
					break

				case 0x0005:
					if cpu.RegisterV[cpu.Y()] > cpu.RegisterV[cpu.X()] {
						cpu.RegisterV[0xF] = 0
					} else {
						cpu.RegisterV[0xF] = 1
					}

					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] - cpu.RegisterV[cpu.Y()]

					cpu.ProgramCounter += 2
					break

				case 0x0006:
					cpu.RegisterV[0xF] = cpu.RegisterV[cpu.X()] & 0x1
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] >> 1

					cpu.ProgramCounter += 2
					break

				case 0x0007:
					if cpu.RegisterV[cpu.X()] > cpu.RegisterV[cpu.Y()] {
						cpu.RegisterV[0xF] = 0
					} else {
						cpu.RegisterV[0xF] = 1
					}

					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.Y()] - cpu.RegisterV[cpu.X()]

					cpu.ProgramCounter += 2
					break
				
				case 0x000E:
					cpu.RegisterV[0xF] = cpu.RegisterV[cpu.X()] >> 7
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] << 1
			
					cpu.ProgramCounter += 2
					break
			}

		case 0x9000:
			if cpu.RegisterV[cpu.X()] != cpu.RegisterV[cpu.Y()] {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}

			break
	
		case 0xA000:
			cpu.Index = cpu.NNN()

			cpu.ProgramCounter += 2
			break
		
		case 0xB000:
			cpu.ProgramCounter = cpu.NNN() + uint16(cpu.RegisterV[0])
			break
	
		case 0xC000:
			cpu.RegisterV[cpu.X()] = uint8(rand.Int()) & cpu.NN()

			cpu.ProgramCounter += 2
			break

		case 0xD000:
			var XCord = cpu.RegisterV[cpu.X()] % 64
			var YCord = cpu.RegisterV[cpu.Y()] % 32
			var Height = cpu.N()
			cpu.RegisterV[0xF] = 0

			for y := 0; y < int(Height); y++ {
			    sprite := cpu.Memory[int(cpu.Index)+y]

			    for x := range 8 {
			        if sprite&(0x80>>x) != 0 {
			            screenX := (int(XCord) + x) % 64
			            screenY := (int(YCord) + y) % 32
			            index := screenX + screenY*64

			            if cpu.Display[index] {
			                cpu.RegisterV[0xF] = 1
			            }

			            cpu.Display[index] = !cpu.Display[index]
			        }
		    	}
			}

			cpu.ProgramCounter += 2
			break
			
	}
}