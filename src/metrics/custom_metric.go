package metrics

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type CustomMetric struct {
	namespace  string
	metricName string
	dimensions []*cloudwatch.Dimension
}

func (c *CustomMetric) buildInput(value float64) *cloudwatch.PutMetricDataInput {
	return &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(c.namespace),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String(c.metricName),
				Dimensions: c.dimensions,
				Value:      &value,
			},
		},
	}
}

func (c *CustomMetric) emitMetric(svc cloudwatchiface.CloudWatchAPI, value float64) (*cloudwatch.PutMetricDataOutput, error) {
	inp := c.buildInput(value)
	return svc.PutMetricData(inp)
}

func CollectMetrics(metric *CustomMetric, svc *cloudwatch.CloudWatch, bufferLen int) (chan float64, chan chan struct{}) {
	ch := make(chan float64, bufferLen)
	doneCh := make(chan chan struct{}, 1)

	go func() {
		for {
			select {
			case ackCh := <-doneCh:
				ackCh <- struct{}{}
				return
			case value := <-ch:
				_, err := metric.emitMetric(svc, value)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	return ch, doneCh
}

func NewCustomMetric(namespace, metricName string, dimensions []*cloudwatch.Dimension) *CustomMetric {
	return &CustomMetric{
		namespace:  namespace,
		metricName: metricName,
		dimensions: dimensions,
	}
}
