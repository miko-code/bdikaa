package bdikaa

import "github.com/fsouza/go-dockerclient"

//Continer interface
type Continer interface {
	CreatDockerConfig() *docker.Config
	CreatDockerHostConfig() *docker.HostConfig
	CreateContiner(c *docker.Client) (interface{}, string, error)
	ConectToStorage(c *docker.Client, cid string) (interface{}, error)
	RemoveContiner(c *docker.Client, cid string) error
}
