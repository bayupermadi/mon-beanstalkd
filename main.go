package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/iwanbk/gobeanstalk"
	"github.com/spf13/viper"
	ses "github.com/srajelli/ses-go"
)

func statsTube(tubes string) {
	host := viper.Get("app.host").(string)
	conn, err := gobeanstalk.Dial(host)
	if err != nil {
		log.Fatal(err)
	}
	tubeMap := make(map[string]interface{})
	burried, err := conn.StatsTube(tubes)
	if err := yaml.Unmarshal(burried, &tubeMap); err != nil {

	}
	thresholdJobs := viper.Get("app.max-buried-job").(int)
	currentBuried, ok := tubeMap["current-jobs-buried"]
	if ok {
		fmt.Println(tubes+" has buried jobs: ", currentBuried.(int))
		if currentBuried.(int) > thresholdJobs {
			message := tubes + " has buried jobs: " + string(currentBuried.(int))
			alert(message)
			fmt.Println(message)
		}
	}

}

func sendEmail(body string) {
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

func sesAws(to string, body string) {
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

func alert(msg string) {
	if viper.GetBool("app.smtp.enabled") == true {
		if viper.GetBool("app.smtp.ses.enabled") == true {
			to := viper.Get("app.smtp.recipient").(string)
			dest := strings.Split(to, ", ")
			start := 0
			for i := 0; i < len(dest); i++ {
				start += i
				sesAws(dest[start], msg)
			}
		} else {
			sendEmail(msg)
		}
	}
}

func logSize(path string, maxSize int64) (int64, error) {
	var size int64
	adjSize := func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	}
	err := filepath.Walk(path, adjSize)

	sizeInMB := size / 1024 / 1024

	if sizeInMB > maxSize {
		message := "Total log size beanstalkd: " + strconv.FormatInt(sizeInMB, 10) + "MB"
		alert(message)
		fmt.Println(message)
	}

	return sizeInMB, err
}

func main() {
	// config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	logPath := viper.Get("app.log.dir").(string)
	logMaxSize := viper.Get("app.log.max-size").(int)

	for {
		tubes := viper.Get("app.tube").(string)
		tubeList := strings.Split(tubes, ", ")
		start := 0
		for i := 0; i < len(tubeList); i++ {
			start = i
			tubeName := string(tubeList[start])
			statsTube(tubeName)
		}

		logSize(logPath, int64(logMaxSize))

		<-time.After(time.Second * 30)
	}

}
