package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/qml.v1"
)

var (
	grid Grid
	file = flag.String("file", "", "json map to load at startup")
)

func init() {
	flag.Parse()
}

func main() {
	if err := qml.Run(run); err != nil {
		log.Fatalf("error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	engine := qml.NewEngine()

	component, err := engine.LoadFile("qml/nurikabe.qml")
	if err != nil {
		return err
	}

	tileComponent, err := engine.LoadFile("qml/tile.qml")
	if err != nil {
		return err
	}

	context := engine.Context()
	context.SetVar("grid", &grid)

	win := component.CreateWindow(nil)

	grid.RowCount = 5
	grid.ColCount = 5
	grid.Grid = win.Root().ObjectByName("grid")
	grid.StatusText = win.Root().ObjectByName("statusText")
	grid.TileComp = &Tile{Object: tileComponent}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		grid.LoadGrid(os.Stdin)
	} else {
		grid.BuildGrid(3, 3)
	}

	win.Show()
	win.Wait()

	return nil
}
