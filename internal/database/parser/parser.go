package parser

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gabrielluciano/liondb/internal/database/engine"
)

type ParseError struct {
	msg string
}

type ParsedCommand struct {
	operation string
	entity    string
	id        uint
	data      *engine.Data
}

func (err *ParseError) Error() string {
	return err.msg
}

func ParseCommand(cmd string) (*ParsedCommand, error) {
	parts := getParts(cmd)

	if len(parts) < 2 {
		return nil, &ParseError{"invalid command"}
	}

	operation, err := getOperation(parts)
	if err != nil {
		return nil, err
	}

	entity, err := getEntity(parts)
	if err != nil {
		return nil, err
	}

	id, err := getId(parts)
	if err != nil {
		return nil, err
	}

	data, err := getData(parts)
	if err != nil {
		return nil, err
	}

	parsedCommand := &ParsedCommand{
		operation: operation,
		entity:    entity,
		id:        id,
		data:      data,
	}

	return parsedCommand, nil
}

func getParts(cmd string) []string {
	cmd = strings.TrimSpace(cmd)
	return strings.Split(cmd, " ")
}

func getOperation(parts []string) (string, error) {
	operations := []string{"NEW", "UPD", "GET", "DEL"}
	operation := strings.ToUpper(parts[0])
	validOperation := slices.Contains(operations, operation)
	if !validOperation {
		return "", &ParseError{msg: "invalid operation"}
	}
	return operation, nil
}

func getEntity(parts []string) (string, error) {
	parts = strings.Split(parts[1], ":")
	if len(parts) < 1 {
		return "", &ParseError{"invalid entity"}
	}
	entity := parts[0]
	return entity, nil
}

func getId(parts []string) (uint, error) {
	parts = strings.Split(parts[1], ":")
	if len(parts) < 2 {
		return 0, &ParseError{"id not found in command"}
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil || id < 1 {
		return 0, &ParseError{"invalid id, must be a positive number"}
	}
	return uint(id), nil
}

func getData(parts []string) (*engine.Data, error) {
	parts = parts[2:]
	if len(parts) == 0 {
		return nil, nil
	}
	// rebuild strings previously splitted
	parts = rebuildStrings(parts)

	if len(parts)%2 != 0 {
		return nil, &ParseError{"invalid data, different number of attributes and values"}
	}
	data := &engine.Data{}
	for i := 0; i < len(parts); i += 2 {
		(*data)[parts[i]] = parts[i+1]
	}
	return data, nil
}

func rebuildStrings(parts []string) []string {
	fixedParts := make([]string, 0, len(parts))
	var i int
	for i = 0; i < len(parts); i++ {
		// If find begining of string delimited by single quotes
		if strings.HasPrefix(parts[i], "'") {
			buffer := ""
			j := i
			// while not find the end of string delimited by single quotes
			for !strings.HasSuffix(parts[j], "'") {
				buffer += parts[j] + " "
				j++
			}
			buffer += parts[j]
			buffer = strings.ReplaceAll(buffer, "'", "")
			fixedParts = append(fixedParts, buffer)
			i += j - 1
		} else {
			parts[i] = strings.ReplaceAll(parts[i], "'", "")
			fixedParts = append(fixedParts, parts[i])
		}
	}
	return fixedParts
}
