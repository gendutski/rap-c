package entity

const (
	// path
	WebLoginPath          string = "/login"
	WebPasswordChangePath string = "/password-change"
	WebLogoutPath         string = "/logout"
	WebLogoutMethod       string = "post"
	WebResetPasswordPath  string = "/reset-password"

	// web route
	LoginRouteName               string = "loginPage"
	GuestLoginRouteName          string = "guestLoginPage"
	PostLoginRouteName           string = "submitLogin"
	PostLogoutRouteName          string = "submitLogout"
	ProfileRouteName             string = "profile"
	RenewPasswordRouteName       string = "renewPassword"
	RequestResetPasswordName     string = "requestReset"
	PostRequestResetPasswordName string = "submitRequestReset"
	ResetPasswordName            string = "resetPassword"
	SubmitResetPasswordName      string = "submitResetPassword"

	// api route

	// route will call when user is authoriezed
	DefaultAuthorizedRouteRedirect string = ProfileRouteName
)
