package main

import (
	"fmt"
	"strings"

	"github.com/bayupermadi/mon-beanstalkd/monitor"
	"github.com/spf13/viper"
)

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

	tubes := viper.Get("app.tube").(string)
	tubeList := strings.Split(tubes, ", ")
	start := 0
	for i := 0; i < len(tubeList); i++ {
		start = i
		tubeName := string(tubeList[start])
		monitor.StatsTube(tubeName)
	}

	monitor.LogSize(logPath, int64(logMaxSize))

}
