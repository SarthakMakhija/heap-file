package field

type Uint32FieldType struct {
}

func (uint32FieldType Uint32FieldType) UnMarshalBinary(buffer []byte) Field {
	return NewUint32Field(littleEndian.Uint32(buffer))
}
