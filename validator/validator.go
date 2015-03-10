package validator

import (
	"github.com/ostlerc/nurikabe/tile"
)

type GridValidator interface {
	CheckWin(Tiles []*tile.Tile, row, col int) bool
}
