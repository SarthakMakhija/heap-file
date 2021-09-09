package index

import (
	"os"
	"testing"
)

func deleteFile(indexFile *IndexFile) {
	_ = os.Remove(indexFile.file.Name())
}

func TestCreateAnIndexFileWithPageSize(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := Open(options)
	defer deleteFile(indexFile)

	expectedPageSize := options.PageSize
	actualPageSize := indexFile.pageSize

	if actualPageSize != expectedPageSize {
		t.Fatalf("Expected page size to be %v, received %v", expectedPageSize, actualPageSize)
	}
}

func TestCreateAnIndexFileWithPageCount(t *testing.T) {
	options := Options{
		PageSize: os.Getpagesize(),
		FileName: "./test",
	}
	CreateATestFileWithEmptyPage(options.FileName, options.PageSize)

	indexFile, _ := Open(options)
	defer deleteFile(indexFile)

	expectedPageCount := 1
	actualPageCount := indexFile.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected page count to be %v, received %v", expectedPageCount, actualPageCount)
	}
}

func TestAllocate5Pages(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := Open(options)
	defer deleteFile(indexFile)

	_ = indexFile.Allocate(5)
	expectedPageCount := 5
	actualPageCount := indexFile.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected page count to be %v, received %v", expectedPageCount, actualPageCount)
	}
}

func TestAllocationOf5PagesShouldIncreaseTheFileSize(t *testing.T) {
	options := DefaultOptions()
	indexFile, _ := Open(options)
	defer deleteFile(indexFile)

	_ = indexFile.Allocate(5)
	expectedFileSize := int64(5 * os.Getpagesize())
	actualFileSize := indexFile.size

	if actualFileSize != expectedFileSize {
		t.Fatalf("Expected file size to be %v, received %v", expectedFileSize, actualFileSize)
	}
}

func CreateATestFileWithEmptyPage(fileName string, pageSize int) {
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	_, _ = file.Write(make([]byte, pageSize))
	_ = file.Close()
}
