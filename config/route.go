package config

import (
	"reflect"
)

type routeDetail struct {
	method string
	path   string
}

func (cfg routeDetail) Method() string {
	return cfg.method
}

func (cfg routeDetail) Path() string {
	return cfg.path
}

type Route struct {
	// api
	LoginAPI                routeDetail `method:"POST" path:"/api/login"`
	GuestLoginAPI           routeDetail `method:"POST" path:"/api/guest-login"`
	PasswordMustChangeAPI   routeDetail `method:"PUT" path:"/api/renew-password"`
	RequestResetPasswordAPI routeDetail `method:"POST" path:"/api/request-reset"`
	ResetPasswordAPI        routeDetail `method:"POST" path:"/api/reset-password"`
	DetailUserAPI           routeDetail `method:"GET" path:"/api/user/detail/:username"`
	ListUserAPI             routeDetail `method:"GET" path:"/api/user/list"`
	TotalUserAPI            routeDetail `method:"GET" path:"/api/user/total"`
	CreateUserAPI           routeDetail `method:"POST" path:"/api/user/create"`
	UpdateUserAPI           routeDetail `method:"PUT" path:"/api/user/update"`
	SetStatusUserAPI        routeDetail `method:"PUT" path:"/api/user/active-status"`
	ListUnitAPI             routeDetail `method:"GET" path:"/api/unit/list"`
	TotalUnitAPI            routeDetail `method:"GET" path:"/api/unit/total"`
	CreateUnitAPI           routeDetail `method:"POST" path:"/api/unit/create"`
	DeleteUnitAPI           routeDetail `method:"DELETE" path:"/api/unit/delete"`

	// web
	LoginWebPage              routeDetail `method:"GET" path:"/login"`
	SubmitTokenSessionWebPage routeDetail `method:"POST" path:"/submit-token-session"`
	LogoutWebPage             routeDetail `method:"POST" path:"/logout"`
	PasswordMustChangeWebPage routeDetail `method:"GET" path:"/renew-password"`
	ForgotPasswordWebPage     routeDetail `method:"GET" path:"/forgot-password"`
	ResetPasswordWebPage      routeDetail `method:"GET" path:"/reset-password"`
	DashboardWebPage          routeDetail `method:"GET" path:"/dashboard"`
	ProfileWebPage            routeDetail `method:"GET" path:"/profile"`
}

func (cfg *Route) DefaultAuthorizedWebPage(method, path string) routeDetail {
	if method != "" && path != "" {
		return routeDetail{
			method: method,
			path:   path,
		}
	}
	return cfg.DashboardWebPage
}

func InitRoute() *Route {
	var route Route
	setRouteDetails(&route)
	return &route
}

func setRouteDetails(route *Route) {
	val := reflect.ValueOf(route).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		methodTag := structField.Tag.Get("method")
		pathTag := structField.Tag.Get("path")

		if methodTag != "" && pathTag != "" {
			detail := routeDetail{
				method: methodTag,
				path:   pathTag,
			}
			field.Set(reflect.ValueOf(detail))
		}
	}
}
