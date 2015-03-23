package validator

type GardenMap map[int]int

type GridSolver interface {
	Solve(GridData, GridValidator) []bool
}
