package tuple

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"testing"
)

func TestMarshalsAndUnMarshalsATupleWithSingleFieldOfTypeString(t *testing.T) {
	tuple := NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))

	marshalledTuple := tuple.MarshalBinary()
	fieldTypes := []field.FieldType{field.StringFieldType{}}

	tuple.UnMarshalBinary(marshalledTuple.Buffer(), fieldTypes)
	fieldValue := tuple.fields[0].Value()
	expected := "Database Systems"

	if fieldValue != expected {
		t.Fatalf("Expected field value to be %v, received %v", expected, fieldValue)
	}
}

func TestMarshalsAndUnMarshalsATupleWithStringAndUint16Field(t *testing.T) {
	tuple := NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(3000))

	marshalledTuple := tuple.MarshalBinary()
	fieldTypes := []field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}}

	unmarshalledTuple := NewTuple()
	unmarshalledTuple.UnMarshalBinary(marshalledTuple.Buffer(), fieldTypes)

	stringFieldValue := unmarshalledTuple.fields[0].Value()
	expectedStringFieldValue := "Database Systems"

	if stringFieldValue != expectedStringFieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
	}

	uint16FieldValue := unmarshalledTuple.fields[1].Value()
	expectedUint16FieldValue := uint16(3000)

	if uint16FieldValue != expectedUint16FieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
	}
}

func TestReturnsTheLastFieldAsTheKeyField(t *testing.T) {
	tuple := NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(100))

	keyField := tuple.KeyField()
	expected := uint16(100)

	if keyField.Value() != expected {
		t.Fatalf("Expected field value to be %v, received %v", expected, keyField.Value())
	}
}
