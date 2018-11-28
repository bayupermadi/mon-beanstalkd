package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/viper"
)

func CW(metricname string, unit string, vmetric float64, dimName string, dimValue string) {

	awsKeyID := viper.Get("app.cloudwatch.aws-key-id").(string)
	awsSecretKey := viper.Get("app.cloudwatch.aws-secret-key").(string)
	awsRegion := viper.Get("app.cloudwatch.aws-region").(string)

	setConfiguration(awsKeyID, awsSecretKey, awsRegion)

	sess := startNewSession()

	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)

	_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("Beanstalkd"),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(metricname),
				Unit:       aws.String(unit),
				Value:      aws.Float64(vmetric),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String(dimName),
						Value: aws.String(dimValue),
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error adding metrics:", err.Error())
		return
	}

}
