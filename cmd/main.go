package main

import (
	"custommetrics"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"fmt"
)

func main() {
	sess := session.Must(session.NewSession())
	svc := cloudwatch.New(sess)
	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("DimensionName"),
			Value: aws.String("DimensionValue"),
		},
	}

	metric := custommetrics.NewCustomMetric("Namespace", "Metric", dimensions)
	ch, doneCh := custommetrics.CollectMetrics(metric, svc, 10)
	ch <- 1.0
	ch<- 2.0
	ch<-3.0
	fmt.Println("Sleep")
	time.Sleep(1 * time.Second)
	fmt.Println("Done")

	ackCh := make(chan struct{}, 1)
	doneCh<-ackCh
	<-ackCh

}
