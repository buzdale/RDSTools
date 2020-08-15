package main

// Commandline to describe or list RDS instances
// AWS_REGION shoudl be set
import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	var holdall string
	nflag := flag.String("d", "all", "Which database to describe (Default:All)")
	listflag := flag.Bool("l", false, "List databases only")
	flag.Parse()
	holdall = *nflag
	if holdall == "all" {
		holdall = ""
	}

	if *listflag == true {
		listDBs()
	} else {
		fmt.Println(describeDB(holdall))
	}
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
