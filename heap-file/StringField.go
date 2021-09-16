package heap_file

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

func (stringField StringField) Value() string {
	return string(stringField.value)
}

func (stringField StringField) MarshalBinary() []byte {
	offset := 0
	buffer := make([]byte, len(stringField.value)+stringValueLengthSize)

	stringField.writeValueLength(buffer)
	offset = offset + stringValueLengthSize
	stringField.writeValueAt(buffer, offset)

	return buffer
}

func (stringField *StringField) UnMarshalBinary(buffer []byte) {
	length := stringField.readValueLengthFrom(buffer)
	stringField.value = stringField.readValue(buffer, stringValueLengthSize, length)
}

func (stringField StringField) writeValueLength(buffer []byte) {
	littleEndian.PutUint16(buffer, uint16(len(stringField.value)))
}

func (stringField StringField) readValueLengthFrom(buffer []byte) uint16 {
	return littleEndian.Uint16(buffer)
}

func (stringField StringField) writeValueAt(buffer []byte, offset int) int {
	return copy(buffer[offset:], stringField.value)
}

func (stringField StringField) readValue(buffer []byte, offset int, length uint16) []byte {
	endOffset := uint16(offset) + length
	return buffer[offset:endOffset]
}
