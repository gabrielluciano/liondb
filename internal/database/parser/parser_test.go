package parser

import (
	"fmt"
	"testing"

	"github.com/gabrielluciano/liondb/internal/database/engine"
)

func TestParseCommandSingleArgument(t *testing.T) {
	command := "NEW cliente:1 name 'Mary'"

	testParseCommand(command, "NEW", "cliente", 1, engine.Data{
		"name": "Mary",
	}, t)
}

func TestParseCommandNoArgument(t *testing.T) {
	command := "DEL cliente:1"

	testParseCommand(command, "DEL", "cliente", 1, nil, t)
}

func TestParseCommandMultipleArguments(t *testing.T) {
	command := "NEW cliente:1 name 'John Tobias' age 18"

	testParseCommand(command, "NEW", "cliente", 1, engine.Data{
		"name": "John Tobias",
		"age":  "18",
	}, t)
}

func testParseCommand(cmd string, expectedOperation string, expectedEntity string,
	expectedId uint, expectedData engine.Data, t *testing.T) {

	parsedCommand, err := ParseCommand(cmd)

	if err != nil {
		t.Errorf("Incorrect result, expected error to be <nil>, got %v", err)
	}

	if parsedCommand.operation != expectedOperation {
		t.Errorf("Incorrect result, expected operation to be %v, got %v", expectedOperation, parsedCommand.operation)
	}

	if parsedCommand.entity != expectedEntity {
		t.Errorf("Incorrect result, expected entity to be %v, got %v", expectedEntity, parsedCommand.entity)
	}

	if parsedCommand.id != expectedId {
		t.Errorf("Incorrect result, expected id to be %v, got %v", expectedId, parsedCommand.id)
	}

	if expectedData != nil {
		for key, value := range *parsedCommand.data {
			if value != expectedData[key] {
				t.Errorf("Incorrect result, expected %s to be %v, got %v", key, expectedData[key], value)
			}
		}
	}
}

func TestRebuildStrings(t *testing.T) {
	expectedName := "John Silva Santos"
	expectedSex := "male"
	expectedLength := 4
	parts := []string{"name", "'John", "Silva", "Santos'", "sex", "'male'"}

	fixedParts := rebuildStrings(parts)
	name := fixedParts[1]
	sex := fixedParts[3]

	if name != expectedName {
		t.Errorf("Incorrect result, expected name to be %v, got %v", expectedName, name)
	}

	if sex != expectedSex {
		t.Errorf("Incorrect result, expected sex to be %v, got %v", expectedSex, sex)
	}

	if len(fixedParts) != expectedLength {
		t.Errorf("Incorrect result, expected len to be %v, got %v", expectedLength, len(fixedParts))
	}

	fmt.Printf("parts: %v\n", parts)
	fmt.Printf("fixedParts: %v\n", fixedParts)
}
