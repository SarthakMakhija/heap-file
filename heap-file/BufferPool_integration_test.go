package heap_file

import (
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/page"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"os"
	"testing"
)

func TestWritesAndReadsSlottedPage(t *testing.T) {
	file := createTestFile("./heap.db")
	bufferPool := NewBufferPool(file, 4096)
	defer deleteFile(file)

	slottedPage := fillASlottedPage()

	_ = bufferPool.Write(slottedPage)
	readSlottedPage, _ := bufferPool.Read(slottedPage.PageId())

	aTuple := readSlottedPage.GetAt(1)

	stringFieldValue := aTuple.AllFields()[0].Value()
	expectedStringFieldValue := "Database Systems"

	if stringFieldValue != expectedStringFieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
	}
	uint16FieldValue := aTuple.AllFields()[1].Value()
	expectedUint16FieldValue := uint16(100)

	if uint16FieldValue != expectedUint16FieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedUint16FieldValue, uint16FieldValue)
	}
}

func fillASlottedPage() *page.SlottedPage {
	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint16Field(uint16(100)))

	slottedPage := page.NewSlottedPage(0, os.Getpagesize())
	slottedPage.Put(aTuple.MarshalBinary())

	return slottedPage
}
