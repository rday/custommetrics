// The metricis package allows the user to emit metrics to Cloudwatch
// over a channel. This simplifies the process of emitting metrics by
// hiding the complexity of network connections and credentials.
// All the user needs to do is pass the channel to a function, and that
// function can emit values to Cloudwatch.
package metrics

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
)

// CustomMetric is composed of a single value for now, but could
// be easily expanded.
type CustomMetric struct {
	metric *cloudwatch.PutMetricDataInput
}

func (c *CustomMetric) emitMetric(svc cloudwatchiface.CloudWatchAPI, value float64) (*cloudwatch.PutMetricDataOutput, error) {
	c.metric.MetricData[0].Value = &value
	return svc.PutMetricData(c.metric)
}

// CollectMetrics starts the collector for a particular metric on a particular Cloudwatch service. This
// collector will emit the metric with the value written to its value channel. The collector can be
// cleanly shutdown with a write to its done channel.
// If an error channel is passed, errors will be written to that channel. The caller is responsible for
// reading this channel, otherwise the collector WILL hang.
func CollectMetrics(metric *CustomMetric, svc cloudwatchiface.CloudWatchAPI, errCh chan error, bufferLen int) (chan float64, chan chan struct{}) {
	valueCh := make(chan float64, bufferLen)
	doneCh := make(chan chan struct{}, 1)

	go func() {
		var ackCh chan struct{}
		done := false

		for {
			select {
			case ackCh = <-doneCh:
				done = true
			case value, ok := <-valueCh:
				// If the caller has written to the done channel, and the value channel has been closed,
				// then acknowlege we are done and exit.
				if !ok && done {
					ackCh <- struct{}{}
					return
				}

				_, err := metric.emitMetric(svc, value)
				if err != nil && errCh != nil {
					errCh <- err
				}
			}
		}
	}()

	return valueCh, doneCh
}

func NewCustomMetric(metric *cloudwatch.PutMetricDataInput) *CustomMetric {
	return &CustomMetric{
		metric: metric,
	}
}
