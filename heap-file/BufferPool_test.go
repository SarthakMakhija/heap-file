package heap_file

import (
	"os"
	"testing"
)

func TestReturnsThePageCountInAFile(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	defer deleteFile(file)

	expectedPageCount := 0
	actualPageCount := bufferPool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected page count to be %v, received %v", expectedPageCount, actualPageCount)
	}
}

func TestReturnsTrueGivenFileContainsZeroPages(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	defer deleteFile(file)

	containsZeroPages := bufferPool.ContainsZeroPages()

	if containsZeroPages != true {
		t.Fatalf("Expected zero pages to be true")
	}
}

func TestReturnsFalseGivenFileContainsMoreThanZeroPages(t *testing.T) {
	file := createTestFile("./heap.db")

	writeToATestFileWithEmptyPage(file.Name(), 4096)
	bufferPool := NewBufferPool(file, 4096)

	defer deleteFile(file)

	containsZeroPages := bufferPool.ContainsZeroPages()

	if containsZeroPages != false {
		t.Fatalf("Expected more than zero pages to be contained")
	}
}

func TestAllocates5Pages(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	defer deleteFile(file)

	_, _ = bufferPool.Allocate(5)
	expectedPageCount := 5
	actualPageCount := bufferPool.pageCount

	if actualPageCount != expectedPageCount {
		t.Fatalf("Expected page count to be %v, received %v", expectedPageCount, actualPageCount)
	}
}

func TestReturnsTheCurrentPageIdAndAllocates5Pages(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	defer deleteFile(file)

	pageId, _ := bufferPool.Allocate(5)
	expectedPageId := 0

	if pageId != expectedPageId {
		t.Fatalf("Expected page id to be %v, received %v", expectedPageId, pageId)
	}
}

func TestReturnsTheNextPageIdAfterAllocating5Pages(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	defer deleteFile(file)

	_, _ = bufferPool.Allocate(5)
	expectedPageId := 5
	pageId := bufferPool.pageCount

	if pageId != expectedPageId {
		t.Fatalf("Expected page id to be %v, received %v", expectedPageId, pageId)
	}
}

func TestAllocationOf5PagesShouldIncreaseTheFileSize(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	defer deleteFile(file)

	_, _ = bufferPool.Allocate(5)
	fileSize := bufferPool.fileSize()
	expectedFileSize := 5 * 4096

	if expectedFileSize != fileSize {
		t.Fatalf("Expected file size to be %v, received %v", expectedFileSize, fileSize)
	}
}

func createTestFile(name string) *os.File {
	file, _ := os.Create(name)
	return file
}

func deleteFile(file *os.File) {
	_ = os.Remove(file.Name())
}

func writeToATestFileWithEmptyPage(fileName string, pageSize int) {
	writeToAATestFileAtOffset(fileName, make([]byte, pageSize), 0)
}

func writeToAATestFileAtOffset(fileName string, content []byte, offset int64) {
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	_, _ = file.Seek(offset, 0)
	_, _ = file.Write(content)
	_ = file.Close()
}
