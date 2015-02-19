package main

import (
	"io"
	"log"
	"strings"
	"testing"

	"gopkg.in/qml.v1"
)

type gridTest struct {
	testNum int
	g       *Grid
	block   bool
	wall    bool
	gardens bool
	win     bool
}

var gridTests []gridTest
var tileComponent qml.Object
var gridComponent qml.Object
var engine *qml.Engine

func init() {
	qml.SetupTesting()
}

func loadGrid(input io.Reader, closed []int) *Grid {
	g := &Grid{
		TileComp: &Tile{Object: tileComponent},
	}

	g.Grid = gridComponent.Create(nil)
	g.LoadGrid(input)
	setClosed(closed, g)
	return g
}

func setClosed(idx []int, g *Grid) {
	for _, i := range idx {
		g.Tiles[i].SetType(1)
	}
}

func TestSetup(t *testing.T) {
	engine = qml.NewEngine()

	var err error
	tileComponent, err = engine.LoadFile("qml/tile.qml")
	if err != nil {
		log.Fatal("Couldn't load tile.qml")
	}
	gridComponent, err = engine.LoadString("test.qml", "import QtQuick 2.0\nGrid {}")
	if err != nil {
		log.Fatal("could not load from string")
	}

	gridTests = []gridTest{
		{1, loadGrid(strings.NewReader(`{"rows":2,"cols":2}`), []int{0, 1, 2, 3}), true, true, true, false},
		{2, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{1, 2, 3, 4}), false, true, true, true},
		{3, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{2, 3, 5, 6}), false, true, true, true},
		{4, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{1, 2, 4, 5}), true, true, true, false},
		{5, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{3, 4, 6, 7}), true, true, true, false},
		{6, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{4, 5, 7, 8}), true, true, true, false},
		{7, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{4, 5, 7}), false, true, true, true},
		{8, loadGrid(strings.NewReader(`{"rows":3,"cols":3}`), []int{0, 2, 5}), false, false, true, false},
		{9, loadGrid(strings.NewReader(`{"rows":3,"cols":3,"tiles":[{"count":2,"index":0},{"count":3,"index":5}]}`), []int{}), false, false, false, false},
		{10, loadGrid(strings.NewReader(`{"rows":3,"cols":3,"tiles":[{"count":2,"index":0},{"count":3,"index":5}]}`), []int{1, 4, 6, 7}), false, true, true, true},
		{11, loadGrid(strings.NewReader(`{"rows":3,"cols":3,"tiles":[{"count":2,"index":0},{"count":3,"index":5}]}`), []int{1, 4, 6}), false, false, false, false},
	}
}

func TestHasBlock(t *testing.T) {
	for i, gt := range gridTests {
		if gt.block != gt.g.hasBlock() {
			t.Fatal("hasBlock invalid for test", i, "(", gt.testNum, ")")
		}
	}
}

func TestWall(t *testing.T) {
	for i, gt := range gridTests {
		if gt.wall != gt.g.singleWall() {
			t.Fatal("wall invalid for test", i, "(", gt.testNum, ")")
		}
	}
}

func TestGarden(t *testing.T) {
	for i, gt := range gridTests {
		if gt.gardens != gt.g.gardensAreCorrect() {
			t.Fatal("gardens invalid for test", i, "(", gt.testNum, ")")
		}
	}
}

func TestWinner(t *testing.T) {
	for i, gt := range gridTests {
		if gt.win != gt.g.CheckWin() {
			t.Fatal("win invalid for test", i, "(", gt.testNum, ")")
		}
	}
}
