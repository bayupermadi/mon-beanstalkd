package alert

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/spf13/viper"
)

func SendEmail(body string) {
	from := viper.Get("app.smtp.user").(string)
	pass := viper.Get("app.smtp.password").(string)
	port := viper.Get("app.smtp.port").(string)
	server := viper.Get("app.smtp.server").(string)
	to := viper.Get("app.smtp.recipient").(string)
	subject := viper.Get("app.smtp.subject").(string)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(server+":"+port,
		smtp.PlainAuth("", from, pass, server),
		from, strings.Split(to, ", "), []byte(msg))

	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return
	}

}
