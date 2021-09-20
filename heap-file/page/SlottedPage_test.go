package page

import (
	"github.com/SarthakMakhija/heap-file/heap-file/field"
	"github.com/SarthakMakhija/heap-file/heap-file/tuple"
	"os"
	"testing"
)

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingPageId(t *testing.T) {
	slottedPage := NewSlottedPage(100, os.Getpagesize(), twoFieldTestTupleDescriptor)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint32Field(3000))

	tupleId := slottedPage.Put(aTuple.MarshalBinary())
	expectedPageId := uint32(100)

	if expectedPageId != tupleId.PageId {
		t.Fatalf("Expected Page Id in tuple id to be %v, received %v", expectedPageId, tupleId.PageId)
	}
}

func TestPutsATupleInASlottedPageAndReturnsTupleIdContainingSlotNo(t *testing.T) {
	slottedPage := NewSlottedPage(100, os.Getpagesize(), twoFieldTestTupleDescriptor)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint32Field(3000))

	tupleId := slottedPage.Put(aTuple.MarshalBinary())
	expectedSlotNo := 1

	if expectedSlotNo != tupleId.SlotNo {
		t.Fatalf("Expected slot no in tuple id to be %v, received %v", expectedSlotNo, tupleId.SlotNo)
	}
}

func TestPutsATupleInASlottedPageAndReadsItBack(t *testing.T) {
	slottedPage := NewSlottedPage(100, os.Getpagesize(), twoFieldTestTupleDescriptor)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint32Field(3000))

	slottedPage.Put(aTuple.MarshalBinary())
	readTuple := slottedPage.GetAt(1)

	stringFieldValue := readTuple.AllFields()[0].Value()
	expectedStringFieldValue := "Database Systems"

	if stringFieldValue != expectedStringFieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
	}

	uint32FieldValue := readTuple.AllFields()[1].Value()
	expectedUint32FieldValue := uint32(3000)

	if uint32FieldValue != expectedUint32FieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedUint32FieldValue, uint32FieldValue)
	}
}

func TestReturnsTheSizeAvailableInAPage(t *testing.T) {
	slottedPage := NewSlottedPage(100, os.Getpagesize(), twoFieldTestTupleDescriptor)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint32Field(3000))

	slottedPage.Put(aTuple.MarshalBinary())

	availableSize := slottedPage.SizeAvailable()
	expectedSize := uint16(4096) - uint16(pageIdSize) - uint16(slotSize) - uint16(aTuple.Size())

	if availableSize != expectedSize {
		t.Fatalf("Expected page available size to be %v, recevied %v", availableSize, expectedSize)
	}
}

func TestReturnsFalseGivenPageIsNotLarEnoughForTheTuple(t *testing.T) {
	slottedPage := NewSlottedPage(100, 40, twoFieldTestTupleDescriptor)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint32Field(3000)) //20 byte tuple

	slottedPage.Put(aTuple.MarshalBinary())

	isPageLargeEnough := slottedPage.HasSizeLargeEnoughToHold(aTuple.MarshalBinary())
	if isPageLargeEnough != false {
		t.Fatalf("Expected page to not be large enough for the tuple")
	}
}

func TestReturnsTrueGivenPageIsLarEnoughForTheTuple(t *testing.T) {
	slottedPage := NewSlottedPage(100, 40, twoFieldTestTupleDescriptor)

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint32Field(3000)) //20 byte tuple

	isPageLargeEnough := slottedPage.HasSizeLargeEnoughToHold(aTuple.MarshalBinary())
	if isPageLargeEnough != true {
		t.Fatalf("Expected page to be large enough for the tuple")
	}
}
