package validator

import (
	"fmt"

	"github.com/ostlerc/nurikabe/tile"
)

type nurikabe struct {
	tiles []*tile.Tile
	row   int
	col   int
}

func NewNurikabe() GridValidator {
	return &nurikabe{}
}

func (n *nurikabe) CheckWin(Tiles []*tile.Tile, row, col int) bool {
	n.tiles = Tiles
	n.row = row
	n.col = col
	if !n.hasBlock() && n.singleWall() && n.gardensAreCorrect() {
		fmt.Println("WINNER!")
		return true
	}
	return false
}

// This function help detect quad blocks
func (n *nurikabe) hasBlock() bool {
	for i, _ := range n.tiles {
		if i/n.col == n.row-1 || // bottom of grid
			i%n.col == n.col-1 || // right side of grid
			n.openAt(i) ||
			n.openAt(i+1) ||
			n.openAt(i+n.col) ||
			n.openAt(i+n.col+1) {
			continue
		}
		return true
	}
	return false
}

func (n *nurikabe) openAt(i int) bool {
	return n.tiles[i].Open()
}

func (n *nurikabe) gardensAreCorrect() bool {
	for i, _ := range n.tiles {
		if c := n.tiles[i].Count(); c > 0 {
			openTiles := make(map[int]bool)
			if x := n.markOpen(i, openTiles); x != c {
				return false
			}
		}
	}
	return true
}

func (n *nurikabe) singleWall() bool {
	firstWall := -1
	wallCount := 0
	for i, _ := range n.tiles {
		if !n.openAt(i) {
			if firstWall == -1 {
				firstWall = i
			}
			wallCount++
		}
	}

	if firstWall == -1 || wallCount == 0 {
		return false
	}

	found := make(map[int]bool)

	return n.markClosed(firstWall, found) == wallCount
}

func (n *nurikabe) markOpen(i int, found map[int]bool) int {
	if i < 0 || i >= len(n.tiles) {
		return 0
	}

	if _, ok := found[i]; ok || !n.openAt(i) {
		return 0
	}

	found[i] = true
	ret := 1

	if i/n.col != n.row-1 { // not bottom of grid
		ret += n.markOpen(i+n.col, found)
	}

	if i >= n.col { // not top of grid
		ret += n.markOpen(i-n.col, found)
	}

	if i%n.col != n.row-1 { // not right side of grid
		ret += n.markOpen(i+1, found)
		ret += n.markOpen(i+n.col+1, found)
		ret += n.markOpen(i-n.col+1, found)
	}

	if i%n.col != 0 { // not left side of grid
		ret += n.markOpen(i-1, found)
		ret += n.markOpen(i+n.col-1, found)
		ret += n.markOpen(i-n.col-1, found)
	}

	return ret
}

func (n *nurikabe) markClosed(i int, found map[int]bool) int {
	if i < 0 || i >= len(n.tiles) {
		return 0
	}

	if _, ok := found[i]; ok || n.openAt(i) {
		return 0
	}

	found[i] = true
	ret := 1

	ret += n.markClosed(i+1, found)
	ret += n.markClosed(i-1, found)
	ret += n.markClosed(i+n.col, found)
	ret += n.markClosed(i-n.col, found)

	return ret
}
