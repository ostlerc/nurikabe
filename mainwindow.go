package main

import (
	"fmt"
	"os"

	"github.com/ostlerc/nurikabe/grid"
	"github.com/ostlerc/nurikabe/stats"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

const (
	StatsFile = ".stats.json"
)

type window struct {
	g       *grid.Grid
	v       validator.GridValidator
	tiles   []qml.Object
	records *stats.Records

	tileComponent qml.Object
	comp          *qml.Window
	currentBoard  string
}

func (w *window) TileChecked(i int) {
	w.qStepsText().Set("moves", w.qStepsText().Int("moves")+1)
	w.g.Toggle(i)
	if w.v.CheckWin(w.g) {
		w.status("Winner!")
		w.records.Log(w.currentBoard, w.qStepsText().Int("moves"), w.qTimeText().Int("seconds"))
		w.records.Save(StatsFile)
	} else {
		w.status("Nurikabe")
	}
	if *verbose {
		w.g.Print()
	}
}

func (w *window) setGameVisibility(show bool) {
	w.qMenuBtn().Set("visible", show)
	w.qStepsText().Set("visible", show)
	w.qTimeText().Set("visible", show)
}

func (w *window) MainMenuPressed() {
	w.status("Nurikabe")
	w.setGameVisibility(false)
	w.currentBoard = ""
	if w.qLoader().String("source") != "qml/main.qml" {
		w.source("qml/main.qml")
	} else {
		w.source("qml/game.qml")
		w.setup("")
	}
}

func (w *window) Level(s string) {
	w.source("qml/game.qml")
	w.currentBoard = s
	w.setup("json/" + s + ".json")
}

func (w *window) setup(file string) {
	if file != "" {
		r, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to load json file "+file)
			w.MainMenuPressed()
			return
		} else {
			w.g, err = grid.FromJson(r)
			if err != nil {
				panic(err)
			}
			end := 7
			if file[8] == '0' {
				end = 8
			}
			w.status("Nurikabe - " + file[5:end])
			w.setGameVisibility(true)
			w.qMenuBtn().Set("visible", true)
			w.qStepsText().Set("visible", true)
		}
	}
	l := w.g.Rows() * w.g.Columns()
	w.qGrid().Set("columns", w.g.Columns())

	w.tiles = make([]qml.Object, l, l)
	for i := 0; i < l; i++ {
		w.tiles[i] = w.tileComponent.Create(nil)
		w.tiles[i].Set("parent", w.qGrid())
		w.tiles[i].Set("index", i)
		w.tiles[i].Set("count", w.g.Count(i))
	}
}

func (w *window) status(s string) {
	w.qStatus().Set("text", s)
}

func (w *window) source(page string) {
	w.qLoader().Set("source", page)
}

func (w *window) qStatus() qml.Object {
	return w.comp.Root().ObjectByName("statusText")
}
func (w *window) qGrid() qml.Object {
	return w.comp.Root().ObjectByName("grid")
}

func (w *window) qLoader() qml.Object {
	return w.comp.Root().ObjectByName("pageLoader")
}

func (w *window) qMenuBtn() qml.Object {
	return w.comp.Root().ObjectByName("menuBtn")
}

func (w *window) qStepsText() qml.Object {
	return w.comp.Root().ObjectByName("movesText")
}

func (w *window) qTimeText() qml.Object {
	return w.comp.Root().ObjectByName("timerText")
}

func RunNurikabe(engine *qml.Engine) error {
	context := engine.Context()

	nurikabeComponent, err := engine.LoadFile("qml/nurikabe.qml")
	if err != nil {
		return err
	}

	window := &window{
		v:    validator.NewNurikabe(),
		comp: nurikabeComponent.CreateWindow(nil),
	}

	window.tileComponent, err = engine.LoadFile("qml/tile.qml")
	if err != nil {
		return err
	}
	context.SetVar("window", window)
	window.records, err = stats.Load(StatsFile)
	if err != nil {
		window.records = stats.New()
	}
	window.MainMenuPressed()

	window.comp.Show()
	window.comp.Wait()
	window.records.Save(StatsFile)
	return nil
}
