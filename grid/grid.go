package grid

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (g *Grid) BuildGrid(rows, cols int) {
	g.rows = rows
	g.cols = cols

	size := g.rows * g.cols
	g.tiles = make([]*tile, size, size)
	for n := 0; n < size; n++ {
		g.tiles[n] = &tile{open: true}
	}
}

func (g *Grid) reset() {
	for _, t := range g.tiles {
		t.open = true
		t.count = 0
	}
}

//TODO: add difficulty parameter
func (g *Grid) Generate(v validator.GridValidator, minGardens, gardenSize, base int) {
	tileMap := make(mapset, len(g.tiles))
	for {
		g.reset()
		for i := 0; i < len(g.tiles); i++ {
			tileMap[i] = closed
		}

		c := 0
		for ; g.placeGarden(R.Intn(gardenSize)+base, tileMap); c++ {
		}

		if c < minGardens {
			continue
		}

		for i, t := range g.tiles {
			t.open = tileMap[i] == opened
		}

		if v.CheckWin(g) {
			return
		}
	}
}

func (g *Grid) Solve(v validator.GridValidator) {
	sol := validator.Solve(g, v)
	for i, t := range g.tiles {
		t.open = !sol[i]
	}
}

func (g *Grid) Print() {
	for i := 0; i < len(g.tiles); i += g.cols {
		for j := 0; j < g.cols; j++ {
			if c := g.tiles[i+j].count; c > 0 {
				fmt.Print(c, " ")
			} else if g.tiles[i+j].open {
				fmt.Print("o ")
			} else {
				fmt.Print("x ")
			}
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
	g.tiles[i].open = true
	g.tiles[i].count = len(tiles)

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
		if !t.open {
			c++
		}
	}
	return c
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
