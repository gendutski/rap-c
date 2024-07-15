package mailusecase

import (
	"net/http"
	"net/url"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"
)

const (
	welcomeSubject string = "Selamat datang di Rap-C"
	resetSubject   string = "Permintaan reset pasword di Rap-C"
)

func NewUsecase(cfg *config.Config, router *config.Route) contract.MailUsecase {
	return &usecase{cfg, router}
}

type usecase struct {
	cfg    *config.Config
	router *config.Route
}

func (uc *usecase) Welcome(user *databaseentity.User, password string) error {
	// init hermes
	h := uc.initHermes()

	// encode email & password
	params := url.Values{}
	params.Add("email", user.Email)
	params.Add("password", password)

	// init hermes email
	email := hermes.Email{
		Body: hermes.Body{
			Greeting: "Hai",
			Name:     user.FullName,
			Intros: []string{
				"Selamat datang di `rap-c`, aplikasi yang dirancang untuk memudahkan pengelolaan resep, menghitung harga pokok penjualan, mengelola stok bahan baku, dan menyimpan catatan transaksi dalam general ledger sederhana.",
				"Anda telah diajukan menjadi user di aplikasi ini, dengan detail sebagai berikut:",
			},
			Dictionary: []hermes.Entry{
				{Key: "Nama Lengkap", Value: user.FullName},
				{Key: "Username", Value: user.Username},
				{Key: "Password", Value: password},
			},
			Actions: []hermes.Action{
				{
					Instructions: "Silahkan klik tombol dibawah untuk login:",
					Button: hermes.Button{
						Text: "Login menggunakan akun anda",
						Link: uc.cfg.URL(uc.router.LoginWebPage.Path()) + "?" + params.Encode(),
					},
				},
			},
			Signature: "Hormat Kami",
		},
	}

	// generate html email
	resHtml, err := h.GenerateHTML(email)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.MailUsecaseGenerateHTMLError, err.Error()),
		}
	}

	// generate text html
	resText, err := h.GeneratePlainText(email)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.MailUsecaseGeneratePlainTextError, err.Error()),
		}
	}

	return uc.send(user.Email, welcomeSubject, resText, resHtml)
}

func (uc *usecase) ResetPassword(user *databaseentity.User, token *databaseentity.PasswordResetToken) error {
	// init hermes
	h := uc.initHermes()

	// encode email & password
	params := url.Values{}
	params.Add("email", user.Email)
	params.Add("token", token.Token)

	// init hermes email
	email := hermes.Email{
		Body: hermes.Body{
			Greeting: "Hai",
			Name:     user.FullName,
			Intros: []string{
				"Anda menerima email ini karena permintaan reset password untuk akun anda telah diterima.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Silahkan klik tombol dibawah untuk reset password Anda:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset your password",
						Link:  uc.cfg.URL(uc.router.ResetPasswordWebPage.Path()) + "?" + params.Encode(),
					},
				},
			},
			Outros: []string{
				"Jika Anda tidak meminta reset password, Anda tidak perlu melakukan tindakan lebih lanjut.",
			},
			Signature: "Hormat Kami",
		},
	}

	// generate html email
	resHtml, err := h.GenerateHTML(email)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.MailUsecaseGenerateHTMLError, err.Error()),
		}
	}

	// generate text html
	resText, err := h.GeneratePlainText(email)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.MailUsecaseGeneratePlainTextError, err.Error()),
		}
	}

	return uc.send(user.Email, resetSubject, resText, resHtml)
}
