package index

import (
	"os"
	"testing"
)

func TestCreatesABPlusTreeByPreAllocatingPagesAlongWithMetaPageAndRootPage(t *testing.T) {
	options := Options{
		PageSize:                 os.Getpagesize(),
		FileName:                 "./test",
		PreAllocatedPagePoolSize: 6,
	}
	tree, _ := Create(options)
	defer deleteFile(tree.pagePool.indexFile)

	expectedPageCount := options.PreAllocatedPagePoolSize + metaPageCount + rootPageCount
	actualPageCount := tree.pagePool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected %v page count, received %v page count", expectedPageCount, actualPageCount)
	}
}
