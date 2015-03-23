package grid

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

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

func (g *Grid) LoadGrid(input io.Reader) error {
	r := bufio.NewReader(input)
	dat, err := r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return errors.New("error reading " + err.Error())
	}
	var jgrid jsonGrid
	err = json.Unmarshal(dat, &jgrid)
	if err != nil {
		return errors.New("error unmarshalling " + err.Error())
	}
	g.BuildGrid(jgrid.Rows, jgrid.Cols)
	for _, t := range jgrid.Tiles {
		g.tiles[t.Index].open = true
		g.tiles[t.Index].count = t.Count
	}
	return nil
}

func (g *Grid) Json() ([]byte, error) {
	jTiles := make([]jsonTile, 0)
	for i, t := range g.tiles {
		if t.count > 0 {
			jTiles = append(jTiles, jsonTile{
				Count: t.count,
				Index: i,
			})
		}
	}
	if len(jTiles) == 0 {
		jTiles = nil
	}
	jGrid := &jsonGrid{
		Rows:  g.rows,
		Cols:  g.cols,
		Tiles: jTiles,
	}
	return json.Marshal(jGrid)
}
