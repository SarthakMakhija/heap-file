package heap_file

type Field interface {
	MarshalBinary() []byte
	UnMarshalBinary([]byte)
}
