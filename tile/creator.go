package tile

import (
	"log"

	"gopkg.in/qml.v1"
)

var TileCreator Creator

type QMLObjectCreator struct {
	o qml.Object
}

func (q *QMLObjectCreator) Create() PropertyHolder {
	return q.o.Create(nil)
}

type Creator interface {
	Create() PropertyHolder
}

func SetupGui(engine *qml.Engine, path string) {
	tileComponent, err := engine.LoadFile(path)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	TileCreator = &QMLObjectCreator{o: tileComponent}
}
