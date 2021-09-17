package tuple

import (
	"testing"
)

func TestMarshalsAndUnMarshalsATupleIdContainingPageId(t *testing.T) {
	tupleId := TupleId{
		PageId: 100,
		SlotNo: 10,
	}
	buffer := tupleId.MarshalBinary()

	unmarshalledTupleId := &TupleId{}
	unmarshalledTupleId.UnMarshalBinary(buffer)

	expected := uint32(100)
	if unmarshalledTupleId.PageId != expected {
		t.Fatalf("Expected pageid in unmarshalled tuple id to be %v, received %v", expected, unmarshalledTupleId.PageId)
	}
}

func TestMarshalsAndUnMarshalsATupleIdContainingSlotNo(t *testing.T) {
	tupleId := TupleId{
		PageId: 100,
		SlotNo: 10,
	}
	buffer := tupleId.MarshalBinary()

	unmarshalledTupleId := &TupleId{}
	unmarshalledTupleId.UnMarshalBinary(buffer)

	expected := 10
	if unmarshalledTupleId.SlotNo != expected {
		t.Fatalf("Expected slotno in unmarshalled tuple id to be %v, received %v", expected, unmarshalledTupleId.SlotNo)
	}
}
