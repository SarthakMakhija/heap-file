package index

import (
	"bytes"
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

func TestReadsAPageIdentifiedByPageId0(t *testing.T) {
	options := Options{
		PageSize: 10,
		FileName: "./test",
	}
	createATestFileWith(options.FileName, []byte("helloAdam0"))

	indexFile, _ := Open(options)
	pagePool := New(indexFile, options)
	defer deleteFile(indexFile)

	pageId := 0
	content, _ := pagePool.Read(pageId)
	expected := []byte("helloAdam0")

	if !bytes.Equal(content, expected) {
		t.Fatalf("Expected page %v to be %v, received %v", pageId, string(expected), string(content))
	}
}

func TestReadsAPageIdentifiedByPageId1(t *testing.T) {
	options := Options{
		PageSize: 10,
		FileName: "./test",
	}
	createATestFileWith(options.FileName, []byte("helloAdam0helloBrad1"))

	indexFile, _ := Open(options)
	pagePool := New(indexFile, options)
	defer deleteFile(indexFile)

	pageId := 1
	content, _ := pagePool.Read(pageId)
	expected := []byte("helloBrad1")

	if !bytes.Equal(content, expected) {
		t.Fatalf("Expected page %v to be %v, received %v", pageId, string(expected), string(content))
	}
}

func createATestFileWithEmptyPage(fileName string, pageSize int) {
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	_, _ = file.Write(make([]byte, pageSize))
	_ = file.Close()
}

func createATestFileWith(fileName string, content []byte) {
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	_, _ = file.Write(content)
	_ = file.Close()
}
