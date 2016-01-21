package sql

import (
	"database/sql"
	"fmt"
	"github.com/Centny/dbm"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
)

var Default = &dbm.MDbs{}
var DbL = map[string]*dbm.MDbs{}

func Db() *sql.DB {
	return Default.Db().(*sql.DB)
}

func DbBy(key string) *sql.DB {
	if mdb, ok := DbL[key]; ok {
		return mdb.Db().(*sql.DB)
	} else {
		panic("database is not found by name " + key)
	}
}

func AddDefault(driver, url string, idle, max int) error {
	mdb, err := dbm.NewMDb(NewSQL_H(driver, url, idle, max))
	if err == nil {
		Default.Add(mdb)
	}
	return err
}
func AddDefault2(driver, url string) error {
	return AddDefault(driver, url, util.CPU(), util.CPU()*2)
}

func AddDbL(key, driver, url string, idle, max int) error {
	mdbs, err := dbm.NewMDbs(NewSQL_H(driver, url, idle, max))
	if err == nil {
		DbL[key] = mdbs
	}
	return err
}
func AddDbL2(key, driver, url string) error {
	return AddDbL(key, driver, url, util.CPU(), util.CPU()*2)
}

type SQL_H struct {
	Driver string
	Url    string
	Idle   int
	Max    int
}

func NewSQL_H(driver, url string, idle, max int) *SQL_H {
	return &SQL_H{
		Driver: driver,
		Url:    url,
		Idle:   idle,
		Max:    max,
	}
}

func (s *SQL_H) Ping(db interface{}) error {
	sdb := db.(*sql.DB)
	_, err := dbutil.DbQueryI(sdb, "select 1")
	return err
}
func (s *SQL_H) Create() (interface{}, error) {
	log.D("SQL_H start connect to %v", s)
	db, err := sql.Open(s.Driver, s.Url)
	if err == nil {
		db.SetMaxIdleConns(s.Idle)
		db.SetMaxOpenConns(s.Max)
		log.D("SQL_H connect to %v success", s)
		return db, nil
	} else {
		log.D("SQL_H connect to %v error->%v", s, err)
		return nil, err
	}
}
func (s *SQL_H) String() string {
	return fmt.Sprintf("SQL(Driver:%v,Url:%v)", s.Driver, s.Url)
}
