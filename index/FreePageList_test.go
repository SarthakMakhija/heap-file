package index

import (
	"reflect"
	"testing"
)

func TestInitializesFreePageList(t *testing.T) {
	freePageList := InitializeFreePageList(5, 10)

	freePageIds := freePageList.pageIds
	expected := []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	if !reflect.DeepEqual(expected, freePageIds) {
		t.Fatalf("Expected freePageIds to be %v, received %v", expected, freePageIds)
	}
}
