package heap_file

import (
	"reflect"
	"testing"
)

func TestInitializesFreePageList(t *testing.T) {
	freePageList := InitializeFreePageList(5, 10)

	freePageIds := freePageList.pageIds
	expected := []uint32{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	if !reflect.DeepEqual(expected, freePageIds) {
		t.Fatalf("Expected freePageIds to be %v, received %v", expected, freePageIds)
	}
}

func TestAllocates3Pages(t *testing.T) {
	freePageList := InitializeFreePageList(5, 10)
	_, startingPageId := freePageList.allocateAndUpdate(3)
	expected := 5

	if startingPageId != 5 {
		t.Fatalf("Expected first free page id to be %v, received %v", expected, startingPageId)
	}
}

func TestAllocates3PagesAndUpdatesFreePageList(t *testing.T) {
	freePageList := InitializeFreePageList(5, 10)
	freePageList.allocateAndUpdate(3)

	freePageIds := freePageList.pageIds
	expected := []uint32{8, 9, 10, 11, 12, 13, 14}

	if !reflect.DeepEqual(expected, freePageIds) {
		t.Fatalf("Expected freePageIds to be %v, received %v", expected, freePageIds)
	}
}

func TestAllocates4Pages(t *testing.T) {
	freePageList := InitializeFreePageList(5, 4)
	_, startingPageId := freePageList.allocateAndUpdate(4)
	expected := 5

	if startingPageId != 5 {
		t.Fatalf("Expected first free page id to be %v, received %v", expected, startingPageId)
	}
}

func TestAllocates4PagesAndUpdatesFreePageList(t *testing.T) {
	freePageList := InitializeFreePageList(5, 4)
	freePageList.allocateAndUpdate(4)

	freePageIds := freePageList.pageIds
	expected := 0

	if len(freePageIds) != 0 {
		t.Fatalf("Expected length of freePageIds to be %v, received %v", expected, freePageIds)
	}
}

func TestDoesNotAllocate4Pages(t *testing.T) {
	freePageList := InitializeFreePageList(5, 3)
	isPageAvailable, _ := freePageList.allocateAndUpdate(4)

	if isPageAvailable != false {
		t.Fatalf("Expected isPageAvailable to be false")
	}
}
