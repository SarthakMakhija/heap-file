package heap_file

import "testing"

func TestShouldMarshalAndUnmarshalAnUint16Field(t *testing.T) {
	var field = NewUint16Field(1000)
	fieldBuffer := field.MarshalBinary()

	unmarshalledField := &Uint16Field{}
	unmarshalledField.UnMarshalBinary(fieldBuffer)

	expected := uint16(1000)
	if expected != unmarshalledField.Value() {
		t.Fatalf("Expected value after unmarshalling %v, received %v", expected, unmarshalledField.Value())
	}
}
