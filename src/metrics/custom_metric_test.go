package metrics

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"testing"
)

type mockClient struct {
	cloudwatchiface.CloudWatchAPI
	Resp cloudwatch.PutMetricDataOutput
}

func (m mockClient) PutMetricData(inp *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	// Only need to return mocked response output
	resp := &cloudwatch.PutMetricDataOutput{}
	return resp, nil
}

func TestCollectMetrics(t *testing.T) {
	svc := &mockClient{}
	errMetricData := &cloudwatch.PutMetricDataInput{
		Namespace: aws.String("YourService"),
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String("ServiceErrors"),
			},
		},
	}
	metric := NewCustomMetric(errMetricData)
	_, _ = metric.emitMetric(svc, 1.0)
}
