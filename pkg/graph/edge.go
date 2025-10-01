package graph

type Edge struct {
	ID         string
	From, To   string
	Weight     float64
	Properties map[string]any
}
