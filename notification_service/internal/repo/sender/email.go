package sender

import (
	gomail "gopkg.in/mail.v2"
)

type EmailSender struct {
	from     string
	smtpHost string
	smtpPort string
	password string
}

func NewEmailSender(from string, password string) *EmailSender {
	return &EmailSender{
		from:     from,
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
		password: password,
	}
}

func (s EmailSender) SendEmail(to string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetBody("text/plain", message)
	d := gomail.NewDialer("smtp.gmail.com", 587, s.from, s.password)
	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}
