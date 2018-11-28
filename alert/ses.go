package alert

import (
	"fmt"

	"github.com/spf13/viper"
	ses "github.com/srajelli/ses-go"
)

func SesAws(to string, body string) {
	from := viper.Get("app.smtp.user").(string)
	subject := viper.Get("app.smtp.subject").(string)
	awsKeyID := viper.Get("app.smtp.ses.aws-key-id").(string)
	awsSecretKey := viper.Get("app.smtp.ses.aws-secret-key").(string)
	awsRegion := viper.Get("app.smtp.ses.aws-region").(string)

	ses.SetConfiguration(awsKeyID, awsSecretKey, awsRegion)

	emailData := ses.Email{
		To:      to,
		From:    from,
		Text:    body,
		Subject: subject,
		ReplyTo: from,
	}

	resp := ses.SendEmail(emailData)

	fmt.Println(resp)
}
