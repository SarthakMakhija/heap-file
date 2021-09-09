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

func TestCreatesABPlusTreeWithARootPage(t *testing.T) {
	options := DefaultOptions()
	tree, _ := Create(options)
	defer deleteFile(tree.pagePool.indexFile)

	if tree.rootPage == nil {
		t.Fatalf("Expected root page to be non-nil received nil")
	}
}

func TestCreatesABPlusTreeByCachingRootPage(t *testing.T) {
	options := DefaultOptions()
	tree, _ := Create(options)
	defer deleteFile(tree.pagePool.indexFile)

	rootPageId := tree.rootPage.id
	rootPage := tree.pageById[rootPageId]

	if rootPage == nil {
		t.Fatalf("Expected root page in page cache to be non-nil received nil")
	}
}
