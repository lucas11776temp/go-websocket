package router

import (
	"errors"
	"reflect"
	"strings"
	"websocket/server/connection"
	"websocket/server/request"
)

type WebMethod func()

type WsMethod func(req *request.Request, ws *connection.Connection)

type Middleware func(req any, res any, next func())

type Group func(route *Route)

type RouteParameters map[string]string

type Route struct {
	path       string
	method     string
	middleware []Middleware
	callback   reflect.Value
	Routes     *GroupRoutes
	options    RouteOption
	parameters RouteParameters
}

type RouteOption struct {
	route *Route
}

type Routes []Route

type GroupRoutes struct {
	webRoutes Routes
	wsRoutes  Routes
}

// Comment
func (ctx *GroupRoutes) getRoute(routes *Routes, method string, path string) (*Route, error) {
	pathArr := strings.Split(strings.Trim(path, "/"), "/")

	for _, route := range *routes {
		if strings.ToUpper(route.method) != strings.ToUpper(method) {
			continue
		}

		routePathArr := strings.Split(route.path, "/")

		if len(pathArr) != len(routePathArr) {
			continue
		}

		for i, segment := range pathArr {
			// must do params search {id}
			if segment != routePathArr[i] {
				continue
			}
		}

		return &route, nil
	}

	return nil, errors.New("Route " + path + " is not found")
}

// Comment
func (ctx *GroupRoutes) WebRoute(method string, path string) (*Route, error) {
	return ctx.getRoute(&ctx.webRoutes, method, path)
}

// Comment
func (ctx *GroupRoutes) WsRoute(path string) (*Route, error) {
	return ctx.getRoute(&ctx.wsRoutes, "GET", path)
}

// Comment
func (ctx *RouteOption) Middlewares(middleware ...Middleware) *RouteOption {
	ctx.route.middleware = append(ctx.route.middleware, middleware...)
	return ctx
}

// Comment
func (ctx *Route) addWebRoute(method string, path string, callback WebMethod) Route {
	route := ctx.route(method, path, reflect.ValueOf(callback))

	ctx.Routes.webRoutes = append(ctx.Routes.webRoutes, route)

	return route
}

// Comment
func (ctx *Route) addWsRoute(method string, path string, callback WsMethod) Route {
	route := ctx.route(method, path, reflect.ValueOf(callback))

	ctx.Routes.wsRoutes = append(ctx.Routes.wsRoutes, route)

	return route
}

// Comment
func (ctx *Route) Ws(path string, method WsMethod) RouteOption {
	return ctx.addWsRoute("GET", path, method).options
}

// comment
func JoinPath(path ...string) string {
	paths := []string{}

	for _, p := range path {
		if p == "" || p == "/" {
			continue
		}
		paths = append(paths, strings.Trim(p, "/"))
	}

	return strings.Join(paths, "/")
}

// Comment
func (ctx *Route) route(method string, path string, value reflect.Value) Route {
	route := Route{
		path:       JoinPath(ctx.path, path),
		method:     strings.ToUpper(method),
		middleware: ctx.middleware,
		callback:   value,
	}

	route.options = RouteOption{route: &route}

	return route
}

// Comment
func (ctx *Route) Get(path string, method WebMethod) RouteOption {
	return ctx.addWebRoute("GET", path, method).options
}

// Comment
func (ctx *Route) Post(path string, method WebMethod) RouteOption {
	return ctx.addWebRoute("POST", path, method).options
}

// Comment
func (ctx *Route) Patch(path string, method WebMethod) RouteOption {
	return ctx.addWebRoute("PATCH", path, method).options
}

// Comment
func (ctx *Route) Delete(path string, method WebMethod) RouteOption {
	return ctx.addWebRoute("DELETE", path, method).options
}

// Comment
func (ctx *Route) Middlewares(middleware ...Middleware) *Route {
	ctx.middleware = append(ctx.middleware, middleware...)
	return ctx
}

// Comment
func (ctx *Route) Call(parameters ...any) {
	values := []reflect.Value{}

	for _, param := range parameters {
		values = append(values, reflect.ValueOf(param))
	}

	ctx.callback.Call(values)
}

// Comment
func (ctx *Route) Group(prefix string, group Group) {
	group(
		&Route{
			path:       JoinPath(ctx.path, prefix),
			middleware: ctx.middleware,
			Routes:     ctx.Routes,
		},
	)
}
