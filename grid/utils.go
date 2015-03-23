package grid

import "fmt"

type mapset map[int]int

func (m mapset) Print(cols int) {
	for i := 0; i < len(m); i += cols {
		for j := 0; j < cols; j++ {
			fmt.Print(m[i+j], " ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func removeAt(i int, a []int) []int {
	a[i], a[len(a)-1], a = a[len(a)-1], 0, a[:len(a)-1]
	return a
}

func remove(v int, a []int) []int {
	for i, x := range a {
		if x == v {
			return removeAt(i, a)
		}
	}
	return a
}

func (g *Grid) reset() {
	for _, t := range g.tiles {
		t.open = true
		t.count = 0
	}
}

func (g *Grid) Print() {
	for i := 0; i < len(g.tiles); i += g.cols {
		for j := 0; j < g.cols; j++ {
			if c := g.tiles[i+j].count; c > 0 {
				fmt.Print(c, " ")
			} else if g.tiles[i+j].open {
				fmt.Print("o ")
			} else {
				fmt.Print("x ")
			}
		}
		fmt.Println()
	}
}
