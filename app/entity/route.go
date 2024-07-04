package entity

const (
	WebLoginPath string = "/login"

	LoginRouteName      string = "loginPage"
	PostLoginRouteName  string = "submitLogin"
	PostLogoutRouteName string = "submitLogout"
	ProfileRouteName    string = "profile"

	// route will call when user is authoriezed
	DefaultAuthorizedRouteRedirect string = ProfileRouteName
)
