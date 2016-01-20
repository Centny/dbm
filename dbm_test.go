package dbm

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"runtime"
	"testing"
	"time"
)

type MockDb struct {
	err    bool
	ping   int
	create int
}

func (m *MockDb) Ping(db interface{}) error {
	if m.err {
		return util.Err("error")
	}
	defer func() {
		m.ping += 1
	}()
	switch m.ping {
	case 0:
		return nil
	case 1:
		return util.Err("error")
	case 2:
		return Closed
	default:
		return nil
	}
}
func (m *MockDb) Create() (interface{}, error) {
	if m.create > 0 && m.err {
		return nil, util.Err("error")
	}
	defer func() {
		m.create += 1
	}()
	switch m.create {
	case 0:
		return "db", nil
	case 1:
		return nil, util.Err("mock error")
	case 2:
		return nil, util.Err("mock error")
	default:
		return "new", nil
	}
}
func (m *MockDb) String() string {
	return "mock db"
}

func TestMdbs(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	mdbs, _ := NewMDbs(&MockDb{})
	mdb, _ := NewMDb(&MockDb{})
	mdbs.Add(mdb)
	if mdbs.Db() == nil {
		t.Error("error")
		return
	}
	rundb := func() {
		if mdbs.Db() == nil {
			t.Error("error")
		}
	}
	go func() {
		for {
			for i := 0; i < 3; i++ {
				go rundb()
			}
			time.Sleep(2 * time.Second)
		}
	}()
	time.Sleep(15 * time.Second)
}

func TestTimeout(t *testing.T) {
	mdbs, _ := NewMDbs(&MockDb{err: true})
	mdbs.Timeout = 2000
	wait := make(chan int)
	go func() {
		defer func() {
			fmt.Println(recover())
			wait <- 0
		}()
		mdbs.Db()
	}()
	<-wait
	go func() {
		defer func() {
			fmt.Println(recover())
			wait <- 0
		}()
		mdbs := &MDbs{}
		mdbs.SelMDb()
	}()
	<-wait
}
