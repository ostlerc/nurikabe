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

func (g *Grid) CheckWin() {
	if !g.hasBlock() && g.singleWall() && g.gardensAreCorrect() {
		fmt.Println("WINNER")
	}
}

// This function help detect quad blocks
func (g *Grid) hasBlock() bool {
	for i, _ := range g.Tiles {
		if i/g.ColCount == g.RowCount-1 || //bottom of grid
			i%g.ColCount == g.ColCount-1 || //right side of grid
			g.openAt(i) ||
			g.openAt(i+1) ||
			g.openAt(i+g.ColCount) ||
			g.openAt(i+g.ColCount+1) {
			continue
		}
		return true
	}
	return false
}

func (g *Grid) gardensAreCorrect() bool {
	for i, _ := range g.Tiles {
		if c := g.Tiles[i].Count(); c > 0 {
			openTiles := make(map[int]bool)
			if x := g.markOpen(i, openTiles); x != c {
				return false
			}
		}
	}
	return true
}

func (g *Grid) singleWall() bool {
	firstWall := -1
	wallCount := 0
	for i, _ := range g.Tiles {
		if !g.openAt(i) {
			if firstWall == -1 {
				firstWall = i
			}
			wallCount++
		}
	}

	if firstWall == -1 || wallCount == 0 {
		return false
	}

	found := make(map[int]bool)

	return g.markClosed(firstWall, found) == wallCount
}

func (g *Grid) markOpen(i int, found map[int]bool) int {
	if i < 0 || i >= len(g.Tiles) {
		return 0
	}

	if _, ok := found[i]; ok || !g.openAt(i) {
		return 0
	}

	found[i] = true
	ret := 1

	if i/g.ColCount != g.RowCount-1 { // not bottom of grid
		ret += g.markOpen(i+g.ColCount, found)
	}

	if i >= g.ColCount { // not top of grid
		ret += g.markOpen(i-g.ColCount, found)
	}

	if i%g.ColCount != g.RowCount-1 { // not right side of grid
		ret += g.markOpen(i+1, found)
		ret += g.markOpen(i+g.ColCount+1, found)
		ret += g.markOpen(i-g.ColCount+1, found)
	}

	if i%g.ColCount != 0 { // not left side of grid
		ret += g.markOpen(i-1, found)
		ret += g.markOpen(i+g.ColCount-1, found)
		ret += g.markOpen(i-g.ColCount-1, found)
	}

	return ret
}

func (g *Grid) markClosed(i int, found map[int]bool) int {
	if i < 0 || i >= len(g.Tiles) {
		return 0
	}

	if _, ok := found[i]; ok || g.openAt(i) {
		return 0
	}

	found[i] = true
	ret := 1

	ret += g.markClosed(i+1, found)
	ret += g.markClosed(i-1, found)
	ret += g.markClosed(i+g.ColCount, found)
	ret += g.markClosed(i-g.ColCount, found)

	return ret
}

func (g *Grid) openAt(i int) bool {
	return g.Tiles[i].Open()
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

func (g *Grid) Clear() {
	for _, v := range g.Tiles {
		v.Object.Set("type", 0)
	}
}
