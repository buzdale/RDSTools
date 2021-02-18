// test to run with known instance

package getCPUUtilization

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
)

func getCPUUtilization(instance string) string {
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBInstanceNotFoundFault:
				fmt.Println(rds.ErrCodeDBInstanceNotFoundFault, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println("Error", err.Error())
		}
	}
	fmt.Printf("CPUUtilization of %s \n", instance)
	fmt.Println(result.Datapoints)

	return "done"
}