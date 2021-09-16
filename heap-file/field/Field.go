package field

type Field interface {
	MarshalBinary() []byte
	MarshalSize() int
	Value() interface{}
}
