package route

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
)

type Route interface {
	// Get a unique hash value for the route.
	//
	// Route implementations must ensure Hash is always unique for two different Routes.
	Hash() string
	// Get the length of the route.
	//
	// Route implementations may determine how to represent their own length.
	Length() int
	// Attach middleware to the route.
	//
	// Route implementations may define the order that middleware is handled.
	Attach(middleware.Middleware)
	// Match a request and update its context.
	//
	// Route implementations must return nil if a request does not match the Route, but may otherwise define any return behavior.
	MatchAndUpdateContext(*http.Request) *http.Request
}

// Create a new Route based on a string expression.
func New(method, expr string, confs ...ConfigFunc) (Route, error) {
	// Determine route type
	var r Route
	var err error
	if isPartialRouteExpr(expr) {
		r, err = build_partialRoute(method, expr)
	} else {
		r, err = build_defaultRoute(method, expr)
	}
	if err != nil {
		return nil, err
	}
	for _, conf := range confs {
		err = conf(r)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Create a new Route based on a string expression, and panic if this fails.
// You should not use this unless you are creating a route on program start and do not intend to modify the route after the fact.
func Declare(method, expr string, confs ...ConfigFunc) Route {
	r, err := New(method, expr, confs...)
	if err != nil {
		panic(err)
	}
	for _, conf := range confs {
		err = conf(r)
		if err != nil {
			panic(err)
		}
	}
	return r
}
