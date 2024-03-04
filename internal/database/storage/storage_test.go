package storage

import (
	"testing"

	"github.com/gabrielluciano/liondb/internal/testutil"
)

func TestInsertRecord_ShouldReturnTrue(t *testing.T) {
	// Arrange
	r := &Record{
		Id: 1,
		Data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	personsStorage := New("persons")

	// Act
	inserted := personsStorage.InsertRecord(r)

	// Assert
	testutil.AssertTrue(t, inserted, "inserted")

	insertedRecord, found := personsStorage.GetRecord(r.Id)
	testutil.AssertTrue(t, found, "found")
	testutil.AssertEquals(t, insertedRecord, r, "insertedRecord")
}

func TestInsertRecord_ShouldReturnFalse(t *testing.T) {
	// Arrange
	r := &Record{
		Id: 1,
		Data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	personsStorage := New("persons")
	personsStorage.InsertRecord(r)

	// Act
	inserted := personsStorage.InsertRecord(r)

	// Assert
	testutil.AssertFalse(t, inserted, "inserted")
}

func TestUpdateRecord_ShouldReturnTrue(t *testing.T) {
	// Arrange
	expectedAge := 23

	original := &Record{
		Id: 1,
		Data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	replacement := &Record{
		Id: 1,
		Data: &Data{
			"age": expectedAge,
		},
	}
	personsStorage := New("persons")
	personsStorage.InsertRecord(original)

	// Act
	replaced := personsStorage.UpdateRecord(replacement)

	// Assert
	testutil.AssertTrue(t, replaced, "replaced")

	updatedRecord, found := personsStorage.GetRecord(original.Id)
	testutil.AssertTrue(t, found, "found")

	updatedAge := (*updatedRecord.Data)["age"]
	testutil.AssertEquals(t, expectedAge, updatedAge, "age")
}

func TestUpdateRecord_ShouldReturnFalse(t *testing.T) {
	// Arrange
	replacement := &Record{
		Id: 1,
		Data: &Data{
			"age": 23,
		},
	}
	personsStorage := New("persons")

	// Act
	replaced := personsStorage.UpdateRecord(replacement)

	// Assert
	testutil.AssertFalse(t, replaced, "replaced")
}

func TestGetRecord_ShouldFindRecord(t *testing.T) {
	// Arrange
	r := &Record{
		Id: 1,
		Data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	personsStorage := New("persons")
	personsStorage.InsertRecord(r)

	// Act
	foundRecord, founded := personsStorage.GetRecord(r.Id)

	// Assert
	testutil.AssertEquals(t, r, foundRecord, "foundRecord")
	testutil.AssertTrue(t, founded, "founded")
}

func TestGetRecord_ShouldNotFindRecord(t *testing.T) {
	// Arrange
	personsStorage := New("persons")

	// Act
	foundRecord, founded := personsStorage.GetRecord(1)

	// Assert
	testutil.AssertNil(t, foundRecord, "foundRecord")
	testutil.AssertFalse(t, founded, "founded")
}

func TestDeleteRecord_ShouldReturnTrue(t *testing.T) {
	// Arrange
	r := &Record{
		Id: 5,
		Data: &Data{
			"name": "Jonh",
			"age":  75,
		},
	}
	personsStorage := New("persons")
	personsStorage.InsertRecord(r)

	// Act
	deletedRecord, deleted := personsStorage.DeleteRecord(r.Id)

	// Assert
	testutil.AssertEquals(t, r, deletedRecord, "deletedRecord")
	testutil.AssertTrue(t, deleted, "deleted")
}

func TestDeleteRecord_ShouldReturnFalse(t *testing.T) {
	// Arrange
	personsStorage := New("persons")

	// Act
	_, deleted := personsStorage.DeleteRecord(1)

	// Assert
	testutil.AssertFalse(t, deleted, "deleted")
}

func TestIterateOverRecords(t *testing.T) {
	// Arrange
	data1 := &Record{
		Id: 1,
		Data: &Data{
			"name": "Mike",
			"age":  23,
		},
	}
	data2 := &Record{
		Id: 2,
		Data: &Data{
			"name": "Jane",
			"age":  34,
		},
	}
	data3 := &Record{
		Id: 3,
		Data: &Data{
			"name": "John",
			"age":  71,
		},
	}

	dataStorage := New("data")
	dataStorage.InsertRecord(data1)
	dataStorage.InsertRecord(data2)
	dataStorage.InsertRecord(data3)

	// Act
	records := make([]*Record, 0)
	dataStorage.IterateOverRecords(func(record *Record) bool {
		records = append(records, record)
		return true
	})

	// Assert
	testutil.AssertEquals(t, dataStorage.Len(), len(records), "len(records)")
}

func TestIterateOverRecordsEmptyStorage(t *testing.T) {
	// Arrange
	dataStorage := New("data")

	// Act
	records := make([]*Record, 0)
	dataStorage.IterateOverRecords(func(record *Record) bool {
		records = append(records, record)
		return true
	})

	// Assert
	testutil.AssertEquals(t, dataStorage.Len(), len(records), "len(records)")
}

func TestGetAllRecordsAscend(t *testing.T) {
	// Arrange
	data1 := &Record{
		Id: 1,
		Data: &Data{
			"name": "Mike",
			"age":  23,
		},
	}
	data2 := &Record{
		Id: 2,
		Data: &Data{
			"name": "Jane",
			"age":  34,
		},
	}
	data3 := &Record{
		Id: 3,
		Data: &Data{
			"name": "John",
			"age":  71,
		},
	}
	dataStorage := New("data")
	dataStorage.InsertRecord(data1)
	dataStorage.InsertRecord(data2)
	dataStorage.InsertRecord(data3)

	// Act
	records := dataStorage.GetAllRecords(false)

	// Assert
	testutil.AssertEquals(t, dataStorage.Len(), len(records), "len(records)")
	testutil.AssertEquals(t, (*data1.Data)["name"], (*records[0].Data)["name"], "first record name")
	testutil.AssertEquals(t, (*data3.Data)["name"], (*records[2].Data)["name"], "last record name")
}

func TestGetAllRecordsDescend(t *testing.T) {
	// Arrange
	data1 := &Record{
		Id: 1,
		Data: &Data{
			"name": "Mike",
			"age":  23,
		},
	}
	data2 := &Record{
		Id: 2,
		Data: &Data{
			"name": "Jane",
			"age":  34,
		},
	}
	data3 := &Record{
		Id: 3,
		Data: &Data{
			"name": "John",
			"age":  71,
		},
	}
	dataStorage := New("data")
	dataStorage.InsertRecord(data1)
	dataStorage.InsertRecord(data2)
	dataStorage.InsertRecord(data3)

	// Act
	records := dataStorage.GetAllRecords(true)

	// Assert
	testutil.AssertEquals(t, dataStorage.Len(), len(records), "len(records)")
	testutil.AssertEquals(t, (*data3.Data)["name"], (*records[0].Data)["name"], "first record name")
	testutil.AssertEquals(t, (*data1.Data)["name"], (*records[2].Data)["name"], "last record name")
}

func TestGetAllRecordsEmptyStorage(t *testing.T) {
	// Arrange
	dataStorage := New("data")

	// Act
	records := dataStorage.GetAllRecords(false)

	// Assert
	testutil.AssertEquals(t, dataStorage.Len(), len(records), "len(records)")
}
