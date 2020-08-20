package main

// Commandline to describe or list RDS instances
// AWS_REGION should be set
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

func main() {
	var holdall string
	nflag := flag.String("d", "all", "Which database to describe (Default:All)")
	listflag := flag.Bool("l", false, "List databases only")
	freeflag := flag.Bool("f", false, "List free storage per database instance")
	flag.Parse()
	holdall = *nflag
	if holdall == "all" {
		holdall = ""
	}
	// Print the list of DBs Maybe these should be case statements.
	if *listflag == true {
		listOfDBs := listDBs()
		for _, n := range listOfDBs {
			fmt.Println(n)
		}
	} else if *freeflag {
		listforStorage := listDBs()
		for _, n := range listforStorage {
			gbytesStorage := *getFreeStorage(n).Datapoints[0].Average / 1024 / 1024 / 1024
			fmt.Println(" Free storage for database ", n, " is ", gbytesStorage, " GBytes")
		}
	} else {
		fmt.Println(describeDB(holdall))
	}
}

func getFreeStorage(instance string) *cloudwatch.GetMetricStatisticsOutput {
	// statistics needs to be a slice of string - We only need one entry though.
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
}

//  Grabs results from describeDB and returns databases.
func listDBs() []string {
	result := describeDB("")
	var list []string
	for _, n := range result.DBInstances {
		// Aurora databases don't give FreeStorage which breaks looking for that metric.
		if *n.Engine != "aurora" {
			list = append(list, *n.DBInstanceIdentifier)
		}
	}
	return list
}
