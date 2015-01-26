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
