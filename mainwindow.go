package main

import (
	"fmt"
	"io/ioutil"
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
	btns    []qml.Object
	records *stats.Records

	tileComponent qml.Object
	btnComponent  qml.Object
	comp          *qml.Window
	currentBoard  string
}

func (w *window) TileChecked(i int) {
	w.qStepsText().Set("moves", w.qStepsText().Int("moves")+1)
	w.g.Toggle(i)
	if w.v.CheckWin(w.g) {
		w.records.Log(w.currentBoard, w.qStepsText().Int("moves"), w.qTimeText().Int("seconds"))
		w.records.Save(StatsFile)
		w.setGameMode(true)
		w.status("Winner!")
	} else {
		w.setGameMode(true)
	}
}

func (w *window) setGameMode(show bool) {
	if show {
		w.status("Nurikabe - " + w.currentBoard)
		w.qRecordText().Set("text", w.records.String(w.currentBoard))

	} else {
		w.currentBoard = ""
		w.status("Nurikabe")
	}
	w.qMenuBtn().Set("visible", show)
	w.qStepsText().Set("visible", show)
	w.qTimeText().Set("visible", show)
	w.qRecordText().Set("visible", show)
}

func (w *window) MainMenuPressed() {
	w.source("qml/main.qml")
	w.setGameMode(false)
	w.buildBtns()
}

func (w *window) Level(s string) {
	w.source("qml/game.qml")
	w.currentBoard = s
	w.loadLevel("levels/" + s + ".json")
}

func (w *window) loadLevel(file string) {
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
			w.setGameMode(true)
		}
	}
	w.buildGrid()
}

func (w *window) buildGrid() {
	l := w.g.Rows() * w.g.Columns()
	w.qGameGrid().Set("columns", w.g.Columns())

	w.tiles = make([]qml.Object, l, l)
	for i := 0; i < l; i++ {
		w.tiles[i] = w.tileComponent.Create(nil)
		w.tiles[i].Set("parent", w.qGameGrid())
		w.tiles[i].Set("index", i)
		w.tiles[i].Set("count", w.g.Count(i))
	}
}

func (w *window) buildBtns() {
	w.qBtnGrid().Set("columns", 5)

	files, err := ioutil.ReadDir("levels/")
	if err != nil {
		panic(err)
	}

	l := len(files)
	w.btns = make([]qml.Object, l, l)

	for i, f := range files {
		name := f.Name()
		name = name[:len(name)-5]
		_, ok := w.records.Stats[name]
		w.btns[i] = w.btnComponent.Create(nil)
		w.btns[i].Set("parent", w.qBtnGrid())
		w.btns[i].Set("text", name)
		w.btns[i].Set("completed", ok)
		w.btns[i].Set("color", "silver")
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

func (w *window) qGameGrid() qml.Object {
	return w.comp.Root().ObjectByName("grid")
}

func (w *window) qBtnGrid() qml.Object {
	return w.comp.Root().ObjectByName("btnGrid")
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

func (w *window) qRecordText() qml.Object {
	return w.comp.Root().ObjectByName("recordText")
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

	window.btnComponent, err = engine.LoadFile("qml/btn.qml")
	if err != nil {
		return err
	}

	context.SetVar("window", window)
	window.records, err = stats.Load(StatsFile)
	if err != nil {
		fmt.Println("Error loading stats", err)
		window.records = stats.New()
	}
	window.MainMenuPressed()

	window.comp.Show()
	window.comp.Wait()
	window.records.Save(StatsFile)
	return nil
}
