# bdikaa

TODO: Write a project description
## Installation

```
go get github.com/miko-code/bdikaa
```
## Usage
```
    //create the docker Api client.
    c, err := bdikaa.GetClinet()
	if err != nil {
		fmt.Println(err.Error())
	}
	  //create  Mysql with default  parameters.
	m := bdikaa.NewMysql()
	//create the mysql container and returning  the container ID  and SQL db instance .
	db, cid, err := m.CreatDockerMysqlContainer(c)
	if err != nil {
		fmt.Println(err.Error())
	}

	db /// do stuf on db

//clean up and remove the continer 
	bdikaa.RemoveContiner(c, cid)
	if err != nil {
		fmt.Println(err.Error())
	}
```

## Credits
This project is based on the amazing  [go-dockerclient](https://github.com/fsouza/go-dockerclient) by @fsouza.
## License
TODO: Write license