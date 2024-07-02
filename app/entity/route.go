package entity

const (
	LoginRouteName      string = "loginPage"
	PostLoginRouteName  string = "submitLogin"
	PostLogoutRouteName string = "submitLogout"
	ProfileRouteName    string = "profile"

	// route will call when user is authoriezed
	DefaultAuthorizedRouteRedirect string = ProfileRouteName
)
