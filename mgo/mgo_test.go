package mgo

import (
	"fmt"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"runtime"
	"sync"
	"testing"
	"time"
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
	// Default.SelMDb().Close()
	// Default.SelMDb().Close()
	// C("abc").Insert(bson.M{"a": 1, "b": 2})
	time.Sleep(2 * time.Second)
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
	time.Sleep(2 * time.Second)
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
	var added = map[int64]bool{}
	var lck = sync.RWMutex{}
	C(Sequence).RemoveAll(nil)
	fmt.Println(Next2("abc", 1))
	fmt.Println(Next2("abc", 1))
	used, err := tutil.DoPerf(2000, "", func(i int) {
		err := C("abc").Insert(bson.M{"a": 1, "b": 2})
		if err != nil {
			t.Error(err.Error())
			return
		}
		_, nv, err := Next2("abc", 1)
		if err != nil {
			t.Error(err.Error())
			return
		}
		_, nv2, err := Next(Sequence, "abc", 1)
		if err != nil {
			t.Error(err.Error())
			return
		}
		lck.Lock()
		if added[nv] {
			panic("exit")
		}
		if added[nv2] {
			panic("exit")
		}
		added[nv] = true
		added[nv2] = true
		lck.Unlock()
	})
	fmt.Println(used, err)
	fmt.Println(Next2("abc", 0))
}

func TestErr(t *testing.T) {
	wait := make(chan int)
	go func() {
		defer func() {
			fmt.Println(recover())
			wait <- 0
		}()
		DbBy("kkksks")
	}()
	<-wait
	NewMGO_H("xxx", "name").Create()
}
