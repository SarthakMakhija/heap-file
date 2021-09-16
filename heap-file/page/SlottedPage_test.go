package page

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"testing"
)

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingPageId(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	tupleId := slottedPage.Put(aTuple)
	expectedPageId := uint32(100)

	if expectedPageId != tupleId.PageId {
		t.Fatalf("Expected Page Id in tuple id to be %v, received %v", expectedPageId, tupleId.PageId)
	}
}

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingSlotNo(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	tupleId := slottedPage.Put(aTuple)
	expectedSlotNo := 1

	if expectedSlotNo != tupleId.SlotNo {
		t.Fatalf("Expected slot no in tuple id to be %v, received %v", expectedSlotNo, tupleId.SlotNo)
	}
}

func TestPutsATupleInASlottedPageAndReadsItBack(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	slottedPage.Put(aTuple)
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

func TestReturnsTheSizeAvailableInAPage(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	slottedPage.Put(aTuple)

	availableSize := slottedPage.SizeAvailable()
	expectedSize := uint16(4096) - uint16(pageIdSize) - uint16(slotSize) - uint16(aTuple.Size())

	if availableSize != expectedSize {
		t.Fatalf("Expected page available size to be %v, recevied %v", availableSize, expectedSize)
	}
}
