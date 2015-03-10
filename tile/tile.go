package tile

import (
	"log"

	"gopkg.in/qml.v1"
)

var TileCreator Creator

type Creator interface {
	Create(*qml.Context) qml.Object
}

type PropertyHolder interface {
	Int(string) int
	Set(string, interface{})
	Destroy()
}

type Tile struct {
	Properties PropertyHolder
	x          int
	y          int
}

func Setup(engine *qml.Engine) {
	tileComponent, err := engine.LoadFile("qml/tile.qml")
	if err != nil {
		log.Fatal(err)
	}
	TileCreator = tileComponent
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
