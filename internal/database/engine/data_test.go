package engine

import "testing"

func TestSetIntExistentKey(t *testing.T) {
	data := &Data{
		"key": 0,
	}

	newValue := 10
	data.SetInt("key", newValue)
	value := (*data)["key"]

	if value != newValue {
		t.Errorf("Incorrect result, expected value to be %v, got %v", newValue, value)
	}
}

func TestSetIntNonExistentKey(t *testing.T) {
	data := &Data{}

	newValue := 10
	data.SetInt("key", newValue)
	value := (*data)["key"]

	if value != newValue {
		t.Errorf("Incorrect result, expected value to be %v, got %v", newValue, value)
	}
}

func TestGetIntExistentKey(t *testing.T) {
	expected := 10
	data := &Data{
		"key": expected,
	}

	value, err := data.GetInt("key")

	if err != nil {
		t.Errorf("Incorrect result, expected error to be <nil>, got %v", err)
	}

	if value != expected {
		t.Errorf("Incorrect result, expected value to be %v, got %v", expected, value)
	}
}

func TestGetIntNonExistentKey(t *testing.T) {
	expected := 0
	data := &Data{}

	value, err := data.GetInt("key")

	err, ok := err.(*KeyNotFoundError)

	if !ok {
		t.Errorf("Incorrect result, expected error to be KeyNotFoundError, got %v", err)
	}

	if value != expected {
		t.Errorf("Incorrect result, expected value to be %v, got %v", expected, value)
	}
}

func TestGetIntDifferentKeyType(t *testing.T) {
	expected := 0
	data := &Data{"key": "text"}

	value, err := data.GetInt("key")

	err, ok := err.(*IncorrectKeyTypeError)

	if !ok {
		t.Errorf("Incorrect result, expected error to be IncorretKeyTypeError, got %v", err)
	}

	if value != expected {
		t.Errorf("Incorrect result, expected value to be %v, got %v", expected, value)
	}
}
