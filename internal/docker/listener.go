package docker

import (
	"context"
	"gocd/internal/labels"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

func StartListener() {
	// Listen for new containers
	// When a new container is created, call the function below
	events, _ := cli.Events(context.Background(), types.EventsOptions{})

	for event := range events {
		if event.Action == "start" {
			logrus.Debugf("%s started", event.ID)

			logrus.Debugf("Container started: %s", event.ID)

			container, err := cli.ContainerInspect(context.Background(), event.ID)
			if err != nil {
				logrus.Error(err)
				continue
			}

			// log labels
			labels := labels.MapToGoCDLabels(container.Config.Labels)
			logrus.Debugf("Labels: %+v", labels)

			AddContainer(&Container{
				ID:      event.ID,
				Labels:  labels,
				Inspect: container,
			})
		}
		// When container is deleted, remove it from the list
		if event.Action == "destroy" {
			logrus.Debugf("Container destroyed: %s", event.ID)
			RemoveContainer(event.ID)
		}
	}
}
