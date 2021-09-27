package heap_file

import (
	"encoding/binary"
	"unsafe"
)

var littleEndian = binary.LittleEndian
var pageIdSize = unsafe.Sizeof(OverflowPage{}.id)

var slotSize = unsafe.Sizeof(Slot{}.tupleOffset) + unsafe.Sizeof(Slot{}.tupleSize) + unsafe.Sizeof(Slot{}.overflowPageSlotId)
var tupleOffsetSize = unsafe.Sizeof(Slot{}.tupleOffset)
var tupleSize = unsafe.Sizeof(Slot{}.tupleSize)

type OverflowPage struct {
	id        uint32
	buffer    []byte
	pageSize  int
	slotCount int
	nextPage  *OverflowPage
}

type Slot struct {
	id                 int
	tupleOffset        uint16
	tupleSize          uint16
	overflowPageSlotId uint16
}

func (slot Slot) MarshalBinary() []byte {
	buffer := make([]byte, slotSize)

	offset := 0
	littleEndian.PutUint16(buffer, slot.tupleOffset)
	offset = offset + 2
	littleEndian.PutUint16(buffer[offset:], slot.tupleSize)
	offset = offset + 2
	littleEndian.PutUint16(buffer[offset:], slot.overflowPageSlotId)

	return buffer
}

func (slot *Slot) UnMarshalBinary(buffer []byte, offset int) {
	slot.tupleOffset = littleEndian.Uint16(buffer[offset:])
	slot.tupleSize = littleEndian.Uint16(buffer[offset+int(tupleOffsetSize):])
	slot.overflowPageSlotId = littleEndian.Uint16(buffer[offset+int(tupleOffsetSize)+int(tupleSize):])
}

func NewOverflowPage(id uint32, pageSize int) *OverflowPage {
	overflowPage := &OverflowPage{
		pageSize:  pageSize,
		buffer:    make([]byte, pageSize),
		id:        id,
		slotCount: 0,
	}
	overflowPage.writePageId(id)
	return overflowPage
}

func NewReadonlyOverflowPage(buffer []byte) *OverflowPage {
	overflowPage := &OverflowPage{
		buffer: buffer,
		id:     littleEndian.Uint32(buffer),
	}
	return overflowPage
}

func (currentOverflowPage *OverflowPage) Put(buffer []byte) *Slot {
	isLargeEnoughToHold, sizeAvailable := currentOverflowPage.isLargeEnoughToHold(buffer)
	if isLargeEnoughToHold == 1 {
		slot := currentOverflowPage.put(buffer)
		currentOverflowPage.addSlot(slot)
		currentOverflowPage.increaseSlotCount()
		return slot
	} else if isLargeEnoughToHold == 0 {
		slot := currentOverflowPage.put(buffer[0:sizeAvailable])
		newOverflowPage := NewOverflowPage(currentOverflowPage.id+1, currentOverflowPage.pageSize)
		currentOverflowPage.nextPage = newOverflowPage
		currentOverflowPage.writeOverflowPageId(newOverflowPage.id)

		newOverflowPageSlot := newOverflowPage.Put(buffer[sizeAvailable:])
		slot.overflowPageSlotId = uint16(newOverflowPageSlot.id)

		currentOverflowPage.addSlot(slot)
		currentOverflowPage.increaseSlotCount()
		return slot
	} else {
		newOverflowPage := NewOverflowPage(currentOverflowPage.id+1, currentOverflowPage.pageSize)
		currentOverflowPage.writeOverflowPageId(newOverflowPage.id)
		currentOverflowPage.nextPage = newOverflowPage
		return newOverflowPage.Put(buffer)
	}
}

func (currentOverflowPage *OverflowPage) GetAt(slotNo int) ([]byte, uint16) {
	slot := currentOverflowPage.getSlot(slotNo)
	if slot == nil {
		return []byte{}, 0
	}
	if slot.overflowPageSlotId == 0 {
		return currentOverflowPage.buffer[slot.tupleOffset : slot.tupleOffset+slot.tupleSize], slot.overflowPageSlotId
	} else {
		return currentOverflowPage.buffer[slot.tupleOffset:], slot.overflowPageSlotId
	}
}

func (currentOverflowPage *OverflowPage) NextOverflowPageId() uint32 {
	return littleEndian.Uint32(currentOverflowPage.buffer[pageIdSize:])
}

func (currentOverflowPage OverflowPage) isLargeEnoughToHold(buffer []byte) (int, uint16) {
	sizeAvailable := currentOverflowPage.sizeAvailable()
	if sizeAvailable > uint16(len(buffer)) {
		return 1, sizeAvailable
	}
	if sizeAvailable > 0 {
		return 0, sizeAvailable
	}
	return -1, sizeAvailable
}

func (currentOverflowPage OverflowPage) sizeAvailable() uint16 {
	size := uint16(0)
	for slotNo := 1; slotNo <= currentOverflowPage.slotCount; slotNo++ {
		slot := currentOverflowPage.getSlot(slotNo)
		size = size + slot.tupleSize + uint16(slotSize)
	}
	return uint16(currentOverflowPage.pageSize) - uint16(slotSize) - size - uint16(pageIdSize) - uint16(pageIdSize) //next overflow page id
}

func (currentOverflowPage OverflowPage) getSlot(slotNo int) *Slot {
	if slotNo <= 0 {
		return nil
	}
	slot := &Slot{}
	slot.UnMarshalBinary(currentOverflowPage.buffer, currentOverflowPage.slotOffset(slotNo))
	return slot
}

func (currentOverflowPage OverflowPage) slotOffset(slotNo int) int {
	return int(pageIdSize) + int(pageIdSize) + (slotNo-1)*int(slotSize) //pageId + nextOverflowPageId
}

func (currentOverflowPage *OverflowPage) put(buffer []byte) *Slot {
	latestOccupiedSlot := currentOverflowPage.getSlot(currentOverflowPage.slotCount)

	tupleStartingOffset := uint16(currentOverflowPage.pageSize)
	if latestOccupiedSlot == nil {
		tupleStartingOffset = tupleStartingOffset - uint16(len(buffer))
	} else {
		tupleStartingOffset = latestOccupiedSlot.tupleOffset - uint16(len(buffer))
	}

	copy(currentOverflowPage.buffer[tupleStartingOffset:], buffer)
	return &Slot{id: currentOverflowPage.slotCount + 1, tupleOffset: tupleStartingOffset, tupleSize: uint16(len(buffer))}
}

func (currentOverflowPage *OverflowPage) addSlot(slot *Slot) {
	offset := currentOverflowPage.slotCount*int(slotSize) + int(pageIdSize) + int(pageIdSize) //overflow page id
	copy(currentOverflowPage.buffer[offset:], slot.MarshalBinary())
}

func (currentOverflowPage *OverflowPage) increaseSlotCount() {
	currentOverflowPage.slotCount = currentOverflowPage.slotCount + 1
}

func (currentOverflowPage *OverflowPage) writePageId(id uint32) {
	littleEndian.PutUint32(currentOverflowPage.buffer, id)
}

func (currentOverflowPage *OverflowPage) writeOverflowPageId(id uint32) {
	littleEndian.PutUint32(currentOverflowPage.buffer[pageIdSize:], id)
}
