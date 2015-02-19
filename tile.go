package main

import "gopkg.in/qml.v1"

type Tile struct {
	Object qml.Object
	x      int
	y      int
}

type JSONTile struct {
	Count int `json:"count,omitempty"`
	Index int `json:"index,omitempty"`
}

func (t *Tile) Open() bool {
	return t.Object.Int("type") == 0 //open spot
}

func (t *Tile) Count() int {
	return t.Object.Int("count")
}

func (t *Tile) SetType(_t int) {
	t.Object.Set("type", _t)
}
