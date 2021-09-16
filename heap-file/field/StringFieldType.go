package field

type StringFieldType struct {
}

func (stringFieldType StringFieldType) UnMarshalBinary(buffer []byte) Field {
	length := stringFieldType.readValueLengthFrom(buffer)
	return NewStringField(string(stringFieldType.readValue(buffer, stringValueLengthSize, length)))
}

func (stringFieldType StringFieldType) readValueLengthFrom(buffer []byte) uint16 {
	return littleEndian.Uint16(buffer)
}

func (stringFieldType StringFieldType) readValue(buffer []byte, offset int, length uint16) []byte {
	endOffset := uint16(offset) + length
	return buffer[offset:endOffset]
}
