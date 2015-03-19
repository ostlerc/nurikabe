package grid

import (
	"io"
	"strings"
	"testing"

	"github.com/ostlerc/nurikabe/tile"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

type gridTest struct {
	testNum int
	g       *Grid
	win     bool
}

var engine *qml.Engine

func init() {
	qml.SetupTesting()
	tile.SetupTesting()
}

func loadGrid(input io.Reader, closed []int) *Grid {
	g := New(validator.NewNurikabe(), nil)

	g.grid = tile.Fake()
	g.LoadGrid(input)
	setClosed(closed, g)
	return g
}

func setClosed(idx []int, g *Grid) {
	for _, i := range idx {
		g.tiles[i].SetType(1)
	}
}

func TestWinner(t *testing.T) {
	engine = qml.NewEngine()

	gridTests := []gridTest{
		{1, loadGrid(strings.NewReader(`{"rows":2,"cols":2}`), []int{0, 1, 2, 3}), false},
		{2, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{1, 2, 3, 4}), true},
		{3, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{2, 3, 5, 6}), true},
		{4, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{1, 2, 4, 5}), false},
		{5, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{3, 4, 6, 7}), false},
		{6, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{4, 5, 7, 8}), false},
		{7, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{4, 5, 7}), true},
		{8, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{0, 2, 5}), false},
		{9, loadGrid(strings.NewReader(`{"rows":3,"cols":3,"tiles":[{"count":2,"index":0},{"count":3,"index":5}]}`), []int{}), false},
		{10, loadGrid(strings.NewReader(`{"rows":3,"cols":3,"tiles":[{"count":2,"index":0},{"count":3,"index":5}]}`), []int{1, 4, 6, 7}), true},
		{11, loadGrid(strings.NewReader(`{"rows":3,"cols":3,"tiles":[{"count":2,"index":0},{"count":3,"index":5}]}`), []int{1, 4, 6}), false},
	}

	for i, gt := range gridTests {
		if gt.win != gt.g.CheckWin() {
			t.Fatal("win invalid for test", i, "(", gt.testNum, ")")
		}
	}
}

func TestBuildGrid(t *testing.T) {
	g := &Grid{grid: tile.Fake()}
	g.BuildGrid(4, 6)
	if g.cols != 6 {
		t.Fatal("Invalid columns ", g.cols)
	}
	if g.rows != 4 {
		t.Fatal("Invalid rows ", g.rows)
	}
}
