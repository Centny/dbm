package main

import (
	"fmt"
	"github.com/Centny/dbm"
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/util"
	// xmgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"runtime"
	"time"
)

func tmgo() {
	runtime.GOMAXPROCS(util.CPU())
	dbm.ShowLog = true
	// for {
	// 	ss, err := xmgo.Dial("cny:123@loc.w:27017/cny")
	// 	fmt.Println(ss, err)
	// 	if err == nil {
	// 		time.Sleep(1000 * time.Second)
	// 	}
	// }
	// err := mgo.AddDefault("cny:123@10.211.55.3:27017/cny", "cny")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// time.Sleep(500 * time.Second)
	mgo.AddDefault2("cny:123@loc.w:27017/cny*8")
	fmt.Println("add done...")
	time.Sleep(10000 * time.Second)
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
