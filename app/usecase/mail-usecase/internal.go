package mailusecase

import (
	"fmt"
	"net/mail"
	"rap-c/config"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/matcornic/hermes/v2"
)

const (
	defaultEmailLogo string = "https://github.com/gendutski/rap-c/blob/main/storage/public-asset/images/logo-s.png?raw=true"
)

func (uc *usecase) initHermes() *hermes.Hermes {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      config.AppName,
			Link:      uc.cfg.AppURL,
			Logo:      defaultEmailLogo,
			Copyright: fmt.Sprintf("Copyright © %s %s. All rights reserved.", time.Now().Format("2006"), config.AppName),
		},
	}
	h.Theme = new(hermes.Default)
	return &h
}

func (uc *usecase) send(to, subject, txtBody, htmlBody string) error {
	from := mail.Address{
		Name:    uc.cfg.MailSenderName,
		Address: uc.cfg.MailSenderAddress,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(uc.cfg.MailHost, uc.cfg.MailPort, uc.cfg.MailUser, uc.cfg.MailPassword)
	return d.DialAndSend(m)
}
