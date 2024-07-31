package sender

import (
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

const DefaultEmailMessage = "current rate is"

type EmailSender struct {
	from     string
	smtpHost string
	smtpPort string
	password string
	logger   *zap.Logger
}

func NewEmailSender(from string, password string, logger *zap.Logger) *EmailSender {
	return &EmailSender{
		from:     from,
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
		password: password,
		logger:   logger.With(zap.String("service", "EmailSender")),
	}
}

func (s EmailSender) SendEmail(to string, message string) error {
	logger := s.logger.With(zap.String("method", "SendEmail"))

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetBody("text/plain", message)
	d := gomail.NewDialer("smtp.gmail.com", 587, s.from, s.password)
	err := d.DialAndSend(m)
	if err != nil {
		logger.Error("can't send email", zap.Error(err))
		return err
	}

	return nil
}
