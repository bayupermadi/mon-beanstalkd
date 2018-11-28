package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bayupermadi/mon-beanstalkd/aws"
)

func LogSize(path string, maxSize int64) (int64, error) {
	var size int64

	adjSize := func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += info.Size()

		return err
	}
	err := filepath.Walk(path, adjSize)

	sizeInMB := size / 1000 / 1000
	message := "Total log size beanstalkd: " + strconv.FormatInt(sizeInMB, 10) + "MB"

	fmt.Println(message)
	aws.CW("Logsize", "Megabytes", float64(sizeInMB), "LogDir", path)

	if sizeInMB > maxSize {
		alert(message)
	}

	return sizeInMB, err
}
