package heap_file

import (
	"encoding/binary"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"os"
	"unsafe"
)

var pageSize = os.Getpagesize()
var littleEndian = binary.LittleEndian
var pageIdSize = unsafe.Sizeof(SlottedPage{}.id)
var slotSize = unsafe.Sizeof(Slot{})

type SlottedPage struct {
	id        uint32
	buffer    []byte
	slotCount int
}

//revisit data types of SlottedPage

type Slot struct {
	tupleOffset uint16
	tupleSize   uint16
}

func NewSlottedPage(id uint32) *SlottedPage {
	slottedPage := &SlottedPage{
		buffer:    make([]byte, pageSize),
		id:        id,
		slotCount: 0,
	}
	slottedPage.writePageId(id)
	return slottedPage
}

func (slottedPage *SlottedPage) Put(tuple *Tuple) TupleId {
	slot := slottedPage.put(tuple)
	slottedPage.addSlot(slot)
	slottedPage.increaseSlotCount()

	return TupleId{
		pageId: slottedPage.id,
		slotNo: slottedPage.slotCount,
	}
}

func (slottedPage *SlottedPage) Get(slotNo int) *Tuple {
	//assume slotNo > 0
	slotStartingOffset := int(pageIdSize) + (slotNo-1)*int(slotSize)
	tupleSizeOffset := slotStartingOffset + 2

	tupleOffset := littleEndian.Uint16(slottedPage.buffer[slotStartingOffset:])
	tupleSize := littleEndian.Uint16(slottedPage.buffer[tupleSizeOffset:])

	tuple := NewTuple()
	tuple.UnMarshalBinary(
		slottedPage.buffer[tupleOffset:tupleOffset+tupleSize],
		[]field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}},
	)
	return tuple
}

func (slottedPage *SlottedPage) put(tuple *Tuple) Slot {
	buffer, size := tuple.MarshalBinary()
	startingOffset := pageSize - size
	copy(slottedPage.buffer[startingOffset:], buffer)

	return Slot{tupleOffset: uint16(startingOffset), tupleSize: uint16(size)}
}

func (slottedPage *SlottedPage) addSlot(slot Slot) {
	offset := pageIdSize
	littleEndian.PutUint16(slottedPage.buffer[offset:], slot.tupleOffset)
	offset = offset + 2
	littleEndian.PutUint16(slottedPage.buffer[offset:], slot.tupleSize)
}

func (slottedPage *SlottedPage) increaseSlotCount() {
	slottedPage.slotCount = slottedPage.slotCount + 1
}

func (slottedPage *SlottedPage) writePageId(id uint32) {
	littleEndian.PutUint32(slottedPage.buffer, id)
}
