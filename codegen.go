package genus

//go:generate stringer -type PlanType types.go
//go:generate go-bindata -pkg schema -o ./cmd/genus/schema/schema.go cmd/genus/schema
