package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	width         int32
	height        int32
	cellWidth     int32
	aliveColor    rl.Color
	hoverColor    rl.Color
	deadColor     rl.Color
	currentMatrix [][]Cell
	previousCellX int32
	previousCellY int32
	started       bool
}

type Cell struct {
	state         string
	previousState string
}

func main() {

	game := Game{}
	game.Init()
	rl.InitWindow(game.width*game.cellWidth, game.height*game.cellWidth, "Conway's game of life")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		for i := 0; i < len(game.currentMatrix); i++ {
			for j := 0; j < len(game.currentMatrix[i]); j++ {
				currentColor := game.deadColor
				if game.currentMatrix[i][j].state == "alive" {
					currentColor = game.aliveColor
				} else if game.currentMatrix[i][j].state == "hover" {
					currentColor = game.hoverColor
				}
				rl.DrawRectangle(int32(j*int(game.cellWidth)), int32(i*int(game.cellWidth)), game.cellWidth, game.cellWidth, currentColor)
			}
		}

		game.Update()
		rl.ClearBackground(rl.White)
		rl.EndDrawing()
	}
}

func (g *Game) Init() {
	g.cellWidth = 20
	g.height = 50
	g.width = 50
	g.aliveColor = rl.White
	g.deadColor = rl.Black
	g.hoverColor = rl.Gray
	g.started = false

	matrix := make([][]Cell, g.height)
	for i := range matrix {
		matrix[i] = make([]Cell, g.width)
		for j := 0; j < int(g.width); j++ {
			matrix[i][j] = Cell{state: "dead"}
		}
	}
	g.currentMatrix = matrix
}

func (g *Game) Update() {

	if !g.started {
		y := rl.GetMousePosition().Y
		x := rl.GetMousePosition().X
		currentCellY := int32(y) / g.cellWidth
		currentCellX := int32(x) / g.cellWidth
		if rl.IsKeyPressed(rl.KeySpace) {
			g.started = true
			g.currentMatrix[currentCellY][currentCellX].state = g.currentMatrix[currentCellY][currentCellX].previousState

			fmt.Println(countNeighbours(&g.currentMatrix, 3, 3))
		}

		if currentCellY != g.previousCellY || currentCellX != g.previousCellX {
			g.currentMatrix[g.previousCellY][g.previousCellX].state = g.currentMatrix[g.previousCellY][g.previousCellX].previousState

			g.currentMatrix[currentCellY][currentCellX].previousState = g.currentMatrix[currentCellY][currentCellX].state

			g.currentMatrix[currentCellY][currentCellX].state = "hover"

			g.previousCellY = currentCellY
			g.previousCellX = currentCellX
		}
		if rl.IsMouseButtonPressed(0) {
			currentCellY := int32(y) / g.cellWidth
			currentCellX := int32(x) / g.cellWidth
			if g.currentMatrix[currentCellY][currentCellX].previousState == "alive" {
				g.currentMatrix[currentCellY][currentCellX].state = "dead"
				g.currentMatrix[currentCellY][currentCellX].previousState = "dead"
			} else {
				g.currentMatrix[currentCellY][currentCellX].state = "alive"
				g.currentMatrix[currentCellY][currentCellX].previousState = "alive"
			}
		}
	}

	if g.started {
		currentIndex := rand.Intn(7)
		switch currentIndex {
		case 0:
			g.aliveColor = rl.Green
		case 1:
			g.aliveColor = rl.Red
		case 2:
			g.aliveColor = rl.Blue
		case 3:
			g.aliveColor = rl.Yellow
		case 4:
			g.aliveColor = rl.Pink
		case 5:
			g.aliveColor = rl.Purple
		case 6:
			g.aliveColor = rl.White
		}
		nextGeneration := make([][]Cell, g.height)
		for i := range nextGeneration {
			nextGeneration[i] = make([]Cell, g.width)
			for j := 0; j < int(g.width); j++ {
				nextGeneration[i][j] = Cell{state: "dead"}
			}
		}

		for y := 0; y < len(g.currentMatrix); y++ {
			for x := 0; x < len(g.currentMatrix[0]); x++ {
				neighbours := countNeighbours(&g.currentMatrix, x, y)

				currentCell := g.currentMatrix[y][x]

				if currentCell.state == "alive" {
					if neighbours < 2 || neighbours > 3 {
						nextGeneration[y][x].state = "dead"
					} else {
						nextGeneration[y][x].state = "alive"
					}
				} else {
					if neighbours == 3 {
						nextGeneration[y][x].state = "alive"
					}
				}
			}
		}
		g.currentMatrix = nextGeneration
	}
}

func countNeighbours(g *[][]Cell, x, y int) int {
	count := 0
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	grid := *g
	rows := len(grid)
	cols := len(grid[0])
	for _, dir := range directions {
		newX := x + dir[0]
		newY := y + dir[1]
		if newX >= 0 && newX < rows && newY >= 0 && newY < cols {
			if grid[newY][newX].state == "alive" {
				count++
			}
		}
	}

	return count
}
