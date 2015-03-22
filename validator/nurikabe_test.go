package validator

import (
	"testing"

	"github.com/ostlerc/nurikabe/tile"
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

func BuildNurikabe(v *vTest) *nurikabe {
	tiles := make([]*tile.Tile, v.rows*v.cols, v.rows*v.cols)
	for i := 0; i < v.rows*v.cols; i++ {
		tiles[i] = tile.New(nil)
	}
	for _, i := range v.closed {
		tiles[i].Properties.Set("type", 1)
	}
	for k, v := range v.count {
		tiles[k].Properties.Set("count", v)
	}
	return &nurikabe{
		tiles: tiles,
		row:   v.rows,
		col:   v.cols,
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
