[![Go Reference](https://pkg.go.dev/badge/github.com/flarebuild/buildkite-agent-metrics-adapter.svg)](https://pkg.go.dev/github.com/flarebuild/buildkite-agent-metrics-adapter)

# buildkite-agent-metrics-adapter

This application implements Kubernetes External Metrics API for Buildkite agent metrics.

## Deployment

* Build docker image
* Deploy to kubernetes

```bash
# docker build -t buildkite-agent-metrics-adapter:latest -f deployDockerfile .
# kubectl apply -f deploy/adapter.yaml
```

## Metrics

There are total counters prefixed with `total_` and queue counters prefixed with `queue_`.

Queue counters should be filtered by metric label `queue` (see [usage examples](#usage-examples) below)

Metric names are taken from Buildkite's [collector package](https://pkg.go.dev/github.com/buildkite/buildkite-agent-metrics@v1.6.1-0.20200922085734-c8145a178990/collector#pkg-constants) and converted to `snake_case`.

F.e. counter for total number of running jobs will be `total_` + `ToSnakeCase(RunningJobsCount)` = `total_running_jobs_count`

## Usage examples

* Autoscaling based on total number of waiting jobs:

    ```yaml
    apiVersion: autoscaling/v2beta2
    kind: HorizontalPodAutoscaler
    metadata:
      name: buildkite-autoscaler
    spec:
      scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: buildkite
      minReplicas: 1
      maxReplicas: 10
      metrics:
      - type: External
        external:
          metric:
            name: total_waiting_jobs_count
          target:
            type: AverageValue
            averageValue: 1m
    ```

* Autoscaling based on number of waiting jobs in specific queue:

    ```yaml
    apiVersion: autoscaling/v2beta2
    kind: HorizontalPodAutoscaler
    metadata:
      name: buildkite-autoscaler
    spec:
      scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: buildkite
      minReplicas: 1
      maxReplicas: 10
      metrics:
      - type: External
        external:
          metric:
            name: queue_waiting_jobs_count
            selector:
              matchLabels:
                queue: default
          target:
            type: AverageValue
            averageValue: 1m
    ```
