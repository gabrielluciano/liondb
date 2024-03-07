package engine

import (
	"testing"

	"github.com/gabrielluciano/liondb/internal/database/parser"
	"github.com/gabrielluciano/liondb/internal/database/storage"
	"github.com/gabrielluciano/liondb/internal/testutil"
)

func TestInsertRecord(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "NEW",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
		Data: &storage.Data{
			"name": "bmw",
		},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertEquals(t, "1", string(result), "result")
}

func TestInsertRecordDuplicated(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "NEW",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
		Data: &storage.Data{
			"name": "bmw",
		},
	}
	executeOperation(parsedCommand)

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertEquals(t, "0", string(result), "result")
}

func TestInsertRecordInvalidId(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "NEW",
		Entity:    "car",
		Id:        parser.Id{Lower: 0, Upper: 1},
		Data: &storage.Data{
			"name": "bmw",
		},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertContains(t, string(result), "invalid id")
}

func TestUpdateRecord(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "NEW",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
		Data: &storage.Data{
			"name": "bmw",
		},
	}
	executeOperation(parsedCommand)

	updateParsedCommand := &parser.ParsedCommand{
		Operation: "UPD",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
		Data: &storage.Data{
			"name": "mercedes",
		},
	}

	// Act
	result := executeOperation(updateParsedCommand)

	// Assert
	testutil.AssertEquals(t, "1", string(result), "result")
}

func TestUpdateInexistentRecord(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	updateParsedCommand := &parser.ParsedCommand{
		Operation: "UPD",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
		Data: &storage.Data{
			"name": "mercedes",
		},
	}

	// Act
	result := executeOperation(updateParsedCommand)

	// Assert
	testutil.AssertEquals(t, "0", string(result), "result")
}

func TestUpdateInvalidId(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	updateParsedCommand := &parser.ParsedCommand{
		Operation: "UPD",
		Entity:    "car",
		Id:        parser.Id{Lower: 0, Upper: 1},
		Data: &storage.Data{
			"name": "mercedes",
		},
	}

	// Act
	result := executeOperation(updateParsedCommand)

	// Assert
	testutil.AssertContains(t, string(result), "invalid id")
}

func TestGetSingleRecord(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	storages["car"].InsertRecord(&storage.Record{
		Id: 1,
		Data: &storage.Data{
			"name": "'mercedes'",
		},
	})
	parsedCommand := &parser.ParsedCommand{
		Operation: "GET",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertContains(t, string(result), "id 1")
	testutil.AssertContains(t, string(result), "name 'mercedes'")
}

func TestGetSingleRecordEmpty(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "GET",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertEquals(t, "0", string(result), "result")
}

func TestGetAllRecords(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	storages["car"].InsertRecord(&storage.Record{
		Id: 1,
		Data: &storage.Data{
			"name": "'mercedes'",
		},
	})
	storages["car"].InsertRecord(&storage.Record{
		Id: 2,
		Data: &storage.Data{
			"name": "'bmw'",
		},
	})
	parsedCommand := &parser.ParsedCommand{
		Operation: "GET",
		Entity:    "car",
		Id:        parser.Id{Lower: 0, Upper: 0},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertContains(t, string(result), "id 1")
	testutil.AssertContains(t, string(result), "name 'mercedes'")
	testutil.AssertContains(t, string(result), "id 2")
	testutil.AssertContains(t, string(result), "name 'bmw'")
}

func TestGetAllRecordsEmpty(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "GET",
		Entity:    "car",
		Id:        parser.Id{Lower: 0, Upper: 0},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertEquals(t, "0", string(result), "result")
}

func TestDeleteRecord(t *testing.T) {
	initializeStorage()
	storages["car"] = storage.New("car")
	storages["car"].InsertRecord(&storage.Record{
		Id: 1,
		Data: &storage.Data{
			"name": "'mercedes'",
		},
	})
	parsedCommand := &parser.ParsedCommand{
		Operation: "DEL",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertEquals(t, "1", string(result), "result")
	_, found := storages["car"].GetRecord(1)
	testutil.AssertFalse(t, found, "found")
}

func TestDeleteRecordEmpty(t *testing.T) {
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "DEL",
		Entity:    "car",
		Id:        parser.Id{Lower: 1, Upper: 1},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertEquals(t, "0", string(result), "result")
}

func TestDeleteRecordInvalidId(t *testing.T) {
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "DEL",
		Entity:    "car",
		Id:        parser.Id{Lower: 0, Upper: 1},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertContains(t, string(result), "invalid id")
}

func TestInvalidOperation(t *testing.T) {
	// Arrange
	initializeStorage()
	storages["car"] = storage.New("car")
	parsedCommand := &parser.ParsedCommand{
		Operation: "APT",
		Entity:    "car",
		Id:        parser.Id{Lower: 0, Upper: 1},
	}

	// Act
	result := executeOperation(parsedCommand)

	// Assert
	testutil.AssertContains(t, string(result), "invalid operation")
}

func TestMessageHandler(t *testing.T) {
	// Arrange
	initializeStorage()

	// Act
	result := messageHandler("GET cliente:1")

	// Assert
	testutil.AssertEquals(t, string(result), "0", "result")
}

func TestMessageHandlerInvalidCommand(t *testing.T) {
	// Arrange
	initializeStorage()

	// Act
	result := messageHandler("GET cliente abc")

	// Assert
	testutil.AssertContains(t, string(result), "Error processing command")
}
