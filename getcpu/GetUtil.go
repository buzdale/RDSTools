// Package getcpu takes a string and prints CPU Utilization - returns "Done"
package getcpu

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Utilization takes instance to get CPU Utilization and prints it then returns "Done"
func Utilization(instance string) string {
	// statistics needs to be a slice of string - We only need one entry though.
	type cpuPercentage struct {
		Average     float64
		Timestamp   time.Time
		Unit        string
		Description string
	}
	//var cpuJson []cpuPercentage
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
	x0 := result.Datapoints[0].Average
	// fmt.Println((aws.Float64(*x0)))
	cpuOoutput := (aws.Float64(*x0))
	// json.Unmarshal([]float64(cpuOoutput), &cpuJson)
	fmt.Println(cpuOoutput)
	return "done"
}
