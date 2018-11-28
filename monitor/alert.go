package monitor

import (
	"strings"

	a "github.com/bayupermadi/mon-beanstalkd/alert"
	"github.com/spf13/viper"
)

func alert(msg string) {
	if viper.GetBool("app.smtp.enabled") == true {
		if viper.GetBool("app.smtp.ses.enabled") == true {
			to := viper.Get("app.smtp.recipient").(string)
			dest := strings.Split(to, ", ")
			start := 0
			for i := 0; i < len(dest); i++ {
				start += i
				a.SesAws(dest[start], msg)
			}
		} else {
			a.SendEmail(msg)
		}
	}
}
