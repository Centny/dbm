dabase manager tool 
===
dabase manager tool for mongo/mysql(golang)

###Features
* auto try ping to server when geting database connection
* distribute connection on multi link to the same database server


### Install
```
go get github.com/Centny/dbm
```

###Useage
* for mongodb [mgo](https://github.com/Centny/dbm/tree/master/mgo)
* for sql [sql](https://github.com/Centny/dbm/tree/master/sql)
* for other database impl `dbm.DbH` for yourself.

###Run Test
* `go get github.com/Centny/dbm`
* `cd $GOPATH/src/github.com/Centny/dbm/dbmtest`
* edit mysql connect url on `tsql.go` or edit mgo connect url on `tmgo.go`
* `go install github.com/Centny/dbm/dbmtest`
* `$GOPATH/bin/dbmtest tsql` to run mysql test
* `$GOPATH/bin/dbmtest tmgo` to run mgo test



###Example
mgo

```
	runtime.GOMAXPROCS(util.CPU())
	err := mgo.AddDefault("cny:123@10.211.55.3:27017/cny", "cny")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// time.Sleep(500 * time.Second)
	err = mgo.AddDefault("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("add done...")
	rundb := func(v int) {
		fmt.Println("running->", v)
		err := mgo.C("abc").Insert(bson.M{"a": 1, "b": 2})
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("done->", v)
	}
	for {
		for i := 0; i < 5; i++ {
			go rundb(i)
		}
		time.Sleep(3 * time.Second)
	}
	fmt.Println("all done...")
```

sql

```
	runtime.GOMAXPROCS(util.CPU())
	err := sql.AddDefault("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// time.Sleep(500 * time.Second)
	err = sql.AddDefault("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("add done...")
	rundb := func(v int) {
		fmt.Println("running->", v)
		iv, err := dbutil.DbQueryI(sql.Db(), "select 1")
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("done->", iv, v)
	}
	for {
		for i := 0; i < 5; i++ {
			go rundb(i)
		}
		time.Sleep(3 * time.Second)
	}
	fmt.Println("all done...")
```