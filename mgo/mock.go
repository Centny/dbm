package mgo

// import (
// 	"errors"
// 	tmgo "gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// 	"time"
// )

// var Mock = false
// var MckC = map[string]int{}
// var mckc = map[string]int{}
// var MckV = map[string]interface{}{}
// var MckE = errors.New("Mock Error")

// func SetMckC(key string, c int) {
// 	MckC[key] = c
// }

// func SetMckV(key string, c int, v interface{}) {
// 	MckC[key] = c
// 	MckV[key] = v
// }

// func chk_mock(key string) bool {
// 	tc, ok := MckC[key]
// 	if ok {
// 		eq := tc == mckc[key]
// 		mckc[key] += 1
// 		return eq
// 	} else {
// 		return false
// 	}
// }

// type Collection interface {
// 	Bulk() *tmgo.Bulk
// 	Count() (n int, err error)
// 	Create(info *tmgo.CollectionInfo) error
// 	DropCollection() error
// 	DropIndex(key ...string) error
// 	DropIndexName(name string) error
// 	EnsureIndex(index tmgo.Index) error
// 	EnsureIndexKey(key ...string) error
// 	Indexes() (indexes []tmgo.Index, err error)
// 	Insert(docs ...interface{}) error
// 	NewIter(session *tmgo.Session, firstBatch []bson.Raw, cursorId int64, err error) *tmgo.Iter
// 	Remove(selector interface{}) error
// 	RemoveAll(selector interface{}) (info *tmgo.ChangeInfo, err error)
// 	RemoveId(id interface{}) error
// 	Repair() *tmgo.Iter
// 	Update(selector interface{}, update interface{}) error
// 	UpdateAll(selector interface{}, update interface{}) (info *tmgo.ChangeInfo, err error)
// 	UpdateId(id interface{}, update interface{}) error
// 	Upsert(selector interface{}, update interface{}) (info *tmgo.ChangeInfo, err error)
// 	UpsertId(id interface{}, update interface{}) (info *tmgo.ChangeInfo, err error)
// 	With(s *tmgo.Session) *tmgo.Collection
// 	Find(query interface{}) Query
// 	FindId(id interface{}) Query
// 	Pipe(pipeline interface{}) Pipe
// }

// type MGO_Collection struct {
// 	*tmgo.Collection
// }

// func (m *MGO_Collection) Count() (n int, err error) {
// 	if Mock && chk_mock("Collection-Count") {
// 		if val, ok := MckV["Collection-Count"]; ok {
// 			return val.(int), nil
// 		} else {
// 			return 0, MckE
// 		}
// 	}
// 	return m.Collection.Count()
// }
// func (m *MGO_Collection) Create(info *tmgo.CollectionInfo) error {
// 	if Mock && chk_mock("Collection-Create") {
// 		return MckE
// 	}
// 	return m.Collection.Create(info)
// }
// func (m *MGO_Collection) DropCollection() error {
// 	if Mock && chk_mock("Collection-DropCollection") {
// 		return MckE
// 	}
// 	return m.Collection.DropCollection()
// }
// func (m *MGO_Collection) DropIndex(key ...string) error {
// 	if Mock && chk_mock("Collection-DropIndex") {
// 		return MckE
// 	}
// 	return m.Collection.DropIndex(key...)
// }
// func (m *MGO_Collection) DropIndexName(name string) error {
// 	if Mock && chk_mock("Collection-DropIndexName") {
// 		return MckE
// 	}
// 	return m.Collection.DropIndexName(name)
// }
// func (m *MGO_Collection) EnsureIndex(index tmgo.Index) error {
// 	if Mock && chk_mock("Collection-EnsureIndex") {
// 		return MckE
// 	}
// 	return m.Collection.EnsureIndex(index)
// }
// func (m *MGO_Collection) EnsureIndexKey(key ...string) error {
// 	if Mock && chk_mock("Collection-EnsureIndexKey") {
// 		return MckE
// 	}
// 	return m.Collection.EnsureIndexKey(key...)
// }
// func (m *MGO_Collection) Indexes() (indexes []tmgo.Index, err error) {
// 	if Mock && chk_mock("Collection-Indexes") {
// 		if val, ok := MckV["Collection-Indexes"]; ok {
// 			return val.([]tmgo.Index), nil
// 		} else {
// 			return nil, MckE
// 		}
// 	}
// 	return m.Collection.Indexes()
// }
// func (m *MGO_Collection) Insert(docs ...interface{}) error {
// 	if Mock && chk_mock("Collection-Insert") {
// 		return MckE
// 	}
// 	return m.Collection.Insert(docs...)
// }
// func (m *MGO_Collection) Remove(selector interface{}) error {
// 	if Mock && chk_mock("Collection-Remove") {
// 		return MckE
// 	}
// 	return m.Collection.Remove(selector)
// }
// func (m *MGO_Collection) RemoveAll(selector interface{}) (info *tmgo.ChangeInfo, err error) {
// 	if Mock && chk_mock("Collection-RemoveAll") {
// 		if val, ok := MckV["Collection-RemoveAll"]; ok {
// 			return val.(*tmgo.ChangeInfo), nil
// 		} else {
// 			return nil, MckE
// 		}
// 	}
// 	return m.Collection.RemoveAll(selector)
// }
// func (m *MGO_Collection) RemoveId(id interface{}) error {
// 	if Mock && chk_mock("Collection-RemoveId") {
// 		return MckE
// 	}
// 	return m.Collection.RemoveId(id)
// }
// func (m *MGO_Collection) Update(selector interface{}, update interface{}) error {
// 	if Mock && chk_mock("Collection-Update") {
// 		return MckE
// 	}
// 	return m.Collection.Update(selector, update)
// }
// func (m *MGO_Collection) UpdateAll(selector interface{}, update interface{}) (info *tmgo.ChangeInfo, err error) {
// 	if Mock && chk_mock("Collection-UpdateAll") {
// 		if val, ok := MckV["Collection-UpdateAll"]; ok {
// 			return val.(*tmgo.ChangeInfo), nil
// 		} else {
// 			return nil, MckE
// 		}
// 	}
// 	return m.Collection.UpdateAll(selector, update)
// }
// func (m *MGO_Collection) UpdateId(id interface{}, update interface{}) error {
// 	if Mock && chk_mock("Collection-UpdateId") {
// 		return MckE
// 	}
// 	return m.Collection.UpdateId(id, update)
// }
// func (m *MGO_Collection) Upsert(selector interface{}, update interface{}) (info *tmgo.ChangeInfo, err error) {
// 	if Mock && chk_mock("Collection-Upsert") {
// 		if val, ok := MckV["Collection-Upsert"]; ok {
// 			return val.(*tmgo.ChangeInfo), nil
// 		} else {
// 			return nil, MckE
// 		}
// 	}
// 	return m.Collection.Upsert(selector, update)
// }
// func (m *MGO_Collection) UpsertId(id interface{}, update interface{}) (info *tmgo.ChangeInfo, err error) {
// 	if Mock && chk_mock("Collection-UpsertId") {
// 		if val, ok := MckV["Collection-UpsertId"]; ok {
// 			return val.(*tmgo.ChangeInfo), nil
// 		} else {
// 			return nil, MckE
// 		}
// 	}
// 	return m.Collection.UpsertId(id, update)
// }
// func (m *MGO_Collection) Find(query interface{}) Query {
// 	return &MGO_Query{Query: m.Collection.Find(query)}
// }
// func (m *MGO_Collection) FindId(id interface{}) Query {
// 	return &MGO_Query{Query: m.Collection.FindId(id)}
// }
// func (m *MGO_Collection) Pipe(pipeline interface{}) Pipe {
// 	return &MGO_Pipe{Pipe: m.Collection.Pipe(pipeline)}
// }

// type Pipe interface {
// 	All(result interface{}) error
// 	AllowDiskUse() Pipe
// 	Batch(n int) Pipe
// 	Explain(result interface{}) error
// 	Iter() *tmgo.Iter
// 	One(result interface{}) error
// }

// type MGO_Pipe struct {
// 	*tmgo.Pipe
// }

// func (m *MGO_Pipe) All(result interface{}) error {
// 	if Mock && chk_mock("Pipe-All") {
// 		return MckE
// 	}
// 	return m.Pipe.All(result)
// }
// func (m *MGO_Pipe) AllowDiskUse() Pipe {
// 	return &MGO_Pipe{Pipe: m.Pipe.AllowDiskUse()}
// }
// func (m *MGO_Pipe) Batch(n int) Pipe {
// 	return &MGO_Pipe{Pipe: m.Pipe.Batch(n)}
// }
// func (m *MGO_Pipe) Explain(result interface{}) error {
// 	if Mock && chk_mock("Pipe-Explain") {
// 		return MckE
// 	}
// 	return m.Pipe.Explain(result)
// }
// func (m *MGO_Pipe) One(result interface{}) error {
// 	if Mock && chk_mock("Pipe-One") {
// 		return MckE
// 	}
// 	return m.Pipe.One(result)
// }

// type Query interface {
// 	All(result interface{}) error
// 	Apply(change tmgo.Change, result interface{}) (info *tmgo.ChangeInfo, err error)
// 	Batch(n int) Query
// 	Comment(comment string) Query
// 	Count() (n int, err error)
// 	Distinct(key string, result interface{}) error
// 	Explain(result interface{}) error
// 	For(result interface{}, f func() error) error
// 	Hint(indexKey ...string) Query
// 	Iter() *tmgo.Iter
// 	Limit(n int) Query
// 	LogReplay() Query
// 	MapReduce(job *tmgo.MapReduce, result interface{}) (info *tmgo.MapReduceInfo, err error)
// 	One(result interface{}) (err error)
// 	Prefetch(p float64) Query
// 	Select(selector interface{}) Query
// 	SetMaxScan(n int) Query
// 	SetMaxTime(d time.Duration) Query
// 	Skip(n int) Query
// 	Snapshot() Query
// 	Sort(fields ...string) Query
// 	Tail(timeout time.Duration) *tmgo.Iter
// }

// type MGO_Query struct {
// 	*tmgo.Query
// }

// func (m *MGO_Query) All(result interface{}) error {
// 	if Mock && chk_mock("Query-All") {
// 		return MckE
// 	}
// 	return m.Query.All(result)
// }
// func (m *MGO_Query) Apply(change tmgo.Change, result interface{}) (info *tmgo.ChangeInfo, err error) {
// 	if Mock && chk_mock("Query-Apply") {
// 		if val, ok := MckV["Query-Apply"]; ok {
// 			return val.(*tmgo.ChangeInfo), nil
// 		} else {
// 			return nil, MckE
// 		}
// 	}
// 	return m.Query.Apply(change, result)
// }
// func (m *MGO_Query) Batch(n int) Query {
// 	m.Query.Batch(n)
// 	return m
// }
// func (m *MGO_Query) Comment(comment string) Query {
// 	m.Query.Comment(comment)
// 	return m
// }
// func (m *MGO_Query) Count() (n int, err error) {
// 	if Mock && chk_mock("Query-Count") {
// 		if val, ok := MckV["Query-Count"]; ok {
// 			return val.(int), nil
// 		} else {
// 			return 0, MckE
// 		}
// 	}
// 	return m.Query.Count()
// }
// func (m *MGO_Query) Distinct(key string, result interface{}) error {
// 	if Mock && chk_mock("Query-Distinct") {
// 		return MckE
// 	}
// 	return m.Query.Distinct(key, result)
// }
// func (m *MGO_Query) Explain(result interface{}) error {
// 	if Mock && chk_mock("Query-Explain") {
// 		return MckE
// 	}
// 	return m.Query.Explain(result)
// }
// func (m *MGO_Query) For(result interface{}, f func() error) error {
// 	if Mock && chk_mock("Query-For") {
// 		return MckE
// 	}
// 	return m.Query.For(result, f)
// }
// func (m *MGO_Query) Hint(indexKey ...string) Query {
// 	m.Query.Hint(indexKey...)
// 	return m
// }
// func (m *MGO_Query) Limit(n int) Query {
// 	m.Query.Limit(n)
// 	return m
// }
// func (m *MGO_Query) LogReplay() Query {
// 	m.Query.LogReplay()
// 	return m
// }
// func (m *MGO_Query) MapReduce(job *tmgo.MapReduce, result interface{}) (info *tmgo.MapReduceInfo, err error) {
// 	if Mock && chk_mock("Query-MapReduce") {
// 		if val, ok := MckV["Query-MapReduce"]; ok {
// 			return val.(*tmgo.MapReduceInfo), nil
// 		} else {
// 			return nil, MckE
// 		}
// 	}
// 	return m.Query.MapReduce(job, result)
// }
// func (m *MGO_Query) One(result interface{}) (err error) {
// 	if Mock && chk_mock("Query-One") {
// 		return MckE
// 	}
// 	return m.Query.One(result)
// }
// func (m *MGO_Query) Prefetch(p float64) Query {
// 	m.Query.Prefetch(p)
// 	return m
// }
// func (m *MGO_Query) Select(selector interface{}) Query {
// 	m.Query.Select(selector)
// 	return m
// }
// func (m *MGO_Query) SetMaxScan(n int) Query {
// 	m.Query.SetMaxScan(n)
// 	return m
// }
// func (m *MGO_Query) SetMaxTime(d time.Duration) Query {
// 	m.Query.SetMaxTime(d)
// 	return m
// }
// func (m *MGO_Query) Skip(n int) Query {
// 	m.Query.Skip(n)
// 	return m
// }
// func (m *MGO_Query) Snapshot() Query {
// 	m.Query.Snapshot()
// 	return m

// }
// func (m *MGO_Query) Sort(fields ...string) Query {
// 	m.Query.Sort(fields...)
// 	return m
// }
