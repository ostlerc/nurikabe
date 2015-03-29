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
	statsFile = ".stats.json"
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
	currentMode       gameMode
}

type gameMode int

const (
	mainMenu = iota
	difficultySelect
	levelSelect
	nurikabe
)

const (
	MenuStart = "Start"
	MenuStats = "Stats"
	MenuRules = "Rules"
	MenuExit  = "Exit"
)

var MenuItems = []string{MenuStart, MenuStats, MenuRules, MenuExit}

func NewMainWindow(engine *qml.Engine) (*window, error) {
	windowComponent, err := engine.LoadFile("qml/window.qml")
	if err != nil {
		return nil, err
	}

	window := &window{
		v:            validator.NewNurikabe(),
		winComponent: windowComponent.CreateWindow(nil),
	}

	window.tileComponent, err = engine.LoadFile("qml/tile.qml")
	if err != nil {
		return nil, err
	}

	window.btnComponent, err = engine.LoadFile("qml/button.qml")
	if err != nil {
		return nil, err
	}

	return window, nil
}

func (w *window) setGameMode(mode gameMode) {
	w.currentMode = mode
	w.clearGrid()
	w.setSource("qml/game.qml") //reload screen

	switch mode {
	case mainMenu:
		w.buildMainMenu()
		break

	case difficultySelect:
		w.buildDifficultySelect()
		break

	case levelSelect:
		w.buildLevelSelect()
		break

	case nurikabe:
		w.buildNurikabeGrid()
		break
	}

	w.qMenuBtn().Set("visible", mode != mainMenu)
	w.qStepsText().Set("visible", mode == nurikabe)
	w.qTimeText().Set("visible", mode == nurikabe)
	w.qRecordText().Set("visible", mode == nurikabe)
}

func (w *window) MainMenuClicked() {
	switch w.currentMode {
	case difficultySelect:
		w.setGameMode(mainMenu)
		break
	case levelSelect:
		w.setGameMode(difficultySelect)
		break
	case nurikabe:
		w.setGameMode(levelSelect)
		break
	}
}

func (w *window) OnBtnClicked(data string) {
	switch w.currentMode {
	case mainMenu:
		switch data {
		case MenuStart:
			w.setGameMode(difficultySelect)
			break
		case MenuStats:
			fmt.Println("Not implemented yet")
			break
		case MenuRules:
			fmt.Println("Not implemented yet")
			break
		case MenuExit:
			w.records.Save(statsFile)
			os.Exit(0)
			break
		}
		break
	case difficultySelect:
		w.currentDifficulty = data
		w.setGameMode(levelSelect)
	case levelSelect:
		w.currentBoard = data
		w.loadLevel("levels/" + w.currentDifficulty + "/" + data)
	case nurikabe: //This is handled by TileChecked
		panic("Err")
	}
}

func (w *window) TileChecked(i int) {
	w.qStepsText().Set("moves", w.qStepsText().Int("moves")+1)
	w.g.Toggle(i)
	if w.v.CheckWin(w.g) {
		w.records.Log(w.currentBoard, w.qStepsText().Int("moves"), w.qTimeText().Int("seconds"))
		w.records.Save(statsFile)
		w.setStatus("Winner!")
	}
}

func (w *window) loadLevel(file string) {
	if file != "" {
		r, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to load json file "+file)
			w.setGameMode(levelSelect)
			return
		} else {
			w.g, err = grid.FromJson(r)
			if err != nil {
				panic(err)
			}
			w.setGameMode(nurikabe)
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
	w.setStatus("Nurikabe - " + w.currentDifficulty[2:] + " " + w.currentBoard)
	w.qRecordText().Set("text", w.records.String(w.currentBoard))
	w.qGameGrid().Set("spacing", 1)
	w.qMenuBtn().Set("text", "Back")

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

func (w *window) buildMainMenu() {
	w.setStatus("Nurikabe - Main Menu")
	w.qGameGrid().Set("spacing", 15)
	w.qGameGrid().Set("columns", 1)

	w.objs = make([]qml.Object, len(MenuItems), len(MenuItems))
	for i, name := range MenuItems {
		w.objs[i] = w.btnComponent.Create(nil)
		w.objs[i].Set("parent", w.qGameGrid())
		w.objs[i].Set("text", name)
		w.objs[i].Set("data", name)
		w.objs[i].Set("alignCenter", true)
		w.objs[i].Set("width", w.winComponent.Root().Int("width")-150)
	}
}

func (w *window) buildDifficultySelect() {
	w.setStatus("Nurikabe - Select Difficulty")
	w.qGameGrid().Set("spacing", 15)
	w.qGameGrid().Set("columns", 1)
	w.qMenuBtn().Set("text", "Menu")

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
	w.currentBoard = ""
	w.setStatus("Nurikabe - " + w.currentDifficulty[2:])
	w.qGameGrid().Set("spacing", 15)
	w.qMenuBtn().Set("text", "Back")
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

func (w *window) loadStats() {
	var err error
	w.records, err = stats.Load(statsFile)
	if err != nil {
		fmt.Println("Error loading stats", err)
		w.records = stats.New()
	}
}

func (w *window) setStatus(s string) {
	w.qStatus().Set("text", s)
}

func (w *window) setSource(page string) {
	w.qLoader().Set("source", page)
}

func (w *window) obj(name string) qml.Object {
	return w.winComponent.Root().ObjectByName(name)
}

func (w *window) qStatus() qml.Object {
	return w.obj("statusText")
}

func (w *window) qGameGrid() qml.Object {
	return w.obj("grid")
}

func (w *window) qLoader() qml.Object {
	return w.obj("pageLoader")
}

func (w *window) qMenuBtn() qml.Object {
	return w.obj("menuBtn")
}

func (w *window) qStepsText() qml.Object {
	return w.obj("movesText")
}

func (w *window) qTimeText() qml.Object {
	return w.obj("timerText")
}

func (w *window) qRecordText() qml.Object {
	return w.obj("recordText")
}

func RunNurikabe(engine *qml.Engine) error {
	context := engine.Context()

	window, err := NewMainWindow(engine)
	if err != nil {
		return err
	}

	context.SetVar("window", window)
	window.loadStats()
	window.setGameMode(mainMenu)

	window.winComponent.Show()
	window.winComponent.Wait()
	window.records.Save(statsFile)
	return nil
}
