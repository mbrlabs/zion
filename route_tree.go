package hodor

import (
	"strings"
)

// ============================================================================
//                              struct node
// ============================================================================
type node struct {
	part  string
	route *Route
	nodes []*node
}

// ============================================================================
//                              struct RouteTree
// ============================================================================
type RouteTree struct {
	root *node
}

func NewRouteTree() *RouteTree {
	return &RouteTree{root: &node{part: "", route: nil}}
}

func (t *RouteTree) InsertRoute(route *Route) {
	parts := strings.Split(route.GetPath(), "/")

	currentNode := t.root
	for _, part := range parts {
		if currentNode.part == part {

		}
	}

	// TODO implement
}

func (t *RouteTree) GetRoute(pattern string) *Route {
	// TODO implement
	return nil
}
