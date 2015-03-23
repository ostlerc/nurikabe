package grid

import (
	"io"
	"strings"
	"testing"

	"github.com/ostlerc/nurikabe/validator"
)

type gridTest struct {
	testNum int
	json    string
	closed  []int
	win     bool
}

var v = validator.NewNurikabe()

func loadGrid(input io.Reader, closed []int) *Grid {
	g := New()

	g.LoadGrid(input)
	setClosed(closed, g)
	return g
}

func setClosed(idx []int, g *Grid) {
	for _, i := range idx {
		g.tiles[i].open = false
	}
}

var gridTests = []gridTest{
	{1, `{"rows":2,"cols":2}`, []int{0, 1, 2, 3}, false},
	{2, `{"rows":3,"cols":3}`, []int{1, 2, 3, 4}, false},
	{3, `{"rows":3,"cols":3}`, []int{2, 3, 5, 6}, false},
	{4, `{"rows":3,"cols":3}`, []int{1, 2, 4, 5}, false},
	{5, `{"rows":3,"cols":3}`, []int{3, 4, 6, 7}, false},
	{6, `{"rows":3,"cols":3}`, []int{4, 5, 7, 8}, false},
	{7, `{"rows":3,"cols":3}`, []int{4, 5, 7}, false},
	{8, `{"rows":3,"cols":3}`, []int{0, 2, 5}, false},
	{9, `{"rows":3,"cols":3,"tiles":[{"count":2},{"count":3,"index":5}]}`, []int{}, false},
	{10, `{"rows":3,"cols":3,"tiles":[{"count":2},{"count":3,"index":5}]}`, []int{1, 4, 6, 7}, true},
	{11, `{"rows":3,"cols":3,"tiles":[{"count":2},{"count":3,"index":5}]}`, []int{1, 4, 6}, false},
}

func TestWinner(t *testing.T) {
	for i, gt := range gridTests {
		grid := loadGrid(strings.NewReader(gt.json), gt.closed)
		if gt.win != v.CheckWin(grid) {
			t.Fatal("win invalid for test", i, "(", gt.testNum, ")")
		}
	}
}

func TestJson(t *testing.T) {
	for i, gt := range gridTests {
		grid := loadGrid(strings.NewReader(gt.json), gt.closed)
		if json, err := grid.Json(); err != nil || string(json) != gt.json {
			t.Fatal("Invalid json", i, string(json), "(", gt.testNum, ")")
		}
	}
}

func TestBuildGrid(t *testing.T) {
	g := &Grid{}
	g.BuildGrid(4, 6)
	if g.cols != 6 {
		t.Fatal("Invalid columns ", g.cols)
	}
	if g.rows != 4 {
		t.Fatal("Invalid rows ", g.rows)
	}
}
