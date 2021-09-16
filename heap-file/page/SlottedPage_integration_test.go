package page

import (
	heap_file "github.com/SarthakMakhija/b-plus-tree/heap-file"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"strconv"
	"testing"
)

func TestPutsMultipleTuplesInASlottedPageAndReadsThemBack(t *testing.T) {
	slottedPage := NewSlottedPage(100)

	tupleIds := add5Tuples(slottedPage)
	tuples := readTuples(tupleIds, slottedPage)

	for index, tuple := range tuples {
		stringFieldValue := tuple.AllFields()[0].Value()
		expectedStringFieldValue := "Database Systems" + strconv.Itoa(index)

		if stringFieldValue != expectedStringFieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
		}
		uint16FieldValue := tuple.AllFields()[1].Value()
		expectedUint16FieldValue := uint16(index)

		if uint16FieldValue != expectedUint16FieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
		}
	}
}

func add5Tuples(slottedPage *SlottedPage) []heap_file.TupleId {
	tupleIds := make([]heap_file.TupleId, 5)

	for index := 0; index < 5; index++ {
		tuple := heap_file.NewTuple()
		tuple.AddField(field.NewStringField("Database Systems" + strconv.Itoa(index)))
		tuple.AddField(field.NewUint16Field(uint16(index)))

		tupleIds[index] = slottedPage.Put(tuple)
	}
	return tupleIds
}

func readTuples(tupleIds []heap_file.TupleId, slottedPage *SlottedPage) []*heap_file.Tuple {
	tuples := make([]*heap_file.Tuple, len(tupleIds))

	for index := 0; index < 5; index++ {
		tupleId := tupleIds[index]
		tuples[index] = slottedPage.Get(tupleId.SlotNo)
	}
	return tuples
}
