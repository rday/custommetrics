package metrics

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"testing"
)

type mockClient struct {
	cloudwatchiface.CloudWatchAPI
	Resp       cloudwatch.PutMetricDataOutput
	ValueCount float64
}

func (m *mockClient) PutMetricData(inp *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	// Only need to return mocked response output
	m.ValueCount += *inp.MetricData[0].Value
	resp := &cloudwatch.PutMetricDataOutput{}
	return resp, nil
}

func TestCollectMetrics(t *testing.T) {
	svc := &mockClient{}
	testMetricData := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("YourService"),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String("ServiceErrors"),
			},
		},
	}

	testMetric := NewCustomMetric(testMetricData)
	testMetricCh, doneCh := CollectMetrics(testMetric, svc, nil, 1)
	ackCh := make(chan struct{}, 1)

	testMetricCh <- 1.0
	testMetricCh <- 2.0
	testMetricCh <- 3.0
	testMetricCh <- 4.0
	testMetricCh <- 5.0
	testMetricCh <- 6.0
	testMetricCh <- 7.0
	testMetricCh <- 8.0
	testMetricCh <- 9.0
	testMetricCh <- 10.0
	doneCh <- ackCh
	close(testMetricCh)

	<-ackCh

	if svc.ValueCount != 55.0 {
		t.Error("Expected 55.0, found %v", svc.ValueCount)
	}
}
