package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"github.com/SarthakMakhija/b-plus-tree/index"
	"os"
	"strconv"
	"testing"
)

func TestPutsAndGets1000Tuples(t *testing.T) {
	options := DbOptions{
		HeapFileOptions: HeapFileOptions{
			PageSize:                 os.Getpagesize(),
			FileName:                 "./heap.db",
			PreAllocatedPagePoolSize: 6,
			TupleDescriptor: tuple.TupleDescriptor{
				FieldTypes: []field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}},
			},
		},
		IndexOptions: index.DefaultOptions(),
	}
	db, _ := Open(options)
	defer deleteFile(db.bufferPool.file)
	defer deleteFileByName(options.IndexOptions.FileName)

	tupleIds := make([]tuple.TupleId, 1000)
	for iterator := 0; iterator < 1000; iterator++ {
		aTuple := tuple.NewTuple()
		aTuple.AddField(field.NewStringField("Database Systems" + strconv.Itoa(iterator)))
		aTuple.AddField(field.NewUint16Field(uint16(iterator)))

		tupleId, _ := db.Put(aTuple)
		tupleIds[iterator] = tupleId
	}

	for iterator, tupleId := range tupleIds {
		readTuple := db.GetBy(tupleId)

		stringFieldValue := readTuple.AllFields()[0].Value()
		expectedStringFieldValue := "Database Systems" + strconv.Itoa(iterator)

		if stringFieldValue != expectedStringFieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
		}

		uint16FieldValue := readTuple.AllFields()[1].Value()
		expectedUint16FieldValue := uint16(iterator)

		if uint16FieldValue != expectedUint16FieldValue {
			t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
		}
	}
}
