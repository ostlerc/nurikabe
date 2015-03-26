package validator

import (
	"fmt"
	"testing"
)

type vTest struct {
	closed []int
	count  map[int]int
	cols   int
	rows   int
	block  bool
	garden bool
	wall   bool
}

type tTile struct {
	open  bool
	count int
}

var tests = []*vTest{
	&vTest{[]int{}, map[int]int{}, 2, 2, false, true, false},
	&vTest{[]int{0, 1, 2, 3}, map[int]int{}, 2, 2, true, true, true},
	&vTest{[]int{0, 1, 3, 4}, map[int]int{}, 3, 3, true, true, true},
	&vTest{[]int{1, 2, 4, 5}, map[int]int{}, 3, 3, true, true, true},
	&vTest{[]int{3, 4, 6, 7}, map[int]int{}, 3, 3, true, true, true},
	&vTest{[]int{4, 5, 7, 8}, map[int]int{}, 3, 3, true, true, true},
	&vTest{[]int{1, 4, 7}, map[int]int{}, 3, 3, false, true, true},
	&vTest{[]int{0, 1, 2}, map[int]int{}, 3, 3, false, true, true},
	&vTest{[]int{4}, map[int]int{}, 3, 3, false, true, true},

	&vTest{[]int{1}, map[int]int{4: 9}, 3, 3, false, false, true},
	&vTest{[]int{0, 1, 2, 3, 8, 9, 10, 11}, map[int]int{4: 4}, 4, 3, false, true, false},
	&vTest{[]int{0, 1, 2, 3, 4, 8, 9, 10, 11}, map[int]int{4: 3}, 4, 3, false, false, true},
	&vTest{[]int{2, 5, 6, 7, 8, 9, 11, 13, 16, 19, 22, 24, 25, 26, 27, 28, 29, 33},
		map[int]int{10: 3, 17: 2, 18: 2, 21: 4, 30: 3, 34: 2}, 6, 6, false, true, false},
}

type fakeGridData struct {
	counts map[int]int
	closed map[int]bool
	rows   int
	cols   int
}

func (f *fakeGridData) Open(i int) bool {
	_, ok := f.closed[i]
	return !ok
}

func (f *fakeGridData) Count(i int) int {
	return f.counts[i]
}

func (f *fakeGridData) Rows() int {
	return f.rows
}

func (f *fakeGridData) Columns() int {
	return f.cols
}

func BuildNurikabe(v *vTest) *nurikabe {
	closed := make(map[int]bool, len(v.closed))
	for _, i := range v.closed {
		closed[i] = true
	}
	counts := make(map[int]int, len(v.count))
	for k, v := range v.count {
		counts[k] = v
	}
	data := &fakeGridData{
		rows:   v.rows,
		cols:   v.cols,
		closed: closed,
		counts: counts,
	}
	return &nurikabe{
		d: data,
		l: v.rows * v.cols,
	}
}

func TestHasBlock(t *testing.T) {
	for _, vtest := range tests {
		if BuildNurikabe(vtest).hasBlock() != vtest.block {
			t.Fatal("Failed")
		}
	}
}

func TestGarden(t *testing.T) {
	for _, vtest := range tests {
		if BuildNurikabe(vtest).gardensAreCorrect() != vtest.garden {
			t.Fatal("Failed", vtest)
		}
	}
}

func TestWall(t *testing.T) {
	for _, vtest := range tests {
		if BuildNurikabe(vtest).singleWall() != vtest.wall {
			t.Fatal("Failed", vtest)
		}
	}
}

func TestGardenPermutation(t *testing.T) {
	expected := map[int]int{2: 2, 3: 5, 4: 12}

	for c, v := range expected {
		n := &nurikabeSolver{
			rows: 5,
			cols: 5,
		}

		g := &gardenSolver{
			i:         0,
			c:         c,
			workchan:  make(chan bool),
			readychan: make(chan bool),
			tileMap:   make(map[int]int, 100),
			hash:      make(map[string]bool, 50000),
		}
		go func() {
			n.gardenPermutations(g)
			close(g.workchan)
		}()

		count := 0
		for {
			_, ok := <-g.workchan
			if !ok {
				break
			}
			count++
			g.readychan <- true
		}

		if count != v {
			t.Fatal("Invalid count", c, v, count)
		}
	}
}

func TestHash(t *testing.T) {
	fmt.Println(hash(map[int]int{32: 8, 50: 12, 15: 1, 2: 2, 3: 3, 5: 5, 4: 4, 1: 0, 100: 32}))
}

func BenchmarkHash(b *testing.B) {
	m := map[int]int{32: 8, 50: 12, 15: 1, 2: 2, 3: 3, 5: 5, 4: 4, 1: 0, 100: 32}
	for i := 0; i < b.N; i++ {
		hash(m)
	}
}

func TestGardenSolve(t *testing.T) {
	n := &nurikabe{}
	s := &nurikabeSolver{
		gardens: make(map[int]int, 25),
		tiles:   make([]bool, 25, 25),
		v:       n,
		rows:    5,
		cols:    5,
		tileMap: make(map[int]int, 100),
		hash:    make(map[string]bool, 10000),
	}
	s.gardens[1] = 5
	s.gardens[9] = 2
	s.gardens[21] = 4
	s.gardens[23] = 2
	gardens := []int{1, 9, 21, 23}
	if !s.gardenSolve(gardens) {
		t.Fatal("Failed to solve correctly")
	}
	Print(s)
}

func BenchmarkGardenSolve(b *testing.B) {
	n := &nurikabe{}
	s := &nurikabeSolver{
		gardens: make(map[int]int, 25),
		tiles:   make([]bool, 25, 25),
		v:       n,
		rows:    5,
		cols:    5,
		tileMap: make(map[int]int, 100),
		hash:    make(map[string]bool, 10000),
	}

	s.gardens[1] = 5
	s.gardens[9] = 2
	s.gardens[21] = 4
	s.gardens[23] = 2
	gardens := []int{1, 9, 21, 23}

	for i := 0; i < b.N; i++ {
		s.hash = make(map[string]bool, 10000)
		s.tileMap = make(map[int]int, 100)
		if !s.gardenSolve(gardens) {
			b.Fatal("Failed to solve correctly")
		}
	}
}

func TestPerms(t *testing.T) {
	if p := perms([]int{-1, 1, 5, -5}, make(map[string]bool, 100)); len(p) != 15 {
		t.Fatal("Incorrect perm count", len(p), p)
	}
}
