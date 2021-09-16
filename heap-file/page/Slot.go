package page

import "unsafe"

var slotSize = unsafe.Sizeof(Slot{})
var tupleOffsetSize = unsafe.Sizeof(Slot{}.tupleOffset)

type Slot struct {
	tupleOffset uint16
	tupleSize   uint16
}

func (slot Slot) MarshalBinary() []byte {
	buffer := make([]byte, slotSize)

	offset := 0
	littleEndian.PutUint16(buffer, slot.tupleOffset)
	offset = offset + 2
	littleEndian.PutUint16(buffer[offset:], slot.tupleSize)

	return buffer
}

func (slot *Slot) UnMarshalBinary(buffer []byte, offset int) {
	slot.tupleOffset = littleEndian.Uint16(buffer[offset:])
	slot.tupleSize = littleEndian.Uint16(buffer[offset+int(tupleOffsetSize):])
}
