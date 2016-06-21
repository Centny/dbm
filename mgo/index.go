package mgo

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	tmgo "gopkg.in/mgo.v2"
)

func ChkIdx(C func(string) *tmgo.Collection, indexes map[string]map[string]tmgo.Index) error {
	for cname, index := range indexes {
		tc := C(cname)
		if _, err := tc.Count(); err != nil {
			return err
		}
		log.D("ChkIdx checking index on collection(%v)...", cname)
		idx_l, err := tc.Indexes()
		if err != nil {
			if qerr, ok := err.(*tmgo.QueryError); !ok || qerr.Code != 26 {
				err = util.Err("ChkIdx list indexes fail with error(%v) on collection(%v)", err, cname)
				log.E("%v", err)
				return err
			}
			log.D("ChkIdx the collection(%v) is not found, it will create empty one...", cname)
			err = tc.Create(&tmgo.CollectionInfo{})
			if err != nil {
				err = util.Err("ChkIdx create collection(%v) fail with error(%v)", cname, err)
				log.E("%v", err)
				return err
			}
		}
		exists := map[string]tmgo.Index{}
		for _, idx := range idx_l {
			exists[idx.Name] = idx
		}
		for iname, idx := range index {
			if _, ok := exists[iname]; ok {
				continue
			}
			idx.Name = iname
			err = C(cname).EnsureIndex(idx)
			if err != nil {
				err = util.Err("ChkIdx ensure index by keys(%v),name(%v) fail with error(%v) on collection(%v)", idx.Key, idx.Name, err, cname)
				log.E("%v", err)
				return err
			}
			log.D("ChkIdx ensure index(%v) on collection(%v) success", iname, cname)
		}
	}
	return nil
}
