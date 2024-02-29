package engine

import (
	"sync"

	"github.com/google/btree"
)

type Record struct {
	mu   sync.Mutex
	id   uint
	data *Data
}

type Storage struct {
	name    string
	records btree.BTreeG[*Record]
}

func (r *Record) Less(than btree.Item) bool {
	other, ok := than.(*Record)
	if !ok {
		panic("Type not allowed in this method. Only *Record is alowed.")
	}
	return r.id < other.id
}

func New(name string) *Storage {
	return &Storage{
		name: name,
		records: *btree.NewG(64, func(a, b *Record) bool {
			return a.id < b.id
		}),
	}
}

func (s *Storage) InsertOrUpdateRecord(r *Record) (*Record, bool) {
	return s.records.ReplaceOrInsert(r)
}

func (s *Storage) GetRecord(id uint) (*Record, bool) {
	return s.records.Get(&Record{id: id})
}

func (s *Storage) IterateOverRecords(iterator func(record *Record) bool) {
	s.records.Ascend(iterator)
}

func (s *Storage) DeleteRecord(id uint) (*Record, bool) {
	return s.records.Delete(&Record{id: id})
}

func (s *Storage) Len() int {
	return s.records.Len()
}
