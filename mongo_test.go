package bdikaa

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func TestMongoNoData(t *testing.T) {

	client, err := GetClinet()
	assert.Nil(t, err)

	tests := []Continer{&Mongo{"latest"}}

	for _, m := range tests {

		i, cid, err := m.CreateContiner(client)
		assert.Nil(t, err)
		defer m.RemoveContiner(client, cid)
		session := i.(*mgo.Session)
		assert.NotNil(t, session)

		c := session.DB("test").C("people")
		err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
			&Person{"Cla", "+55 53 8402 8510"})
		assert.Nil(t, err)

		result := Person{}
		err = c.Find(bson.M{"name": "Ale"}).One(&result)
		assert.Nil(t, err)

		fmt.Println("Phone:", result.Phone)

	}
}
