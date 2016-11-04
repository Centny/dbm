package dbm

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
)

var ShowLog bool = false
var ShowLogTime int64 = 8000

func slog(format string, args ...interface{}) {
	if ShowLog {
		log.D_(-1, format, args...)
	}
}

type DbH interface {
	Ping(db interface{}) error
	Create() (interface{}, error)
	// Close(db interface{})
	String() string
}

var Closed = errors.New("Closed explicitly")

type MDb struct {
	H      DbH
	DB     interface{}
	Active bool
	Hited  uint64

	lck      *sync.RWMutex
	ping     int32
	last     int64
	log_time int64
}

func NewMDb(h DbH) (*MDb, error) {
	db, err := h.Create()
	mdb := &MDb{
		H:      h,
		DB:     db,
		Active: true,
		Hited:  0,
		ping:   0,
		lck:    &sync.RWMutex{},
	}
	return mdb, err
}

func (m *MDb) Db() interface{} {
	return m.DB
}

func (m *MDb) TPing() {
	var showlog = ShowLog && (util.Now()-m.log_time > ShowLogTime)
	if showlog {
		m.log_time = util.Now()
		log.D("MDb start ping to %v ", m.String())
	}
	err := m.H.Ping(m.DB)
	if err == nil || err.Error() != "Closed explicitly" {
		if err == nil {
			if showlog {
				log.D("MDb ping to %v success", m.String())
			}
		} else {
			log.E("MDb ping to %v error->%v, will mark to not active", m.String(), err)
		}
		m.lck.Lock()
		m.Active = err == nil
		m.ping = 0
		m.lck.Unlock()
		return
	}
	// m.H.Close(m.DB)
	//do reconnect
	log.E("MDb ping to %v error->%v, will try reconnect", m.String(), err)
	for {
		db, err := m.H.Create()
		if err == nil {
			log.D("MDb connect to %v success, will mark to active", m.String())
			m.lck.Lock()
			m.DB = db
			m.ping = 0
			m.Active = true
			m.lck.Unlock()
			break
		} else {
			log.E("MDb connect to %v error->%v, will retry after 5s", m.String(), err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (m *MDb) IsActive() bool {
	m.lck.RLock()
	defer m.lck.RUnlock()
	return m.Active
}

func (m *MDb) String() string {
	return fmt.Sprintf("DB(%v),Active(%v),Hited(%v)", m.H.String(), m.Active, m.Hited)
}

type MDbs struct {
	Dbs     []*MDb
	onum    uint32
	Timeout int64
	Delay   int64
	Running bool
}

func NewMDbs2() *MDbs {
	var mdbs = &MDbs{
		Timeout: 30000,
		Delay:   3000,
	}
	// mdbs.StartLoop()
	return mdbs
}

func NewMDbs(h DbH) (*MDbs, error) {
	mdbs := &MDbs{
		Timeout: 30000,
		Delay:   3000,
	}
	// mdbs.StartLoop()
	mdb, err := NewMDb(h)
	if err == nil {
		mdbs.Add(mdb)
	}
	return mdbs, err
}

func (m *MDbs) Add(mdb ...*MDb) {
	m.Dbs = append(m.Dbs, mdb...)
	m.StartLoop()
}

func (m *MDbs) Db() interface{} {
	return m.SelMDb().Db()
}

func (m *MDbs) SelMDb() *MDb {
	all := len(m.Dbs)
	if all < 1 {
		panic("database session list is empty, please add at last one")
	}
	tidx := atomic.AddUint32(&m.onum, 1)
	bidx := int(tidx % uint32(all))
	beg := util.Now()
	for {
		for i := 0; i < all; i++ {
			mdb := m.Dbs[(bidx+i)%all]
			if mdb.IsActive() {
				atomic.AddUint64(&mdb.Hited, 1)
				return mdb
			}
		}
		log.W("MDbs all session is not active, it will retry after 1s")
		time.Sleep(time.Second)
		if util.Now()-beg > m.Timeout {
			break
		}
	}
	panic(fmt.Sprintf("MDbs wait database active timeout(%vms)", m.Timeout))
}

func (m *MDbs) LoopPing() {
	log.D("MDbs loop ping is started...")
	for m.Running {
		for _, mdb := range m.Dbs {
			if mdb.ping < 1 {
				mdb.ping = 1
				go mdb.TPing()
			}
		}
		time.Sleep(time.Duration(m.Delay) * time.Millisecond)
	}
	log.D("MDbs loop ping is stopped...")
}

func (m *MDbs) StartLoop() {
	if m.Running {
		return
	}
	m.Running = true
	go m.LoopPing()
}
