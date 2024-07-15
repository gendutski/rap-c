package entity

const (
	// path
	// web page auth path
	WebLoginPath              string = "/login"
	WebLogoutPath             string = "/logout"
	WebRequestResetPath       string = "/request-reset"
	WebResetPasswordPath      string = "/reset-password"
	WebPasswordMustChangePath string = "/renew-password"
	// profile page path
	WebProfilePath string = "/profile"
	// default authorized path
	WebDefaultAuthorizedPath = WebProfilePath

	// api route
	ApiLoginRouteName      string = "apiLogin"
	ApiGuestLoginRouteName string = "apiGuestLogin"
)
