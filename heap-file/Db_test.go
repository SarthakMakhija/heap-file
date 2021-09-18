package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"github.com/SarthakMakhija/b-plus-tree/index"
	"os"
	"reflect"
	"testing"
)

func TestCreatesADbByPreAllocatingPages(t *testing.T) {
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

	expectedPageCount := options.PreAllocatedPagePoolSize()
	actualPageCount := db.bufferPool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected %v page count, received %v page count", expectedPageCount, actualPageCount)
	}
}

func TestCreatesABPlusTreeWhenOpened(t *testing.T) {
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

	tree := db.bPlusTree
	if tree == nil {
		t.Fatalf("Expected bPlusTree to be initialized when db is opened")
	}
}

func TestCreatesADbWithFreePageListAndUsesTheFirstPageForHeapFile(t *testing.T) {
	options := DbOptions{
		HeapFileOptions: HeapFileOptions{
			PageSize:                 os.Getpagesize(),
			FileName:                 "./heap.db",
			PreAllocatedPagePoolSize: 6,
		},
		IndexOptions: index.DefaultOptions(),
	}
	db, _ := Open(options)
	defer deleteFile(db.bufferPool.file)
	defer deleteFileByName(options.IndexOptions.FileName)

	expected := []uint32{1, 2, 3, 4, 5}
	freePageIds := db.freePageList.pageIds

	if !reflect.DeepEqual(expected, freePageIds) {
		t.Fatalf("Expected free pageIds to be %v, received %v", expected, freePageIds)
	}
}

func TestPutsAndGetsATupleByTupleId(t *testing.T) {
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

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	tupleId, _ := db.Put(aTuple)

	readTuple := db.GetByTupleId(tupleId)

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

func TestPutsAndGetsATupleByKey(t *testing.T) {
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

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	_, _ = db.Put(aTuple)
	readTuple, _ := db.GetByKey(aTuple.KeyField())

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
