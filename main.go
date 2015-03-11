package main

import (
	"log"
	"os"

	"github.com/ostlerc/nurikabe/tile"

	"gopkg.in/qml.v1"
)

func main() {
	if err := qml.Run(run); err != nil {
		log.Fatalf("error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	engine := qml.NewEngine()

	tile.Setup(engine, "qml/tile.qml")
	MainWindow(engine)

	return nil
}
