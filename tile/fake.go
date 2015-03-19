package tile

import (
	"gopkg.in/qml.v1"
)

type fakeCreator struct{}

func (f *fakeCreator) Create(*qml.Context) PropertyHolder {
	return &fakePropertyHolder{m: make(map[string]int)}
}

type fakePropertyHolder struct {
	m map[string]int
}

func (f *fakePropertyHolder) Int(s string) int {
	if v, ok := f.m[s]; ok {
		return v
	}
	return 0
}
func (f *fakePropertyHolder) Set(s string, i interface{}) { f.m[s], _ = i.(int) }
func (f *fakePropertyHolder) Destroy()                    {}

func Fake() PropertyHolder {
	return &fakePropertyHolder{m: make(map[string]int)}
}

func SetupTesting() {
	TileCreator = &fakeCreator{}
}
