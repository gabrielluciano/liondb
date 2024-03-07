package engine

import (
	"fmt"

	"github.com/gabrielluciano/liondb/internal/database/parser"
	"github.com/gabrielluciano/liondb/internal/database/server"
	"github.com/gabrielluciano/liondb/internal/database/storage"
)

var storages map[string]*storage.Storage

func Start() {
	initializeStorage()
	initializeServer()
}

func initializeStorage() {
	storages = make(map[string]*storage.Storage)
}

func initializeServer() {
	server := server.New("7123")
	server.SetMessageHandler(messageHandler)
	server.Listen()
}

func messageHandler(command string) []byte {
	parsedCommand, err := parser.ParseCommand(command)
	if err != nil {
		return []byte(fmt.Sprintf("Error processing command: %v", err))
	}

	_, ok := storages[parsedCommand.Entity]
	if !ok {
		storages[parsedCommand.Entity] = storage.New(parsedCommand.Entity)
	}
	return executeOperation(parsedCommand)
}

func executeOperation(parsedCommand *parser.ParsedCommand) []byte {
	switch parsedCommand.Operation {
	case "NEW":
		return insertRecord(parsedCommand)
	case "UPD":
		return updateRecord(parsedCommand)
	case "GET":
		return getRecords(parsedCommand)
	case "DEL":
		return deleteRecord(parsedCommand)
	default:
		return []byte("Error processing command: invalid operation")
	}
}

func insertRecord(parsedCommand *parser.ParsedCommand) []byte {
	s := storages[parsedCommand.Entity]
	if parsedCommand.Id.Lower != parsedCommand.Id.Upper || parsedCommand.Id.Lower == uint(0) {
		return []byte("Error processing command: invalid id")
	}
	inserted := s.InsertRecord(&storage.Record{
		Id:   parsedCommand.Id.Lower,
		Data: parsedCommand.Data,
	})
	if !inserted {
		return []byte("0")
	}
	return []byte("1")
}

func updateRecord(parsedCommand *parser.ParsedCommand) []byte {
	s := storages[parsedCommand.Entity]
	if parsedCommand.Id.Lower != parsedCommand.Id.Upper || parsedCommand.Id.Lower == uint(0) {
		return []byte("Error processing command: invalid id")
	}
	updated := s.UpdateRecord(&storage.Record{
		Id:   parsedCommand.Id.Lower,
		Data: parsedCommand.Data,
	})
	if !updated {
		return []byte("0")
	}
	return []byte("1")
}

func getRecords(parsedCommand *parser.ParsedCommand) []byte {
	s := storages[parsedCommand.Entity]
	if parsedCommand.Id.Lower == 0 && parsedCommand.Id.Upper == 0 {
		return getAllRecords(s)
	}
	if parsedCommand.Id.Lower == parsedCommand.Id.Upper {
		return getSingleRecord(parsedCommand.Id.Lower, s)
	}
	return []byte("Internal error")
}

func getSingleRecord(id uint, s *storage.Storage) []byte {
	record, found := s.GetRecord(id)
	if !found {
		return []byte("0")
	}

	response, err := SerializeRecord(record)
	if err != nil {
		return []byte(err.Error())
	}
	return response
}

func getAllRecords(s *storage.Storage) []byte {
	records := s.GetAllRecords(false)
	if len(records) == 0 {
		return []byte("0")
	}
	response := make([]byte, 0)
	for _, record := range records {
		serialized, err := SerializeRecord(record)
		if err != nil {
			return []byte("Error deserializing data")
		}
		response = append(response, serialized...)
		response = append(response, '\n')
	}
	return response[:len(response)-1]
}

func deleteRecord(parsedCommand *parser.ParsedCommand) []byte {
	if parsedCommand.Id.Lower == 0 && parsedCommand.Id.Upper == 0 {
		return []byte("Error processing command: invalid id")
	}
	if parsedCommand.Id.Lower != parsedCommand.Id.Upper {
		return []byte("Error processing command: invalid id")
	}
	s := storages[parsedCommand.Entity]
	_, deleted := s.DeleteRecord(parsedCommand.Id.Lower)
	if !deleted {
		return []byte("0")
	}

	return []byte("1")
}
