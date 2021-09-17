package heap_file

import (
	"os"
	"reflect"
	"testing"
)

func TestCreatesADbByPreAllocatingPages(t *testing.T) {
	options := Options{
		PageSize:                 os.Getpagesize(),
		FileName:                 "./heap.db",
		PreAllocatedPagePoolSize: 6,
	}
	db, _ := Open(options)
	defer deleteFile(db.bufferPool.file)

	expectedPageCount := options.PreAllocatedPagePoolSize
	actualPageCount := db.bufferPool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected %v page count, received %v page count", expectedPageCount, actualPageCount)
	}
}

func TestCreatesADbWithFreePageList(t *testing.T) {
	options := Options{
		PageSize:                 os.Getpagesize(),
		FileName:                 "./heap.db",
		PreAllocatedPagePoolSize: 6,
	}
	db, _ := Open(options)
	defer deleteFile(db.bufferPool.file)

	expected := []uint32{0, 1, 2, 3, 4, 5}
	freePageIds := db.freePageList.pageIds

	if !reflect.DeepEqual(expected, freePageIds) {
		t.Fatalf("Expected free pageIds to be %v, received %v", expected, freePageIds)
	}
}
