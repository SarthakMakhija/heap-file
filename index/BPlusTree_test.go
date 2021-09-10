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
	tree, _ := CreateBPlusTree(options)
	defer deleteFile(tree.pagePool.indexFile)

	expectedPageCount := options.PreAllocatedPagePoolSize + metaPageCount + rootPageCount
	actualPageCount := tree.pagePool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected %v page count, received %v page count", expectedPageCount, actualPageCount)
	}
}

func TestCreatesABPlusTreeWithARootPage(t *testing.T) {
	options := DefaultOptions()
	tree, _ := CreateBPlusTree(options)
	defer deleteFile(tree.pagePool.indexFile)

	if tree.pageHierarchy.rootPage == nil {
		t.Fatalf("Expected root page to be non-nil received nil")
	}
}

func TestCreatesABPlusTreeByCachingRootPage(t *testing.T) {
	options := DefaultOptions()
	tree, _ := CreateBPlusTree(options)
	defer deleteFile(tree.pagePool.indexFile)

	rootPageId := tree.pageHierarchy.RootPageId()
	rootPage := tree.pageHierarchy.PageById(rootPageId)

	if rootPage == nil {
		t.Fatalf("Expected root page in page cache to be non-nil received nil")
	}
}

func TestDoesNotGetByKeyAsSearchedKeyDoesNotExist(t *testing.T) {
	options := DefaultOptions()
	tree, _ := CreateBPlusTree(options)
	defer deleteFile(tree.pagePool.indexFile)

	tree.pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{key: []byte("A")},
		{key: []byte("B")},
	}

	getResult := tree.Get([]byte("C"))

	if getResult.found != false && getResult.err != nil {
		t.Fatalf("Expected found to be false received %v, and error to be nil, received err %v", getResult.found, getResult.err)
	}
}

func TestGetsByKeyGivenKeyIsFoundInTheNonLeafPage(t *testing.T) {
	writeLeftPageToFile := func(fileName string, pageSize int) {
		leftPage := Page{
			id: 1,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("A"),
					value: []byte("Storage"),
				},
			},
		}
		writeToAATestFileAtOffset(fileName, leftPage.MarshalBinary(), int64(pageSize*leftPage.id))
	}
	writeRightPageToFile := func(fileName string, pageSize int) {
		rightPage := Page{
			id: 2,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("B"),
					value: []byte("Database"),
				},
				{
					key:   []byte("C"),
					value: []byte("Systems"),
				},
			},
		}
		writeToAATestFileAtOffset(fileName, rightPage.MarshalBinary(), int64(pageSize*rightPage.id))
	}

	options := Options{
		PageSize:                 os.Getpagesize(),
		FileName:                 "./test",
		PreAllocatedPagePoolSize: 8,
	}
	tree, _ := CreateBPlusTree(options)
	defer deleteFile(tree.pagePool.indexFile)

	tree.pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("B"),
			value: []byte("Database"),
		},
	}
	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	tree.pageHierarchy.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("B"),
		value: []byte("Database"),
	}
	getResult := tree.Get([]byte("B"))

	if !expectedKeyValuePair.Equals(getResult.keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, getResult.keyValuePair)
	}
}
