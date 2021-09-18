package page

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"os"
	"strconv"
	"testing"
)

func TestPutsMultipleTuplesInASlottedPageAndReadsThemBack(t *testing.T) {
	slottedPage := NewSlottedPage(100, os.Getpagesize())

	tupleIds := add5Tuples(slottedPage)
	tuples := readTuples(tupleIds, slottedPage)

	for index, aTuple := range tuples {
		stringFieldValue := aTuple.AllFields()[0].Value()
		expectedStringFieldValue := "Database Systems" + strconv.Itoa(index)

		if stringFieldValue != expectedStringFieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
		}
		uint16FieldValue := aTuple.AllFields()[1].Value()
		expectedUint16FieldValue := uint16(index)

		if uint16FieldValue != expectedUint16FieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
		}
	}
}

func add5Tuples(slottedPage *SlottedPage) []tuple.TupleId {
	tupleIds := make([]tuple.TupleId, 5)

	for index := 0; index < 5; index++ {
		aTuple := tuple.NewTuple()
		aTuple.AddField(field.NewStringField("Database Systems" + strconv.Itoa(index)))
		aTuple.AddField(field.NewUint16Field(uint16(index)))

		tupleIds[index] = slottedPage.Put(aTuple)
	}
	return tupleIds
}

func readTuples(tupleIds []tuple.TupleId, slottedPage *SlottedPage) []*tuple.Tuple {
	tuples := make([]*tuple.Tuple, len(tupleIds))

	for index := 0; index < 5; index++ {
		tupleId := tupleIds[index]
		tuples[index] = slottedPage.GetAt(tupleId.SlotNo)
	}
	return tuples
}
