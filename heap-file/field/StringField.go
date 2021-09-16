package field

import "encoding/binary"

const stringValueLengthSize = 2

var littleEndian = binary.LittleEndian

type StringField struct {
	value []byte
}

func NewStringField(value string) StringField {
	return StringField{
		value: []byte(value),
	}
}

func (stringField StringField) Value() interface{} {
	return string(stringField.value)
}

func (stringField StringField) MarshalBinary() []byte {
	offset := 0
	buffer := make([]byte, stringField.MarshalSize())

	stringField.writeValueLength(buffer)
	offset = offset + stringValueLengthSize
	stringField.writeValueAt(buffer, offset)

	return buffer
}

func (stringField StringField) MarshalSize() int {
	return len(stringField.value) + stringValueLengthSize
}

func (stringField StringField) writeValueLength(buffer []byte) {
	littleEndian.PutUint16(buffer, uint16(len(stringField.value)))
}

func (stringField StringField) writeValueAt(buffer []byte, offset int) int {
	return copy(buffer[offset:], stringField.value)
}
