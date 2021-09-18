package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"os"
	"strconv"
	"testing"
)

func TestPutsAndGets1000Tuples(t *testing.T) {
	options := DbOptions{
		PageSize:                 os.Getpagesize(),
		FileName:                 "./heap.db",
		PreAllocatedPagePoolSize: 6,
		TupleDescriptor: tuple.TupleDescriptor{
			FieldTypes: []field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}},
		},
	}
	db, _ := Open(options)
	defer deleteFile(db.bufferPool.file)

	tupleIds := make([]tuple.TupleId, 1000)
	for index := 0; index < 1000; index++ {
		aTuple := tuple.NewTuple()
		aTuple.AddField(field.NewStringField("Database Systems" + strconv.Itoa(index)))
		aTuple.AddField(field.NewUint16Field(uint16(index)))

		tupleId, _ := db.Put(aTuple)
		tupleIds[index] = tupleId
	}

	for index, tupleId := range tupleIds {
		readTuple := db.GetBy(tupleId)

		stringFieldValue := readTuple.AllFields()[0].Value()
		expectedStringFieldValue := "Database Systems" + strconv.Itoa(index)

		if stringFieldValue != expectedStringFieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
		}

		uint16FieldValue := readTuple.AllFields()[1].Value()
		expectedUint16FieldValue := uint16(index)

		if uint16FieldValue != expectedUint16FieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
		}
	}
}
