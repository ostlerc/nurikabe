package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ostlerc/nurikabe/grid"
	"github.com/ostlerc/nurikabe/stats"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

const (
	statsFile = ".stats.json"
	levelDir  = "levels/"
)

type window struct {
	g       *grid.Grid
	v       validator.GridValidator
	objs    []qml.Object
	records *stats.Records

	tileComponent qml.Object
	btnComponent  qml.Object
	txtComponent  qml.Object
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
	rulesPage
	statsPage
	nurikabePage
)

const (
	MenuPlay  = "Play"
	MenuStats = "Records"
	MenuRules = "Rules"
	MenuExit  = "Exit"
)

const rulesText = `Each puzzle consists of a grid containing clues in various places.` +
	` The object is to create islands by partitioning between clues with walls so:` +
	` Each island contains exactly one clue.` +
	` The number of squares in each island equals the value of the clue.` +
	` All islands are isolated from each other horizontally and vertically.` +
	` There are no wall areas of 2x2 or larger.` +
	` When completed, all walls form a continuous path.`

var MenuItems = []string{MenuPlay, MenuStats, MenuRules, MenuExit}

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

	window.txtComponent, err = engine.LoadFile("qml/text.qml")
	if err != nil {
		return nil, err
	}

	return window, nil
}

func (w *window) setGameMode(mode gameMode) {
	w.currentMode = mode
	w.clearGrid()
	w.setSource("qml/game.qml") //reload screen

	w.qToolBtn().Set("visible", mode != mainMenu)
	w.qStepsText().Set("visible", mode == nurikabePage)
	w.qTimeText().Set("visible", mode == nurikabePage)
	w.qRecordText().Set("visible", mode == nurikabePage)

	switch mode {
	case mainMenu:
		w.buildMainMenu()
	case difficultySelect:
		w.buildDifficultySelect()
	case levelSelect:
		w.buildLevelSelect()
	case nurikabePage:
		w.buildNurikabeGrid()
	case rulesPage:
		w.buildRules()
	case statsPage:
		w.buildStats()
	}
}

func (w *window) ToolButtonClicked() {
	switch w.currentMode {
	case rulesPage:
		w.setGameMode(mainMenu)
	case statsPage:
		w.setGameMode(mainMenu)
	case difficultySelect:
		w.setGameMode(mainMenu)
	case levelSelect:
		w.setGameMode(difficultySelect)
	case nurikabePage:
		w.setGameMode(levelSelect)
	}
}

func (w *window) OnBtnClicked(data string) {
	switch w.currentMode {
	case mainMenu:
		switch data {
		case MenuPlay:
			w.setGameMode(difficultySelect)
		case MenuStats:
			w.setGameMode(statsPage)
		case MenuRules:
			w.setGameMode(rulesPage)
		case MenuExit:
			w.records.Save(statsFile)
			os.Exit(0)
		}
	case difficultySelect:
		w.currentDifficulty = data
		w.setGameMode(levelSelect)
	case levelSelect:
		w.currentBoard = data
		w.loadLevel(levelDir + w.currentDifficulty + "/" + data)
	case nurikabePage: //This is handled by TileChecked
		panic("Err")
	case rulesPage:
		panic("Err")
	}
}

func (w *window) TileChecked(i int) {
	w.qStepsText().Set("moves", w.qStepsText().Int("moves")+1)
	w.g.Toggle(i)
	if w.v.CheckWin(w.g) {
		w.records.Log(w.currentDifficulty, levelInt(w.currentBoard), w.qStepsText().Int("moves"), w.qTimeText().Int("seconds"))
		w.records.Save(statsFile)
		w.setStatus("Winner!")
	}
}

func levelInt(file string) int {
	ret, err := strconv.Atoi(levelStr(file))
	if err != nil {
		panic(err)
	}
	return ret
}

func levelStr(file string) string {
	return file[:len(file)-5]
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
			w.setGameMode(nurikabePage)
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
	w.setStatus("Nurikabe - " + w.currentDifficulty[2:] + " " + levelStr(w.currentBoard))
	w.qRecordText().Set("text", w.records.String(w.currentDifficulty, levelInt(w.currentBoard)))
	w.qGameGrid().Set("spacing", 1)
	w.qToolBtn().Set("text", "Back")

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

func (w *window) buildLevelSelect() {
	w.currentBoard = ""
	w.setStatus("Nurikabe - " + w.currentDifficulty[2:])
	w.qGameGrid().Set("spacing", 15)
	w.qToolBtn().Set("text", "Back")
	w.qGameGrid().Set("columns", 4)

	names := files(levelDir + w.currentDifficulty)
	w.objs = make([]qml.Object, len(names), len(names))
	for i, name := range names {
		_, ok := w.records.Level(w.currentDifficulty, levelInt(name))
		w.objs[i] = w.btnComponent.Create(nil)
		w.objs[i].Set("parent", w.qGameGrid())
		w.objs[i].Set("text", levelStr(name)) //remove '.json' from name
		w.objs[i].Set("data", name)
		w.objs[i].Set("showstar", true)
		w.objs[i].Set("completed", ok)
		w.objs[i].Set("width", 50)
	}
}

func (w *window) buildDifficultySelect() {
	w.setStatus("Nurikabe - Select Difficulty")
	w.qGameGrid().Set("spacing", 15)
	w.qGameGrid().Set("columns", 1)
	w.qToolBtn().Set("text", "Menu")

	names := dirs(levelDir)
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

func (w *window) buildRules() {
	w.currentBoard = ""
	w.setStatus("Nurikabe - Rules")
	w.qGameGrid().Set("spacing", 15)
	w.qToolBtn().Set("text", "Back")
	w.qGameGrid().Set("columns", 4)

	w.objs = make([]qml.Object, 1, 1)
	w.objs[0] = w.txtComponent.Create(nil)
	w.objs[0].Set("parent", w.qGameGrid())
	w.objs[0].Set("text", rulesText)
	w.objs[0].Set("width", w.winComponent.Root().Int("width")-50)
}

func (w *window) buildStats() {
	w.currentBoard = ""
	w.setStatus("Nurikabe - Records")
	w.qGameGrid().Set("spacing", 15)
	w.qToolBtn().Set("text", "Back")
	w.qGameGrid().Set("columns", 4)

	l := w.records.Length()
	w.objs = make([]qml.Object, 0, l+1)

	buildTxtBox := func(s string) {
		obj := w.txtComponent.Create(nil)
		obj.Set("parent", w.qGameGrid())
		obj.Set("text", s)
		w.objs = append(w.objs, obj)
	}

	headers := []string{"Difficulty", "Level", "Steps", "Seconds"}
	for _, txt := range headers {
		buildTxtBox(txt)
	}

	for _, rec := range w.records.All() {
		buildTxtBox(rec.Difficulty[2:])
		buildTxtBox(strconv.Itoa(rec.Lvl))
		buildTxtBox(strconv.Itoa(rec.Steps))
		buildTxtBox(strconv.Itoa(rec.Seconds))
	}
}

func (w *window) loadStats() {
	var err error
	d := dirs(levelDir)
	sorter := make(map[string]int, len(d))
	for _, f := range dirs(levelDir) {
		sorter[f[2:]] = int(f[0] - '0')

	}
	w.records, err = stats.Load(statsFile, sorter)
	if err != nil {
		fmt.Println("Error loading stats", err)
		w.records = stats.New(sorter)
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

func (w *window) qToolBtn() qml.Object {
	return w.obj("toolBtn")
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
