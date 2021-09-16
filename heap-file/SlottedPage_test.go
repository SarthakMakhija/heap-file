package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"testing"
)

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingPageId(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	tuple := NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(3000))

	tupleId := slottedPage.Put(tuple)
	expectedPageId := uint32(100)

	if expectedPageId != tupleId.pageId {
		t.Fatalf("Expected Page Id in tuple id to be %v, received %v", expectedPageId, tupleId.pageId)
	}
}

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingSlotNo(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	tuple := NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(3000))

	tupleId := slottedPage.Put(tuple)
	expectedSlotNo := 1

	if expectedSlotNo != tupleId.slotNo {
		t.Fatalf("Expected slot no in tuple id to be %v, received %v", expectedSlotNo, tupleId.slotNo)
	}
}

func TestPutsATupleInASlottedPageAndReadsItBack(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	tuple := NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(3000))

	slottedPage.Put(tuple)
	readTuple := slottedPage.Get(1)

	stringFieldValue := readTuple.fields[0].Value()
	expectedStringFieldValue := "Database Systems"

	if stringFieldValue != expectedStringFieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
	}

	uint16FieldValue := readTuple.fields[1].Value()
	expectedUint16FieldValue := uint16(3000)

	if uint16FieldValue != expectedUint16FieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
	}
}
