package docker

import (
	"context"
	"gocd/internal/labels"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

var cli *client.Client

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Docker client initialized")

	// Fill up containers
	c, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		logrus.Fatal(err)
	}

	for _, container := range c {
		logrus.Infof("Container found: %s", container.ID)
		labels := labels.MapToGoCDLabels(container.Labels)

		inspect, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			logrus.Error(err)
		}

		AddContainer(&Container{
			ID:      container.ID,
			Labels:  labels,
			Inspect: inspect,
		})
	}
}
