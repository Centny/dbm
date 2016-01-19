package mgo

import (
	"fmt"
	"github.com/Centny/gwf/log"
	tmgo "gopkg.in/mgo.v2"
	"sync/atomic"
	"time"
)

var Default = &MDbs{}
var DbL = map[string]*MDbs{}

func Db() *tmgo.Database {
	return Default.Db()
}
func C(name string) *tmgo.Collection {
	return Default.C(name)
}

func DbBy(key string) *tmgo.Database {
	if mdb, ok := DbL[key]; ok {
		return mdb.Db()
	} else {
		panic("database is not found by name " + key)
	}
}

func CBy(key string, name string) *tmgo.Collection {
	if mdb, ok := DbL[key]; ok {
		return mdb.C(name)
	} else {
		panic("database is not found by name " + key)
	}
}

func AddDefault(url, name string) error {
	mdb, err := NewMDb(url, name)
	if err == nil {
		Default.Add(mdb)
	}
	return err
}

func AddDbL(key, url, name string) error {
	mdbs, err := NewMDbs(url, name)
	if err == nil {
		DbL[key] = mdbs
	}
	return err
}

type MDb struct {
	*tmgo.Session
	Url     string
	Name    string
	Active  bool
	Running bool
	Delay   time.Duration
}

func NewMDb(url, name string) (*MDb, error) {
	return NewMDb2(url, name, true)
}

func NewMDb2(url, name string, ping bool) (*MDb, error) {
	log.I("dail to mgo by url(%v),name(%v),ping(%v)", url, name, ping)
	ss, err := tmgo.Dial(url)
	mdb := &MDb{
		Session: ss,
		Url:     url,
		Name:    name,
		Active:  true,
		Running: true,
		Delay:   2000,
	}
	if err == nil && ping {
		go mdb.RunPing()
	}
	return mdb, err
}

func (m *MDb) Db() *tmgo.Database {
	return m.DB(m.Name)
}

func (m *MDb) C(name string) *tmgo.Collection {
	return m.Db().C(name)
}

func (m *MDb) RunPing() {
	for m.Running {
		log.I("MDb start ping to %v ", m.String())
		m.Active = false
		err := m.Ping()
		m.Active = err == nil
		if err == nil {
			log.I("MDb ping to %v success, will retry after %vms", m.String(), int64(m.Delay))
		} else {
			log.E("MDb ping to %v error->%v, will retry after %vms", m.String(), err, int64(m.Delay))
		}
		time.Sleep(m.Delay * time.Millisecond)
	}
}
func (m *MDb) String() string {
	return fmt.Sprintf("Url(%v),Name(%v)", m.Url, m.Name)
}

type MDbs struct {
	Dbs  []*MDb
	onum uint32
}

func NewMDbs(url, name string) (*MDbs, error) {
	mdbs := &MDbs{}
	mdb, err := NewMDb(url, name)
	if err == nil {
		mdbs.Add(mdb)
	}
	return mdbs, err
}

func (m *MDbs) Add(mdb ...*MDb) {
	m.Dbs = append(m.Dbs, mdb...)
}

func (m *MDbs) Db() *tmgo.Database {
	return m.SelMDb().Db()
}

func (m *MDbs) C(name string) *tmgo.Collection {
	return m.SelMDb().C(name)
}

func (m *MDbs) SelMDb() *MDb {
	all := len(m.Dbs)
	if all < 1 {
		panic("database session list is empty, please add at last one")
	}
	tidx := atomic.AddUint32(&m.onum, 1)
	bidx := int(tidx % uint32(all))
	for {
		for i := 0; i < all; i++ {
			mdb := m.Dbs[(bidx+i)%all]
			if mdb.Active {
				return mdb
			}
		}
		log.W("MDbs all session is not active, it will retry after 1s")
		time.Sleep(time.Second)
	}
	panic("never calling to this")

}
