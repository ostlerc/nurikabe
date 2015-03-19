package tile

import (
	"log"

	"gopkg.in/qml.v1"
)

var TileCreator Creator

type QMLObjectCreator struct {
	o qml.Object
}

func (q *QMLObjectCreator) Create(c *qml.Context) PropertyHolder {
	return q.o.Create(c)
}

type Creator interface {
	Create(*qml.Context) PropertyHolder
}

func Setup(engine *qml.Engine, path string) {
	tileComponent, err := engine.LoadFile(path)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	TileCreator = &QMLObjectCreator{o: tileComponent}
}
