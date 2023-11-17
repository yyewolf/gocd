package docker

import (
	"context"
	"errors"
	"fmt"
	"gocd/internal/labels"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
		fmt.Println(list)
		defer mutex.Unlock()
		listCopy := make([]*Container, len(list))
		copy(listCopy, list)

		for _, c := range list {
			logrus.Infof("Updating container %s", c.Inspect.Name)

			// Get current image ID from container
			i, _, err := cli.ImageInspectWithRaw(context.Background(), c.Inspect.Image)
			if err != nil {
				logrus.Error("Failed to update container at image inspect: ", err)
				continue
			}

			// Pull image, delete container, create new container
			out, err := cli.ImagePull(context.Background(), c.Inspect.Config.Image, types.ImagePullOptions{})
			if err != nil {
				logrus.Error("Failed to update container at image pull: ", err)
				continue
			}

			var buf [1024]byte
			io.ReadFull(out, buf[:])
			out.Close()

			outStr := string(buf[:])

			err = cli.ContainerStop(context.Background(), c.ID, container.StopOptions{})
			if err != nil {
				logrus.Error("Failed to update container at container stop: ", err)
			}

			err = cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{})
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

			resp, err := cli.ContainerCreate(context.Background(), c.Inspect.Config, c.Inspect.HostConfig, nil, nil, c.Inspect.Name)
			if err != nil {
				logrus.Error("Failed to update container at container create: ", err)
				continue
			}

			// Put networking config into new container
			for k, v := range nw {
				err = cli.NetworkConnect(context.Background(), k, resp.ID, v)
				if err != nil {
					logrus.Error("Failed to update container at network connect: ", err)
				}
			}

			err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
			if err != nil {
				logrus.Error("Failed to update container at container start: ", err)
				continue
			}

			if !strings.Contains(outStr, "Image is up to date") {
				// Remove old image
				_, err = cli.ImageRemove(context.Background(), i.ID, types.ImageRemoveOptions{})
				if err != nil {
					logrus.Error("Failed to remove old image: ", err)
				}
			}

		}
		containers[token] = listCopy
	}()

	return nil
}
