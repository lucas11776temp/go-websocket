package router

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"websocket/server/connection"
	"websocket/server/request"
	"websocket/server/response"
)

type WebMethod func()

type WsMethod func(req *request.Request, ws *connection.Connection)

type Next func() *response.Response

type Middleware func(req *request.Request, res *response.Response, next Next) *response.Response

type Group func(route *Route)

type RouteParams map[string]string

type Middlewares []Middleware

type Route struct {
	path        string
	method      string
	middlewares Middlewares
	callback    reflect.Value
	option      *RouteOption
	parameters  RouteParams
	Routes      *GroupRoutes
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
func (ctx *Route) GetMiddlewares() Middlewares {
	return ctx.middlewares
}

// Comment
func matchRoute(routePathArr []string, pathArr []string) (RouteParams, error) {
	parameters := make(RouteParams)

	for i, segment := range routePathArr {
		regex, _ := regexp.Compile("\\{[a-zA-Z_]+\\}")
		param := string(regex.Find([]byte(segment)))

		if param != "" {
			parameters[strings.Trim(strings.Trim(param, "{"), "}")] = pathArr[i]
			continue
		}

		if segment == pathArr[i] {
			continue
		}

		return parameters, fmt.Errorf("Route %s does not exist", strings.Join(pathArr, "/"))
	}

	return parameters, nil
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

		parameters, err := matchRoute(routePathArr, pathArr)

		if err != nil {
			continue
		}

		route.parameters = parameters

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
func (ctx *RouteOption) Middleware(middleware ...Middleware) *RouteOption {
	ctx.route.Middleware(middleware...)

	return ctx
}

// Comment
func (ctx *Route) addWebRoute(method string, path string, callback WebMethod) *Route {
	route := ctx.route(method, path, reflect.ValueOf(callback))

	ctx.Routes.webRoutes = append(ctx.Routes.webRoutes, *route)

	return route
}

// Comment
func (ctx *Route) addWsRoute(method string, path string, callback WsMethod) *Route {
	route := ctx.route(method, path, reflect.ValueOf(callback))

	ctx.Routes.wsRoutes = append(ctx.Routes.wsRoutes, *route)

	return route
}

// Comment
func (ctx *Route) Ws(path string, method WsMethod) *RouteOption {
	return ctx.addWsRoute("GET", path, method).option
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
func (ctx *Route) route(method string, path string, value reflect.Value) *Route {
	route := Route{
		path:        JoinPath(ctx.path, path),
		method:      strings.ToUpper(method),
		middlewares: ctx.middlewares,
		callback:    value,
	}

	route.option = &RouteOption{route: &route}

	return &route
}

// Comment
func (ctx *Route) Get(path string, method WebMethod) *RouteOption {
	return ctx.addWebRoute("GET", path, method).option
}

// Comment
func (ctx *Route) Post(path string, method WebMethod) *RouteOption {
	return ctx.addWebRoute("POST", path, method).option
}

// Comment
func (ctx *Route) Patch(path string, method WebMethod) *RouteOption {
	return ctx.addWebRoute("PATCH", path, method).option
}

// Comment
func (ctx *Route) Delete(path string, method WebMethod) *RouteOption {
	return ctx.addWebRoute("DELETE", path, method).option
}

// Comment
func (ctx *Route) Middleware(middleware ...Middleware) *Route {
	ctx.middlewares = append(ctx.middlewares, middleware...)
	return ctx
}

// Comment
func (ctx *Route) Parameters() RouteParams {
	return ctx.parameters
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
			path:        JoinPath(ctx.path, prefix),
			middlewares: ctx.middlewares,
			Routes:      ctx.Routes,
		},
	)
}
