package bdikaa

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Citys struct {
	City  string
	State string
}

func TestMongoNoData(t *testing.T) {

	client, err := GetClinet()
	assert.Nil(t, err)

	tests := []Continer{NewMongoDB()}

	for _, m := range tests {

		i, cid, err := m.CreateContiner(client)
		assert.Nil(t, err)
		defer m.RemoveContiner(client, cid)
		session := i.(*mgo.Session)
		assert.NotNil(t, session)

		c := session.DB("Country").C("Citys")
		err = c.Insert(&Citys{"GRANBY", "MA"},
			&Citys{"WILBRAHAM", "MA"})
		assert.Nil(t, err)

		result := Citys{}
		err = c.Find(bson.M{"city": "GRANBY"}).One(&result)
		assert.Nil(t, err)

		fmt.Println("State:", result.State)

	}

}

func TestMongoWithData(t *testing.T) {
	client, err := GetClinet()
	assert.Nil(t, err)
	dir, err := os.Getwd()
	assert.Nil(t, err)
	seeds := strings.Replace(dir, " ", "\\ ", -1) + "/data/mongo"
	tests := []Continer{&Mongo{"latest", seeds, "zips.json", "country", "citys"}}
	for _, m := range tests {
		i, cid, err := m.CreateContiner(client)
		assert.Nil(t, err)
		defer m.RemoveContiner(client, cid)
		session := i.(*mgo.Session)
		assert.NotNil(t, session)

		result := Citys{}
		time.Sleep(5 * time.Second)
		c := session.DB("country").C("citys")
		err = c.Find(bson.M{"city": "GRANBY"}).One(&result)
		assert.Nil(t, err)
		assert.Equal(t, "MA", result.State)
	}
}
