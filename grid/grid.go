package grid

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/ostlerc/nurikabe/tile"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

type grid struct {
	grid       qml.Object
	statusText qml.Object
	valid      validator.GridValidator

	tiles []*tile.Tile

	cols int
	rows int
}

// json member variables must be external for unmarshalling
type jsonGrid struct {
	Rows  int        `json:"rows"`
	Cols  int        `json:"cols"`
	Tiles []jsonTile `json:"tiles"`
}

type jsonTile struct {
	Count int `json:"count,omitempty"`
	Index int `json:"index,omitempty"`
}

func New(v validator.GridValidator, g qml.Object, status qml.Object) *grid {
	return &grid{
		valid:      v,
		grid:       g,
		statusText: status,
	}
}

func (g *grid) CheckWin() bool {
	return g.valid.CheckWin(g.tiles, g.rows, g.cols)
}

func (g *grid) LoadGrid(input io.Reader) error {
	r := bufio.NewReader(input)
	dat, err := r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return errors.New("error reading " + err.Error())
	}
	var newg jsonGrid
	err = json.Unmarshal(dat, &newg)
	if err != nil {
		return errors.New("error unmarshalling " + err.Error())
	}
	g.BuildGrid(newg.Rows, newg.Cols)
	for _, t := range newg.Tiles {
		g.tiles[t.Index].Properties.Set("type", 0)
		g.tiles[t.Index].Properties.Set("count", t.Count)
	}
	return nil
}

func (g *grid) BuildGrid(rows, cols int) {
	for _, b := range g.tiles {
		b.Properties.Set("visible", false)
		b.Properties.Destroy()
	}
	g.rows = rows
	g.cols = cols
	g.grid.Set("columns", g.cols)

	size := g.rows * g.cols
	g.tiles = make([]*tile.Tile, size, size)
	for n := 0; n < size; n++ {
		g.tiles[n] = tile.New(g.grid)
		g.tiles[n].Properties.Set("index", n)
	}
}
