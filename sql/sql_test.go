package sql

import (
	"fmt"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"sync"
	"testing"
)

func TestDefault(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	err := AddDefault2("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = AddDefault2("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
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
		iv, err := dbutil.DbQueryInt(Db(), "select 1")
		if err != nil {
			t.Error(err.Error())
		}
		wg.Done()
		fmt.Println(iv)
	}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go rundb()
	}
	wg.Wait()
}

func TestDbL(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	err := AddDbL2("a1", "mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = AddDbL2("a2", "mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
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
		_, err = dbutil.DbQueryF(DbBy("a1"), "select 1")
		if err != nil {
			t.Error(err.Error())
		}
		_, err = dbutil.DbQueryF(DbBy("a2"), "select 1")
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
	err := AddDefault2("mysql", "cny:123@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
	if err != nil {
		t.Error(err.Error())
		return
	}
	used, err := tutil.DoPerf(5000, "", func(i int) {
		_, err = dbutil.DbQueryF(Db(), "select 1")
		if err != nil {
			t.Error(err.Error())
		}
	})
	fmt.Println(used, err)
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
	NewSQL_H("xxx", "name", 1, 3).Create()
}
