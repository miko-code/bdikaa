package bdikaa

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/fsouza/go-dockerclient"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

type Mysql struct {
	RootPass string
	DbName   string
	UserName string
	Pass     string
	DataDir  string
	Tag      string
}

//NewMysql returning defult  values for the mysql configoration.
func NewMysql() *Mysql {
	return &Mysql{
		RootPass: "root",
		DbName:   "dbname",
		UserName: "root",
		Tag:      "latest",
	}
}

//

func (m *Mysql) CreatDockerConfig() *docker.Config {

	conf := &docker.Config{
		Image: fmt.Sprintf("mysql:%s", m.Tag),
		Env: []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", m.RootPass),
			fmt.Sprintf("MYSQL_DATABASE=%s", m.DbName),
		},
	}

	if m.UserName != "" {
		conf.Env = append(conf.Env, fmt.Sprintf("MYSQL_USER=%s", m.UserName))
	}
	if m.Pass != "" {
		conf.Env = append(conf.Env, fmt.Sprintf("MYSQL_PASSWORD=%s", m.Pass))
	}

	return conf
}

// func creatDockerConfig(m *Mysql) *docker.Config {
// 	conf := &docker.Config{
// 		Image: fmt.Sprintf("mysql:%s", m.Tag),
// 		Env: []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", m.RootPass),
// 			fmt.Sprintf("MYSQL_DATABASE=%s", m.DbName),
// 		},
// 	}

// 	if m.UserName != "" {
// 		conf.Env = append(conf.Env, fmt.Sprintf("MYSQL_USER=%s", m.UserName))
// 	}
// 	if m.Pass != "" {
// 		conf.Env = append(conf.Env, fmt.Sprintf("MYSQL_PASSWORD=%s", m.Pass))
// 	}

// 	return conf
// }

func (m *Mysql) CreatDockerHostConfig() *docker.HostConfig {
	var dh *docker.HostConfig
	if m.DataDir != "" {
		dh = &docker.HostConfig{Binds: []string{m.DataDir + ":/docker-entrypoint-initdb.d"}}
	}

	return dh
}

//check if container db is responsive.
func checkIfAlive(m *Mysql, client *docker.Client, cid string) (*sql.DB, error) {
	dc, err := client.InspectContainer(cid)
	ip := dc.NetworkSettings.IPAddress
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.UserName, m.RootPass, ip, 3306, m.DbName)
	log.Println(url)
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println("sql.Open  ", err.Error())
	}

	for i := 0; i < RETRY; i++ {

		log.Println("try to  conect continer DB  #", i)
		time.Sleep(5 * time.Second)
		err := db.Ping()
		if err != nil {
			log.Println("db ping error ", err)
			continue
		}
		break
	}
	return db, err
}

//CreatDockerMysqlContainer creating the  mysql container and returning  the container ID  and SQL db instance .
func (m *Mysql) CreateContiner(client *docker.Client) (interface{}, string, error) {

	err := GetImageIfNotExsit(client, "mysql", m.Tag)
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

	err = client.StartContainer(c.ID, nil)
	if err != nil {
		log.Println("enable to create Start ", err.Error())
		return nil, "", err
	}

	db, err := checkIfAlive(m, client, c.ID)

	if err != nil {
		log.Println("enable to to conecnte  DB ", err.Error())
	}

	return db, c.ID, nil
}
func (m *Mysql) RemoveContiner(c *docker.Client, cid string) error {
	err := RemoveContinerID(c, cid)
	if err != nil {
		return err
	}
	return nil
}
