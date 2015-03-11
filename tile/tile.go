package tile

type Tile struct {
	Properties PropertyHolder
	x          int
	y          int
}

func New(parent interface{}) *Tile {
	tile := &Tile{
		Properties: TileCreator.Create(nil),
	}
	tile.Properties.Set("parent", parent)
	return tile
}

func (t *Tile) Open() bool {
	return t.Properties.Int("type") == 0 //open spot
}

func (t *Tile) Count() int {
	return t.Properties.Int("count")
}

func (t *Tile) SetType(_t int) {
	t.Properties.Set("type", _t)
}
