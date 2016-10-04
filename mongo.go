package bdikaa

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"

	"golang.org/x/net/context"

	"github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

type Mongo struct {
	Tag        string
	Seeds      string
	FileName   string
	DbName     string
	Collection string
}

//NewMongoDB ,create a pontier of Mongo struct with defualt values.
func NewMongoDB() *Mongo {
	return &Mongo{"latest", "", "", "", ""}
}

//CreatDockerConfig set continer properties.
func (m *Mongo) CreatDockerConfig() *docker.Config {
	conf := &docker.Config{
		Image:        fmt.Sprintf("mongo:%s", m.Tag),
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}
	return conf
}

//CreatDockerHostConfig set host config properties.
func (m *Mongo) CreatDockerHostConfig() *docker.HostConfig {
	var dh *docker.HostConfig
	if m.Seeds != "" {
		dh = &docker.HostConfig{Binds: []string{m.Seeds + ":/data/seeds"}}
	}
	return dh
}

//CreateContiner  , create and atart the continer and returning the continer ID and Mongo session.
func (m *Mongo) CreateContiner(client *docker.Client) (interface{}, string, error) {
	err := GetImageIfNotExsit(client, "mongo", m.Tag)
	if err != nil {
		log.Println("enable to create Continer ", err.Error())
		return nil, "", err
	}
	conf := m.CreatDockerConfig()
	hostConf := m.CreatDockerHostConfig()
	netConf := &docker.NetworkingConfig{}

	name := fmt.Sprintf("bdika_%s", uuid.New())
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
	i, err := m.ConnectToStorage(client, c.ID)
	if err != nil {
		log.Println("enable to  Start ", err.Error())
		return nil, "", err
	}

	return i, c.ID, nil
}

//ConnectToStorage , check if the mongo continer is ready and import data if nedded .
func (m *Mongo) ConnectToStorage(client *docker.Client, cid string) (interface{}, error) {

	dc, err := client.InspectContainer(cid)
	ip := dc.NetworkSettings.IPAddress

	session, err := mgo.Dial(fmt.Sprintf("%s:%d", ip, 27017))
	if err != nil {
		log.Println("enable to  connect mongo continer  ", err.Error())
		return nil, err
	}
	for i := 0; i < RETRY; i++ {

		log.Println("try to  connect continer DB  %d", i)
		time.Sleep(5 * time.Second)
		err := session.Ping()
		if err != nil {
			log.Println("db ping error ", err)
			continue
		}
		if m.Seeds != "" {
			err := importData(m, cid, client)
			if err != nil {
				log.Println("importData error ", err)
				return nil, err
			}

		}

		break
	}

	return session, nil
}

//RemoveContiner by continer ID.
func (m *Mongo) RemoveContiner(c *docker.Client, cid string) error {
	err := RemoveContinerID(c, cid)
	if err != nil {
		return err
	}
	return nil
}

//importData running a mongo import cmd , by using an exsiteng json file and seting the DbName and collection.
func importData(m *Mongo, cid string, client *docker.Client) error {
	cmd := fmt.Sprintf("mongoimport -d %s -c %s /data/seeds/%s", m.DbName, m.Collection, m.FileName)

	opts := docker.CreateExecOptions{

		Container:    cid,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bash", "-c", cmd},
	}
	fmt.Println("trying to load %s", cmd)

	execID, err := client.CreateExec(opts)
	if err != nil {
		return err
	}

	success := make(chan struct{})
	config := docker.StartExecOptions{
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		InputStream:  os.Stdin,
		RawTerminal:  true,
		Success:      success,
	}

	go func() error {
		err := client.StartExec(execID.ID, config)
		if err != nil {
			fmt.Println("errr: %s", err.Error)
			return err
		}
		return nil
	}()

	<-config.Success
	return nil
}
