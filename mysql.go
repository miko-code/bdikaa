package bdikaa

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/fsouza/go-dockerclient"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pborman/uuid"
)

type Mysql struct {
	rootPass string
	dbName   string
	userName string
	pass     string
	dataDir  string
	tag      string
}

func NewMysql() *Mysql {
	return &Mysql{
		rootPass: "root",
		dbName:   "dbname",
		userName: "root",
		tag:      "latest",
	}
}

func creatDockerConfig(m *Mysql) *docker.Config {
	conf := &docker.Config{
		Image: fmt.Sprintf("mysql:%s", m.tag),
		Env: []string{fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", m.rootPass),
			fmt.Sprintf("MYSQL_DATABASE=%s", m.dbName),
		},
	}

	if m.userName != "" {
		conf.Env = append(conf.Env, fmt.Sprintf("MYSQL_USER=%s", m.userName))
	}
	if m.pass != "" {
		conf.Env = append(conf.Env, fmt.Sprintf("MYSQL_PASSWORD=%s", m.pass))
	}

	return conf
}

func creatDockerHostConfig(m *Mysql) *docker.HostConfig {
	var dh *docker.HostConfig
	if m.dataDir != "" {
		dh = &docker.HostConfig{Binds: []string{m.dataDir + ":/docker-entrypoint-initdb.d"}}
	}

	return dh
}

//check if container db is responsive.
func checkIfAlive(m *Mysql, client *docker.Client, cid string) (*sql.DB, error) {
	dc, err := client.InspectContainer(cid)
	ip := dc.NetworkSettings.IPAddress
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.userName, m.rootPass, ip, 3306, m.dbName)
	log.Println(url)
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println("sql.Open  ", err.Error())
	}

	for i := 0; i < RETRY; i++ {

		log.Println("try to  conect continer DB  #", i)
		time.Sleep(5 * time.Second)
		err = db.Ping()
		if err != nil {
			log.Println("db ping error ", err.Error())
			continue
		}

		break
	}

	return db, err
}

//create the mysql container and returning  the container ID  and SQL db instance .
func (m *Mysql) CreatDockerMysqlContainer(client *docker.Client) (*sql.DB, string, error) {

	err := GetImageIfNotExsit(client, "mysql", m.tag)
	if err != nil {
		log.Println("enable to create Continer ", err.Error())
		return nil, "", err
	}

	conf := creatDockerConfig(m)
	hostConf := creatDockerHostConfig(m)

	name := "bdika_" + uuid.New()
	opts := docker.CreateContainerOptions{name, conf, hostConf}

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
