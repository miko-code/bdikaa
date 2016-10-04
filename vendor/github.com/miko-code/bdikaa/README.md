# bdikaa
Bdikaa is a framework to test your database action against docker container instead of local databases.
It will pull the correct image and start the container after the test is done it will be removed.
In order to use it you need the get the client  and an instance of the db struct.

 -   Without parameters you will get default Database and you will need to create and add  your own data. 
 
 - You can add a path to existing sql file and it will loaded to the container DB.

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
	 
	 //create  Mysql with custome  parameters and Data.
	dataDir:="PATH TO THE SQL FILE DIR"
	m := bdikaa.NewMysql()
	m.DataDir = dataDir
	
	
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
This project is based on the amazing  [go-dockerclient](https://github.com/fsouza/go-dockerclient) by [@fsouza] ().
## License
The MIT License (MIT)

Copyright (c) 2014 Arshad Chummun

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
Status API Training Shop Blog About
