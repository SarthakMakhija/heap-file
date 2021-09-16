package heap_file

import (
	"testing"
)

func TestShouldMarshalAndUnmarshalAStringField(t *testing.T) {
	var field = NewStringField("Database Storage Systems")
	fieldBuffer := field.MarshalBinary()

	unmarshalledField := &StringField{}
	unmarshalledField.UnMarshalBinary(fieldBuffer)

	expected := "Database Storage Systems"
	if expected != unmarshalledField.Value() {
		t.Fatalf("Expected value after unmarshalling %v, received %v", expected, unmarshalledField.Value())
	}
}
