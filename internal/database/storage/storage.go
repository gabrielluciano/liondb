package storage

import (
	"sync"

	"github.com/google/btree"
)

type Data map[string]interface{}

type Record struct {
	Mu   sync.Mutex
	Id   uint
	Data *Data
}

type Storage struct {
	name    string
	records btree.BTreeG[*Record]
}

func New(name string) *Storage {
	return &Storage{
		name:    name,
		records: *btree.NewG(64, lessFunc),
	}
}

func lessFunc(a, b *Record) bool {
	return a.Id < b.Id
}

func (s *Storage) GetRecord(id uint) (*Record, bool) {
	return s.records.Get(&Record{Id: id})
}

func (s *Storage) IterateOverRecords(iterator func(record *Record) bool) {
	s.records.Ascend(iterator)
}

func (s *Storage) GetAllRecords(descend bool) []*Record {
	records := make([]*Record, 0, s.Len())
	if descend {
		s.records.Descend(func(record *Record) bool {
			records = append(records, record)
			return true
		})
	} else {
		s.records.Ascend(func(record *Record) bool {
			records = append(records, record)
			return true
		})
	}
	return records
}

func (s *Storage) InsertRecord(r *Record) bool {
	_, found := s.GetRecord(r.Id)
	if found {
		return false
	}
	s.records.ReplaceOrInsert(r)
	return true
}

func (s *Storage) UpdateRecord(r *Record) bool {
	savedRecord, found := s.GetRecord(r.Id)
	if !found {
		return false
	}
	savedRecord.Mu.Lock()
	defer savedRecord.Mu.Unlock()
	for k, v := range *r.Data {
		(*savedRecord.Data)[k] = v
	}
	return true
}

func (s *Storage) DeleteRecord(id uint) (*Record, bool) {
	return s.records.Delete(&Record{Id: id})
}

func (s *Storage) Len() int {
	return s.records.Len()
}
