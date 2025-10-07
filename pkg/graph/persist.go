package graph

type Persister interface {
	Save(g *Graph) error
	Load() (*Graph, error)
}
