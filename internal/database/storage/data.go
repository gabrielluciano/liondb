package storage

import (
	"fmt"
)

type KeyNotFoundError struct {
	key string
}

type IncorrectKeyTypeError struct {
	key          string
	expectedType string
}

func (err *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Key: '%s'", err.key)
}

func (err *IncorrectKeyTypeError) Error() string {
	return fmt.Sprintf("Key: '%s' doesn't have type '%s'", err.key, err.expectedType)
}

type Data map[string]interface{}

func (d *Data) GetInt(key string) (int, error) {
	genericValue, ok := (*d)[key]

	if !ok {
		return 0, &KeyNotFoundError{key: key}
	}

	value, ok := genericValue.(int)
	if !ok {
		return 0, &IncorrectKeyTypeError{key: key, expectedType: "int"}
	}

	return value, nil
}

func (d *Data) SetInt(key string, value int) {
	(*d)[key] = value
}
