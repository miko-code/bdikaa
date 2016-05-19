package bdikaa

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMysqlNoData(t *testing.T) {
	client, err := GetClinet()
	assert.Nil(t, err)

	tests := []*mysql{&mysql{"root", "dbname", "root", "", "", "", "5.6"},
		&mysql{"root", "dbname", "root", "", "", "", "latest"}}

	for _, m := range tests {

		db, cid, err := m.CreatDockerMysqlContainer(client)
		db.Ping()
		assert.Nil(t, err)
		db.Close()
		err = RemoveContiner(client, cid)
		assert.Nil(t, err)
	}
}

func TestMysqlWithData(t *testing.T) {
	client, err := GetClinet()
	assert.Nil(t, err)
	// geting the curnnet dir.
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataDir := strings.Replace(dir, " ", "\\ ", -1) + "/data"
	tests := []*mysql{&mysql{"root", "world", "root", "", "", dataDir, "latest"}}
	for _, m := range tests {

		db, cid, err := m.CreatDockerMysqlContainer(client)
		rows, err := db.Query("SELECT *  FROM  City")
		assert.True(t, rows.Next(), "expected true got %s ", err)
		db.Close()
		err = RemoveContiner(client, cid)
		assert.Nil(t, err)
	}

}
