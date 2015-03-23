package validator

type GridSolver interface {
	Solve(GridData, GridValidator) []bool
}
