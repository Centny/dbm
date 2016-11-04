package mgo

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/Centny/dbm"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Indexes = map[string]map[string]mgo.Index{
	"abc": map[string]mgo.Index{
		"abc_a": mgo.Index{
			Key: []string{"a"},
		},
		"abc_b": mgo.Index{
			Key: []string{"b"},
		},
	},
}

var Indexes2 = map[string]map[string]mgo.Index{
	"abc": map[string]mgo.Index{
		"abc_a": mgo.Index{
			Key: []string{"a"},
		},
		"abc_c": mgo.Index{
			Key: []string{"c"},
		},
	},
}

func TestDefault(t *testing.T) {
	dbm.ShowLog = true
	time.Sleep(time.Second)
	runtime.GOMAXPROCS(util.CPU())
	Default = dbm.NewMDbs2()
	err := AddDefault("cny:123@loc.m:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = AddDefault("cny:123@loc.m:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if Db() == nil {
		t.Error("error")
		return
	}
	Db().C("abc").DropCollection()
	//
	err = ChkIdx(C, Indexes)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = ChkIdx(C, Indexes2)
	if err != nil {
		t.Error(err.Error())
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

func TestDefault2(t *testing.T) {
	Default = dbm.NewMDbs2()
	AddDefault2("cny:123@loc.m:27017/cny*5")
	if len(Default.Dbs) != 5 {
		t.Error("error")
		return
	}
	AddDefault2("cny:123@loc.m:27017/cny*5;cny:123@loc.m:27017/cny;cny:123@loc.m:27017/cny*3")
	if len(Default.Dbs) != 14 {
		t.Error("error")
		return
	}
}

func TestDbL(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	Default = dbm.NewMDbs2()
	err := AddDbL("a1", "cny:123@loc.m:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = AddDbL("a2", "cny:123@loc.m:27017/cny", "cny")
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
	Default = dbm.NewMDbs2()
	err := AddDefault("cny:123@loc.m:27017/cny", "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("xxxx->")
	var added = map[int64]bool{}
	var lck = sync.RWMutex{}
	C(Sequence).RemoveAll(nil)
	fmt.Println(Next2("abc", 1))
	fmt.Println(Next2("abc", 10))
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
	fmt.Println("used->", used, err)
	fmt.Println(Next2("abc", 0))
	//
	// added = map[int64]bool{}
	used, err = tutil.DoPerf(20000, "", func(i int) {
		nv := WaitNext("abc")
		if err != nil {
			t.Error(err.Error())
			return
		}
		nv2 := WaitNext("abc")
		if err != nil {
			t.Error(err.Error())
			return
		}
		lck.Lock()
		if added[nv] {
			fmt.Println(nv)
			fmt.Println(pool)
			panic("exists")
		}
		if added[nv2] {
			fmt.Println(nv2)
			panic("exists")
		}
		added[nv] = true
		added[nv2] = true
		lck.Unlock()
	})
	fmt.Println("used->", used, err)
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

func TestSequenc(t *testing.T) {

}
