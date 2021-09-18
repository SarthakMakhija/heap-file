package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"testing"
)

func TestCreatesAHeapFileWithACurrentPage(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	_, _ = bufferPool.Allocate(5)
	defer deleteFile(file)

	heapFile := NewHeapFile(bufferPool, InitializeFreePageList(0, 5))
	currentPageId := heapFile.currentPage.PageId()
	expected := uint32(0)

	if currentPageId != expected {
		t.Fatalf("Expected current page id to be %v, received %v", expected, currentPageId)
	}
}

func TestPutsAndGetsATupleInAPage(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	_, _ = bufferPool.Allocate(5)
	defer deleteFile(file)

	heapFile := NewHeapFile(bufferPool, InitializeFreePageList(0, 5))

	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(3000))

	tupleId := heapFile.Put(aTuple)
	readTuple := heapFile.GetAt(tupleId.SlotNo)

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
