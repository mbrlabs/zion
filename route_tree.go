package hodor

import (
	"net/http"
	"strings"
)

// ============================================================================
//                              struct node
// ============================================================================

type node struct {
	part  string
	route *route
	nodes []*node
}

// Returns one of the child nodes, that match the given pattern.
// If such a node (route) does not exist, it will be inserted and returned.
func (n *node) getOrInsertNode(part string) *node {
	// panic if tying to add a wildcard to child nodes if parent already has child nodes
	isWildcard := strings.HasPrefix(part, "*")
	if isWildcard && len(n.nodes) > 0 {
		panic("Ambigious mapping found: wildcards can't be mixed with static urls or named parameters")
	}

	// try to find existing match
	isParam := strings.HasPrefix(part, ":")
	for _, childNode := range n.nodes {
		// panic if a different named parameter is added to a list of child nodes, which already
		// have another named parameter
		if childNode.isParam() && isParam && childNode.part != part {
			panic("Ambigious mapping of named parameters found.")
		}

		if childNode.part == part {
			return childNode
		}
	}

	// create new node if nothing found
	newNode := &node{part: part, route: nil}
	n.nodes = append(n.nodes, newNode)
	return newNode
}

func (n *node) isParam() bool {
	return strings.HasPrefix(n.part, ":")
}

func (n *node) isWildcard() bool {
	return strings.HasPrefix(n.part, "*")
}

// ============================================================================
//                              struct routeTree
// ============================================================================

type routeTree struct {
	treeRoots map[string]*node
}

func newRouteTree() routeTree {
	roots := make(map[string]*node)
	roots[http.MethodGet] = &node{part: "", route: nil}
	roots[http.MethodHead] = &node{part: "", route: nil}
	roots[http.MethodPost] = &node{part: "", route: nil}
	roots[http.MethodPut] = &node{part: "", route: nil}
	roots[http.MethodDelete] = &node{part: "", route: nil}
	roots[http.MethodOptions] = &node{part: "", route: nil}
	return routeTree{treeRoots: roots}
}

// TODO check for named params & wildcard params witch same name
// Tries to insert a route into the route tree. Panics if mappings are ambigious
func (rt *routeTree) insert(r *route) {
	parts := strings.Split(r.getPattern(), "/")
	// get tree root, corresponding to the http method
	root := rt.treeRoots[r.method]
	if root == nil {
		panic("Unsupported http method: " + r.method)
	}

	// check if wildcard is at leaf. If not, panic.
	for i, part := range parts {
		if i < len(parts)-1 && strings.HasPrefix(part, "*") {
			panic("Wildcards must be a leaf node")
		}
	}

	// handle the root pattern ("")
	if len(parts) == 1 && len(parts[0]) == 0 {
		if root.route == nil {
			root.route = r
			return
		}
		panic("Abmigious mapping for the root pattern")
	}

	// handle other patterns
	currentNode := root
	for i, part := range parts {
		currentNode = currentNode.getOrInsertNode(part)
		// end of pattern reached. we need to assign the route here
		if i == len(parts)-1 {
			if currentNode.route == nil {
				currentNode.route = r
				return
			}
			panic("Abmigious mapping for: " + r.pattern)
		}
	}

	panic("Something went terribly wrong while mapping the route: " + r.pattern)
}

// Returns a route and sets the url parameters of the context.
func (rt *routeTree) get(ctx *Context) *route {
	// get tree root, corresponding to the http method
	root := rt.treeRoots[ctx.Request.Method]
	if root == nil {
		return nil
	}

	path := strings.Trim(ctx.Request.URL.Path, "/")

	// deny everything that contains a colon or star
	if strings.Contains(path, ":") || strings.Contains(path, "*") {
		return nil
	}

	// handle root path
	if len(path) == 0 {
		return root.route
	}

	// handle everything else
	parts := strings.Split(path, "/")
	currentNode := root
	var namedParam *node
	foundPart := false
	for _, part := range parts {
		foundPart = false
		namedParam = nil
		for _, childNode := range currentNode.nodes {
			// found a wildcard match
			if childNode.isWildcard() {
				key := strings.TrimLeft(currentNode.part, "*")
				ctx.UrlParams[key] = part
				return childNode.route
			}
			// found a static match
			if childNode.part == part {
				currentNode = childNode
				foundPart = true
				break
			}
			// found a possible match in a named param
			if childNode.isParam() {
				namedParam = childNode
			}
		}

		// found match in named param, since it's set to a non-nil value.
		// set as current node, extract value and set in context.
		if namedParam != nil && !foundPart {
			currentNode = namedParam
			foundPart = true
			key := strings.TrimLeft(currentNode.part, ":")
			ctx.UrlParams[key] = part
		}

		// if the url part can't be found, the route does not exist
		if !foundPart {
			return nil
		}
	}

	return currentNode.route
}
