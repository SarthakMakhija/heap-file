package field

import "unsafe"

const intSize = unsafe.Sizeof(Uint16Field{}.value)

type Uint16Field struct {
	value uint16
}

func NewUint16Field(value uint16) Uint16Field {
	return Uint16Field{
		value: value,
	}
}

func (uint16Field Uint16Field) Value() uint16 {
	return uint16Field.value
}

func (uint16Field Uint16Field) MarshalBinary() []byte {
	buffer := make([]byte, intSize)

	littleEndian.PutUint16(buffer, uint16Field.value)
	return buffer
}
