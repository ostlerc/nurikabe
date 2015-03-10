package main

import (
	"log"
	"os"

	"github.com/ostlerc/nurikabe/grid"
	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

func MainWindow(engine *qml.Engine) {
	component, err := engine.LoadFile("qml/nurikabe.qml")
	if err != nil {
		log.Fatal(err)
	}

	win := component.CreateWindow(nil)
	g := grid.New(validator.NewNurikabe(),
		win.Root().ObjectByName("grid"),
		win.Root().ObjectByName("statusText"))

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
	context.SetVar("grid", g)

	win.Show()
	win.Wait()
}
