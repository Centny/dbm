package mgo

import (
	"fmt"
	"github.com/Centny/dbm"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	tmgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Default = dbm.NewMDbs2()
var DbL = map[string]*dbm.MDbs{}
var Sequence = "sequence"

func Db() *tmgo.Database {
	return Default.Db().(*tmgo.Database)
}
func C(name string) *tmgo.Collection {
	return Db().C(name)
}
func Next(name, id string, increase int64) (oldv, newv int64, err error) {
	return NextV(C(name), id, increase)
}
func Next2(id string, increase int64) (oldv, newv int64, err error) {
	return NextV(C(Sequence), id, increase)
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

func NextBy(key, name, id string, increase int64) (oldv, newv int64, err error) {
	return NextV(CBy(key, name), id, increase)
}

func NextBy2(key, id string, increase int64) (oldv, newv int64, err error) {
	return NextV(CBy(key, Sequence), id, increase)
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
	_, err := mdb.CollectionNames()
	if err != nil && (err.Error() == "Closed explicitly" || err.Error() == "EOF") {
		return dbm.Closed
	} else {
		return err
	}
}
func (m *MGO_H) Create() (interface{}, error) {
	if len(m.Url) < 1 || len(m.Name) < 1 {
		return nil, util.Err("the database con/name is empty")
	}
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

func NextV(c *tmgo.Collection, id string, increase int64) (oldv int64, newv int64, err error) {
	var res = bson.M{}
	info, err := c.Find(bson.M{"_id": id}).Select(bson.M{"sequence": 1}).Apply(tmgo.Change{
		Update:    bson.M{"$inc": bson.M{"sequence": increase}},
		ReturnNew: true,
		Upsert:    true,
	}, &res)
	if err != nil {
		err = util.Err("require sequence(increase:%v) for %v fail, err:%v ", increase, id, err)
		return
	}
	if info.Updated != 1 && info.UpsertedId == nil {
		err = util.Err("require sequence(increase:%v) for %v fail, %v updated ", increase, id, info.Updated)
		return
	}
	newv = res["sequence"].(int64)
	oldv = newv - increase
	return
}
