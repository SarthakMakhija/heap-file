package page

import (
	"encoding/binary"
	"github.com/SarthakMakhija/b-plus-tree/heap-file"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"os"
	"unsafe"
)

var pageSize = os.Getpagesize()
var littleEndian = binary.LittleEndian
var pageIdSize = unsafe.Sizeof(SlottedPage{}.id)

type SlottedPage struct {
	id        uint32
	buffer    []byte
	slotCount int
}

//revisit data types of SlottedPage

func NewSlottedPage(id uint32) *SlottedPage {
	slottedPage := &SlottedPage{
		buffer:    make([]byte, pageSize),
		id:        id,
		slotCount: 0,
	}
	slottedPage.writePageId(id)
	return slottedPage
}

func (slottedPage *SlottedPage) Put(tuple *heap_file.Tuple) heap_file.TupleId {
	slot := slottedPage.put(tuple)
	slottedPage.addSlot(slot)
	slottedPage.increaseSlotCount()

	return heap_file.TupleId{
		PageId: slottedPage.id,
		SlotNo: slottedPage.slotCount,
	}
}

func (slottedPage *SlottedPage) Get(slotNo int) *heap_file.Tuple {
	slotStartingOffset := int(pageIdSize) + (slotNo-1)*int(slotSize)

	slot := &Slot{}
	slot.UnMarshalBinary(slottedPage.buffer, slotStartingOffset)

	tuple := heap_file.NewTuple()
	tuple.UnMarshalBinary(
		slottedPage.buffer[slot.tupleOffset:slot.tupleOffset+slot.tupleSize],
		[]field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}},
	)
	return tuple
}

func (slottedPage *SlottedPage) put(tuple *heap_file.Tuple) Slot {
	buffer, size := tuple.MarshalBinary()
	startingOffset := pageSize - size
	copy(slottedPage.buffer[startingOffset:], buffer)

	return Slot{tupleOffset: uint16(startingOffset), tupleSize: uint16(size)}
}

func (slottedPage *SlottedPage) addSlot(slot Slot) {
	offset := pageIdSize
	copy(slottedPage.buffer[offset:], slot.MarshalBinary())
}

func (slottedPage *SlottedPage) increaseSlotCount() {
	slottedPage.slotCount = slottedPage.slotCount + 1
}

func (slottedPage *SlottedPage) writePageId(id uint32) {
	littleEndian.PutUint32(slottedPage.buffer, id)
}
