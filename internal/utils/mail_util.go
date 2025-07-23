package utils

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"strings"
)

// SendEmail send email
func SendEmail(to []string, cc []string, subject string, body string) error {
	// use gomail to send email
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("email.user"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	for _, s := range cc {
		c := strings.Split(s, ",")
		m.SetAddressHeader("Cc", c[0], c[1])
	}

	d := gomail.NewDialer(viper.GetString("email.host"), viper.GetInt("email.port"), viper.GetString("email.user"), viper.GetString("email.password"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
