package monitor

import (
	"fmt"
	"log"

	"github.com/bayupermadi/mon-beanstalkd/aws"
	"github.com/go-yaml/yaml"
	"github.com/iwanbk/gobeanstalk"
	"github.com/spf13/viper"
)

func StatsTube(tubes string) {
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
		aws.CW("Tubes Name", "Count", float64(currentBuried.(int)), "Tubes", tubes)
		if currentBuried.(int) > thresholdJobs {
			message := tubes + " has buried jobs: " + string(currentBuried.(int))
			alert(message)
			fmt.Println(message)
		}
	}

}
