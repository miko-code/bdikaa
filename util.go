package bdikaa

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/fsouza/go-dockerclient"
)

//RETRY count for conecting continer DB
const RETRY = 5

//GetClinet get the api clinet for linux or OSX.
func GetClinet() (*docker.Client, error) {

	if runtime.GOOS == "linux" {
		endpoint := "unix:///var/run/docker.sock"
		return docker.NewClient(endpoint)
	}
	return docker.NewClientFromEnv()
}

//RemoveContiner by the continer ID.
func RemoveContinerID(client *docker.Client, cid string) error {
	err := client.StopContainer(cid, 5)
	if err != nil {
		log.Println("err %s", err.Error())
	}

	return client.RemoveContainer(docker.RemoveContainerOptions{ID: cid})
}

//GetImageIfNotExsit pull the correct image and tag.
func GetImageIfNotExsit(client *docker.Client, image string, tag string) error {

	//	var buf bytes.Buffer
	opts := docker.PullImageOptions{
		Repository:        image,
		Tag:               tag,
		OutputStream:      os.Stdout,
		InactivityTimeout: time.Duration(time.Minute * 1),
	}
	//	opts := docker.PullImageOptions{image, tag, "", &buf, false, time.Duration(time.Minute * 5)}
	auth := docker.AuthConfiguration{}
	log.Printf("pulling image if nedded {%s}:{%s}", image, tag)
	return client.PullImage(opts, auth)
}
