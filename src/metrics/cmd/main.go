package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"metrics"
	"time"
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

	errMetricData := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("YourService"),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String("ServiceErrors"),
				Dimensions: dimensions,
			},
		},
	}

	durationMetricData := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("YourService"),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String("ServiceExecutionTime"),
				Dimensions: dimensions,
				Unit:       aws.String("Milliseconds"),
			},
		},
	}

	errMetric := metrics.NewCustomMetric(errMetricData)
	errMetricCh, errMetricDoneCh, errMetricErrCh := metrics.CollectMetrics(errMetric, svc, 10)

	durationMetric := metrics.NewCustomMetric(durationMetricData)
	durationMetricCh, durationMetricDoneCh, durationMetricErrCh := metrics.CollectMetrics(durationMetric, svc, 10)

	go func() {
		for {
			select {
			case err := <-errMetricErrCh:
				fmt.Println("ErrMetric failed: ", err)
			case err := <-durationMetricErrCh:
				fmt.Println("DurationMetric failed: ", err)
			}
		}
	}()

	errMetricCh <- 1.0
	errMetricCh <- 2.0
	errMetricCh <- 3.0
	durationMetricCh <- 1024.0
	durationMetricCh <- 1211.0
	durationMetricCh <- 1132.0

	fmt.Println("Sleep")
	time.Sleep(1 * time.Second)
	fmt.Println("Done")

	ackCh := make(chan struct{}, 1)
	durationMetricDoneCh <- ackCh
	<-ackCh
	errMetricDoneCh <- ackCh
	<-ackCh
}
