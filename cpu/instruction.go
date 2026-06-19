package cpu

const (
	OpcodeCLS  uint16 = 0x00E0
	OpcodeRET  uint16 = 0x00EE
	OpcodeSYS  uint16 = 0x0000
	OpcodeJP   uint16 = 0x1000
	OpcodeCALL uint16 = 0x2000
	OpcodeSEByte uint16 = 0x3000
	OpcodeSNEByte uint16 = 0x4000
	OpcodeSEReg uint16 = 0x5000
	OpcodeLDByte uint16 = 0x6000
	OpcodeADDByte uint16 = 0x7000
	OpcodeLDReg uint16 = 0x8000
	OpcodeOR    uint16 = 0x8001
	OpcodeAND   uint16 = 0x8002
	OpcodeXOR   uint16 = 0x8003
	OpcodeADD   uint16 = 0x8004
	OpcodeSUB   uint16 = 0x8005
	OpcodeSHR   uint16 = 0x8006
	OpcodeSUBN  uint16 = 0x8007
	OpcodeSHL   uint16 = 0x800E
	OpcodeSNEReg uint16 = 0x9000
	OpcodeLDI uint16 = 0xA000
	OpcodeJPV0 uint16 = 0xB000
	OpcodeRND uint16 = 0xC000
	OpcodeDRW uint16 = 0xD000
	OpcodeSKP  uint16 = 0x009E
	OpcodeSKNP uint16 = 0x00A1
	OpcodeLDVxDT uint16 = 0x0007
	OpcodeLDVxK  uint16 = 0x000A
	OpcodeLDDT   uint16 = 0x0015
	OpcodeLDST   uint16 = 0x0018
	OpcodeADDI   uint16 = 0x001E
	OpcodeLDF    uint16 = 0x0029
	OpcodeLDB    uint16 = 0x0033
	OpcodeLDIVx  uint16 = 0x0055
	OpcodeLDVxI  uint16 = 0x0065
)