package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/LazyCode2/CHIP-8/cpu"
	"github.com/hajimehoshi/ebiten/v2"
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

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if g.cpu.Display[x+y*W] {
				for dy := 0; dy < SCALE; dy++ {
					for dx := 0; dx < SCALE; dx++ {
						screen.Set(
							x*SCALE+dx,
							y*SCALE+dy,
							color.White,
						)
					}
				}
			}
		}
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return W * SCALE, H * SCALE
}

func main() {
	c := &cpu.Chip8{}
	c.LoadROM("./rom/danm8ku.ch8")

	game := &Game{cpu: c}

	ebiten.SetWindowSize(W*SCALE, H*SCALE)
	ebiten.SetWindowTitle("CHIP-8")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}