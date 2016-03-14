package main

import (
	"fmt"
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"runtime"
	"time"
)

func tmgo() {
	runtime.GOMAXPROCS(util.CPU())
	// err := mgo.AddDefault("cny:123@10.211.55.3:27017/cny", "cny")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// time.Sleep(500 * time.Second)
	err := mgo.AddDefault("cny:123@192.168.2.57:27017/cny", "cny")
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
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("all done...")
}
