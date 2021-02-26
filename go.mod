module github.com/flarebuild/buildkite-agent-metrics-adapter

go 1.15

require (
	github.com/buildkite/buildkite-agent-metrics v1.6.1-0.20200922085734-c8145a178990
	github.com/kubernetes-sigs/custom-metrics-apiserver v0.0.0-20210223164403-718972b7a2f4
	github.com/stretchr/testify v1.6.1
	k8s.io/apimachinery v0.20.4
	k8s.io/component-base v0.20.0
	k8s.io/klog/v2 v2.4.0
	k8s.io/metrics v0.20.4
)
