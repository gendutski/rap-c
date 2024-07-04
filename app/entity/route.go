package entity

const (
	// path
	WebLoginPath          string = "/login"
	WebPasswordChangePath string = "/password-change"

	// route
	LoginRouteName      string = "loginPage"
	PostLoginRouteName  string = "submitLogin"
	PostLogoutRouteName string = "submitLogout"
	ProfileRouteName    string = "profile"

	// route will call when user is authoriezed
	DefaultAuthorizedRouteRedirect string = ProfileRouteName
)
