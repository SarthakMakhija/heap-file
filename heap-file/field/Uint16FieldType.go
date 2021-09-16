package field

type Uint16FieldType struct {
}

func (uint16FieldType Uint16FieldType) UnMarshalBinary(buffer []byte) Field {
	return NewUint16Field(littleEndian.Uint16(buffer))
}
