package main

// Commandline to describe or list RDS instances
// AWS_REGION should be set
import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
	getcpu "github.com/buzdale/RDSTools/getCPU"
)

func main() {
	var holdall string
	nflag := flag.String("d", "all", "Which database to describe (Default:All)")
	listflag := flag.Bool("l", false, "List databases only")
	freeflag := flag.Bool("f", false, "List free storage per database instance")
	cpuflag := flag.Bool("c", false, "Calculate average CPU usage last 5 minutes")
	flag.Parse()
	holdall = *nflag
	if holdall == "all" {
		holdall = ""
	}
	// Print the list of DBs Maybe these should be case statements.
	if *listflag == true {
		listOfDBs := listDBs(true)
		for _, n := range listOfDBs {
			fmt.Println(n)
		}
	} else if *freeflag {
		listforStorage := listDBs(false)
		for _, n := range listforStorage {
			gbytesStorage := *getFreeStorage(n).Datapoints[0].Average / 1024 / 1024 / 1024
			gbytesAllocatedStorage := getAllocatedStorage(n)
			fmt.Printf("Free storage for database %s:\t\t%4.2fGB of %s GB \n", n, gbytesStorage, gbytesAllocatedStorage)
		}
	} else if *cpuflag {
		cpuAverage := listDBs(true)
		for _, n := range cpuAverage {
			cpuPercentage := getcpu.Utilization(n)
			fmt.Println(cpuPercentage)
			// cpustring := cpuPercentage.GoString()
			// fmt.Printf(" CPU utilization for database %s is %s %% \n", n, cpustring)
		}

	} else {
		fmt.Println(describeDB(holdall))
	}
}

// Gets the average free storage for the last 5 minutes.
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
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return result

}

// Gets the allocated storage for the database instance.
func getAllocatedStorage(instance string) string {
	svc := rds.New(session.New())
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instance),
	}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	n := result.DBInstances[0]
	AllocatedStorage := *n.AllocatedStorage

	return strconv.FormatInt(AllocatedStorage, 10)
}

// Describes one or all databases in Json format
func describeDB(instance string) *rds.DescribeDBInstancesOutput {
	svc := rds.New(session.New())
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(instance),
	}

	result, err := svc.DescribeDBInstances(input)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return result
}

//  Grabs results from describeDB and returns databases.
func listDBs(listFlag bool) []string {
	result := describeDB("")
	var list []string
	for _, n := range result.DBInstances {
		if listFlag == true {
			temp1 := *n.DBInstanceIdentifier
			temp2 := *n.DBInstanceClass
			temp := temp1 + " " + temp2
			list = append(list, temp) //*n.DBInstanceIdentifier)
		} else {
			// // Aurora databases don't give FreeStorage  neither do stopped databases - which breaks looking for that metric.
			if *n.Engine != "aurora" && *n.DBInstanceStatus != "stopped" {
				list = append(list, *n.DBInstanceIdentifier)
			}
		}
	}
	return list
}
