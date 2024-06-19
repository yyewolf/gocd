package docker

import (
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
}
