package main

import (
	"flag"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"

	"k8s.io/klog/v2"

	basecmd "github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/cmd"

	buildkiteProvider "github.com/flarebuild/buildkite-agent-metrics-adapter/pkg/provider"
)

type adapter struct {
    basecmd.AdapterBase

    Token string
    Interval time.Duration
    Endpoint string
}

func (cmd *adapter) addFlags() {
    cmd.Flags().StringVar(&cmd.Token, "buildkite-agent-token", cmd.Token, "A Buildkite Agent Registration Token")
    cmd.Flags().StringVar(&cmd.Endpoint, "buildkite-agent-api-endpoint", "https://agent.buildkite.com/v3", "A custom Buildkite Agent API endpoint")
    cmd.Flags().DurationVar(&cmd.Interval, "buildkite-agent-metrics-update-interval", time.Second*10, "Update metrics every interval")
}

func main() {
    logs.InitLogs()
    defer logs.FlushLogs()

    // initialize the flags, with one custom flag for the message
    cmd := &adapter{}
    cmd.addFlags()
    cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure you get the klog flags
    cmd.Flags().Parse(os.Args)

    if cmd.Token == "" {
        // try to get from env
        cmd.Token = os.Getenv("BUILDKITE_AGENT_TOKEN")
        if cmd.Token == "" {
            klog.Fatalf("buildkite agent token is not set, use env BUILDKITE_AGENT_TOKEN or --buildkite-agent-token")
        }
    }

    cmd.WithExternalMetrics(buildkiteProvider.NewProvider(cmd.Token, cmd.Endpoint, cmd.Interval))

    if err := cmd.Run(wait.NeverStop); err != nil {
        klog.Fatalf("unable to run buildkite metrics adapter: %v", err)
    }
}
