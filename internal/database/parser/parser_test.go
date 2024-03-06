package parser

import (
	"testing"

	"github.com/gabrielluciano/liondb/internal/database/storage"
	"github.com/gabrielluciano/liondb/internal/testutil"
)

func TestGetParts(t *testing.T) {
	// Arrange
	expectedLen := 6
	command := "NEW client:1 name 'John Sanchez' age 17"

	// Act
	parts, err := getParts(command)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, expectedLen, len(parts), "len(parts)")
	testutil.AssertEquals(t, "'John Sanchez'", parts[3], "name")
	testutil.AssertEquals(t, "17", parts[5], "age")
}

func TestGetPartsEmptyString(t *testing.T) {
	// Arrange
	expectedLen := 0
	command := ""

	// Act
	parts, err := getParts(command)

	// Assert
	testutil.AssertNotNil(t, err, "error")
	testutil.AssertEquals(t, expectedLen, len(parts), "len(parts)")
}

func TestGetEntityNoId(t *testing.T) {
	// Arrange
	part := "client"

	// Act
	entity, err := getEntity(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, "client", entity, "entity")
}

func TestGetEntitySingleId(t *testing.T) {
	// Arrange
	part := "car:2"

	// Act
	entity, err := getEntity(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, "car", entity, "entity")
}

func TestGetEntityRangeId(t *testing.T) {
	// Arrange
	part := "car[1:10]"

	// Act
	entity, err := getEntity(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, "car", entity, "entity")
}

func TestGetEntityEmptyArgument(t *testing.T) {
	// Arrange
	part := ""

	// Act
	entity, err := getEntity(part)

	// Assert
	testutil.AssertNotNil(t, err, "error")
	testutil.AssertEquals(t, "", entity, "entity")
}

func TestGetIdSingleId(t *testing.T) {
	// Arrange
	part := "cliente:1"

	// Act
	ids, err := getIds(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, uint(1), ids.Lower, "lower")
	testutil.AssertEquals(t, uint(1), ids.Upper, "upper")
}

func TestGetIdNoId(t *testing.T) {
	// Arrange
	part := "cliente"

	// Act
	ids, err := getIds(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, uint(0), ids.Lower, "lower")
	testutil.AssertEquals(t, uint(0), ids.Upper, "upper")
}

func TestGetIdRangeMinToValue(t *testing.T) {
	// Arrange
	part := "cliente[:10]"

	// Act
	ids, err := getIds(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, uint(0), ids.Lower, "lower")
	testutil.AssertEquals(t, uint(10), ids.Upper, "upper")
}

func TestGetIdRangeValueToMax(t *testing.T) {
	// Arrange
	part := "cliente[10:]"

	// Act
	ids, err := getIds(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, uint(10), ids.Lower, "lower")
	testutil.AssertEquals(t, uint(0), ids.Upper, "upper")
}

func TestGetIdRangeValueToValue(t *testing.T) {
	// Arrange
	part := "cliente[10:35]"

	// Act
	ids, err := getIds(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, uint(10), ids.Lower, "lower")
	testutil.AssertEquals(t, uint(35), ids.Upper, "upper")
}

func TestGetIdRangeEmpty(t *testing.T) {
	// Arrange
	part := "cliente[:]"

	// Act
	ids, err := getIds(part)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, uint(0), ids.Lower, "lower")
	testutil.AssertEquals(t, uint(0), ids.Upper, "upper")
}

func TestGetDataBooleanTrue(t *testing.T) {
	// Arrange
	parts := []string{"key", "true"}

	// Act
	data, err := getData(parts)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, (*data)["key"], true, "value")
}

func TestGetDataBooleanFalse(t *testing.T) {
	// Arrange
	parts := []string{"key", "false"}

	// Act
	data, err := getData(parts)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, (*data)["key"], false, "value")
}

func TestGetDataInteger(t *testing.T) {
	// Arrange
	parts := []string{"key", "10"}

	// Act
	data, err := getData(parts)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, (*data)["key"], 10, "value")
}

func TestGetDataFloat(t *testing.T) {
	// Arrange
	parts := []string{"key", "10.10"}

	// Act
	data, err := getData(parts)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, (*data)["key"], 10.10, "value")
}

func TestGetDataString(t *testing.T) {
	// Arrange
	parts := []string{"key", "'John Silva'"}

	// Act
	data, err := getData(parts)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, (*data)["key"], "'John Silva'", "value")
}

func TestGetDataStringDefault(t *testing.T) {
	// Arrange
	parts := []string{"key", "word"}

	// Act
	data, err := getData(parts)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, (*data)["key"], "'word'", "value")
}

func TestParseCommandSingleArgument(t *testing.T) {
	command := "NEW cliente:1 name 'Mary'"

	testParseCommand_ShouldParse(command, "NEW",
		"cliente",
		Id{Lower: 1, Upper: 1},
		storage.Data{
			"name": "'Mary'",
		}, t)
}

func TestParseCommandMultipleArguments(t *testing.T) {
	command := "UPD product[1:10] category condiment price 17.8"

	testParseCommand_ShouldParse(command, "UPD",
		"product",
		Id{Lower: 1, Upper: 10},
		storage.Data{
			"category": "'condiment'",
			"price":    17.8,
		}, t)
}

func testParseCommand_ShouldParse(cmd string, expectedOperation string, expectedEntity string,
	expectedId Id, expectedData storage.Data, t *testing.T) {
	// Act
	parsedCommand, err := ParseCommand(cmd)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertEquals(t, expectedOperation, parsedCommand.Operation, "operation")
	testutil.AssertEquals(t, expectedEntity, parsedCommand.Entity, "entity")
	testutil.AssertEquals(t, expectedId.Lower, parsedCommand.Id.Lower, "id")
	testutil.AssertEquals(t, expectedId.Upper, parsedCommand.Id.Upper, "id")
	if expectedData != nil {
		for key, value := range *parsedCommand.Data {
			testutil.AssertEquals(t, expectedData[key], value, key)
		}
	}
}

func TestParseCommandInvalidNumberOfParts(t *testing.T) {
	command := "UPD product:1 category"

	testParseCommand_ShouldError(command, t)
}

func TestParseCommandInsufficientParts(t *testing.T) {
	command := "UPD"

	testParseCommand_ShouldError(command, t)
}

func testParseCommand_ShouldError(cmd string, t *testing.T) {
	// Act
	_, err := ParseCommand(cmd)

	// Assert
	testutil.AssertNotNil(t, err, "error")

	_, ok := err.(*ParseError)
	testutil.AssertTrue(t, ok, "ParseError")
}
