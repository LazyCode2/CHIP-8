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



