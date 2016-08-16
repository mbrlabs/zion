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

func (n *node) nextNode(part string) *node {
	isParam := strings.HasPrefix(part, ":")
	for _, childNode := range n.nodes {
		// panic if a different named parameter is added to a list of child nodes, which already
		// have another named parameter
		if strings.HasPrefix(childNode.part, ":") && isParam && childNode.part != part {
			panic("Ambigious mapping of named parameters found.")
		}

		if childNode.part == part {
			return childNode
		}
	}

	newNode := &node{part: part, route: nil}
	n.nodes = append(n.nodes, newNode)
	return newNode
}

// ============================================================================
//                              struct RouteTree
// ============================================================================
type RouteTree struct {
	root *node
}

func NewRouteTree() RouteTree {
	return RouteTree{root: &node{part: "", route: nil}}
}

func (t *RouteTree) InsertRoute(route *Route) {
	parts := strings.Split(route.GetPath(), "/")

	// handle the root pattern ("")
	if len(parts) == 1 && parts[0] == "" {
		if t.root.route == nil {
			t.root.route = route
			return
		} else {
			panic("Abmigious mapping for the root pattern")
		}
	}

	// handle other patterns
	currentNode := t.root
	for i, part := range parts {
		currentNode = currentNode.nextNode(part)
		// end of pattern reached. we need to assign the route here
		if i == len(parts)-1 {
			if currentNode.route == nil {
				currentNode.route = route
				return
			} else {
				panic("Abmigious mapping for: " + route.path)
			}
		}
	}

	panic("Something went terribly wrong while mapping the route: " + route.path)
}

func (t *RouteTree) GetRoute(pattern string) *Route {
	// TODO implement
	return nil
}
