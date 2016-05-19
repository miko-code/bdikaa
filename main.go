package bdikaa

//import (
//	"log"

//	_ "github.com/go-sql-driver/mysql"
//)

//func main() {

//	client, err := GetClinet()
//	if err != nil {
//		log.Println("enable to create clinet %s", err)

//	}

//	m := newMysql()

//	db, cid, err := m.CreatDockerMysqlContainer(client)
//	if err != nil {
//		RemoveContiner(client, cid)
//		log.Println("enable to CreatDockerMysqlContainer %s", err)

//	}
//	db.Ping()

//	err = RemoveContiner(client, cid)
//	if err != nil {

//		log.Println("enable to CreatDockerMysqlContainer %s", err)

//	}

//}
