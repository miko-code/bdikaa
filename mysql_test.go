package bdikaa

import (
	"os"
	"strings"
	"testing"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMysqlNoData(t *testing.T) {
	client, err := GetClinet()
	assert.Nil(t, err)

	tests := []Continer{&Mysql{"root", "dbname", "root", "", "", "5.6"},
		&Mysql{"root", "dbname", "root", "", "", "latest"}}

	for _, m := range tests {

		i, cid, err := m.CreateContiner(client)
		assert.Nil(t, err)
		db := i.(*sql.DB)
		err = db.Ping()
		assert.Nil(t, err)
		db.Close()
		err = m.RemoveContiner(client, cid)
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
	dataDir := strings.Replace(dir, " ", "\\ ", -1) + "/data/sql"
	tests := []Continer{&Mysql{"root", "dbname", "root", "", dataDir, "5.6"},
		&Mysql{"root", "dbname", "root", "", dataDir, "latest"}}

	for _, m := range tests {

		i, cid, err := m.CreateContiner(client)
		assert.Nil(t, err)
		db := i.(*sql.DB)
		err = db.Ping()
		assert.Nil(t, err)
		rows, err := db.Query("SELECT *  FROM  City")
		assert.True(t, rows.Next(), "expected true got  ", err)

		db.Close()
		err = m.RemoveContiner(client, cid)
		assert.Nil(t, err)

	}

}
