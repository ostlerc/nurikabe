package validator

import (
	"fmt"
)

type nurikabe struct {
	d GridData
	l int
}

func NewNurikabe() GridValidator {
	return &nurikabe{}
}

var Verbose = false

type nurikabeSolver struct {
	gardens map[int]int
	tiles   []bool
	v       GridValidator
	rows    int
	cols    int

	tileMap map[int]int
}

type gardenSolver struct {
	i, c      int
	workchan  chan bool
	readychan chan bool
	tileMap   map[int]int
	done      bool
}

func Solve(d GridData, v GridValidator, smart bool) []bool {
	l := d.Rows() * d.Columns()
	s := &nurikabeSolver{
		gardens: make(map[int]int, l),
		tiles:   make([]bool, l, l),
		v:       v,
		rows:    d.Rows(),
		cols:    d.Columns(),
		tileMap: make(map[int]int),
	}

	for i := 0; i < l; i++ {
		s.gardens[i] = d.Count(i)
	}

	if smart {
		gardenIndecies := make([]int, 0, l)
		for k, v := range s.gardens {
			if v > 0 {
				gardenIndecies = append(gardenIndecies, k)
			}
		}
		if s.gardenSolve(gardenIndecies) {
			return s.tiles
		}
	} else if s.dumbSolve(0) {
		return s.tiles
	}

	return nil
}

func (n *nurikabeSolver) Open(i int) bool {
	return !n.tiles[i]
}

func (n *nurikabeSolver) Count(i int) int {
	return n.gardens[i]
}

func (n *nurikabeSolver) Rows() int {
	return n.rows
}

func (n *nurikabeSolver) Columns() int {
	return n.cols
}

func Print(d GridData) {
	l := d.Rows() * d.Columns()
	for i := 0; i < l; i += d.Columns() {
		for j := 0; j < d.Columns(); j++ {
			if c := d.Count(i + j); c > 0 {
				fmt.Print(c, " ")
			} else if d.Open(i + j) {
				fmt.Print("o ")
			} else {
				fmt.Print("x ")
			}
		}
		fmt.Println()
	}
}

func (n *nurikabeSolver) dumbSolve(i int) bool {
	if i == len(n.tiles)-1 {
		n.tiles[i] = true
		if n.v.CheckWin(n) {
			return true
		}
		n.tiles[i] = false
		if n.v.CheckWin(n) {
			return true
		}
		return false
	}
	n.tiles[i] = true
	if n.dumbSolve(i + 1) {
		return true
	}
	n.tiles[i] = false
	if n.dumbSolve(i + 1) {
		return true
	}
	return false
}

var counter = 0

func (n *nurikabeSolver) gardenSolve(gardens []int) bool {
	if len(gardens) == 0 {

		for i := 0; i < len(n.tiles); i++ {
			n.tiles[i] = true
		}

		for k, _ := range n.tileMap {
			n.tiles[k] = false
		}

		if n.v.CheckWin(n) {
			return true
		}

		return false
	}

	g := &gardenSolver{
		i:         gardens[0],
		c:         n.Count(gardens[0]),
		workchan:  make(chan bool),
		readychan: make(chan bool),
		tileMap:   n.tileMap,
	}
	go func() {
		n.gardenPermutations(g)
		close(g.workchan)
		close(g.readychan)
	}()

	for {
		_, ok := <-g.workchan
		if !ok {
			break
		}
		if n.gardenSolve(gardens[1:]) {
			g.done = true
			g.readychan <- true
			return true

		}
		g.readychan <- true
	}

	return false
}

// Find all possible garden permutations for garden at index i, where there is still c count possibilities left.
// a bool will be sent on the workchan when tileMap contains the key indecies of a garden permutation.
func (n *nurikabeSolver) gardenPermutations(g *gardenSolver) {
	if g.done {
		return
	}

	if _, ok := g.tileMap[g.i]; ok {
		return
	}

	g.tileMap[g.i] = g.c
	g.c--
	defer func() {
		delete(g.tileMap, g.i)
		g.c++
	}()

	if g.c == 0 {
		g.workchan <- true
		<-g.readychan
		return
	}

	steps := make([]int, 0, 4)

	if g.i/n.Columns() != n.Rows()-1 { // not bottom of grid
		steps = append(steps, n.Columns())
	}

	if g.i >= n.Columns() { // not top of grid
		steps = append(steps, -n.Columns())
	}

	if g.i%n.Columns() != n.Columns()-1 { // not right side of grid
		steps = append(steps, 1)
	}

	if g.i%n.Columns() != 0 { // not left side of grid
		steps = append(steps, -1)
	}

	for _, step := range steps {
		g.i += step
		n.gardenPermutations(g)
		g.i -= step
	}
	return
}

func (n *nurikabe) CheckWin(d GridData) bool {
	n.d = d
	n.l = d.Rows() * d.Columns()
	if !n.hasBlock() && n.singleWall() && n.gardensAreCorrect() && n.openCountCorrect() {
		return true
	}
	return false
}

// This function detects quad blocks
func (n *nurikabe) hasBlock() bool {
	for i := 0; i < n.l; i++ {
		if i/n.d.Columns() == n.d.Rows()-1 || // bottom of grid
			i%n.d.Columns() == n.d.Columns()-1 || // right side of grid
			n.d.Open(i) ||
			n.d.Open(i+1) ||
			n.d.Open(i+n.d.Columns()) ||
			n.d.Open(i+n.d.Columns()+1) {
			continue
		}
		if Verbose {
			fmt.Println("Block err")
		}
		return true
	}
	return false
}

func (n *nurikabe) openCountCorrect() bool {
	open := 0
	expected := 0
	for i := 0; i < n.l; i++ {
		if n.d.Open(i) {
			open++
		}
		expected += n.d.Count(i)
	}
	if open != expected && Verbose {
		fmt.Println("open", open, "!=", expected)
	}
	return open == expected
}

// This function counts 4-connected open squares at each garden count spot
func (n *nurikabe) gardensAreCorrect() bool {
	for i := 0; i < n.l; i++ {
		if c := n.d.Count(i); c > 0 {
			openTiles := make(map[int]bool)
			if x := n.mark(i, openTiles, true); x != c {
				if Verbose {
					fmt.Println("gardens", x, "!=", c)
				}
				return false
			}
		}
	}
	return true
}

// This function determines if there is one contiguous 4-connected wall
func (n *nurikabe) singleWall() bool {
	firstWall := -1
	wallCount := 0
	for i := 0; i < n.l; i++ {
		if !n.d.Open(i) {
			if firstWall == -1 {
				firstWall = i
			}
			wallCount++
		}
	}

	if firstWall == -1 || wallCount == 0 {
		if Verbose {
			fmt.Println("early wall")
		}
		return false
	}

	found := make(map[int]bool)

	c := n.mark(firstWall, found, false)
	if c != wallCount && Verbose {
		fmt.Println("wall", c, "!=", wallCount)
	}
	return c == wallCount
}

func (n *nurikabe) mark(i int, found map[int]bool, open bool) int {
	if i < 0 || i >= n.l {
		return 0
	}

	if _, ok := found[i]; ok || n.d.Open(i) != open {
		return 0
	}

	found[i] = true
	ret := 1

	if i/n.d.Columns() != n.d.Rows()-1 { // not bottom of grid
		ret += n.mark(i+n.d.Columns(), found, open)
	}

	if i >= n.d.Columns() { // not top of grid
		ret += n.mark(i-n.d.Columns(), found, open)
	}

	if i%n.d.Columns() != n.d.Columns()-1 { // not right side of grid
		ret += n.mark(i+1, found, open)
	}

	if i%n.d.Columns() != 0 { // not left side of grid
		ret += n.mark(i-1, found, open)
	}

	return ret
}
