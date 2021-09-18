package page

import (
	"encoding/binary"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/field"
	"github.com/SarthakMakhija/b-plus-tree/heap-file/tuple"
	"unsafe"
)

var littleEndian = binary.LittleEndian
var pageIdSize = unsafe.Sizeof(SlottedPage{}.id)

type SlottedPage struct {
	pageSize  int
	id        uint32
	buffer    []byte
	slotCount int
}

//revisit data types of SlottedPage

func NewSlottedPage(id uint32, pageSize int) *SlottedPage {
	slottedPage := &SlottedPage{
		pageSize:  pageSize,
		buffer:    make([]byte, pageSize),
		id:        id,
		slotCount: 0,
	}
	slottedPage.writePageId(id)
	return slottedPage
}

func NewReadonlySlottedPageFrom(buffer []byte) *SlottedPage {
	slottedPage := &SlottedPage{
		buffer: buffer,
		id:     littleEndian.Uint32(buffer),
	}
	return slottedPage
}

func (slottedPage *SlottedPage) Put(aTuple tuple.MarshalledTuple) tuple.TupleId {
	slot := slottedPage.put(aTuple)
	slottedPage.addSlot(slot)
	slottedPage.increaseSlotCount()

	return tuple.TupleId{
		PageId: slottedPage.id,
		SlotNo: slottedPage.slotCount,
	}
}

func (slottedPage *SlottedPage) GetAt(slotNo int) *tuple.Tuple {
	aTuple := tuple.NewTuple()
	slot := slottedPage.getSlot(slotNo)
	if slot == nil {
		return aTuple
	}
	aTuple.UnMarshalBinary(
		slottedPage.buffer[slot.tupleOffset:slot.tupleOffset+slot.tupleSize],
		[]field.FieldType{field.StringFieldType{}, field.Uint16FieldType{}},
	)
	return aTuple
}

func (slottedPage SlottedPage) SizeAvailable() uint16 {
	size := uint16(pageIdSize)
	for slotNo := 1; slotNo <= slottedPage.slotCount; slotNo++ {
		slot := slottedPage.getSlot(slotNo)
		size = size + slot.tupleSize + uint16(slotSize)
	}
	return uint16(slottedPage.pageSize) - size
}

func (slottedPage SlottedPage) PageId() uint32 {
	return slottedPage.id
}

func (slottedPage SlottedPage) Buffer() []byte {
	return slottedPage.buffer
}

func (slottedPage *SlottedPage) put(marshalledTuple tuple.MarshalledTuple) Slot {
	latestOccupiedSlot := slottedPage.getSlot(slottedPage.slotCount)

	tupleStartingOffset := uint16(slottedPage.pageSize)
	if latestOccupiedSlot == nil {
		tupleStartingOffset = tupleStartingOffset - uint16(marshalledTuple.Size())
	} else {
		tupleStartingOffset = latestOccupiedSlot.tupleOffset - uint16(marshalledTuple.Size())
	}

	copy(slottedPage.buffer[tupleStartingOffset:], marshalledTuple.Buffer())
	return Slot{tupleOffset: tupleStartingOffset, tupleSize: uint16(marshalledTuple.Size())}
}

func (slottedPage *SlottedPage) getSlot(slotNo int) *Slot {
	if slotNo <= 0 {
		return nil
	}
	slot := &Slot{}
	slot.UnMarshalBinary(slottedPage.buffer, slottedPage.slotOffset(slotNo))
	return slot
}

func (slottedPage SlottedPage) slotOffset(slotNo int) int {
	return int(pageIdSize) + (slotNo-1)*int(slotSize)
}

func (slottedPage *SlottedPage) addSlot(slot Slot) {
	offset := slottedPage.slotCount*int(slotSize) + int(pageIdSize)
	copy(slottedPage.buffer[offset:], slot.MarshalBinary())
}

func (slottedPage *SlottedPage) increaseSlotCount() {
	slottedPage.slotCount = slottedPage.slotCount + 1
}

func (slottedPage *SlottedPage) writePageId(id uint32) {
	littleEndian.PutUint32(slottedPage.buffer, id)
}
