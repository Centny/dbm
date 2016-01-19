package main

import (
	"fmt"
	"github.com/Centny/dbm/mgo"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"runtime"
	"sync"
	"time"
)

func tmgo() {
	runtime.GOMAXPROCS(util.CPU())
	err := mgo.AddDefault("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = mgo.AddDefault("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("add done...")
	time.Sleep(5 * time.Second)
	wg := &sync.WaitGroup{}
	rundb := func() {
		fmt.Println("running...")
		err := mgo.C("abc").Insert(bson.M{"a": 1, "b": 2})
		if err != nil {
			fmt.Println(err.Error())
		}
		wg.Done()
		fmt.Println("done...")
	}
	wg.Add(200)
	for i := 0; i < 200; i++ {
		go rundb()
	}
	wg.Wait()
	fmt.Println("all done...")
}
