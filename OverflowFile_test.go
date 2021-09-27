package heap_file

import (
	"math/rand"
	"os"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	bytes := make([]byte, n)
	for iterator := range bytes {
		bytes[iterator] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bytes)
}

func TestPutsBufferIntoAnOverflowPageGivenBufferFitsInAPage(t *testing.T) {
	key := randStringBytes(100)
	pageSize := 200

	file := createTestFile("./heap.db")
	writeToATestFileWithEmptyPage(file.Name(), pageSize)

	overflowPool := NewOverflowPool(file, pageSize)
	defer deleteFile(file)

	overflowFile := NewOverflowFile(overflowPool, pageSize)
	slot := overflowFile.Put([]byte(key))

	readBuffer := overflowFile.GetAt(slot.id)
	readString := string(readBuffer)

	if readString != key {
		t.Fatalf("Expected %v, received %v", key, readString)
	}
}

func TestPutsBufferIntoAnOverflowPageGivenBufferDoesNotFitInAPage(t *testing.T) {
	key := randStringBytes(300)
	pageSize := 200

	file := createTestFile("./heap.db")
	writeToATestFileWithNEmptyPages(file.Name(), pageSize, 5)

	overflowPool := NewOverflowPool(file, pageSize)
	defer deleteFile(file)

	overflowFile := NewOverflowFile(overflowPool, pageSize)
	slot := overflowFile.Put([]byte(key))

	readBuffer := overflowFile.GetAt(slot.id)
	readString := string(readBuffer)

	if readString != key {
		t.Fatalf("Expected %v, received %v", key, readString)
	}
}

func createTestFile(name string) *os.File {
	file, _ := os.Create(name)
	return file
}

func deleteFile(file *os.File) {
	_ = os.Remove(file.Name())
}

func writeToATestFileWithNEmptyPages(fileName string, pageSize int, n int) {
	for count := 1; count <= n; count++ {
		writeToAATestFileAtOffset(fileName, make([]byte, pageSize), int64((count-1)*pageSize))
	}
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
