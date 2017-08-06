package custommetrics

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"testing"
)

type mockClient struct {
	cloudwatchiface.CloudWatchAPI
	Resp cloudwatch.PutMetricDataOutput
}

func (m mockClient) PutMetricData(inp *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	// Only need to return mocked response output
	resp := &cloudwatch.PutMetricDataOutput{
		Namespace: inp.Namespace,
	}
	return resp, nil
}

func TestCollectMetrics(t *testing.T) {
	svc := &mockClient{}
	metric := NewCustomMetric("Test Namespace", "Test Metric", nil)
	resp, _ := metric.emitMetric(svc, 1.0)
	t.Error(resp)
}
