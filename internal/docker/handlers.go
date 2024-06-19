package docker

import (
	"context"
	"fmt"
	"gocd/internal/discord"
	"gocd/internal/labels"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
)

type Container struct {
	ID      string
	Labels  labels.GoCDLabels
	Inspect types.ContainerJSON
	Error   error
}

func UpdateContainers(token string) error {
	logrus.Info("Docker client initialized")

	// Fill up containers
	c, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "gocd.token="+token)),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	var containers []*Container

	for _, container := range c {
		logrus.Infof("Container found: %s", container.ID)
		labels := labels.MapToGoCDLabels(container.Labels)

		inspect, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			logrus.Error(err)
		}

		container := &Container{
			ID:      container.ID,
			Labels:  labels,
			Inspect: inspect,
		}

		containers = append(containers, container)
	}

	go func() {
		list := containers

		// Prepare discord's message
		message := fmt.Sprintf("Updating %d container(s):\n", len(list))
		for _, c := range list {
			lbl := labels.MapToGoCDLabels(c.Inspect.Config.Labels)
			if lbl.Repo == "" {
				message += fmt.Sprintf("- **%s**", c.Inspect.Name)
			} else {
				message += fmt.Sprintf("- **[%s](%s)**", c.Inspect.Name, lbl.Repo)
			}
		}

		messageID, _ := discord.SendMessage(message)

		for _, c := range list {
			logrus.Infof("Updating container %s", c.Inspect.Name)

			// Pull image, delete container, create new container
			out, err := cli.ImagePull(context.Background(), c.Inspect.Config.Image, types.ImagePullOptions{})
			if err != nil {
				logrus.Error("Failed to update container at image pull: ", err)
				c.Error = err
			}

			if out == nil {
				out = io.NopCloser(strings.NewReader(""))
			}

			var buf []byte
			// read out till it's exhausted
			buf, _ = io.ReadAll(out)

			outStr := string(buf[:])

			logrus.Debug("Image pull output: ", outStr)

			logrus.Debug("Stopping container")
			err = cli.ContainerStop(context.Background(), c.ID, container.StopOptions{
				Signal: "SIGKILL",
			})
			if err != nil {
				logrus.Error("Failed to update container at container stop: ", err)
				c.Error = err
			}

			logrus.Debug("Removing container")
			err = cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{
				Force: true,
			})
			if err != nil {
				logrus.Error("Failed to update container at container remove: ", err)
				logrus.Error("Attempting to start container anyway")
				c.Error = err
			}
			// Get networking config from old container
			nw := c.Inspect.NetworkSettings.Networks

			logrus.Debug("Creating new container")
			// Create new container
			resp, err := cli.ContainerCreate(context.Background(), c.Inspect.Config, c.Inspect.HostConfig, nil, nil, c.Inspect.Name)
			if err != nil {
				logrus.Error("Failed to update container at container create: ", err)
				c.Error = err
				continue
			}

			logrus.Debug("Connecting new container to network")
			// Put networking config into new container
			for k, v := range nw {
				err = cli.NetworkConnect(context.Background(), k, resp.ID, v)
				if err != nil {
					logrus.Error("Failed to update container at network connect: ", err)
					c.Error = err
				}
			}

			logrus.Debug("Starting new container")
			err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
			if err != nil {
				logrus.Error("Failed to update container at container start: ", err)
				c.Error = err
				continue
			}

			if !strings.Contains(outStr, "Image is up to date") {
				logrus.Debug("Removing old image")
				// Remove old image
				_, err = cli.ImagesPrune(context.Background(), filters.NewArgs())
				if err != nil {
					logrus.Error("Failed to remove old image: ", err)
				}
			}
			logrus.Debug("Finished updating container")
		}

		// Prepare discord's message
		message = fmt.Sprintf("Updated %d container(s):\n", len(list))
		for _, c := range list {
			lbl := labels.MapToGoCDLabels(c.Inspect.Config.Labels)
			if lbl.Repo == "" {
				message += fmt.Sprintf("- **%s**", c.Inspect.Name)
			} else {
				message += fmt.Sprintf("- **[%s](%s)**", c.Inspect.Name, lbl.Repo)
			}

			// Error report
			if c.Error != nil {
				message += fmt.Sprintf(" (%v)\n", c.Error)
			} else {
				message += "\n"
			}
		}

		discord.UpdateMessage(messageID, message)
	}()

	return nil
}
