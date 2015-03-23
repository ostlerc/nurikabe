package grid

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/ostlerc/nurikabe/tile"
	"github.com/ostlerc/nurikabe/validator"
)

var R = rand.New(rand.NewSource(time.Now().UnixNano()))

type Grid struct {
	parent interface{}

	tiles []*tile.Tile

	cols int
	rows int
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

func New(parent interface{}) *Grid {
	return &Grid{
		parent: parent,
		tiles:  make([]*tile.Tile, 0),
		rows:   0,
		cols:   0,
	}
}

func (g *Grid) Open(i int) bool {
	return g.tiles[i].Open()
}

func (g *Grid) Count(i int) int {
	return g.tiles[i].Count()
}

func (g *Grid) Rows() int {
	return g.rows
}

func (g *Grid) Columns() int {
	return g.cols
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
		g.tiles[t.Index].Properties.Set("type", 0)
		g.tiles[t.Index].Properties.Set("count", t.Count)
		g.tiles[t.Index].Properties.Set("index", t.Index)
	}
	return nil
}

func (g *Grid) BuildGrid(rows, cols int) {
	for _, b := range g.tiles {
		b.Properties.Set("visible", false)
		b.Properties.Destroy()
	}
	g.rows = rows
	g.cols = cols

	size := g.rows * g.cols
	g.tiles = make([]*tile.Tile, size, size)
	for n := 0; n < size; n++ {
		g.tiles[n] = tile.New(g.parent)
		g.tiles[n].Properties.Set("index", n)
	}
}

func (g *Grid) reset() {
	for _, t := range g.tiles {
		t.Reset()
	}
}

//TODO: add difficulty parameter
func (g *Grid) Generate(v validator.GridValidator, minGardens, gardenSize, base int) {
	tileMap := make(mapset)
	for {
		g.reset()
		for i, _ := range g.tiles {
			tileMap[i] = closed
		}

		c := 0
		for ; g.placeGarden(R.Intn(gardenSize)+base, tileMap); c++ {
		}

		if c < minGardens {
			continue
		}

		for i, t := range g.tiles {
			if tileMap[i] == opened {
				t.Properties.Set("type", 0)
			} else {
				t.Properties.Set("type", 1)
			}
		}

		if v.CheckWin(g) {
			return
		}
	}
}

func (g *Grid) Solve(v validator.GridValidator) {
	sol := validator.Solve(g, v)
	for i, t := range g.tiles {
		if sol[i] {
			t.Properties.Set("type", 1)
		} else {
			t.Properties.Set("type", 0)
		}
	}
}

func (g *Grid) Print() {
	for i := 0; i < len(g.tiles); i += g.cols {
		for j := 0; j < g.cols; j++ {
			fmt.Print(g.tiles[i+j].Properties.Int("type"), " ")
		}
		fmt.Println()
	}
}

const (
	opened = iota
	closed = iota
	sealed = iota
)

type mapset map[int]int

func (m mapset) Print(cols int) {
	for i := 0; i < len(m); i += cols {
		for j := 0; j < cols; j++ {
			fmt.Print(m[i+j], " ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g *Grid) placeGarden(max int, tileMap mapset) bool {
	i := -1
	for c := 0; c < 10; c++ {
		z := R.Intn(len(tileMap))
		if tileMap[z] == closed {
			i = z
			break
		}
	}
	if i == -1 {
		for k, v := range tileMap {
			if v == closed {
				i = k
				break
			}
		}
		if i == -1 {
			return false
		}
	}
	tiles := g.markOpen(i, max, tileMap)
	if len(tiles) < 2 {
		return false
	}
	g.tiles[i].Properties.Set("type", 0)
	g.tiles[i].Properties.Set("count", len(tiles))

	return true
}

func removeAt(i int, a []int) []int {
	a[i], a[len(a)-1], a = a[len(a)-1], 0, a[:len(a)-1]
	return a
}

func remove(v int, a []int) []int {
	for i, x := range a {
		if x == v {
			return removeAt(i, a)
		}
	}
	return a
}

func (g *Grid) markOpen(i, c int, tileMap mapset) []int {
	if c == 0 || tileMap[i] == sealed || tileMap[i] == opened {
		return []int{}
	}
	steps := []int{1, -1, g.cols, -g.cols}

	if i/g.cols == g.rows-1 { // bottom of grid
		steps = remove(g.cols, steps)
	}

	if i < g.cols { // top of grid
		steps = remove(-g.cols, steps)
	}

	if i%g.cols == g.cols-1 { // right side of grid
		steps = remove(1, steps)
	}

	if i%g.cols == 0 { // left side of grid
		steps = remove(-1, steps)
	}

	remainingSteps := make([]int, len(steps))
	copy(remainingSteps, steps)

	ret := make([]int, 0, c)
	ret = append(ret, i)
	c--
	tileMap[i] = opened
	for c > 0 && len(remainingSteps) > 0 {
		stepIndex := R.Intn(len(remainingSteps))
		v := remainingSteps[stepIndex] + i
		remainingSteps = removeAt(stepIndex, remainingSteps)

		tList := g.markOpen(v, c, tileMap)
		if l := len(tList); l > 0 {
			c -= l
			ret = append(ret, tList...)
		}
	}

	seal := func(x int) {
		if tileMap[x] == closed {
			tileMap[x] = sealed
		}
	}

	//seal up boundaries
	for _, s := range steps {
		seal(s + i)
	}
	return ret
}

func (g *Grid) closedCount() int {
	c := 0
	for _, t := range g.tiles {
		if !t.Open() {
			c++
		}
	}
	return c
}

func (g *Grid) Json() ([]byte, error) {
	jTiles := make([]jsonTile, 0)
	for _, t := range g.tiles {
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
		Rows:  g.rows,
		Cols:  g.cols,
		Tiles: jTiles,
	}
	return json.Marshal(jGrid)
}
