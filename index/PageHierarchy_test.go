package index

import (
	"os"
	"testing"
)

func TestDoesNotGetByKey(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := OpenIndexFile(options)
	pagePool := NewPagePool(indexFile, options)
	pageHierarchy := NewPageHierarchy(pagePool)

	defer deleteFile(pagePool.indexFile)

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("A"),
			value: uint64(100),
		},
		{
			key:   []byte("B"),
			value: uint64(200),
		},
	}

	_, found, err := pageHierarchy.Get([]byte("C"))

	if found != false && err != nil {
		t.Fatalf("Expected found to be false received %v, and error to be nil, received err %v", found, err)
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
	keyValuePair, _, _ := pageHierarchy.Get([]byte("B"))

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
					key:   []byte("B"),
					value: uint64(200),
				},
				{
					key:   []byte("C"),
					value: uint64(300),
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

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("B"),
			value: uint64(200),
		},
	}
	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	pageHierarchy.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("A"),
		value: uint64(100),
	}
	keyValuePair, _, _ := pageHierarchy.Get([]byte("A"))

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
					key:   []byte("B"),
					value: uint64(200),
				},
				{
					key:   []byte("C"),
					value: uint64(300),
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

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("B"),
			value: uint64(200),
		},
	}
	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	pageHierarchy.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("C"),
		value: uint64(300),
	}
	keyValuePair, _, _ := pageHierarchy.Get([]byte("C"))

	if !expectedKeyValuePair.Equals(keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, keyValuePair)
	}
}

func TestGetsByKeyInTheLeafPageWhichIsTheRightChildOfRootPageGivenKeyIsFoundInTheNonLeafPage(t *testing.T) {
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
					key:   []byte("B"),
					value: uint64(200),
				},
				{
					key:   []byte("C"),
					value: uint64(300),
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

	pageHierarchy.rootPage.keyValuePairs = []KeyValuePair{
		{
			key:   []byte("B"),
			value: uint64(200),
		},
	}
	writeLeftPageToFile(options.FileName, options.PageSize)
	writeRightPageToFile(options.FileName, options.PageSize)
	pageHierarchy.rootPage.childPageIds = []int{1, 2}

	expectedKeyValuePair := KeyValuePair{
		key:   []byte("B"),
		value: uint64(200),
	}
	keyValuePair, _, _ := pageHierarchy.Get([]byte("B"))

	if !expectedKeyValuePair.Equals(keyValuePair) {
		t.Fatalf("Expected KeyValuePair to be %v, received %v", expectedKeyValuePair, keyValuePair)
	}
}
