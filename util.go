package bdikaa

import (
	"log"
	"runtime"

	"github.com/fsouza/go-dockerclient"
)

const RETRY = 7

func GetClinet() (*docker.Client, error) {

	if runtime.GOOS == "linux" {
		endpoint := "unix:///var/run/docker.sock"
		return docker.NewClient(endpoint)
	}
	return docker.NewClientFromEnv()
}

func RemoveContiner(client *docker.Client, cid string) error {
	err := client.StopContainer(cid, 5)
	if err != nil {
		log.Println("err %s", err.Error())
	}

	return client.RemoveContainer(docker.RemoveContainerOptions{ID: cid})
}
