package validator

import "fmt"

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
	state   []bool
	v       GridValidator
	rows    int
	cols    int
}

func Solve(d GridData, v GridValidator) []bool {
	ret := make([]bool, d.Rows()*d.Columns(), d.Rows()*d.Columns())
	gardens := make(map[int]int, len(ret))
	for i := 0; i < len(ret); i++ {
		gardens[i] = d.Count(i)
	}
	s := &nurikabeSolver{
		gardens: gardens,
		state:   ret,
		v:       v,
		rows:    d.Rows(),
		cols:    d.Columns(),
	}

	if s.dumbSolve(0) {
		fmt.Println("Solved", s.state)
		return s.state
	}

	return nil
}

func (n *nurikabeSolver) Open(i int) bool {
	return !n.state[i]
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

func (n *nurikabeSolver) dumbSolve(i int) bool {
	if i == len(n.state)-1 {
		n.state[i] = true
		if n.v.CheckWin(n) {
			return true
		}
		n.state[i] = false
		if n.v.CheckWin(n) {
			return true
		}
		return false
	}
	n.state[i] = true
	if n.dumbSolve(i + 1) {
		return true
	}
	n.state[i] = false
	if n.dumbSolve(i + 1) {
		return true
	}
	return false
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
			if x := n.markOpen(i, openTiles); x != c {
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

	c := n.markClosed(firstWall, found)
	if c != wallCount && Verbose {
		fmt.Println("wall", c, "!=", wallCount)
	}
	return c == wallCount
}

func (n *nurikabe) markOpen(i int, found map[int]bool) int {
	if i < 0 || i >= n.l {
		return 0
	}

	if _, ok := found[i]; ok || !n.d.Open(i) {
		return 0
	}

	found[i] = true
	ret := 1

	if i/n.d.Columns() != n.d.Rows()-1 { // not bottom of grid
		ret += n.markOpen(i+n.d.Columns(), found)
	}

	if i >= n.d.Columns() { // not top of grid
		ret += n.markOpen(i-n.d.Columns(), found)
	}

	if i%n.d.Columns() != n.d.Columns()-1 { // not right side of grid
		ret += n.markOpen(i+1, found)
	}

	if i%n.d.Columns() != 0 { // not left side of grid
		ret += n.markOpen(i-1, found)
	}

	return ret
}

func (n *nurikabe) markClosed(i int, found map[int]bool) int {
	if i < 0 || i >= n.l {
		return 0
	}

	if _, ok := found[i]; ok || n.d.Open(i) {
		return 0
	}

	found[i] = true
	ret := 1

	if i/n.d.Columns() != n.d.Rows()-1 { // not bottom of grid
		ret += n.markClosed(i+n.d.Columns(), found)
	}

	if i >= n.d.Columns() { // not top of grid
		ret += n.markClosed(i-n.d.Columns(), found)
	}

	if i%n.d.Columns() != n.d.Columns()-1 { // not right side of grid
		ret += n.markClosed(i+1, found)
	}

	if i%n.d.Columns() != 0 { // not left side of grid
		ret += n.markClosed(i-1, found)
	}

	return ret
}
