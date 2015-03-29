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
	objs    []qml.Object
	records *stats.Records

	tileComponent qml.Object
	btnComponent  qml.Object
	winComponent  *qml.Window

	currentDifficulty string
	currentBoard      string
	currentMode       GameMode
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
	w.currentMode = mode
	w.clearGrid()
	w.setSource("qml/game.qml")

	switch mode {
	case Nurikabe:
		w.status("Nurikabe - " + w.currentDifficulty[2:] + " " + w.currentBoard)
		w.qRecordText().Set("text", w.records.String(w.currentBoard))
		w.qGameGrid().Set("spacing", 1)
		w.buildNurikabeGrid()
		w.qMenuBtn().Set("text", "Back")
		break

	case LevelSelect:
		w.currentBoard = ""
		w.status("Nurikabe - " + w.currentDifficulty[2:])
		w.buildLevelSelect()
		w.qGameGrid().Set("spacing", 15)
		w.qMenuBtn().Set("text", "Menu")
		break

	case DifficultySelect:
		w.currentBoard = ""
		w.currentDifficulty = ""
		w.status("Nurikabe")
		w.buildDifficultySelect()
		w.qGameGrid().Set("spacing", 15)
		break
	}

	w.qMenuBtn().Set("visible", mode != DifficultySelect)
	w.qStepsText().Set("visible", mode == Nurikabe)
	w.qTimeText().Set("visible", mode == Nurikabe)
	w.qRecordText().Set("visible", mode == Nurikabe)
}

func (w *window) MainMenuClicked() {
	switch w.currentMode {
	case DifficultySelect:
	case LevelSelect:
		w.setGameMode(DifficultySelect)
	case Nurikabe:
		w.setGameMode(LevelSelect)
	}
}

func (w *window) OnBtnClicked(data string) {
	switch w.currentMode {
	case DifficultySelect:
		w.currentDifficulty = data
		w.setGameMode(LevelSelect)
	case LevelSelect:
		w.currentBoard = data
		w.loadLevel("levels/" + w.currentDifficulty + "/" + data)
	case Nurikabe: //This is handled by TileChecked
		panic("Err")
	}
}

func (w *window) OnDifficultyClicked(s string) {
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
	for i, _ := range w.objs {
		w.objs[i].Set("visible", false)
		w.objs[i].Destroy()
	}
	w.objs = nil
}

func (w *window) buildNurikabeGrid() {
	l := w.g.Rows() * w.g.Columns()
	w.qGameGrid().Set("columns", w.g.Columns())

	w.clearGrid()
	w.objs = make([]qml.Object, l, l)
	dimension := w.g.Columns()
	if rows := w.g.Rows(); rows > dimension {
		dimension = rows
	}
	windowDim := w.winComponent.Root().Int("width") - 50
	if height := w.winComponent.Root().Int("height"); height < windowDim {
		windowDim = height
	}
	dimension = windowDim / dimension
	for i := 0; i < l; i++ {
		w.objs[i] = w.tileComponent.Create(nil)
		w.objs[i].Set("parent", w.qGameGrid())
		w.objs[i].Set("index", i)
		w.objs[i].Set("count", w.g.Count(i))
		w.objs[i].Set("width", dimension)
		w.objs[i].Set("height", dimension)
	}
}

func (w *window) buildDifficultySelect() {
	w.qGameGrid().Set("columns", 1)

	names := dirs("levels/")
	w.objs = make([]qml.Object, len(names), len(names))
	for i, name := range names {
		w.objs[i] = w.btnComponent.Create(nil)
		w.objs[i].Set("parent", w.qGameGrid())
		w.objs[i].Set("text", name[2:])
		w.objs[i].Set("data", name)
		w.objs[i].Set("alignCenter", true)
		w.objs[i].Set("width", w.winComponent.Root().Int("width")-150)
	}
}

func (w *window) buildLevelSelect() {
	w.qGameGrid().Set("columns", 4)

	names := files("levels/" + w.currentDifficulty)
	w.objs = make([]qml.Object, len(names), len(names))
	for i, name := range names {
		_, ok := w.records.Stats[name]
		w.objs[i] = w.btnComponent.Create(nil)
		w.objs[i].Set("parent", w.qGameGrid())
		w.objs[i].Set("text", name[:len(name)-5]) //remove '.json' from name
		w.objs[i].Set("data", name)
		w.objs[i].Set("showstar", true)
		w.objs[i].Set("completed", ok)
		w.objs[i].Set("width", 50)
	}
}

func dirs(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	l := len(files)
	names := make([]string, 0, l)
	for _, f := range files {
		if f.IsDir() {
			names = append(names, f.Name())
		}
	}
	sort.Strings(names)
	return names
}

func files(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	l := len(files)
	names := make([]string, l, l)
	for i, f := range files {
		names[i] = f.Name()
	}

	sort.Strings(names)
	return names
}

func (w *window) status(s string) {
	w.qStatus().Set("text", s)
}

func (w *window) setSource(page string) {
	w.qLoader().Set("source", page)
}

func (w *window) qStatus() qml.Object {
	return w.winComponent.Root().ObjectByName("statusText")
}

func (w *window) qGameGrid() qml.Object {
	return w.winComponent.Root().ObjectByName("grid")
}

func (w *window) qLoader() qml.Object {
	return w.winComponent.Root().ObjectByName("pageLoader")
}

func (w *window) qMenuBtn() qml.Object {
	return w.winComponent.Root().ObjectByName("menuBtn")
}

func (w *window) qStepsText() qml.Object {
	return w.winComponent.Root().ObjectByName("movesText")
}

func (w *window) qTimeText() qml.Object {
	return w.winComponent.Root().ObjectByName("timerText")
}

func (w *window) qRecordText() qml.Object {
	return w.winComponent.Root().ObjectByName("recordText")
}

func RunNurikabe(engine *qml.Engine) error {
	context := engine.Context()

	nurikabeComponent, err := engine.LoadFile("qml/nurikabe.qml")
	if err != nil {
		return err
	}

	window := &window{
		v:            validator.NewNurikabe(),
		winComponent: nurikabeComponent.CreateWindow(nil),
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
	window.setGameMode(DifficultySelect)

	window.winComponent.Show()
	window.winComponent.Wait()
	window.records.Save(StatsFile)
	return nil
}
