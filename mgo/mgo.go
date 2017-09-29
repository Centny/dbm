package mgo

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"sort"

	"math"

	"github.com/Centny/dbm"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	tmgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Default = dbm.NewMDbs2()
var DbL = map[string]*dbm.MDbs{}
var Sequence = "sequence"
var Increase int64 = 1000
var WaitTimeout time.Duration = 30000

type CF func(name string) *tmgo.Collection

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
func WaitNext(id string) int64 {
	return PoolWaitNextV(C, Sequence, id, Increase, WaitTimeout)
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

func AddDefault2(urls string) {
	for _, murl := range strings.Split(urls, ";") {
		url_m := strings.SplitAfterN(murl, "*", 2)
		var count = 1
		if len(url_m) > 1 {
			tc, err := util.ParseInt(url_m[1])
			if err != nil {
				panic(fmt.Sprintf("parsing uri(%v) fail with error %v", murl, err))
			}
			count = tc
		}
		var url = strings.TrimSuffix(url_m[0], "*")
		url_n := strings.SplitAfterN(strings.SplitN(url, "?", 2)[0], "/", 2)
		if len(url_n) != 2 {
			panic(fmt.Sprintf("invalid db uri(%v)", url))
		}
		for i := 0; i < count; i++ { //connection multi time
			var tempDelay time.Duration
			for {
				mdb, err := dbm.NewMDb(NewMGO_H(url, url_n[1]))
				if err == nil {
					Default.Add(mdb)
					break
				}
				if tempDelay == 0 {
					tempDelay = 100 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 8 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.E("AddDefault2 connection to server(%v) fail with error(%v), retrying after %v", url, err, tempDelay)
				time.Sleep(tempDelay)
			}
		}
	}
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
	errc int
}

func NewMGO_H(url, name string) *MGO_H {
	return &MGO_H{
		Name: name,
		Url:  url,
	}
}
func (m *MGO_H) Ping(db interface{}) (err error) {
	mdb := db.(*tmgo.Database)
	defer func() {
		terr := recover()
		if terr != nil {
			mdb.Session.Close()
			err = dbm.Closed
		}
	}()
	err = mdb.Session.Ping()
	if err == nil {
		m.errc = 0
		return nil
	}
	m.errc += 1
	if m.errc >= 3 || err.Error() == "Closed explicitly" || err.Error() == "EOF" {
		mdb.Session.Close()
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
	}
	log.D("MGO_H dail to %v error->%v", m, err)
	if ss != nil {
		ss.Close()
	}
	return nil, err
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

//
var pool = map[string][]int64{}
var pool_lck = sync.RWMutex{}

func PoolNextV(cf CF, name, id string, increase int64) (val int64, err error) {
	pool_lck.Lock()
	defer pool_lck.Unlock()
	if vals, ok := pool[id]; ok && len(vals) > 0 {
		pool[id] = vals[1:]
		return vals[0], nil
	}
	oldv, _, err := NextV(cf(name), id, increase)
	if err != nil {
		return 0, err
	}
	vals := []int64{}
	for i := int64(1); i < increase; i++ {
		vals = append(vals, oldv+int64(i)+1)
	}
	pool[id] = vals
	return oldv + 1, nil
}

func PoolWaitNextV(cf CF, name, id string, increase int64, timeout time.Duration) int64 {
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		next, err := PoolNextV(cf, name, id, increase)
		if err != nil {
			if tempDelay == 0 {
				tempDelay = 5 * time.Millisecond
			} else {
				tempDelay *= 2
			}
			if tempDelay > timeout {
				break
			}
			log.W("MDbs next sequence error: %v; retrying in %v", err, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		tempDelay = 0
		return next
	}
	panic("MDbs next sequence time out")
}

func PoolWaitNext(cf CF, id string) int64 {
	return PoolWaitNextV(cf, Sequence, id, Increase, WaitTimeout)
}

// func (m *MDbs) LoadSequence(name, id string, expect int64) (oldv int64, err error) {
// 	oldv, _, err = NextV(m.C(name), id, expect)
// 	return
// }

// func (m *MDbs) NextSequence(name, id string, increase int64) (int64, error) {

// }

// func (m *MDbs) NextString(name, pre, id string, increase int64) string {
// 	return fmt.Sprintf("%v%v", pre, m.NextInt(name, id, increase))
// }

type SortD struct {
	bson.D
	Desc bool
}

func (s SortD) Len() int {
	return len(s.D)
}

func (s SortD) Less(i, j int) bool {
	if s.Desc {
		return math.Abs(util.FloatVal(s.D[i].Value)) > math.Abs(util.FloatVal(s.D[j].Value))
	}
	return math.Abs(util.FloatVal(s.D[i].Value)) < math.Abs(util.FloatVal(s.D[j].Value))
}

func (s SortD) Swap(i, j int) {
	s.D[i], s.D[j] = s.D[j], s.D[i]
}

func (s SortD) Val() bson.D {
	var val = bson.D{}
	for _, v := range s.D {
		if util.IntVal(v.Value) > 0 {
			val = append(val, bson.DocElem{
				Name:  v.Name,
				Value: 1,
			})
		} else {
			val = append(val, bson.DocElem{
				Name:  v.Name,
				Value: -1,
			})
		}
	}
	return val
}

func ParseSortD(msort util.Map) bson.D {
	var sd = SortD{}
	for key := range msort {
		sd.D = append(sd.D, bson.DocElem{
			Name:  key,
			Value: msort.IntVal(key),
		})
	}
	sort.Sort(sd)
	return sd.Val()
}
