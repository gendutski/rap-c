package entity

const (
	// path
	WebLoginPath          string = "/login"
	WebPasswordChangePath string = "/password-change"

	// web route
	LoginRouteName         string = "loginPage"
	GuestLoginRouteName    string = "guestLoginPage"
	PostLoginRouteName     string = "submitLogin"
	PostLogoutRouteName    string = "submitLogout"
	ProfileRouteName       string = "profile"
	RenewPasswordRouteName string = "renewPassword"

	// api route

	// route will call when user is authoriezed
	DefaultAuthorizedRouteRedirect string = ProfileRouteName
)
