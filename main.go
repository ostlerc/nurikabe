package main

import (
	"flag"
	"log"
	"os"

	"github.com/ostlerc/nurikabe/validator"

	"gopkg.in/qml.v1"
)

var (
	verbose = flag.Bool("v", false, "Show debug output")
)

func init() {
	flag.Parse()
}

func main() {
	validator.Verbose = *verbose
	if err := qml.Run(run); err != nil {
		log.Fatalf("error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	engine := qml.NewEngine()
	return RunNurikabe(engine)
}
