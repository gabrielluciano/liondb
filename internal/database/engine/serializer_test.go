package engine

import (
	"testing"

	"github.com/gabrielluciano/liondb/internal/database/storage"
	"github.com/gabrielluciano/liondb/internal/testutil"
)

func TestSerializeRecord(t *testing.T) {
	// Arrange
	record := &storage.Record{
		Id: 1,
		Data: &storage.Data{
			"name":   "'John Silva'",
			"age":    45,
			"weight": 75.8,
			"smoker": false,
		},
	}

	// Act
	serialized, err := SerializeRecord(record)

	// Assert
	testutil.AssertNil(t, err, "error")
	testutil.AssertContains(t, string(serialized), "name 'John Silva'")
	testutil.AssertContains(t, string(serialized), "age 45")
	testutil.AssertContains(t, string(serialized), "weight 75.8")
	testutil.AssertContains(t, string(serialized), "smoker false")
}
