package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ostlerc/nurikabe/grid"
	"github.com/ostlerc/nurikabe/validator"
)

var (
	width  = flag.Int("width", 5, "grid width")
	height = flag.Int("height", 5, "grid height")
	min    = flag.Int("min", 3, "minimum gardens")
	growth = flag.Int("growth", 4, "garden growth. base + growth is max garden size")
	base   = flag.Int("base", 2, "Minimum garden size")

	verbose = flag.Bool("v", false, "Verbose")
)

func init() {
	flag.Parse()
}

func main() {
	validator.Verbose = *verbose
	g := grid.New(validator.NewNurikabe(), nil)

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		err := g.LoadGrid(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		g.BuildGrid(*height, *width)
	}

	g.Generate(*min, *growth, *base)
	j, err := g.Json()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
	if !g.CheckWin() {
		panic("Fail")
	}
	g.Print()
}
