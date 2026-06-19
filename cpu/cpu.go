package cpu

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"

	"github.com/LazyCode2/CHIP-8/utils"
)

var Logger = utils.New(utils.DEBUG)

var FontSet = [80]uint8{
    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
    0x20, 0x60, 0x20, 0x20, 0x70, // 1
    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
    0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8 struct {
	Stack[16]	 uint16
	StackPointer uint16  

	Memory[4096]  uint8
	RegisterV[16] uint8

	ProgramCounter uint16
	Opcode 		   uint16
	Index		   uint16

	Display[64 * 32] bool

	Key[16]	uint8

	DelayTimer uint8
	SoundTimer uint8
}

func (cpu *Chip8) Init() {
	cpu.ProgramCounter = 0x200
	cpu.StackPointer = 0

	for i := range FontSet {
		cpu.Memory[i] = FontSet[i]
	}
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
		case OpcodeSYS:
			switch cpu.Opcode & 0x00FF {
				case OpcodeCLS:
					cpu.ExecuteCLS()
					cpu.ProgramCounter += 2
					break
				case OpcodeRET:
					cpu.ExecuteRET()
				    cpu.ProgramCounter += 2
				    break
			}
		case OpcodeJP:
			cpu.ProgramCounter = cpu.Opcode & 0x0FFF
			break

		case OpcodeCALL:
			cpu.Stack[cpu.StackPointer] = cpu.ProgramCounter
			cpu.StackPointer++
			cpu.ProgramCounter = cpu.Opcode & 0xFFF
		
		case OpcodeSEByte:
			if cpu.RegisterV[cpu.X()] == cpu.NN() {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break

		case OpcodeSNEByte:
			if cpu.RegisterV[cpu.X()] != cpu.NN() {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break

		case OpcodeSEReg:
			if cpu.RegisterV[cpu.X()] == cpu.RegisterV[cpu.Y()] {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}
			break
		
		case OpcodeLDByte:
			cpu.RegisterV[cpu.X()] = cpu.NN()
			cpu.ProgramCounter += 2
			break
		
		case OpcodeADDByte:
			result := uint16(cpu.RegisterV[cpu.X()]) + uint16(cpu.NN())
			cpu.RegisterV[cpu.X()] = uint8(result)
			cpu.ProgramCounter += 2
			break
		
		case 0x8000:
			switch cpu.Opcode & 0x000F {
				case OpcodeLDReg:
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break

				case OpcodeOR:
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] | cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break
				case OpcodeAND:
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] & cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break

				case OpcodeXOR:
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] ^ cpu.RegisterV[cpu.Y()]
					cpu.ProgramCounter += 2
					break

				case OpcodeADD:
					sum := uint16(cpu.RegisterV[cpu.X()]) + uint16(cpu.RegisterV[cpu.Y()])
					
					if sum > 255 {
						cpu.RegisterV[0xF] = 1
					} else {
						cpu.RegisterV[0xF] = 0
					}

					cpu.RegisterV[cpu.X()] = uint8(sum)

					cpu.ProgramCounter += 2
					break

				case OpcodeSUB:
					if cpu.RegisterV[cpu.Y()] > cpu.RegisterV[cpu.X()] {
						cpu.RegisterV[0xF] = 0
					} else {
						cpu.RegisterV[0xF] = 1
					}

					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] - cpu.RegisterV[cpu.Y()]

					cpu.ProgramCounter += 2
					break

				case OpcodeSHR:
					cpu.RegisterV[0xF] = cpu.RegisterV[cpu.X()] & 0x1
					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.X()] >> 1

					cpu.ProgramCounter += 2
					break

				case OpcodeSUBN:
					if cpu.RegisterV[cpu.X()] > cpu.RegisterV[cpu.Y()] {
						cpu.RegisterV[0xF] = 0
					} else {
						cpu.RegisterV[0xF] = 1
					}

					cpu.RegisterV[cpu.X()] = cpu.RegisterV[cpu.Y()] - cpu.RegisterV[cpu.X()]

					cpu.ProgramCounter += 2
					break
				
				case OpcodeSHL:
					cpu.RegisterV[0xF] = cpu.RegisterV[cpu.X()] >> 7
					cpu.RegisterV[cpu.X()] <<= 1
					cpu.ProgramCounter += 2
					break
			}

		case OpcodeSNEReg:
			if cpu.RegisterV[cpu.X()] != cpu.RegisterV[cpu.Y()] {
				cpu.ProgramCounter += 4
			} else {
				cpu.ProgramCounter += 2
			}

			break
	
		case OpcodeLDI:
			cpu.Index = cpu.NNN()

			cpu.ProgramCounter += 2
			break
		
		case OpcodeJPV0:
			cpu.ProgramCounter = cpu.NNN() + uint16(cpu.RegisterV[0])
			break
	
		case OpcodeRND:
			cpu.RegisterV[cpu.X()] = uint8(rand.Int()) & cpu.NN()

			cpu.ProgramCounter += 2
			break

		case OpcodeDRW:
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
		case 0xE000:
				switch cpu.Opcode & 0x00FF {
					case OpcodeSKP:
						if cpu.Key[cpu.RegisterV[cpu.X()]] != 0 {
							cpu.ProgramCounter += 4
						} else {
							cpu.ProgramCounter += 2
						}

						break

					case OpcodeSKNP:
						if cpu.Key[cpu.RegisterV[cpu.X()]] == 0 {
							cpu.ProgramCounter += 4
						} else {
							cpu.ProgramCounter += 2
						}

						break
				} 
		
		case 0xF000:
			switch  cpu.Opcode & 0x00FF {
				case OpcodeLDVxDT:
					cpu.RegisterV[cpu.X()] = cpu.DelayTimer

					cpu.ProgramCounter += 2
					break

				case OpcodeLDVxK:
				    keyPressed := false
				    for i := 0; i < 16; i++ {
				        if cpu.Key[i] != 0 {
				            cpu.RegisterV[cpu.X()] = uint8(i)
				            keyPressed = true
				            break
				        }
				    }

				    if !keyPressed {
				        return 
				    }
				    
				    cpu.ProgramCounter += 2
				    break

				case OpcodeLDDT:
					cpu.DelayTimer = cpu.RegisterV[cpu.X()]

					cpu.ProgramCounter += 2
					break
			
				case OpcodeLDST:
					cpu.SoundTimer = cpu.RegisterV[cpu.X()]

					cpu.ProgramCounter += 2
					break

				case OpcodeADDI:
					cpu.Index += uint16(cpu.RegisterV[cpu.X()])

					cpu.ProgramCounter += 2
					break

				case OpcodeLDF:
					cpu.Index = uint16(cpu.RegisterV[cpu.X()]) * 0x5

					cpu.ProgramCounter += 2
					break
				case OpcodeLDB:
					x := cpu.RegisterV[cpu.X()]

					cpu.Memory[cpu.Index+2] = x % 10
					x /= 10
					cpu.Memory[cpu.Index+1] = x % 10
					x /= 10
					cpu.Memory[cpu.Index] = x

					cpu.ProgramCounter += 2
					break

				case OpcodeLDIVx:
					for i := 0; i <= int(cpu.X()); i++ {
						cpu.Memory[int(cpu.Index)+i] = cpu.RegisterV[i]
					}

					cpu.Index += uint16(cpu.X()) + 1

					cpu.ProgramCounter += 2
					break

				case OpcodeLDVxI:
					for i := 0; i <= int(cpu.X()); i++ {
						cpu.RegisterV[i] = cpu.Memory[int(cpu.Index)+i]
					}

					cpu.Index += uint16(cpu.X()) + 1

					cpu.ProgramCounter += 2
					break
				
				default:
					Logger.Error("Unknown opcode [%v]",cpu.Opcode)
			}

			break 
		}

		if cpu.DelayTimer > 0 {
			cpu.DelayTimer--
		}
}

func (cpu *Chip8) LoadROM(path string) {
	cpu.Init()

	file, err := os.Open(path)
	if err != nil {
		Logger.Error("failed to open ROM: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		Logger.Error("failed to read ROM: %v", err)
	}

	for i, b := range data {
		addr := 0x200 + i
		if addr >= len(cpu.Memory) {
			break
		}
		cpu.Memory[addr] = b
	}

	Logger.Info("ROM Loaded!")
}

// Execute Functions

func (cpu *Chip8) ExecuteCLS() {
	for i := range cpu.Display {
		cpu.Display[i] = false
	}
}

func (cpu *Chip8) ExecuteRET() {
    cpu.StackPointer--
    cpu.ProgramCounter = cpu.Stack[cpu.StackPointer]
}

func (cpu *Chip8) DumpMemory(start, end uint16) {
	for addr := start; addr < end; addr += 8 {
		fmt.Printf("%04X: ", addr)

		for i := 0; i < 8 && addr+uint16(i) < end; i++ {
			fmt.Printf("%02X ", cpu.Memory[addr+uint16(i)])
		}

		fmt.Println()
	}
}