package metrics

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

type CustomMetric struct {
	metric *cloudwatch.PutMetricDataInput
}

func (c *CustomMetric) emitMetric(svc cloudwatchiface.CloudWatchAPI, value float64) (*cloudwatch.PutMetricDataOutput, error) {
	c.metric.MetricData[0].Value = &value
	return svc.PutMetricData(c.metric)
}

func CollectMetrics(metric *CustomMetric, svc *cloudwatch.CloudWatch, bufferLen int) (chan float64, chan chan struct{}, chan error) {
	ch := make(chan float64, bufferLen)
	doneCh := make(chan chan struct{}, 1)
	errCh := make(chan error, 1)

	go func() {
		for {
			select {
			case ackCh := <-doneCh:
				ackCh <- struct{}{}
				return
			case value := <-ch:
				_, err := metric.emitMetric(svc, value)
				if err != nil {
					errCh <- err
				}
			}
		}
	}()

	return ch, doneCh, errCh
}

func NewCustomMetric(metric *cloudwatch.PutMetricDataInput) *CustomMetric {
	return &CustomMetric{
		metric: metric,
	}
}
