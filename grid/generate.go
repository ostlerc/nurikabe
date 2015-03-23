package grid

import "github.com/ostlerc/nurikabe/validator"

const (
	opened = iota
	closed = iota
	sealed = iota
)

//TODO: add difficulty parameter
func (g *Grid) Generate(v validator.GridValidator, minGardens, gardenSize, base int) {
	tileMap := make(mapset, len(g.tiles))
	for {
		g.reset()
		for i := 0; i < len(g.tiles); i++ {
			tileMap[i] = closed
		}

		c := 0
		for ; g.placeGarden(R.Intn(gardenSize)+base, tileMap); c++ {
		}

		if c < minGardens {
			continue
		}

		for i, t := range g.tiles {
			t.open = tileMap[i] == opened
		}

		if v.CheckWin(g) {
			return
		}
	}
}

func (g *Grid) placeGarden(max int, tileMap mapset) bool {
	i := -1
	for c := 0; c < 10; c++ {
		z := R.Intn(len(tileMap))
		if tileMap[z] == closed {
			i = z
			break
		}
	}
	if i == -1 {
		for k, v := range tileMap {
			if v == closed {
				i = k
				break
			}
		}
		if i == -1 {
			return false
		}
	}
	tiles := g.markOpen(i, max, tileMap)
	if len(tiles) < 2 {
		return false
	}
	g.tiles[i].open = true
	g.tiles[i].count = len(tiles)

	return true
}

func (g *Grid) markOpen(i, c int, tileMap mapset) []int {
	if c == 0 || tileMap[i] == sealed || tileMap[i] == opened {
		return []int{}
	}
	steps := []int{1, -1, g.cols, -g.cols}

	if i/g.cols == g.rows-1 { // bottom of grid
		steps = remove(g.cols, steps)
	}

	if i < g.cols { // top of grid
		steps = remove(-g.cols, steps)
	}

	if i%g.cols == g.cols-1 { // right side of grid
		steps = remove(1, steps)
	}

	if i%g.cols == 0 { // left side of grid
		steps = remove(-1, steps)
	}

	remainingSteps := make([]int, len(steps))
	copy(remainingSteps, steps)

	ret := make([]int, 0, c)
	ret = append(ret, i)
	c--
	tileMap[i] = opened
	for c > 0 && len(remainingSteps) > 0 {
		stepIndex := R.Intn(len(remainingSteps))
		v := remainingSteps[stepIndex] + i
		remainingSteps = removeAt(stepIndex, remainingSteps)

		tList := g.markOpen(v, c, tileMap)
		if l := len(tList); l > 0 {
			c -= l
			ret = append(ret, tList...)
		}
	}

	seal := func(x int) {
		if tileMap[x] == closed {
			tileMap[x] = sealed
		}
	}

	//seal up boundaries
	for _, s := range steps {
		seal(s + i)
	}
	return ret
}
