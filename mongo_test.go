package bdikaa

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type City struct {
	Name  string
	State string
}

func TestMongoNoData(t *testing.T) {

	client, err := GetClinet()
	assert.Nil(t, err)

	tests := []Continer{&Mongo{"latest", "", "", "", ""}}

	for _, m := range tests {

		i, cid, err := m.CreateContiner(client)
		assert.Nil(t, err)
		defer m.RemoveContiner(client, cid)
		session := i.(*mgo.Session)
		assert.NotNil(t, session)

		c := session.DB("Country").C("Citys")
		err = c.Insert(&City{"GRANBY", "MA"},
			&City{"WILBRAHAM", "MA"})
		assert.Nil(t, err)

		result := City{}
		err = c.Find(bson.M{"name": "GRANBY"}).One(&result)
		assert.Nil(t, err)

		fmt.Println("State:", result.State)

	}

}

func TestMongoNoWithData(t *testing.T) {
	client, err := GetClinet()
	assert.Nil(t, err)
	dir, err := os.Getwd()
	assert.Nil(t, err)
	seeds := strings.Replace(dir, " ", "\\ ", -1) + "/data/mongo"
	tests := []Continer{&Mongo{"latest", seeds, "zips.json", "Country", "Citys"}}
	for _, m := range tests {
		i, _, err := m.CreateContiner(client)
		assert.Nil(t, err)
		//		defer m.RemoveContiner(client, cid)
		session := i.(*mgo.Session)
		assert.NotNil(t, session)

		result := City{}
		c := session.DB("Country").C("Citys")
		err = c.Find(bson.M{"name": "GRANBY"}).One(&result)
		fmt.Println("State:", result.State)
		assert.Equal(t, "MA", result.State)
	}
}
