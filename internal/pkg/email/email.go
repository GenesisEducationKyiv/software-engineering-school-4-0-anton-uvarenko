package email

import (
	gomail "gopkg.in/mail.v2"
)

type EmailSender struct {
	from     string
	smtpHost string
	smtpPort string
	password string
	dialer   dialer
}

type dialer interface {
	DialAndSend(...*gomail.Message) error
}

func NewEmailSender(from string, password string) *EmailSender {
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	return &EmailSender{
		from:     from,
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
		password: password,
		dialer:   d,
	}
}

func (s *EmailSender) SendEmail(to string, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetBody("text/plain", message)
	err := s.dialer.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}
