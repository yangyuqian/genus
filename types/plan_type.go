package types

type PlanType int

const (
	SINGLETON PlanType = iota
	REPEATABLE
)
