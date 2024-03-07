package engine

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/gabrielluciano/liondb/internal/database/storage"
)

type SerializationError struct {
	message string
}

func (err *SerializationError) Error() string {
	return fmt.Sprintf("Serialization error: %v", err.message)
}

func SerializeRecord(record *storage.Record) ([]byte, error) {
	buffer := bytes.Buffer{}
	buffer.WriteString("id ")
	buffer.WriteString(fmt.Sprint(record.Id))
	for key, value := range *record.Data {
		buffer.WriteString(" ")
		buffer.WriteString(key)
		buffer.WriteString(" ")
		serializedValue, err := SerializeValue(value)
		if err != nil {
			return []byte{}, &SerializationError{message: err.Error()}
		}
		buffer.WriteString(serializedValue)
	}
	return buffer.Bytes(), nil
}

func SerializeValue(value interface{}) (string, error) {
	var response string
	var err error
	switch v := value.(type) {
	case float64:
		response = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		response = strconv.Itoa(v)
	case bool:
		response = strconv.FormatBool(v)
	case string:
		response = v
	default:
		err = errors.New("error serializing value")
	}
	return response, err
}
