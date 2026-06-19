package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/LazyCode2/CHIP-8/cpu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	W     = 64
	H     = 32
	SCALE = 10
)

type Game struct {
	cpu *cpu.Chip8
}

func (g *Game) handleKeys() {
	g.cpu.Key = [16]uint8{}

	if ebiten.IsKeyPressed(ebiten.Key1) { g.cpu.Key[0x1] = 1 }
	if ebiten.IsKeyPressed(ebiten.Key2) { g.cpu.Key[0x2] = 1 }
	if ebiten.IsKeyPressed(ebiten.Key3) { g.cpu.Key[0x3] = 1 }
	if ebiten.IsKeyPressed(ebiten.Key4) { g.cpu.Key[0xC] = 1 }

	if ebiten.IsKeyPressed(ebiten.KeyQ) { g.cpu.Key[0x4] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyW) { g.cpu.Key[0x5] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyE) { g.cpu.Key[0x6] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyR) { g.cpu.Key[0xD] = 1 }

	if ebiten.IsKeyPressed(ebiten.KeyA) { g.cpu.Key[0x7] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyS) { g.cpu.Key[0x8] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyD) { g.cpu.Key[0x9] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyF) { g.cpu.Key[0xE] = 1 }

	if ebiten.IsKeyPressed(ebiten.KeyZ) { g.cpu.Key[0xA] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyX) { g.cpu.Key[0x0] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyC) { g.cpu.Key[0xB] = 1 }
	if ebiten.IsKeyPressed(ebiten.KeyV) { g.cpu.Key[0xF] = 1 }
}

func (g *Game) Update() error {
	g.handleKeys()

	for i := 0; i < 8; i++ {
		g.cpu.Emulate()
		fmt.Printf("PC=%04X OPCODE=%04X\n", g.cpu.ProgramCounter, g.cpu.Opcode)

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	const debugWidth = 220
	displayOffsetX := debugWidth

	// CHIP-8 Display
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if g.cpu.Display[x+y*W] {
				for dy := 0; dy < SCALE; dy++ {
					for dx := 0; dx < SCALE; dx++ {
						screen.Set(
							displayOffsetX+x*SCALE+dx,
							y*SCALE+dy,
							color.White,
						)
					}
				}
			}
		}
	}

	// LEFT PANEL
	left := fmt.Sprintf(
		"PC: %04X\n"+
			"SP: %02X\n"+
			"I : %04X\n\n"+
			"V0: %02X\n"+
			"V1: %02X\n"+
			"V2: %02X\n"+
			"V3: %02X\n"+
			"V4: %02X\n"+
			"V5: %02X\n"+
			"V6: %02X\n"+
			"V7: %02X\n"+
			"V8: %02X\n"+
			"V9: %02X\n"+
			"VA: %02X\n"+
			"VB: %02X\n"+
			"VC: %02X\n"+
			"VD: %02X\n"+
			"VE: %02X\n"+
			"VF: %02X",
		g.cpu.ProgramCounter,
		g.cpu.StackPointer,
		g.cpu.Index,
		g.cpu.RegisterV[0], g.cpu.RegisterV[1], g.cpu.RegisterV[2], g.cpu.RegisterV[3],
		g.cpu.RegisterV[4], g.cpu.RegisterV[5], g.cpu.RegisterV[6], g.cpu.RegisterV[7],
		g.cpu.RegisterV[8], g.cpu.RegisterV[9], g.cpu.RegisterV[10], g.cpu.RegisterV[11],
		g.cpu.RegisterV[12], g.cpu.RegisterV[13], g.cpu.RegisterV[14], g.cpu.RegisterV[15],
	)

	ebitenutil.DebugPrintAt(screen, left, 10, 10)

	// RIGHT PANEL
	right := fmt.Sprintf(
		"Opcode: %04X\n"+
			"Delay : %02X\n"+
			"Sound : %02X\n\n"+
			"Stack:\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X\n"+
			"%04X",
		g.cpu.Opcode,
		g.cpu.DelayTimer,
		g.cpu.SoundTimer,
		g.cpu.Stack[0], g.cpu.Stack[1], g.cpu.Stack[2], g.cpu.Stack[3],
		g.cpu.Stack[4], g.cpu.Stack[5], g.cpu.Stack[6], g.cpu.Stack[7],
		g.cpu.Stack[8], g.cpu.Stack[9], g.cpu.Stack[10], g.cpu.Stack[11],
		g.cpu.Stack[12], g.cpu.Stack[13], g.cpu.Stack[14], g.cpu.Stack[15],
	)

	ebitenutil.DebugPrintAt(
		screen,
		right,
		displayOffsetX+W*SCALE+20,
		10,
	)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return W*SCALE + 440, H*SCALE
}

func main() {
	c := &cpu.Chip8{}
	c.LoadROM("./rom/PONG")

	game := &Game{cpu: c}

	ebiten.SetWindowSize(W*SCALE+440, H*SCALE)
	ebiten.SetWindowTitle("CHIP-8")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}