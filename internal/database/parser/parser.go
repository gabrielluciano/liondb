package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gabrielluciano/liondb/internal/database/storage"
)

var splitCommandRegex = regexp.MustCompile("[^\\s\"']+|\"[^\"]*\"|'[^']*'")

type ParseError struct {
	msg string
}

type Id struct {
	Lower uint
	Upper uint
}

type ParsedCommand struct {
	Operation string
	Entity    string
	Id        Id
	Data      *storage.Data
}

func (err *ParseError) Error() string {
	return err.msg
}

func ParseCommand(cmd string) (*ParsedCommand, error) {
	parts, err := getParts(cmd)
	if err != nil {
		return nil, &ParseError{"Error parsing command: " + err.Error()}
	}

	entity, err := getEntity(parts[1])
	if err != nil {
		return nil, &ParseError{"Error parsing entity: " + err.Error()}
	}

	ids, err := getIds(parts[1])
	if err != nil {
		return nil, &ParseError{"Error parsing id: " + err.Error()}
	}

	data, err := getData(parts[2:])
	if err != nil {
		return nil, &ParseError{"Error parsing data: " + err.Error()}
	}

	return &ParsedCommand{
		Operation: strings.ToUpper(parts[0]),
		Entity:    entity,
		Id:        ids,
		Data:      data,
	}, nil
}

func getParts(cmd string) ([]string, error) {
	parts := splitCommandRegex.FindAllString(cmd, -1)
	if len(parts) < 2 {
		return nil, errors.New("invalid command")
	}
	return parts, nil
}

func getEntity(entityPart string) (string, error) {
	splitChar := ":"
	if strings.Contains(entityPart, "[") {
		splitChar = "["
	}
	parts := strings.Split(entityPart, splitChar)
	if parts[0] == "" {
		return "", errors.New("invalid entity")

	}
	return parts[0], nil
}

func getIds(idPart string) (Id, error) {
	if strings.Contains(idPart, "[") {
		idPart = strings.Split(idPart, "[")[1]
		idPart = strings.ReplaceAll(idPart, "[", "")
		idPart = strings.ReplaceAll(idPart, "]", "")
		ids := strings.Split(idPart, ":")
		lower, err := parseId(ids[0])
		if err != nil {
			return Id{}, nil
		}
		upper, err := parseId(ids[1])
		if err != nil {
			return Id{}, nil
		}
		return Id{Lower: lower, Upper: upper}, nil
	} else if strings.Contains(idPart, ":") {
		idString := strings.Split(idPart, ":")[1]
		id, err := parseId(idString)
		if err != nil {
			return Id{}, nil
		}
		return Id{Lower: id, Upper: id}, nil
	}
	return Id{}, nil
}

func parseId(idString string) (uint, error) {
	var id int
	if idString == "" {
		return 0, nil
	}
	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id format")
	}
	return uint(id), nil
}

func getData(parts []string) (*storage.Data, error) {
	if len(parts) == 0 {
		return nil, nil
	}

	if len(parts)%2 != 0 {
		return nil, &ParseError{"invalid data, different number of attributes and values"}
	}

	data := &storage.Data{}
	for i := 0; i < len(parts); i += 2 {
		(*data)[parts[i]] = parseDataTypes(parts[i+1])
	}
	return data, nil
}

func parseDataTypes(part string) interface{} {
	if strings.HasPrefix(part, "'") {
		return part
	}

	if strings.ToUpper(part) == "TRUE" {
		return true
	}

	if strings.ToUpper(part) == "FALSE" {
		return false
	}

	integer, err := strconv.Atoi(part)
	if err == nil {
		return integer
	}

	float, err := strconv.ParseFloat(part, 64)
	if err == nil {
		return float
	}

	return fmt.Sprintf("'%s'", part)
}
