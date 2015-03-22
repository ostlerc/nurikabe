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

func (g *Grid) reset() {
	for _, t := range g.Tiles {
		t.Reset()
	}
}

//TODO: add difficulty parameter
func (g *Grid) Generate(minGardens, gardenSize int) {
	tileMap := make(mapset)
	for {
		g.reset()
		for i, _ := range g.Tiles {
			tileMap[i] = closed
		}

		c := 0
		for ; g.placeGarden(R.Intn(gardenSize)+2, tileMap); c++ {
		}

		if c < minGardens {
			continue
		}

		for i, t := range g.Tiles {
			if tileMap[i] == opened {
				t.Properties.Set("type", 0)
			} else {
				t.Properties.Set("type", 1)
			}
		}

		if g.CheckWin() {
			return
		}
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
	g.Tiles[i].Properties.Set("type", 0)
	g.Tiles[i].Properties.Set("count", len(tiles))

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
	steps := []int{1, -1, g.Cols, -g.Cols}

	if i/g.Cols == g.Rows-1 { // bottom of grid
		steps = remove(g.Cols, steps)
	}

	if i < g.Cols { // top of grid
		steps = remove(-g.Cols, steps)
	}

	if i%g.Cols == g.Rows-1 { // right side of grid
		steps = remove(1, steps)
	}

	if i%g.Cols == 0 { // left side of grid
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
	for _, t := range g.Tiles {
		if !t.Open() {
			c++
		}
	}
	return c
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
