package tile

type Tile struct {
	Properties PropertyHolder
}

func New(parent interface{}) *Tile {
	tile := &Tile{
		Properties: TileCreator.Create(),
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

func (t *Tile) Reset() {
	t.Properties.Set("type", 1) //closed
	t.Properties.Set("count", 0)
}
