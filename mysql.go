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

type mysql struct {
	rootPass string
	dbName   string
	userName string
	pass     string
	confFile string
	dataDir  string
	tag      string
}

func newMysql() *mysql {
	return &mysql{
		rootPass: "root",
		dbName:   "dbname",
		userName: "root",
		tag:      "latest",
	}
}

func creatDockerConfig(m *mysql) *docker.Config {
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

func creatDockerHostConfig(m *mysql) *docker.HostConfig {
	var dh *docker.HostConfig
	if m.dataDir != "" {
		dh = &docker.HostConfig{Binds: []string{m.dataDir + ":/docker-entrypoint-initdb.d"}}
	}

	return dh
}
func checkIfAlive(m *mysql, client *docker.Client, cid string) (*sql.DB, error) {
	dc, err := client.InspectContainer(cid)
	ip := dc.NetworkSettings.IPAddress
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.userName, m.rootPass, ip, 3306, m.dbName)
	log.Println(url)
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Println("sql.Open %s ", err.Error())
	}

	for i := 0; i < RETRY; i++ {

		log.Println("try to  conecet  #%d", i)
		time.Sleep(5 * time.Second)
		err = db.Ping()
		if err != nil {
			log.Println("db ping error %s", err.Error())
			continue
		}

		break
	}

	return db, err
}

func (m *mysql) CreatDockerMysqlContainer(client *docker.Client) (*sql.DB, string, error) {
	conf := creatDockerConfig(m)
	hostConf := creatDockerHostConfig(m)

	name := "bdika_" + uuid.New()
	opts := docker.CreateContainerOptions{name, conf, hostConf}

	c, err := client.CreateContainer(opts)
	if err != nil {
		log.Println("enable to create Continer %s", err)
		return nil, "", err
	}

	err = client.StartContainer(c.ID, nil)
	if err != nil {
		log.Println("enable to create Start %s", err)
		return nil, "", err
	}

	db, err := checkIfAlive(m, client, c.ID)

	if err != nil {
		log.Println("enable to to conecnte  DB %s", err)
	}

	return db, c.ID, nil
}
