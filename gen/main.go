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
	min    = flag.Int("min", 3, "minimum gardens count")
	growth = flag.Int("growth", 4, "garden growth. base + growth is max garden size")
	base   = flag.Int("base", 2, "minimum garden size")

	verbose = flag.Bool("v", false, "Verbose")
)

func init() {
	flag.Parse()
}

func main() {
	validator.Verbose = *verbose
	v := validator.NewNurikabe()
	var g *grid.Grid
	var err error

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		g, err = grid.FromJson(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		g = grid.New(*height, *width)
	}

	g.Generate(v, *min, *growth, *base)
	j, err := g.Json()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
	if !v.CheckWin(g) {
		panic("Fail")
	}
	g.Print()
	g.Solve(v)
	fmt.Println("solved")
	g.Print()
}
