package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"gopkg.in/qml.v1"
)

type Grid struct {
	Grid       qml.Object
	StatusText qml.Object

	TileComp *Tile

	Tiles []*Tile

	ColCount int
	RowCount int
}

type JSONGrid struct {
	Rows  int        `json:"rows"`
	Cols  int        `json:"cols"`
	Tiles []JSONTile `json:"tiles"`
}

func (g *Grid) StatusFromTile(t *Tile) string {
	name := "open"
	switch t.Object.Int("type") {
	case 1:
		name = "wall"
	case 2:
		name = "nest"
	case 3:
		name = "food"
	}
	return fmt.Sprintf("%v %v %v", name, t.x, t.y)
}

func (g *Grid) createTile() *Tile {
	tile := &Tile{
		Object: g.TileComp.Object.Create(nil),
		x:      1,
	}
	tile.Object.Set("parent", g.Grid)
	return tile
}

func (g *Grid) SaveGrid(filename string) {
	filename = filename[7:]
	jg := &JSONGrid{
		Rows: g.RowCount,
		Cols: g.ColCount,
	}
	tiles := make([]JSONTile, 0, jg.Rows*jg.Cols)

	for _, v := range g.Tiles {
		c := v.Object.Int("count")
		if c == 0 { //skip non number node
			continue
		}
		t := JSONTile{
			Count: v.Object.Int("count"),
			Index: v.Object.Int("index"),
		}
		tiles = append(tiles, t)
	}

	jg.Tiles = tiles

	dat, err := json.Marshal(jg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(dat))
	err = ioutil.WriteFile(filename, dat, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully Saved", filename)
}

func (g *Grid) LoadGrid(input io.Reader) {
	r := bufio.NewReader(input)
	dat, err := r.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	var newg JSONGrid
	err = json.Unmarshal(dat, &newg)
	if err != nil {
		fmt.Println(err)
		return
	}
	g.BuildGrid(newg.Rows, newg.Cols)
	for _, t := range newg.Tiles {
		g.Tiles[t.Index].Object.Set("type", 0)
		g.Tiles[t.Index].Object.Set("count", t.Count)
	}
}

func (g *Grid) BuildGrid(rows, cols int) {
	for _, b := range g.Tiles {
		b.Object.Set("visible", false)
		b.Object.Destroy()
	}
	g.RowCount = rows
	g.ColCount = cols
	g.Grid.Set("columns", g.ColCount)

	fmt.Println("Building a", g.RowCount, g.ColCount, "grid")
	size := g.RowCount * g.ColCount
	g.Tiles = make([]*Tile, size, size)
	for n := 0; n < size; n++ {
		g.Tiles[n] = g.createTile()
		g.Tiles[n].Object.Set("index", n)
	}
}

func (g *Grid) ClearGrid() {
	for _, v := range g.Tiles {
		v.Object.Set("type", 0)
	}
}
