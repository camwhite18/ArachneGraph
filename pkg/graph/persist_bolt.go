package graph

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type BoltPersist struct {
	Path string
}

func (bp *BoltPersist) Save(g Graph) error {
	db, err := bolt.Open(bp.Path, 0600, nil)
	if err != nil {
		return fmt.Errorf("could not open bolt database: %v", err)
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		nodes, err := tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			return fmt.Errorf("could not create nodes bucket: %v", err)
		}
		edges, err := tx.CreateBucketIfNotExists([]byte("edges"))
		if err != nil {
			return fmt.Errorf("could not create edges bucket: %v", err)
		}

		for id, node := range g.Nodes() {
			data, err := json.Marshal(node)
			if err != nil {
				return fmt.Errorf("could not serialize node %v: %v", node, err)
			}
			if err := nodes.Put([]byte(id), data); err != nil {
				return fmt.Errorf("could not serialize node %v: %v", node, err)
			}
		}
		for from, outMap := range g.OutEdges() {
			for to, edge := range outMap {
				b, err := json.Marshal(edge)
				if err != nil {
					return fmt.Errorf("could not serialize edge %v: %v", edge, err)
				}
				key := []byte(from + "->" + to)
				if err := edges.Put(key, b); err != nil {
					return fmt.Errorf("could not serialize edge %v: %v", edge, err)
				}
			}
		}
		return nil
	})
}

func (bp *BoltPersist) Load() (Graph, error) {
	g := NewGraph()
	db, err := bolt.Open(bp.Path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open bolt database: %v", err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		nodes := tx.Bucket([]byte("nodes"))
		if nodes == nil {
			return fmt.Errorf("nodes bucket does not exist")
		}
		edges := tx.Bucket([]byte("edges"))
		if edges == nil {
			return fmt.Errorf("edges bucket does not exist")
		}

		err := nodes.ForEach(func(k, v []byte) error {
			var node Node
			if err := json.Unmarshal(v, &node); err != nil {
				return fmt.Errorf("could not deserialize node %s: %v", k, err)
			}
			if err := g.AddNode(&node); err != nil {
				return fmt.Errorf("could not add node %v: %v", node, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("could not deserialize nodes bucket: %v", err)
		}

		err = edges.ForEach(func(k, v []byte) error {
			var edge Edge
			if err := json.Unmarshal(v, &edge); err != nil {
				return fmt.Errorf("could not deserialize edge %s: %v", k, err)
			}
			if err := g.AddEdge(&edge); err != nil {
				return fmt.Errorf("could not add edge %v: %v", edge, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("could not deserialize edges bucket: %v", err)
		}

		return nil
	})
	return g, err
}
