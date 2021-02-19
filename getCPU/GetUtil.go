// Package getCPUUtilizationPKG takes a string and prints CPU Utilization - returns "Done"
package getCPU

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// func getCPUUtilization takes instance to get CPU Utilization and prints it then returns "Done"
func GetCPUUtilization(instance string) string {
	// statistics needs to be a slice of string - We only need one entry though.
	statistics := []string{
		"Average",
	}
	svc := cloudwatch.New(session.New())
	result, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{{Name: aws.String("DBInstanceIdentifier"), Value: aws.String(instance)}},
		EndTime:    aws.Time(time.Now().Add(time.Second * -300)),
		// ExtendedStatistics: aws.StringSlice(statistics),
		MetricName: aws.String("CPUUtilization"),
		Namespace:  aws.String("AWS/RDS"),
		Period:     aws.Int64(300),
		StartTime:  aws.Time(time.Now().Add(time.Second * -600)),
		Statistics: aws.StringSlice(statistics),
		// Unit:       aws.String("Seconds"),
	},
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	fmt.Printf("CPUUtilization of %s \n", instance)
	fmt.Println(result.Datapoints)

	return "done"
}
