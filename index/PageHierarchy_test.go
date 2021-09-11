package index

import (
	"os"
	"reflect"
	"testing"
)

func TestReturnsPageById(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	pageHierarchy := NewPageHierarchy(pagePool)
	pageHierarchy.pageById[0] = &Page{
		id: 0,
	}

	defer deleteFile(pagePool.indexFile)

	page := pageHierarchy.PageById(0)
	if page.id != 0 {
		t.Fatalf("Expected page id to be 0 received %v", page.id)
	}
}

func TestReturnsTheRootPageId(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	pageHierarchy := NewPageHierarchy(pagePool)
	pageHierarchy.rootPage = &Page{id: 100}

	defer deleteFile(pagePool.indexFile)

	rootPageId := pageHierarchy.RootPageId()
	if rootPageId != 100 {
		t.Fatalf("Expected root page id to be 100 received %v", rootPageId)
	}
}

func TestDoesNotGetByKey(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{key: []byte("A")},
		{key: []byte("B")},
	}

	getResult := pageHierarchy.Get([]byte("C"))

	if getResult.found != false && getResult.err != nil {
		t.Fatalf("Expected found to be false received %v, and error to be nil, received err %v", getResult.found, getResult.err)
	}
}

func TestGetsByKeyInRootLeafPage(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: []byte("Database"),
		},
		{
			key:   []byte("B"),
			value: []byte("Systems"),
		},
	}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("B"),
		value: []byte("Systems"),
	}
	getResult := pageHierarchy.Get([]byte("B"))

	if !expectedKeyValuePair.Equals(getResult.keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, getResult.keyValuePair)
	}
}

func TestGetsByKeyInTheLeafPageWhichIsTheLeftChildOfRootPage(t *testing.T) {
	writeLeftPageToFile := func(fileName string, pageSize int) {
		leftPage := Page{
			id: 1,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("A"),
					value: []byte("Database"),
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
					value: []byte("Storage"),
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
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{{key: []byte("B")}}

	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	pageHierarchy.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: []byte("Database"),
	}
	getResult := pageHierarchy.Get([]byte("A"))

	if !expectedKeyValuePair.Equals(getResult.keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, getResult.keyValuePair)
	}
}

func TestGetsByKeyInTheLeafPageWhichIsTheRightChildOfRootPage(t *testing.T) {
	writeLeftPageToFile := func(fileName string, pageSize int) {
		leftPage := Page{
			id: 1,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("A"),
					value: []byte("Database"),
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
					value: []byte("Storage"),
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
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{{key: []byte("B")}}

	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	pageHierarchy.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("C"),
		value: []byte("Systems"),
	}
	getResult := pageHierarchy.Get([]byte("C"))

	if !expectedKeyValuePair.Equals(getResult.keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, getResult.keyValuePair)
	}
}

func TestGetsByKeyInTheLeafPageWhichIsTheRightChildOfRootPageGivenKeyIsFoundInTheNonLeafPage(t *testing.T) {
	writeLeftPageToFile := func(fileName string, pageSize int) {
		leftPage := Page{
			id: 1,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("A"),
					value: []byte("Database"),
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
					value: []byte("Storage"),
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
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{{key: []byte("B")}}

	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	pageHierarchy.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("B"),
		value: []byte("Storage"),
	}
	getResult := pageHierarchy.Get([]byte("B"))

	if !expectedKeyValuePair.Equals(getResult.keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, getResult.keyValuePair)
	}
}

func TestPutsAKeyValuePairInRootLeafPage(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: []byte("Database"),
		},
		{
			key:   []byte("C"),
			value: []byte("Systems"),
		},
	}
	_ = pageHierarchy.Put(KeyValuePair{key: []byte("B"), value: []byte("Storage")})
	expected := []KeyValuePair{
		{key: []byte("A"), value: []byte("Database")},
		{key: []byte("B"), value: []byte("Storage")},
		{key: []byte("C"), value: []byte("Systems")},
	}

	pageKeyValuePairs := pageHierarchy.rootPage.keyValuePairs
	if !reflect.DeepEqual(expected, pageKeyValuePairs) {
		t.Fatalf("Expected Key value pairs to be %v, received %v", expected, pageKeyValuePairs)
	}
}

func TestPutsAKeyValuePairInTheRightPage(t *testing.T) {
	writeLeftPageToFile := func(fileName string, pageSize int) {
		leftPage := Page{
			id: 1,
			keyValuePairs: []KeyValuePair{
				{
					key:   []byte("A"),
					value: []byte("Database"),
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
					value: []byte("Storage"),
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
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{{key: []byte("B")}}

	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	pageHierarchy.rootPage.childPageIds = []int{1, 2}

	_ = pageHierarchy.Put(KeyValuePair{key: []byte("D"), value: []byte("OS")})

	getResult := pageHierarchy.Get([]byte("D"))
	expected := KeyValuePair{key: []byte("D"), value: []byte("OS")}

	if !expected.Equals(getResult.keyValuePair) {
		t.Fatalf("Expected Key value pair to be %v, received %v", expected, getResult.keyValuePair)
	}
}

func TestPutsAKeyValuePairAfterSplittingTheRootPage(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)
	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: []byte("Database"),
		},
		{
			key:   []byte("C"),
			value: []byte("Systems"),
		},
		{
			key:   []byte("E"),
			value: []byte("OS"),
		},
	}

	_ = pageHierarchy.Put(KeyValuePair{key: []byte("D"), value: []byte("File System")})

	getResult := pageHierarchy.Get([]byte("D"))
	expected := KeyValuePair{key: []byte("D"), value: []byte("File System")}

	if !expected.Equals(getResult.keyValuePair) {
		t.Fatalf("Expected Key value pair to be %v, received %v", expected, getResult.keyValuePair)
	}
}

func TestSplitsTheRootPageAndCreatesANewRootWithKeyValuePairs(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)
	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: []byte("Database"),
		},
		{
			key:   []byte("C"),
			value: []byte("Systems"),
		},
		{
			key:   []byte("E"),
			value: []byte("OS"),
		},
	}

	_ = pageHierarchy.Put(KeyValuePair{key: []byte("D"), value: []byte("File System")})

	keyValuePairsOfNewRootPage := pageHierarchy.rootPage.keyValuePairs
	expected := []KeyValuePair{{key: []byte("C"), value: []byte("Systems")}}

	if !reflect.DeepEqual(expected, keyValuePairsOfNewRootPage) {
		t.Fatalf("Expected Key value pair in the new root to be %v, received %v", expected, keyValuePairsOfNewRootPage)
	}
}

func TestSplitsTheRootPageAndWithKeyValuePairsInOldRoot(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)
	existingRootPage := pageHierarchy.rootPage
	existingRootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: []byte("Database"),
		},
		{
			key:   []byte("C"),
			value: []byte("Systems"),
		},
		{
			key:   []byte("E"),
			value: []byte("OS"),
		},
	}

	_ = pageHierarchy.Put(KeyValuePair{key: []byte("D"), value: []byte("File System")})

	keyValuePairs := existingRootPage.keyValuePairs
	expected := []KeyValuePair{{key: []byte("A"), value: []byte("Database")}}

	if !reflect.DeepEqual(expected, keyValuePairs) {
		t.Fatalf("Expected Key value pair in the old root to be %v, received %v", expected, keyValuePairs)
	}
}

func TestSplitsTheRootPageAndWithKeyValuePairsInRightSiblingPage(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	_ = pagePool.Allocate(options.PreAllocatedPagePoolSize)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)
	existingRootPage := pageHierarchy.rootPage
	existingRootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: []byte("Database"),
		},
		{
			key:   []byte("C"),
			value: []byte("Systems"),
		},
		{
			key:   []byte("E"),
			value: []byte("OS"),
		},
	}

	_ = pageHierarchy.Put(KeyValuePair{key: []byte("D"), value: []byte("File System")})
	rightSibling := pageHierarchy.PageById(pageHierarchy.rootPage.childPageIds[1])

	keyValuePairs := rightSibling.keyValuePairs
	expected := []KeyValuePair{{key: []byte("C"), value: []byte("Systems")}, {key: []byte("D"), value: []byte("File System")}, {key: []byte("E"), value: []byte("OS")}}

	if !reflect.DeepEqual(expected, keyValuePairs) {
		t.Fatalf("Expected Key value pair in the right sibling to be %v, received %v", expected, keyValuePairs)
	}
}
