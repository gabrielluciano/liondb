package storage

import (
	"sync"

	"github.com/google/btree"
)

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
		name: name,
		records: *btree.NewG(64, func(a, b *Record) bool {
			return a.Id < b.Id
		}),
	}
}

func (s *Storage) InsertOrUpdateRecord(r *Record) (*Record, bool) {
	return s.records.ReplaceOrInsert(r)
}

func (s *Storage) GetRecord(id uint) (*Record, bool) {
	return s.records.Get(&Record{Id: id})
}

func (s *Storage) IterateOverRecords(iterator func(record *Record) bool) {
	s.records.Ascend(iterator)
}

func (s *Storage) DeleteRecord(id uint) (*Record, bool) {
	return s.records.Delete(&Record{Id: id})
}

func (s *Storage) Len() int {
	return s.records.Len()
}
