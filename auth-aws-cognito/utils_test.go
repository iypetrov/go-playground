package main

import (
	"testing"
)

type TestStruct struct {
	Name string
	Age  int
}

func TestCastStringToStructType(t *testing.T) {
	original := TestStruct{
		Name: "Alice",
		Age:  30,
	}

	encodedValue, err := StructToString(original)
	if err != nil {
		t.Fatalf("Failed to gob encode: %v", err)
	}

	decodedValue, err := StringToStruct(encodedValue, TestStruct{})
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	result, ok := decodedValue.(TestStruct)
	if !ok {
		t.Fatalf("Decoded value is not of type TestStruct")
	}

	if original.Name != result.Name {
		t.Errorf("Expected %v, got %v", original, result)
	}

	if original.Age != result.Age {
		t.Errorf("Expected %v, got %v", original, result)
	}
}
