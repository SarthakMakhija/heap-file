package tuple

import (
	"encoding/binary"
	"unsafe"
)

var pageIdSize = unsafe.Sizeof(TupleId{}.PageId)
var tupleIdSize = unsafe.Sizeof(pageIdSize * 2) //SlotNo serialized as uint32, hence 2 uint32
var littleEndian = binary.LittleEndian

type TupleId struct {
	PageId uint32
	SlotNo int //serialized as uint32
}

func (tupleId TupleId) MarshalBinary() []byte {
	buffer := make([]byte, tupleIdSize)

	offset := 0
	littleEndian.PutUint32(buffer, tupleId.PageId)
	offset = offset + int(pageIdSize)
	littleEndian.PutUint32(buffer[offset:], uint32(tupleId.SlotNo))

	return buffer
}

func (tupleId *TupleId) UnMarshalBinary(buffer []byte) {
	offset := 0
	pageId := littleEndian.Uint32(buffer)
	offset = offset + int(pageIdSize)
	slotNo := littleEndian.Uint32(buffer[offset:])

	tupleId.PageId = pageId
	tupleId.SlotNo = int(slotNo)
}
