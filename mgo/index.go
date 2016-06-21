package mgo

import (
	tmgo "gopkg.in/mgo.v2"
)

func ChkIdx(C func(string) *tmgo.Collection, indexes map[string]map[string]tmgo.Index) error {
	for cname, index := range indexes {
		idx_l, err := C(cname).Indexes()
		if err != nil {
			return err
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
				return err
			}
		}
	}
	return nil
}
