package grid

import (
	"io"
	"strings"
	"testing"

	"github.com/ostlerc/nurikabe/tile"
	"github.com/ostlerc/nurikabe/validator"
)

type gridTest struct {
	testNum int
	json    string
	closed  []int
	win     bool
}

func init() {
	tile.SetupTesting()
}

func loadGrid(input io.Reader, closed []int) *Grid {
	g := New(validator.NewNurikabe(), nil)

	g.LoadGrid(input)
	setClosed(closed, g)
	return g
}

func setClosed(idx []int, g *Grid) {
	for _, i := range idx {
		g.tiles[i].SetType(1)
	}
}

var gridTests = []gridTest{
	{1, `{"rows":2,"cols":2}`, []int{0, 1, 2, 3}, false},
	{2, `{"rows":3,"cols":3}`, []int{1, 2, 3, 4}, true},
	{3, `{"rows":3,"cols":3}`, []int{2, 3, 5, 6}, true},
	{4, `{"rows":3,"cols":3}`, []int{1, 2, 4, 5}, false},
	{5, `{"rows":3,"cols":3}`, []int{3, 4, 6, 7}, false},
	{6, `{"rows":3,"cols":3}`, []int{4, 5, 7, 8}, false},
	{7, `{"rows":3,"cols":3}`, []int{4, 5, 7}, true},
	{8, `{"rows":3,"cols":3}`, []int{0, 2, 5}, false},
	{9, `{"rows":3,"cols":3,"tiles":[{"count":2},{"count":3,"index":5}]}`, []int{}, false},
	{10, `{"rows":3,"cols":3,"tiles":[{"count":2},{"count":3,"index":5}]}`, []int{1, 4, 6, 7}, true},
	{11, `{"rows":3,"cols":3,"tiles":[{"count":2},{"count":3,"index":5}]}`, []int{1, 4, 6}, false},
}

func TestWinner(t *testing.T) {
	for i, gt := range gridTests {
		grid := loadGrid(strings.NewReader(gt.json), gt.closed)
		if gt.win != grid.CheckWin() {
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
	if g.Cols != 6 {
		t.Fatal("Invalid columns ", g.Cols)
	}
	if g.Rows != 4 {
		t.Fatal("Invalid rows ", g.Rows)
	}
}
