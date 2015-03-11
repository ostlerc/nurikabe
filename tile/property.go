package tile

type PropertyHolder interface {
	Int(string) int
	Set(string, interface{})
	Destroy()
}
