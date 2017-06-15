package metrics

import (
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
	metric := NewCustomMetric("Test Namespace", "Test Metric", nil)
	_, _ = metric.emitMetric(svc, 1.0)
}
