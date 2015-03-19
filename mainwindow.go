package main

import (
	"log"
	"os"

	"github.com/ostlerc/nurikabe/grid"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

type window struct {
	g          *grid.Grid
	qmlgrid    qml.Object
	statusText qml.Object
}

func (w *window) TileChecked() {
	if w.g.CheckWin() {
		w.statusText.Set("text", "Winner!")
	} else {
		w.statusText.Set("text", "Nurikabe")
	}
}

func CreateMainWindow(engine *qml.Engine) {
	component, err := engine.LoadFile("qml/nurikabe.qml")
	if err != nil {
		log.Fatal(err)
	}

	comp := component.CreateWindow(nil)
	qmlgrid := comp.Root().ObjectByName("grid")
	g := grid.New(validator.NewNurikabe(), qmlgrid)

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		err := g.LoadGrid(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		g.BuildGrid(3, 3)
	}

	context := engine.Context()
	window := &window{
		g:          g,
		qmlgrid:    qmlgrid,
		statusText: comp.Root().ObjectByName("statusText"),
	}
	context.SetVar("grid", g)
	context.SetVar("window", window)
	window.qmlgrid.Set("columns", g.Cols)

	comp.Show()
	comp.Wait()
}
