package grid

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/ostlerc/nurikabe/tile"
	"github.com/ostlerc/nurikabe/validator"
)

type Grid struct {
	valid  validator.GridValidator
	parent interface{}

	Tiles []*tile.Tile

	Cols int
	Rows int
}

// json member variables must be external for unmarshalling
type jsonGrid struct {
	Rows  int        `json:"rows"`
	Cols  int        `json:"cols"`
	Tiles []jsonTile `json:"tiles,omitempty"`
}

type jsonTile struct {
	Count int `json:"count,omitempty"`
	Index int `json:"index,omitempty"`
}

func New(v validator.GridValidator, parent interface{}) *Grid {
	return &Grid{
		valid:  v,
		parent: parent,
		Tiles:  make([]*tile.Tile, 0),
		Rows:   0,
		Cols:   0,
	}
}

func (g *Grid) CheckWin() bool {
	return g.valid.CheckWin(g.Tiles, g.Rows, g.Cols)
}

func (g *Grid) LoadGrid(input io.Reader) error {
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
		g.Tiles[t.Index].Properties.Set("type", 0)
		g.Tiles[t.Index].Properties.Set("count", t.Count)
		g.Tiles[t.Index].Properties.Set("index", t.Index)
	}
	return nil
}

func (g *Grid) BuildGrid(rows, cols int) {
	for _, b := range g.Tiles {
		b.Properties.Set("visible", false)
		b.Properties.Destroy()
	}
	g.Rows = rows
	g.Cols = cols

	size := g.Rows * g.Cols
	g.Tiles = make([]*tile.Tile, size, size)
	for n := 0; n < size; n++ {
		g.Tiles[n] = tile.New(g.parent)
		g.Tiles[n].Properties.Set("index", n)
	}
}

func (g *Grid) Json() ([]byte, error) {
	jTiles := make([]jsonTile, 0)
	for _, t := range g.Tiles {
		c := t.Properties.Int("count")
		i := t.Properties.Int("index")
		if c != 0 {
			jTiles = append(jTiles, jsonTile{
				Count: c,
				Index: i,
			})
		}
	}
	if len(jTiles) == 0 {
		jTiles = nil
	}
	jGrid := &jsonGrid{
		Rows:  g.Rows,
		Cols:  g.Cols,
		Tiles: jTiles,
	}
	return json.Marshal(jGrid)
}
