package main

import (
	"log"
	"os"
	"time"

	"github.com/ostlerc/nurikabe/grid"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

type window struct {
	g          *grid.Grid
	statusText qml.Object
}

func (w *window) TileChecked() {
	if w.g.CheckWin() {
		go func() {
			w.statusText.Set("text", "Winner!")
			time.Sleep(5 * time.Second)
			w.statusText.Set("text", "Nurikabe")
		}()
	}
}

func CreateMainWindow(engine *qml.Engine) {
	component, err := engine.LoadFile("qml/nurikabe.qml")
	if err != nil {
		log.Fatal(err)
	}

	comp := component.CreateWindow(nil)
	g := grid.New(validator.NewNurikabe(), comp.Root().ObjectByName("grid"))

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
	window := &window{g: g, statusText: comp.Root().ObjectByName("statusText")}
	context.SetVar("grid", g)
	context.SetVar("window", window)

	comp.Show()
	comp.Wait()
}
