package page

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"testing"
)

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingPageId(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	tuple := heap_file.NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(3000))

	tupleId := slottedPage.Put(tuple)
	expectedPageId := uint32(100)

	if expectedPageId != tupleId.PageId {
		t.Fatalf("Expected Page Id in tuple id to be %v, received %v", expectedPageId, tupleId.PageId)
	}
}

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingSlotNo(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	tuple := heap_file.NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(3000))

	tupleId := slottedPage.Put(tuple)
	expectedSlotNo := 1

	if expectedSlotNo != tupleId.SlotNo {
		t.Fatalf("Expected slot no in tuple id to be %v, received %v", expectedSlotNo, tupleId.SlotNo)
	}
}

func TestPutsATupleInASlottedPageAndReadsItBack(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	tuple := heap_file.NewTuple()
	tuple.AddField(field.NewStringField("Database Systems"))
	tuple.AddField(field.NewUint16Field(3000))

	slottedPage.Put(tuple)
	readTuple := slottedPage.Get(1)

	stringFieldValue := readTuple.AllFields()[0].Value()
	expectedStringFieldValue := "Database Systems"

	if stringFieldValue != expectedStringFieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
	}

	uint16FieldValue := readTuple.AllFields()[1].Value()
	expectedUint16FieldValue := uint16(3000)

	if uint16FieldValue != expectedUint16FieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
	}
}
