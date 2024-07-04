package mailmodule

import (
	"net/http"
	"rap-c/app/entity"
	"rap-c/config"

	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"
)

const (
	welcomeSubject string = "Selamat datang di Rap-C"
)

type MailUsecase interface {
	Welcome(user *entity.User, password string) error
}

func NewUsecase(cfg config.Config) MailUsecase {
	return &usecase{cfg}
}

type usecase struct {
	cfg config.Config
}

func (uc *usecase) Welcome(user *entity.User, password string) error {
	// init hermes
	h := uc.initHermes()

	// init hermes email
	email := hermes.Email{
		Body: hermes.Body{
			Name: user.FullName,
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
						Link: uc.cfg.URL(entity.WebLoginPath),
					},
				},
			},
		},
	}

	// generate html email
	resHtml, err := h.GenerateHTML(email)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.GeneratingEmailHTMLError, err.Error()),
		}
	}

	// generate text html
	resText, err := h.GeneratePlainText(email)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.GeneratingEmailPlainTextError, err.Error()),
		}
	}

	return uc.send(user.Email, welcomeSubject, resText, resHtml)
}
