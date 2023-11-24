package docker

import (
	"context"
	"errors"
	"fmt"
	"gocd/internal/discord"
	"gocd/internal/labels"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
)

type Container struct {
	ID      string
	Labels  labels.GoCDLabels
	Inspect types.ContainerJSON
}

var mutex sync.Mutex
var containers = make(map[string][]*Container)

func AddContainer(container *Container) {
	mutex.Lock()
	list, ok := containers[container.Labels.Token]
	if !ok {
		list = make([]*Container, 0)
	}

	list = append(list, container)
	containers[container.Labels.Token] = list
	mutex.Unlock()
}

func RemoveContainer(id string) {
	mutex.Lock()
	for token, list := range containers {
		for i, c := range list {
			if c.ID == id {
				logrus.Infof("Removing container %s", id)
				list = append(list[:i], list[i+1:]...)
				containers[token] = list
			}
		}
	}
	mutex.Unlock()
}

func UpdateContainers(token string) error {
	_, ok := containers[token]
	if !ok {
		return errors.New("no containers for token")
	}

	go func() {
		mutex.Lock()
		list := containers[token]
		defer mutex.Unlock()
		listCopy := make([]*Container, len(list))
		copy(listCopy, list)

		// Prepare discord's message
		message := fmt.Sprintf("Updating %d container(s):\n", len(list))
		for _, c := range list {
			lbl := labels.MapToGoCDLabels(c.Inspect.Config.Labels)
			if lbl.Repo == "" {
				message += fmt.Sprintf("- **%s**\n", c.Inspect.Name)
			} else {
				message += fmt.Sprintf("- **[%s](%s)**\n", c.Inspect.Name, lbl.Repo)
			}
		}

		discord.SendMessage(message)

		for _, c := range list {
			logrus.Infof("Updating container %s", c.Inspect.Name)

			// Pull image, delete container, create new container
			out, err := cli.ImagePull(context.Background(), c.Inspect.Config.Image, types.ImagePullOptions{})
			if err != nil {
				logrus.Error("Failed to update container at image pull: ", err)
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
			}

			logrus.Debug("Removing container")
			err = cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{
				Force: true,
			})
			if err != nil {
				logrus.Error("Failed to update container at container remove: ", err)
				logrus.Error("Attempting to start container anyway")
			}
			// Get networking config from old container
			nw := c.Inspect.NetworkSettings.Networks

			// Delete current container from listCopy
			for i, c2 := range listCopy {
				if c.ID == c2.ID {
					listCopy = append(listCopy[:i], listCopy[i+1:]...)
				}
			}

			logrus.Debug("Creating new container")
			// Create new container
			resp, err := cli.ContainerCreate(context.Background(), c.Inspect.Config, c.Inspect.HostConfig, nil, nil, c.Inspect.Name)
			if err != nil {
				logrus.Error("Failed to update container at container create: ", err)
				continue
			}

			logrus.Debug("Connecting new container to network")
			// Put networking config into new container
			for k, v := range nw {
				err = cli.NetworkConnect(context.Background(), k, resp.ID, v)
				if err != nil {
					logrus.Error("Failed to update container at network connect: ", err)
				}
			}

			logrus.Debug("Starting new container")
			err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
			if err != nil {
				logrus.Error("Failed to update container at container start: ", err)
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
		containers[token] = listCopy
	}()

	return nil
}
