package bdikaa

import "github.com/fsouza/go-dockerclient"

//Continer interface
type Continer interface {
	CreatDockerConfig() *docker.Config
	CreatDockerHostConfig() *docker.HostConfig
	CreateContiner(c *docker.Client) (interface{}, string, error)
	RemoveContiner(c *docker.Client, cid string) error
}
