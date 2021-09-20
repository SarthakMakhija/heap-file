package field

import "testing"

func TestShouldMarshalAndUnmarshalAnUint32Field(t *testing.T) {
	var field = NewUint32Field(1000)
	fieldBuffer := field.MarshalBinary()

	uint32FieldType := &Uint32FieldType{}
	uint32Field := uint32FieldType.UnMarshalBinary(fieldBuffer).(Uint32Field)

	expected := uint32(1000)
	if expected != uint32Field.Value() {
		t.Fatalf("Expected value after unmarshalling %v, received %v", expected, uint32Field.Value())
	}
}
