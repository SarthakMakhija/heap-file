package field

import (
	"testing"
)

func TestShouldMarshalAndUnmarshalAStringField(t *testing.T) {
	var field = NewStringField("Database Storage Systems")
	fieldBuffer := field.MarshalBinary()

	stringFieldType := &StringFieldType{}
	unmarshalledField := stringFieldType.UnMarshalBinary(fieldBuffer).(StringField)

	expected := "Database Storage Systems"
	if expected != unmarshalledField.Value() {
		t.Fatalf("Expected value after unmarshalling %v, received %v", expected, unmarshalledField.Value())
	}
}
