// This package is an exmaple of how to use the custom metric package.
// The custom metric package allows the user to emit metrics to Cloudwatch
// over a channel. This simplifies the process of emitting metrics by
// hiding the complexity of network connections and credentials.
// All the user needs to do is pass the channel to a function, and that
// function can emit values to Cloudwatch.
package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"metrics"
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

	// Create our metric collection pools using the session created above.
	errMetricErrCh := make(chan error, 1)
	errMetric := metrics.NewCustomMetric(errMetricData)
	errMetricCh, errMetricDoneCh := metrics.CollectMetrics(errMetric, svc, errMetricErrCh, 10)

	durationMetricErrCh := make(chan error, 1)
	durationMetric := metrics.NewCustomMetric(durationMetricData)
	durationMetricCh, durationMetricDoneCh := metrics.CollectMetrics(durationMetric, svc, durationMetricErrCh, 10)

	// Handle errors in the background.
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

	// Metrics can be emitted asynchronously. A collection pool handles the network
	// operations. The metric channel is easy to mock in unit tests. All the user
	// needs to do is write a value to a channel.
	errMetricCh <- 1.0
	errMetricCh <- 2.0
	errMetricCh <- 3.0
	durationMetricCh <- 1024.0
	durationMetricCh <- 1211.0
	durationMetricCh <- 1132.0

	close(errMetricCh)
	close(durationMetricCh)

	// Shutdown our metric collection pools.
	ackCh := make(chan struct{}, 1)
	durationMetricDoneCh <- ackCh
	<-ackCh
	errMetricDoneCh <- ackCh
	<-ackCh
}
