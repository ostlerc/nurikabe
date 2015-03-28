package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

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
		w.status("Winner!")
	}
}

type GameMode int

const (
	DifficultySelect = iota
	LevelSelect
	Nurikabe
)

func (w *window) setGameMode(mode GameMode) {
	w.clearGrid()
	w.source("qml/game.qml")
	if mode == Nurikabe {
		w.status("Nurikabe - " + w.currentBoard)
		w.qRecordText().Set("text", w.records.String(w.currentBoard))
		w.qGameGrid().Set("spacing", 1)
		w.buildNurikabeGrid()

	} else {
		w.currentBoard = ""
		w.status("Nurikabe")
		w.buildLevelSelect()
		w.qGameGrid().Set("spacing", 15)
	}

	w.qMenuBtn().Set("visible", mode == Nurikabe)
	w.qStepsText().Set("visible", mode == Nurikabe)
	w.qTimeText().Set("visible", mode == Nurikabe)
	w.qRecordText().Set("visible", mode == Nurikabe)
}

func (w *window) MainMenuPressed() {
	w.setGameMode(LevelSelect)
}

func (w *window) Level(s string) {
	w.currentBoard = s
	w.loadLevel("levels/" + s + ".json")
}

func (w *window) loadLevel(file string) {
	if file != "" {
		r, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to load json file "+file)
			w.setGameMode(LevelSelect)
			return
		} else {
			w.g, err = grid.FromJson(r)
			if err != nil {
				panic(err)
			}
			w.setGameMode(Nurikabe)
		}
	}
}

func (w *window) clearGrid() {
	for i, _ := range w.btns {
		w.btns[i].Set("visible", false)
		w.btns[i].Destroy()
	}
	w.btns = nil
}

func (w *window) buildNurikabeGrid() {
	l := w.g.Rows() * w.g.Columns()
	w.qGameGrid().Set("columns", w.g.Columns())

	w.clearGrid()
	w.btns = make([]qml.Object, l, l)
	dimension := w.g.Columns()
	if rows := w.g.Rows(); rows > dimension {
		dimension = rows
	}
	windowDim := w.comp.Root().Int("width") - 50
	if height := w.comp.Root().Int("height"); height < windowDim {
		windowDim = height
	}
	dimension = windowDim / dimension
	for i := 0; i < l; i++ {
		w.btns[i] = w.tileComponent.Create(nil)
		w.btns[i].Set("parent", w.qGameGrid())
		w.btns[i].Set("index", i)
		w.btns[i].Set("count", w.g.Count(i))
		w.btns[i].Set("width", dimension)
		w.btns[i].Set("height", dimension)
	}
}

func (w *window) buildLevelSelect() {
	w.qGameGrid().Set("columns", 5)

	files, err := ioutil.ReadDir("levels/")
	if err != nil {
		panic(err)
	}

	l := len(files)
	w.btns = make([]qml.Object, l, l)
	names := make([]string, l, l)
	for i, f := range files {
		names[i] = f.Name()[:len(f.Name())-5]
	}

	sort.Strings(names)

	for i, name := range names {
		_, ok := w.records.Stats[name]
		w.btns[i] = w.btnComponent.Create(nil)
		w.btns[i].Set("parent", w.qGameGrid())
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
	window.setGameMode(LevelSelect)

	window.comp.Show()
	window.comp.Wait()
	window.records.Save(StatsFile)
	return nil
}
