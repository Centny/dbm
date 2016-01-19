package mgo

import (
	"fmt"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"runtime"
	"sync"
	"testing"
)

func TestDefault(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	err := AddDefault("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = AddDefault("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if Db() == nil {
		t.Error("error")
		return
	}
	// time.Sleep(5 * time.Second)
	wg := &sync.WaitGroup{}
	rundb := func() {
		err := C("abc").Insert(bson.M{"a": 1, "b": 2})
		if err != nil {
			t.Error(err.Error())
		}
		wg.Done()
	}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go rundb()
	}
	wg.Wait()
}

func TestDbL(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	err := AddDbL("a1", "cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = AddDbL("a2", "cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if DbBy("a1") == nil {
		t.Error("error")
		return
	}
	if DbBy("a2") == nil {
		t.Error("error")
		return
	}
	// time.Sleep(5 * time.Second)
	wg := &sync.WaitGroup{}
	rundb := func() {
		err = CBy("a1", "abc").Insert(bson.M{"a": 1, "b": 2})
		if err != nil {
			t.Error(err.Error())
		}
		err = CBy("a2", "abc").Insert(bson.M{"a": 1, "b": 2})
		if err != nil {
			t.Error(err.Error())
		}
		wg.Done()
	}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go rundb()
	}
	wg.Wait()
	fmt.Println("all done")
}

func TestPerformance(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	err := AddDefault("cny:123@loc.w:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("xxxx->")
	used, err := tutil.DoPerf(5000, "", func(i int) {
		err := C("abc").Insert(bson.M{"a": 1, "b": 2})
		if err != nil {
			t.Error(err.Error())
		}
	})
	fmt.Println(used, err)
}
