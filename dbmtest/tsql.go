package main

import (
	vsql "database/sql"
	"fmt"
	"github.com/Centny/dbm/sql"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"time"
)

func tsql() {
	runtime.GOMAXPROCS(util.CPU())
	err := sql.AddDefault("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local", 5, 8)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// time.Sleep(500 * time.Second)
	err = sql.AddDefault("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local", 5, 8)
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
}

func tsql2() {
	db, err := vsql.Open("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	for {
		_, err := dbutil.DbQueryInt(db, "select 1")
		log.D("%v", err)
		time.Sleep(time.Second)
	}
}
