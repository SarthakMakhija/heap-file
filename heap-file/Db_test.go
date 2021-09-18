package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"os"
	"reflect"
	"testing"
)

func TestCreatesADbByPreAllocatingPages(t *testing.T) {
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

	expectedPageCount := options.PreAllocatedPagePoolSize
	actualPageCount := db.bufferPool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected %v page count, received %v page count", expectedPageCount, actualPageCount)
	}
}

func TestCreatesADbWithFreePageListAndUsesTheFirstPageForHeapFile(t *testing.T) {
	options := DbOptions{
		PageSize:                 os.Getpagesize(),
		FileName:                 "./heap.db",
		PreAllocatedPagePoolSize: 6,
	}
	db, _ := Open(options)
	defer deleteFile(db.bufferPool.file)

	expected := []uint32{1, 2, 3, 4, 5}
	freePageIds := db.freePageList.pageIds

	if !reflect.DeepEqual(expected, freePageIds) {
		t.Fatalf("Expected free pageIds to be %v, received %v", expected, freePageIds)
	}
}

func TestPutsAndGetsATuple(t *testing.T) {
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

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	tupleId, _ := db.Put(aTuple)

	readTuple := db.GetBy(tupleId)

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
