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
	v          validator.GridValidator
	tiles      []qml.Object
}

func (w *window) TileChecked(i int) {
	w.g.Toggle(i)
	if w.v.CheckWin(w.g) {
		w.statusText.Set("text", "Winner!")
	} else {
		w.statusText.Set("text", "Nurikabe")
	}
	if *verbose {
		w.g.Print()
	}
}

func CreateMainWindow(engine *qml.Engine) {
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

	g := grid.New()

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
		v:          validator.NewNurikabe(),
	}
	context.SetVar("grid", g)
	context.SetVar("window", window)
	window.qmlgrid.Set("columns", g.Columns())

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
