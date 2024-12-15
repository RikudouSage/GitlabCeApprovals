package router

import "GitlabCeForcedApprovals/http"

type ControllerMethod func(request *http.Request) (*http.Response, error)

type Route struct {
	Path             string
	ControllerMethod ControllerMethod
}

func NewRoute(path string, controller ControllerMethod) *Route {
	return &Route{Path: path, ControllerMethod: controller}
}

type Router struct {
	Routes []*Route
}

func NewRouter() *Router {
	return &Router{
		Routes: make([]*Route, 0),
	}
}

func (receiver *Router) AddRoute(route *Route) {
	receiver.Routes = append(receiver.Routes, route)
}
