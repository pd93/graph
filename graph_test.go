package graph

import "testing"

func TestGraph_Vertex(t *testing.T) {
	tests := map[string]struct {
		vertices         []int
		expectedVertices []int
	}{
		"graph with 3 vertices": {
			vertices:         []int{1, 2, 3},
			expectedVertices: []int{1, 2, 3},
		},
		"graph with duplicated vertex": {
			vertices:         []int{1, 2, 2},
			expectedVertices: []int{1, 2},
		},
	}

	for name, test := range tests {
		graph := New(IntHash)

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}

		for _, vertex := range test.vertices {
			hash := graph.hash(vertex)
			if _, ok := graph.vertices[hash]; !ok {
				t.Errorf("%s: vertex %v not found in graph: %v", name, vertex, graph.vertices)
			}
		}
	}
}

func TestGraph_Edge(t *testing.T) {
	TestGraph_WeightedEdge(t)
}

func TestGraph_WeightedEdge(t *testing.T) {
	TestGraph_WeightedEdgeByHashes(t)
}

func TestGraph_EdgeByHashes(t *testing.T) {
	TestGraph_WeightedEdgeByHashes(t)
}

func TestGraph_WeightedEdgeByHashes(t *testing.T) {
	tests := map[string]struct {
		vertices      []int
		edgeHashes    [][3]int
		expectedEdges []Edge[int]
		shouldFail    bool
	}{
		"graph with 2 edges": {
			vertices:   []int{1, 2, 3},
			edgeHashes: [][3]int{{1, 2, 10}, {1, 3, 20}},
			expectedEdges: []Edge[int]{
				{Source: 1, Target: 2, Weight: 10},
				{Source: 1, Target: 3, Weight: 20},
			},
		},
		"hashes for non-existent vertices": {
			vertices:   []int{1, 2},
			edgeHashes: [][3]int{{1, 3, 20}},
			shouldFail: true,
		},
	}

	for name, test := range tests {
		graph := New(IntHash)

		for _, vertex := range test.vertices {
			graph.Vertex(vertex)
		}
		for _, edge := range test.edgeHashes {
			err := graph.WeightedEdgeByHashes(edge[0], edge[1], edge[2])

			if test.shouldFail != (err != nil) {
				t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
			}
		}
		for _, expectedEdge := range test.expectedEdges {
			sourceHash := graph.hash(expectedEdge.Source)
			targetHash := graph.hash(expectedEdge.Target)

			edge, ok := graph.edges[sourceHash][targetHash]
			if !ok {
				t.Fatalf("%s: edge with source %v and target %v not found", name, expectedEdge.Source, expectedEdge.Target)
			}

			if edge.Source != expectedEdge.Source {
				t.Errorf("%s: edge sources don't match: expected source %v, got %v", name, expectedEdge.Source, edge.Source)
			}

			if edge.Target != expectedEdge.Target {
				t.Errorf("%s: edge targets don't match: expected target %v, got %v", name, expectedEdge.Target, edge.Target)
			}

			if edge.Weight != expectedEdge.Weight {
				t.Errorf("%s: edge weights don't match: expected weight %v, got %v", name, expectedEdge.Weight, edge.Weight)
			}
		}
	}
}

func TestGraph_GetEdge(t *testing.T) {
	TestGraph_GetEdgeByHashes(t)
}

func TestGraph_GetEdgeByHashes(t *testing.T) {
	tests := map[string]struct {
		graph         *Graph[int, int]
		vertices      []int
		edgeHashes    [2]int
		getEdgeHashes [2]int
		expectedEdge  Edge[int]
		shouldFail    bool
	}{
		"get edge of undirected graph": {
			graph:         New(IntHash),
			vertices:      []int{1, 2, 3},
			edgeHashes:    [2]int{1, 2},
			getEdgeHashes: [2]int{2, 1},
			expectedEdge:  Edge[int]{Source: 1, Target: 2},
		},
		"get non-existent edge of undirected graph": {
			graph:         New(IntHash),
			vertices:      []int{1, 2, 3},
			edgeHashes:    [2]int{1, 2},
			getEdgeHashes: [2]int{1, 3},
			shouldFail:    true,
		},
		"get edge of directed graph": {
			graph:         New(IntHash, Directed()),
			vertices:      []int{1, 2, 3},
			edgeHashes:    [2]int{1, 2},
			getEdgeHashes: [2]int{1, 2},
			expectedEdge:  Edge[int]{Source: 1, Target: 2},
		},
		"get non-existent edge of directed graph": {
			graph:         New(IntHash, Directed()),
			vertices:      []int{1, 2, 3},
			edgeHashes:    [2]int{1, 2},
			getEdgeHashes: [2]int{1, 3},
			shouldFail:    true,
		},
	}

	for name, test := range tests {
		for _, vertex := range test.vertices {
			test.graph.Vertex(vertex)
		}

		test.graph.EdgeByHashes(test.edgeHashes[0], test.edgeHashes[1])

		edge, err := test.graph.GetEdgeByHashes(test.getEdgeHashes[0], test.getEdgeHashes[1])

		if test.shouldFail != (err != nil) {
			t.Fatalf("%s: error expectancy doesn't match: expected %v, got %v (error: %v)", name, test.shouldFail, (err != nil), err)
		}

		if test.graph.properties.isDirected {
			if edge.Source != test.expectedEdge.Source || edge.Target != test.expectedEdge.Target {
				t.Errorf("%s: edges don't match: expected %v, got %v", name, test.expectedEdge, edge)
			}
		} else {
			if edge.Source != test.expectedEdge.Source && edge.Source != test.expectedEdge.Target ||
				edge.Target != test.expectedEdge.Target && edge.Target != test.expectedEdge.Source {
				t.Errorf("%s: edges don't match: expected %v, got %v", name, test.expectedEdge, edge)
			}
		}

	}
}

func slicesAreEqual[T any](a []T, b []T, equals func(a, b T) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aValue := range a {
		found := false
		for _, bValue := range b {
			if equals(aValue, bValue) {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}
