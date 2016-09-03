package bdikaa

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"

	"golang.org/x/net/context"

	"github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

type Mongo struct {
	Tag string
}

func (m *Mongo) CreatDockerConfig() *docker.Config {
	conf := &docker.Config{
		Image: fmt.Sprintf("mongo:%s", m.Tag),
	}
	return conf
}
func (m *Mongo) CreatDockerHostConfig() *docker.HostConfig {
	//add data dir
	conf := &docker.HostConfig{}
	return conf
}

func (m *Mongo) CreateContiner(client *docker.Client) (interface{}, string, error) {
	err := GetImageIfNotExsit(client, "mongo", m.Tag)
	if err != nil {
		log.Println("enable to create Continer ", err.Error())
		return nil, "", err
	}
	conf := m.CreatDockerConfig()
	hostConf := m.CreatDockerHostConfig()
	netConf := &docker.NetworkingConfig{}

	name := "bdika_" + uuid.New()
	opts := docker.CreateContainerOptions{name, conf, hostConf, netConf, context.Background()}

	c, err := client.CreateContainer(opts)
	if err != nil {
		log.Println("enable to create Continer ", err.Error())
		return nil, "", err
	}

	err = client.StartContainer(c.ID, hostConf)
	if err != nil {
		log.Println("enable to  Start ", err.Error())
		return nil, "", err
	}
	i, err := m.ConectToStorage(client, c.ID)
	if err != nil {
		log.Println("enable to  Start ", err.Error())
		return nil, "", err
	}

	return i, c.ID, nil
}
func (m *Mongo) ConectToStorage(client *docker.Client, cid string) (interface{}, error) {
	dc, err := client.InspectContainer(cid)
	ip := dc.NetworkSettings.IPAddress

	session, err := mgo.Dial(fmt.Sprintf("%s:%d", ip, 27017))
	if err != nil {
		log.Println("enable to  conect mongo continer  ", err.Error())
		return nil, err
	}
	for i := 0; i < RETRY; i++ {

		log.Println("try to  conect continer DB  #", i)
		time.Sleep(5 * time.Second)
		err := session.Ping()
		if err != nil {
			log.Println("db ping error ", err)
			continue
		}
		break
	}

	return session, nil
}
func (m *Mongo) RemoveContiner(c *docker.Client, cid string) error {
	err := RemoveContinerID(c, cid)
	if err != nil {
		return err
	}
	return nil
}
