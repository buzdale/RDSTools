package main

// Commandline to describe or list RDS instances
// AWS_REGION shoudl be set
import (
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
)

// type metricResponse struct {
// 	Average int64  `json: "Average"`
// 	Unit    string `json: "Unit"`
// }

func main() {
	var holdall string
	nflag := flag.String("d", "all", "Which database to describe (Default:All)")
	listflag := flag.Bool("l", false, "List databases only")
	freeflag := flag.String("f", "", "List free storage per database instance (Default ignite0000")
	flag.Parse()
	holdall = *nflag
	instance := *freeflag
	if holdall == "all" {
		holdall = ""
	}

	if *listflag == true {
		listDBs()
	} else if instance != "" {
		gbytesStorage := *getFreeStorage(instance).Datapoints[0].Average / 1024 / 1024 / 1024
		fmt.Println(" Free storage for database ", instance, " is ", gbytesStorage, " GBytes")
	} else {
		fmt.Println(describeDB(holdall))
	}
}

func getFreeStorage(instance string) *cloudwatch.GetMetricStatisticsOutput {
	// statistics needs to be a slice of string - We only need on entry.
	statistics := make([]string, 1)
	statistics[0] = "Average"
	svc := cloudwatch.New(session.New())
	result, err := svc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		MetricName: aws.String("FreeStorageSpace"),
		Namespace:  aws.String("AWS/RDS"),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("DBInstanceIdentifier"),
				Value: aws.String(instance),
			},
		},
		EndTime:    aws.Time(time.Now()),
		StartTime:  aws.Time(time.Now().Add(time.Minute * -10)),
		Period:     aws.Int64(300),
		Statistics: aws.StringSlice(statistics),
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

	return result

}

// Describes one or all databases in Json format
func describeDB(instance string) *rds.DescribeDBInstancesOutput {
	svc := rds.New(session.New())
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instance),
	}

	result, err := svc.DescribeDBInstances(input)
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
			fmt.Println(err.Error())
		}

	}
	return result
	//fmt.Println(result)

}

//  Grabs results from describeDB and shows databases.
func listDBs() {
	result := describeDB("")
	for _, n := range result.DBInstances {
		fmt.Println(*n.DBInstanceIdentifier)
	}
}
