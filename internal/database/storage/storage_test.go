package storage

import (
	"sync"
	"testing"
)

func TestInsertRecord(t *testing.T) {
	r := &Record{
		id: 1,
		data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}

	personsStorage := New("persons")
	replacedRecord, replaced := personsStorage.InsertOrUpdateRecord(r)

	if replacedRecord != nil {
		t.Errorf("Incorrect result, expected replacedRecord to be <nil>, got: %v", replacedRecord)
	}

	if replaced {
		t.Errorf("Incorrect result, expected replaced to be false, got: %v", replaced)
	}
}

func TestUpdateRecord(t *testing.T) {
	original := &Record{
		id: 1,
		data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	replacement := &Record{
		id: 1,
		data: &Data{
			"name": "Mark",
			"age":  23,
		},
	}
	personsStorage := New("persons")
	personsStorage.InsertOrUpdateRecord(original)

	replacedRecord, replaced := personsStorage.InsertOrUpdateRecord(replacement)

	if replacedRecord != original {
		t.Errorf("Incorrect result, expected replacedRecord to be %v, got: %v", original, replacedRecord)
	}

	if !replaced {
		t.Errorf("Incorrect result, expected replaced to be true, got: %v", replaced)
	}
}

func TestGetRecordFoundRecord(t *testing.T) {
	r := &Record{
		id: 1,
		data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	personsStorage := New("persons")
	personsStorage.InsertOrUpdateRecord(r)

	foundRecord, founded := personsStorage.GetRecord(r.id)

	if foundRecord != r {
		t.Errorf("Incorrect result, expected foundRecord to be %v, got: %v", r, foundRecord)
	}

	if !founded {
		t.Errorf("Incorrect result, expected founded to be true, got: %v", founded)
	}
}

func TestGetRecordNotFoundRecord(t *testing.T) {
	personsStorage := New("persons")

	foundRecord, founded := personsStorage.GetRecord(1)

	if foundRecord != nil {
		t.Errorf("Incorrect result, expected foundRecord to be <nil>, got: %v", foundRecord)
	}

	if founded {
		t.Errorf("Incorrect result, expected founded to be false, got: %v", founded)
	}
}

func TestDeleteRecord(t *testing.T) {
	r := &Record{
		id: 5,
		data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	personsStorage := New("persons")
	personsStorage.InsertOrUpdateRecord(r)

	deletedRecord, deleted := personsStorage.DeleteRecord(r.id)

	if deletedRecord != r {
		t.Errorf("Incorrect result, expected deletedRecord to be %v, got: %v", r, deletedRecord)
	}

	if !deleted {
		t.Errorf("Incorrect result, expected deleted to be true, got: %v", deleted)
	}
}

func TestUpdateRecordConcurrently(t *testing.T) {
	finalBalance := 500
	client := &Record{
		id: 5,
		data: &Data{
			"name":    "Jonh",
			"balance": 0,
		},
	}

	clientsStorage := New("persons")
	clientsStorage.InsertOrUpdateRecord(client)

	var wg sync.WaitGroup
	for i := 1; i <= finalBalance; i++ {
		wg.Add(1)
		go func() {
			wg.Done()

			client.mu.Lock()
			balance, err := client.data.GetInt("balance")
			if err != nil {
				t.Error("Test failed", err)
			}
			client.data.SetInt("balance", balance+1)
			client.mu.Unlock()
		}()
	}
	wg.Wait()

	foundRecord, founded := clientsStorage.GetRecord(client.id)

	if !founded {
		t.Errorf("Incorrect result, expected founded to be true, got: %v", founded)
	}

	balance, _ := foundRecord.data.GetInt("balance")

	if balance != finalBalance {
		t.Errorf("Incorrect result, expected balance to be %v, got: %v", finalBalance, balance)
	}
}

func TestIterateOverValues(t *testing.T) {
	data1 := &Record{
		id: 1,
		data: &Data{
			"name": "Mike",
			"age":  23,
		},
	}
	data2 := &Record{
		id: 2,
		data: &Data{
			"name": "Jane",
			"age":  34,
		},
	}
	data3 := &Record{
		id: 3,
		data: &Data{
			"name": "John",
			"age":  71,
		},
	}

	dataStorage := New("data")
	dataStorage.InsertOrUpdateRecord(data1)
	dataStorage.InsertOrUpdateRecord(data2)
	dataStorage.InsertOrUpdateRecord(data3)

	records := make([]*Record, 0)

	dataStorage.IterateOverRecords(func(record *Record) bool {
		records = append(records, record)
		return true
	})

	if len(records) != dataStorage.Len() {
		t.Errorf("Incorrect result, expected len(records) to be %v, got: %v", dataStorage.Len(), len(records))
	}
}
