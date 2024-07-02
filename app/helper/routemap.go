package helper

import "github.com/labstack/echo/v4"

type RouteMaper map[string]*echo.Route

func (e RouteMaper) Get(key, field string) string {
	route := e[key]
	if route != nil {
		switch field {
		case "method":
			return route.Method
		case "path":
			return route.Path
		case "name":
			return route.Name
		}
	}
	return ""
}

func RouteMap(payload []*echo.Route) RouteMaper {
	result := make(map[string]*echo.Route)
	for _, r := range payload {
		result[r.Name] = r
	}
	return result
}
