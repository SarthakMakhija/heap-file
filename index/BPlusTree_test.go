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

func TestDoesNotGetByKey(t *testing.T) {
	options := DefaultOptions()
	tree, _ := Create(options)
	defer deleteFile(tree.pagePool.indexFile)

	tree.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: uint64(100),
		},
		{
			key:   []byte("B"),
			value: uint64(200),
		},
	}

	_, found, err := tree.Get([]byte("C"))

	if found != false && err != nil {
		t.Fatalf("Expected found to be false received %v, and error to be nil, received err %v", found, err)
	}
}

func TestGetsByKeyInRootLeafPage(t *testing.T) {
	options := DefaultOptions()
	tree, _ := Create(options)
	defer deleteFile(tree.pagePool.indexFile)

	tree.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: uint64(100),
		},
		{
			key:   []byte("B"),
			value: uint64(200),
		},
	}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("B"),
		value: uint64(200),
	}
	keyValuePair, _, _ := tree.Get([]byte("B"))

	if !expectedKeyValuePair.Equals(keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, keyValuePair)
	}
}

func TestGetsByKeyInTheLeafPageWhichIsTheLeftChildOfRootPage(t *testing.T) {
	writeLeftPageToFile := func(fileName string, pageSize int) {
		leftPage := Page{
			id: 1,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("A"),
					value: uint64(100),
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
					key:   []byte("R"),
					value: uint64(500),
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
	tree, _ := Create(options)
	defer deleteFile(tree.pagePool.indexFile)

	tree.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("C"),
			value: uint64(300),
		},
	}
	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	tree.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: uint64(100),
	}
	keyValuePair, _, _ := tree.Get([]byte("A"))

	if !expectedKeyValuePair.Equals(keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, keyValuePair)
	}
}

func TestGetsByKeyInTheLeafPageWhichIsTheRightChildOfRootPage(t *testing.T) {
	writeLeftPageToFile := func(fileName string, pageSize int) {
		leftPage := Page{
			id: 1,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("A"),
					value: uint64(100),
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
					key:   []byte("R"),
					value: uint64(500),
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
	tree, _ := Create(options)
	defer deleteFile(tree.pagePool.indexFile)

	tree.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("C"),
			value: uint64(300),
		},
	}
	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	tree.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("R"),
		value: uint64(500),
	}
	keyValuePair, _, _ := tree.Get([]byte("R"))

	if !expectedKeyValuePair.Equals(keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, keyValuePair)
	}
}
