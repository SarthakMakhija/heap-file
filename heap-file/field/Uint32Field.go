package field

import "unsafe"

const intSize = unsafe.Sizeof(Uint32Field{}.value)

type Uint32Field struct {
	value uint32
}

func NewUint32Field(value uint32) Uint32Field {
	return Uint32Field{
		value: value,
	}
}

func (uint32Field Uint32Field) Value() interface{} {
	return uint32Field.value
}

func (uint32Field Uint32Field) MarshalBinary() []byte {
	buffer := make([]byte, uint32Field.MarshalSize())

	littleEndian.PutUint32(buffer, uint32Field.value)
	return buffer
}

func (uint32Field Uint32Field) MarshalSize() int {
	return int(intSize)
}
