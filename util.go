package bdikaa

import (
	"bytes"
	"log"
	"runtime"

	"github.com/fsouza/go-dockerclient"
)

//retry count for conecting continer DB
const RETRY = 5

//get the api clinet for linux or OSX.
func GetClinet() (*docker.Client, error) {

	if runtime.GOOS == "linux" {
		endpoint := "unix:///var/run/docker.sock"
		return docker.NewClient(endpoint)
	}
	return docker.NewClientFromEnv()
}

//remove continer by the continer ID.
func RemoveContiner(client *docker.Client, cid string) error {
	err := client.StopContainer(cid, 5)
	if err != nil {
		log.Println("err %s", err.Error())
	}

	return client.RemoveContainer(docker.RemoveContainerOptions{ID: cid})
}

//pull the correct image and tag.
func GetImageIfNotExsit(client *docker.Client, image string, tag string) error {

	var buf bytes.Buffer
	opts := docker.PullImageOptions{image, "base", tag, &buf, false}
	auth := docker.AuthConfiguration{}
	log.Printf("trying to pull image {%s}:{%s}", image, tag)
	return client.PullImage(opts, auth)
}
