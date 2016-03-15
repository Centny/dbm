package dbm

import (
	"errors"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"sync"
	"sync/atomic"
	"time"
)

var ShowLog bool = false

func slog(format string, args ...interface{}) {
	if ShowLog {
		log.D_(-1, format, args...)
	}
}

type DbH interface {
	Ping(db interface{}) error
	Create() (interface{}, error)
	String() string
}

var Closed = errors.New("Closed explicitly")

type MDb struct {
	H       DbH
	DB      interface{}
	Active  bool
	Running bool
	Timeout int64 //ms
	Hited   uint64

	lck  *sync.Cond
	ping int32
}

func NewMDb(h DbH) (*MDb, error) {
	db, err := h.Create()
	mdb := &MDb{
		H:       h,
		DB:      db,
		Active:  true,
		Running: true,
		Timeout: 3000,
		Hited:   0,
		ping:    0,
		lck:     sync.NewCond(&sync.Mutex{}),
	}
	return mdb, err
}

func (m *MDb) Db() interface{} {
	return m.DB
}

func (m *MDb) recon() {
	db, err := m.H.Create()
	if err == nil {
		m.DB = db
		log.D("MDb connect to %v success, will mark to active", m.String())
		m.lck.L.Lock()
		m.ping = 0
		m.lck.Broadcast()
		m.lck.L.Unlock()
	} else {
		log.E("MDb connect to %v error->%v, will retry after 5s", m.String(), err)
		time.Sleep(5 * time.Second)
		go m.recon()
	}
}
func (m *MDb) rping_() {
	slog("MDb start ping to %v ", m.String())
	err := m.H.Ping(m.DB)
	m.lck.L.Lock()
	m.Active = err == nil
	if err == nil {
		slog("MDb ping to %v success", m.String())
		m.ping = 0
	} else if err.Error() == "Closed explicitly" {
		log.E("MDb ping to %v error->%v, will try reconnect", m.String(), err)
		m.ping = 1
		go m.recon()
	} else {
		log.E("MDb ping to %v error->%v, will mark to not active", m.String(), err)
		m.ping = 0
	}
	m.lck.Broadcast()
	m.lck.L.Unlock()
}
func (m *MDb) tping_() {
	go m.rping_()
	time.Sleep(time.Duration(m.Timeout) * time.Millisecond)
	m.lck.L.Lock()
	if m.ping > 0 {
		m.Active = false
		log.W("MDb ping to %v timeout(%vms), will mark to not active", m.String(), m.Timeout)
		m.lck.Broadcast()
	}
	m.lck.L.Unlock()
}
func (m *MDb) TPing() bool {
	m.lck.L.Lock()
	defer m.lck.L.Unlock()
	if m.ping < 1 {
		m.ping = 1
		go m.tping_()
	}
	if m.Active {
		m.lck.Wait()
	}
	return m.Active
}

func (m *MDb) String() string {
	return fmt.Sprintf("DB(%v),Active(%v),Hited(%v)", m.H.String(), m.Active, m.Hited)
}

type MDbs struct {
	Dbs     []*MDb
	onum    uint32
	Timeout int64
}

func NewMDbs2() *MDbs {
	return &MDbs{
		Timeout: 30000,
	}
}

func NewMDbs(h DbH) (*MDbs, error) {
	mdbs := &MDbs{
		Timeout: 30000,
	}
	mdb, err := NewMDb(h)
	if err == nil {
		mdbs.Add(mdb)
	}
	return mdbs, err
}

func (m *MDbs) Add(mdb ...*MDb) {
	m.Dbs = append(m.Dbs, mdb...)
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
			if mdb.TPing() {
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
