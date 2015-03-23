package main

import (
	"log"
	"os"

	"github.com/ostlerc/nurikabe/grid"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

type window struct {
	g     *grid.Grid
	v     validator.GridValidator
	tiles []qml.Object

	status qml.Object
}

func (w *window) TileChecked(i int) {
	w.g.Toggle(i)
	if w.v.CheckWin(w.g) {
		w.status.Set("text", "Winner!")
	} else {
		w.status.Set("text", "Nurikabe")
	}
	if *verbose {
		w.g.Print()
	}
}

func ShowMainWindow(engine *qml.Engine) {
	var g *grid.Grid
	var err error

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		g, err = grid.FromJson(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		g = grid.New(3, 3)
	}

	context := engine.Context()

	tileComponent, err := engine.LoadFile("qml/tile.qml")
	if err != nil {
		panic(err)
	}
	nurikabeComponent, err := engine.LoadFile("qml/nurikabe.qml")
	if err != nil {
		log.Fatal(err)
	}

	comp := nurikabeComponent.CreateWindow(nil)
	qmlgrid := comp.Root().ObjectByName("grid")
	window := &window{
		g:      g,
		status: comp.Root().ObjectByName("statusText"),
		v:      validator.NewNurikabe(),
	}
	context.SetVar("window", window)
	qmlgrid.Set("columns", g.Columns())

	l := g.Rows() * g.Columns()
	tiles := make([]qml.Object, l, l)
	for i := 0; i < l; i++ {
		tiles[i] = tileComponent.Create(nil)
		tiles[i].Set("parent", qmlgrid)
		tiles[i].Set("index", i)
		tiles[i].Set("count", g.Count(i))
	}

	window.tiles = tiles

	comp.Show()
	comp.Wait()
}
