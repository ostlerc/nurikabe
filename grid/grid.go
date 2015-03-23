package grid

import (
	"math/rand"
	"time"

	"github.com/ostlerc/nurikabe/validator"
)

var R = rand.New(rand.NewSource(time.Now().UnixNano()))

type tile struct {
	open  bool
	count int
}

type Grid struct {
	tiles []*tile
	cols  int
	rows  int
}

func New() *Grid {
	return &Grid{
		tiles: nil,
		rows:  0,
		cols:  0,
	}
}

func (g *Grid) Toggle(i int) {
	g.tiles[i].open = !g.tiles[i].open
}

func (g *Grid) Open(i int) bool {
	return g.tiles[i].open
}

func (g *Grid) Count(i int) int {
	return g.tiles[i].count
}

func (g *Grid) Rows() int {
	return g.rows
}

func (g *Grid) Columns() int {
	return g.cols
}

func (g *Grid) BuildGrid(rows, cols int) {
	g.rows = rows
	g.cols = cols

	size := g.rows * g.cols
	g.tiles = make([]*tile, size, size)
	for n := 0; n < size; n++ {
		g.tiles[n] = &tile{open: true}
	}
}

func (g *Grid) Solve(v validator.GridValidator) {
	sol := validator.Solve(g, v)
	for i, t := range g.tiles {
		t.open = !sol[i]
	}
}
