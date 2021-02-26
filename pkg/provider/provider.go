package provider

import (
	"time"

	resource "k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/klog/v2"

	bk "github.com/buildkite/buildkite-agent-metrics/collector"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/flarebuild/buildkite-agent-metrics-adapter/pkg/util"
)

// buildkiteProvider is implementation of provider.ExternalMetricsProvider
type buildkiteProvider struct {
	bk.Collector
	interval time.Duration
	result *bk.Result
	metricsInfo []provider.ExternalMetricInfo
	metrics map[string]external_metrics.ExternalMetricValueList
}

func (p *buildkiteProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	metrics := &external_metrics.ExternalMetricValueList{}
	//klog.Infof("metric: %s, ns: %s, selector: %s", info.Metric, namespace, metricSelector.String())
	for _, metric := range p.metrics[info.Metric].Items {
		if metricSelector.Matches(labels.Set(metric.MetricLabels)) {
			metrics.Items = append(metrics.Items, metric)
		}
	}
	return metrics, nil
}

func (p *buildkiteProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	return p.metricsInfo
}

func (p *buildkiteProvider) updateMetrics() (time.Duration, error) {
	t := time.Now()
	result, err := p.Collect()
	if err != nil {
		return time.Duration(0), err
	}

	metricsInfo := []provider.ExternalMetricInfo{}
	metrics := map[string]external_metrics.ExternalMetricValueList{}

	if p.result != nil {
		for n, i := range p.result.Totals {
			metricName := util.ToSnakeCase("Total"+n)
			metricsInfo = append(metricsInfo, provider.ExternalMetricInfo{Metric: metricName})
			metrics[metricName] = external_metrics.ExternalMetricValueList{
				Items: []external_metrics.ExternalMetricValue{{
					MetricName: metricName,
					Timestamp: meta.Time{Time: t},
					Value: *resource.NewQuantity(int64(i), resource.DecimalSI),
				}},
			}
		}
		for q := range p.result.Queues {
			metricLabels := map[string]string{"queue": q}
			for n, i := range p.result.Queues[q] {
				metricName := util.ToSnakeCase("Queue"+n)
				metricValue := external_metrics.ExternalMetricValue{
					MetricName: metricName,
					MetricLabels: metricLabels,
					Timestamp: meta.Time{Time: t},
					Value: *resource.NewQuantity(int64(i), resource.DecimalSI),
				}
				qMetrics, ok := metrics[metricName]
				if ok { // already saw this metric, but for another queue
					qMetrics.Items = append(qMetrics.Items, metricValue)
				} else { // see this metric for the first time
					metrics[metricName] = external_metrics.ExternalMetricValueList{
						Items: []external_metrics.ExternalMetricValue{metricValue},
					}
					metricsInfo = append(metricsInfo, provider.ExternalMetricInfo{Metric: metricName})
				}
			}
		}
	}

	p.result = result
	p.metricsInfo = metricsInfo
	p.metrics = metrics

	return result.PollDuration, nil
}

func (p *buildkiteProvider) updateMetricsLoop() {
	minPollDuration, err := p.updateMetrics()
	if (err != nil) {
		klog.Errorln(err)
	}
	for {
		waitTime := p.interval
		// Respect the min poll duration returned by the API
		if waitTime < minPollDuration {
			klog.Infof("Increasing poll duration based on rate-limit headers")
			waitTime = minPollDuration
		}
		klog.Infof("Waiting for %v (minimum of %v)", waitTime, minPollDuration)
		time.Sleep(waitTime)
		minPollDuration, err = p.updateMetrics()
		if err != nil {
			klog.Errorln(err)
		}
	}
}

// NewProvider creates external metrics provider
func NewProvider(token, endpoint string, interval time.Duration) provider.ExternalMetricsProvider {
	prov := &buildkiteProvider{
		Collector: bk.Collector{
			Token: token,
			Endpoint: endpoint,
			UserAgent: "buildkite-agent-metrics-adapter",
			Quiet: true,
		},
		interval: interval,
	}
	go prov.updateMetricsLoop()
	return prov
}