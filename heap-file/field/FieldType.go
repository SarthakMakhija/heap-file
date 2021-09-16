package field

type FieldType interface {
	UnMarshalBinary([]byte) Field
}
