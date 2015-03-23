package validator

type GridData interface {
	Open(int) bool
	Count(int) int
	Rows() int
	Columns() int
}

type GridValidator interface {
	CheckWin(GridData) bool
}
