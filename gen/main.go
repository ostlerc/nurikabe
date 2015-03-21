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
)

func init() {
	flag.Parse()
}

func main() {
	g := grid.New(validator.NewNurikabe(), nil)

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		err := g.LoadGrid(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		g.BuildGrid(*width, *height)
	}

	j, err := g.Json()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
