package heap_file

import (
	"github.com/SarthakMakhija/heap-file/heap-file/field"
	"github.com/SarthakMakhija/heap-file/heap-file/page"
	"github.com/SarthakMakhija/heap-file/heap-file/tuple"
	"testing"
)

func TestWritesAndReadsSlottedPage(t *testing.T) {
	file := createTestFile("./heap.db")
	options := DefaultOptions()
	bufferPool := NewBufferPool(file, options)
	defer deleteFile(file)

	slottedPage := fillASlottedPage(options)

	_ = bufferPool.Write(slottedPage)
	readSlottedPage, _ := bufferPool.Read(slottedPage.PageId())

	aTuple := readSlottedPage.GetAt(1)

	stringFieldValue := aTuple.AllFields()[0].Value()
	expectedStringFieldValue := "Database Systems"

	if stringFieldValue != expectedStringFieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedStringFieldValue, stringFieldValue)
	}
	uint32FieldValue := aTuple.AllFields()[1].Value()
	expectedUint32FieldValue := uint32(100)

	if uint32FieldValue != expectedUint32FieldValue {
		t.Fatalf("Expected field value to be %v, received %v", expectedUint32FieldValue, uint32FieldValue)
	}
}

func fillASlottedPage(options DbOptions) *page.SlottedPage {
	aTuple := tuple.NewTuple()
	aTuple.AddField(field.NewStringField("Database Systems"))
	aTuple.AddField(field.NewUint32Field(uint32(100)))

	slottedPage := page.NewSlottedPage(0, options.PageSize(), options.TupleDescriptor())
	slottedPage.Put(aTuple.MarshalBinary())

	return slottedPage
}
