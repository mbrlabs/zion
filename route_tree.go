package hodor

import (
    "strings"
)

// ============================================================================
//                              struct node
// ============================================================================
type node struct {
    part    string
    route   *Route
}

// ============================================================================
//                              struct RouteTree
// ============================================================================
type RouteTree struct {
    root *node
}

func (t *RouteTree) InsertRoute(route *Route) {
    parts := strings.Split(route.GetPath(), "/")
    if len(parts) == 0 {

    }
    
    // TODO implement
}

func (t *RouteTree) GetRoute(pattern string) *Route {
    // TODO implement
    return nil
}

