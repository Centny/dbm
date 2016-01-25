package mgo

import (
	"fmt"
	"github.com/Centny/dbm"
	"github.com/Centny/gwf/log"
	tmgo "gopkg.in/mgo.v2"
)

var Default = &dbm.MDbs{}
var DbL = map[string]*dbm.MDbs{}

func Db() *tmgo.Database {
	return Default.Db().(*tmgo.Database)
}
func C(name string) *tmgo.Collection {
	return Db().C(name)
}

func DbBy(key string) *tmgo.Database {
	if mdb, ok := DbL[key]; ok {
		return mdb.Db().(*tmgo.Database)
	} else {
		panic("database is not found by name " + key)
	}
}

func CBy(key string, name string) *tmgo.Collection {
	return DbBy(key).C(name)
}

func AddDefault(url, name string) error {
	mdb, err := dbm.NewMDb(NewMGO_H(url, name))
	if err == nil {
		Default.Add(mdb)
	}
	return err
}

func AddDbL(key, url, name string) error {
	mdbs, err := dbm.NewMDbs(NewMGO_H(url, name))
	if err == nil {
		DbL[key] = mdbs
	}
	return err
}

type MGO_H struct {
	Name string
	Url  string
}

func NewMGO_H(url, name string) *MGO_H {
	return &MGO_H{
		Name: name,
		Url:  url,
	}
}
func (m *MGO_H) Ping(db interface{}) error {
	mdb := db.(*tmgo.Database)
	err := mdb.Session.Ping()
	if err != nil && err.Error() == "Closed explicitly" {
		return dbm.Closed
	} else {
		return err
	}
}
func (m *MGO_H) Create() (interface{}, error) {
	log.D("MGO_H start dail to %v", m)
	ss, err := tmgo.Dial(m.Url)
	if err == nil {
		log.D("MGO_H dail to %v success", m)
		return ss.DB(m.Name), nil
	} else {
		log.D("MGO_H dail to %v error->%v", m, err)
		return nil, err
	}
}
func (m *MGO_H) String() string {
	return fmt.Sprintf("MGO(Name:%v,Url:%v)", m.Name, m.Url)
}

type MDbs struct {
	*dbm.MDbs
}

func NewMDbs(url, name string) (*MDbs, error) {
	mdbs, err := dbm.NewMDbs(NewMGO_H(url, name))
	return &MDbs{MDbs: mdbs}, err
}

func (m *MDbs) Db() *tmgo.Database {
	return m.MDbs.Db().(*tmgo.Database)
}

func (m *MDbs) C(name string) *tmgo.Collection {
	return m.Db().C(name)
}
