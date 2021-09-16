package field

import "testing"

func TestShouldMarshalAndUnmarshalAnUint16Field(t *testing.T) {
	var field = NewUint16Field(1000)
	fieldBuffer := field.MarshalBinary()

	uint16FieldType := &Uint16FieldType{}
	uint16Field := uint16FieldType.UnMarshalBinary(fieldBuffer).(Uint16Field)

	expected := uint16(1000)
	if expected != uint16Field.Value() {
		t.Fatalf("Expected value after unmarshalling %v, received %v", expected, uint16Field.Value())
	}
}
