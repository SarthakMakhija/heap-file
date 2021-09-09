package index

import (
	"os"
	"testing"
)

func TestReturnsThePageCountInAnIndexFile(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := Open(options)
	pagePool := New(indexFile, options)

	defer deleteFile(indexFile)

	expectedPageCount := 0
	actualPageCount := pagePool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected page count to be %v, received %v", expectedPageCount, actualPageCount)
	}
}

func TestReturnsTrueGivenIndexFileContainsZeroPages(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := Open(options)
	pagePool := New(indexFile, options)

	defer deleteFile(indexFile)

	containsZeroPages := pagePool.ContainsZeroPages()

	if containsZeroPages != true {
		t.Fatalf("Expected zero pages to be true")
	}
}

func TestReturnsFalseGivenIndexFileContainsMoreThanZeroPages(t *testing.T) {
	options := Options{
		PageSize: os.Getpagesize(),
		FileName: "./test",
	}
	createATestFileWithEmptyPage(options.FileName, options.PageSize)

	indexFile, _ := Open(options)
	defer deleteFile(indexFile)
	pagePool := New(indexFile, options)

	defer deleteFile(indexFile)

	containsZeroPages := pagePool.ContainsZeroPages()

	if containsZeroPages != false {
		t.Fatalf("Expected zero pages to be false")
	}
}

func TestAllocates5Pages(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := Open(options)
	pagePool := New(indexFile, options)

	defer deleteFile(indexFile)

	_ = pagePool.Allocate(5)
	expectedPageCount := 5
	actualPageCount := pagePool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected page count to be %v, received %v", expectedPageCount, actualPageCount)
	}
}

func TestAllocationOf5PagesShouldIncreaseTheFileSize(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := Open(options)
	pagePool := New(indexFile, options)

	defer deleteFile(indexFile)

	_ = pagePool.Allocate(5)
	expectedFileSize := int64(5 * os.Getpagesize())
	actualFileSize := pagePool.indexFile.size

	if actualFileSize != expectedFileSize {
		t.Fatalf("Expected file size to be %v, received %v", expectedFileSize, actualFileSize)
	}
}

func createATestFileWithEmptyPage(fileName string, pageSize int) {
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	_, _ = file.Write(make([]byte, pageSize))
	_ = file.Close()
}
