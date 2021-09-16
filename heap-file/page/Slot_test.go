package page

import "testing"

func TestMarshalsAndUnMarshalsASlotWithTupleOffset(t *testing.T) {
	slot := Slot{
		tupleOffset: 4070,
		tupleSize:   15,
	}
	buffer := slot.MarshalBinary()

	readSlot := &Slot{}
	readSlot.UnMarshalBinary(buffer, 0)
	expectedTupleOffset := uint16(4070)

	if expectedTupleOffset != readSlot.tupleOffset {
		t.Fatalf("Expected tuple offset from slot is %v, received %v", expectedTupleOffset, slot.tupleOffset)
	}
}

func TestMarshalsAndUnMarshalsASlotWithTupleSize(t *testing.T) {
	slot := Slot{
		tupleOffset: 4070,
		tupleSize:   15,
	}
	buffer := slot.MarshalBinary()

	readSlot := &Slot{}
	readSlot.UnMarshalBinary(buffer, 0)
	expectedTupleSize := uint16(15)

	if expectedTupleSize != readSlot.tupleSize {
		t.Fatalf("Expected tuple size from slot is %v, received %v", expectedTupleSize, slot.tupleSize)
	}
}
